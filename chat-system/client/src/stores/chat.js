import { defineStore } from 'pinia'
import { ElNotification } from 'element-plus'
import { approveFriendRequest, getConversationHistory, getConversationList, getFriendList, getFriendRequestList, login, register, requestFriend } from '../api/http'
import { createChatSocket } from '../websocket/chatSocket'

const MESSAGE_STATUS_TEXT = {
  created: '已创建',
  sending: '发送中',
  delivered: '已送达',
  read: '已读',
  recalled: '已撤回'
}

function normalizeMessage(item, selfAccount, peer) {
  const peerID = peer?.peer_id || peer?.friend_id
  const peerAccount = peer?.peer?.account || peer?.friend?.account
  const from = item.from || (item.sender_id === peerID ? peerAccount : selfAccount)
  const to = item.to || (item.receiver_id === peerID ? peerAccount : selfAccount)
  return {
    from,
    to,
    content: item.content,
    send_time: item.send_time,
    status: item.status || 'delivered',
    status_text: MESSAGE_STATUS_TEXT[item.status] || MESSAGE_STATUS_TEXT.delivered
  }
}

function conversationPeer(conversation) {
  return conversation?.peer || conversation?.friend || null
}

function conversationPeerID(conversation) {
  return conversation?.peer_id || conversation?.friend_id || conversationPeer(conversation)?.id
}

function conversationPeerAccount(conversation) {
  return conversationPeer(conversation)?.account || ''
}

function conversationPeerNickname(conversation) {
  return conversationPeer(conversation)?.nickname || conversationPeerAccount(conversation)
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    token: localStorage.getItem('chat_token') || '',
    account: localStorage.getItem('chat_account') || '',
    nickname: localStorage.getItem('chat_nickname') || '',
    rememberedAccount: localStorage.getItem('remembered_account') || '',
    lastLoginTime: localStorage.getItem('last_login_time') || '',
    friends: [],
    friendRequests: [],
    conversations: [],
    currentConversation: null,
    messages: {},
    socket: null,
    connected: false
  }),
  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    currentMessages: (state) => {
      const peerAccount = conversationPeerAccount(state.currentConversation)
      if (!peerAccount) return []
      return state.messages[peerAccount] || []
    },
    currentPeer: (state) => conversationPeer(state.currentConversation),
    currentPeerOnline: (state) => {
      const peerAccount = conversationPeerAccount(state.currentConversation)
      const friend = state.friends.find((item) => item.friend.account === peerAccount)
      return Boolean(friend?.online)
    }
  },
  actions: {
    async registerUser(form) {
      return register({
        nickname: form.nickname,
        password: form.password,
        confirm_password: form.confirmPassword
      })
    },
    async loginUser(form) {
      const data = await login({ account: form.account, password: form.password })
      this.token = data.token
      this.account = data.account
      this.nickname = data.nickname
      this.lastLoginTime = data.last_login_time || ''
      localStorage.setItem('chat_token', data.token)
      localStorage.setItem('chat_account', data.account)
      localStorage.setItem('chat_nickname', data.nickname)
      localStorage.setItem('last_login_time', this.lastLoginTime)
      if (form.remember) {
        this.rememberedAccount = form.account
        localStorage.setItem('remembered_account', form.account)
      } else {
        this.rememberedAccount = ''
        localStorage.removeItem('remembered_account')
      }
      await this.refreshFriends()
      await this.refreshFriendRequests()
      await this.refreshConversations()
      this.connectSocket()
    },
    logout() {
      this.socket?.close()
      this.socket = null
      this.connected = false
      this.token = ''
      this.account = ''
      this.nickname = ''
      this.currentConversation = null
      this.friends = []
      this.friendRequests = []
      this.conversations = []
      this.messages = {}
      localStorage.removeItem('chat_token')
      localStorage.removeItem('chat_account')
      localStorage.removeItem('chat_nickname')
    },
    async refreshFriends() {
      this.friends = await getFriendList()
    },
    async refreshConversations() {
      this.conversations = await getConversationList()
      if (this.currentConversation) {
        const peerAccount = conversationPeerAccount(this.currentConversation)
        this.currentConversation = this.conversations.find((item) => conversationPeerAccount(item) === peerAccount) || this.currentConversation
      }
    },
    async refreshFriendRequests() {
      this.friendRequests = await getFriendRequestList()
    },
    async addFriendByAccount(account) {
      await requestFriend(account)
    },
    async acceptFriendRequest(id) {
      await approveFriendRequest(id)
      await this.refreshFriendRequests()
      await this.refreshFriends()
      await this.refreshConversations()
    },
    async startConversation(friend) {
      const existing = this.conversations.find((item) => conversationPeerAccount(item) === friend.friend.account)
      if (existing) {
        await this.selectConversation(existing)
        return
      }
      const conversation = {
        id: `friend-${friend.friend_id}`,
        user_id: this.account,
        peer_id: friend.friend_id,
        peer: friend.friend,
        status: 'normal'
      }
      this.currentConversation = conversation
      await this.loadHistory(conversation)
    },
    async selectConversation(conversation) {
      this.currentConversation = conversation
      const peerAccount = conversationPeerAccount(conversation)
      const friend = this.friends.find((item) => item.friend.account === peerAccount)
      if (friend) friend.unread = 0
      await this.loadHistory(conversation)
    },
    async loadHistory(conversation = this.currentConversation) {
      const peerID = conversationPeerID(conversation)
      const peerAccount = conversationPeerAccount(conversation)
      if (!peerID || !peerAccount) return
      const rows = await getConversationHistory({ friend_id: peerID })
      this.messages[peerAccount] = rows.map((item) => normalizeMessage(item, this.account, conversation))
    },
    connectSocket() {
      if (!this.token || this.socket) return
      this.socket = createChatSocket(this.token, {
        onOpen: () => { this.connected = true },
        onClose: () => { this.connected = false; this.socket = null },
        onMessage: (message) => this.handleSocketMessage(message),
        onError: () => { this.connected = false }
      })
    },
    sendMessage(content) {
      const peerAccount = conversationPeerAccount(this.currentConversation)
      if (!peerAccount || !content.trim()) return false
      return this.socket?.send('chat', {
        from: this.account,
        to: peerAccount,
        content: content.trim()
      })
    },
    handleSocketMessage(message) {
      const data = message.data || {}
      if (message.type === 'chat') {
        const peer = data.from === this.account ? data.to : data.from
        if (!this.messages[peer]) this.messages[peer] = []
        this.messages[peer].push(normalizeMessage(data, this.account, this.friends.find((item) => item.friend.account === peer)))
        const friend = this.friends.find((item) => item.friend.account === peer)
        if (friend && conversationPeerAccount(this.currentConversation) !== peer) {
          friend.unread = (friend.unread || 0) + 1
        }
        this.refreshConversations()
      }
      if (message.type === 'system') {
        ElNotification({ title: '系统消息', message: data.content || '系统消息', type: 'warning' })
      }
      if (message.type === 'message_delivered') {
        ElNotification({ title: '消息状态', message: '消息已送达', type: 'success' })
      }
      if (message.type === 'online' || message.type === 'offline') {
        const friend = this.friends.find((item) => item.friend.account === data.account)
        if (friend) friend.online = message.type === 'online'
      }
    }
  }
})

export { conversationPeer, conversationPeerAccount, conversationPeerNickname }
