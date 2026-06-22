<template>
  <main class="app-shell">
    <section v-if="!store.isLoggedIn" class="auth-view">
      <div class="auth-panel">
        <div class="auth-header">
          <h1>Chat System</h1>
          <span>Web 即时通信</span>
        </div>
        <el-tabs v-model="authMode" stretch>
          <el-tab-pane label="登录" name="login">
            <el-form :model="loginForm" label-position="top" @submit.prevent>
              <el-form-item label="账号">
                <el-input v-model="loginForm.account" placeholder="10000001" clearable />
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="loginForm.password" type="password" show-password @keyup.enter="handleLogin" />
              </el-form-item>
              <div class="auth-row">
                <el-checkbox v-model="loginForm.remember">记住账号</el-checkbox>
                <span v-if="store.lastLoginTime" class="last-login">上次登录 {{ formatTime(store.lastLoginTime) }}</span>
              </div>
              <el-button type="primary" :loading="loading" class="full-button" @click="handleLogin">登录</el-button>
            </el-form>
          </el-tab-pane>
          <el-tab-pane label="注册" name="register">
            <el-form :model="registerForm" label-position="top" @submit.prevent>
              <el-form-item label="昵称">
                <el-input v-model="registerForm.nickname" maxlength="24" clearable />
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="registerForm.password" type="password" show-password />
              </el-form-item>
              <el-form-item label="确认密码">
                <el-input v-model="registerForm.confirmPassword" type="password" show-password @keyup.enter="handleRegister" />
              </el-form-item>
              <el-button type="primary" :loading="loading" class="full-button" @click="handleRegister">注册并生成账号</el-button>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </section>

    <section v-else class="chat-view">
      <aside class="sidebar friend-sidebar">
        <header class="profile-bar">
          <div>
            <strong>{{ store.nickname }}</strong>
            <span>{{ store.account }}</span>
          </div>
          <el-button size="small" plain @click="store.logout">退出</el-button>
        </header>

        <div class="panel-toolbar">
          <span>好友</span>
          <el-button size="small" type="primary" plain @click="addFriendVisible = true">申请好友</el-button>
        </div>

        <section v-if="store.friendRequests.length" class="request-panel">
          <div class="request-title">好友申请</div>
          <div v-for="request in store.friendRequests" :key="request.id" class="request-item">
            <span class="friend-main">
              <strong>{{ request.user.nickname }}</strong>
              <small>{{ request.user.account }}</small>
            </span>
            <el-button size="small" type="primary" @click="handleAcceptRequest(request.id)">同意</el-button>
          </div>
        </section>

        <el-scrollbar class="friend-scroll">
          <button
            v-for="friend in store.friends"
            :key="friend.friend.account"
            class="list-item"
            @click="startConversation(friend)"
          >
            <span class="presence" :class="friend.online ? 'online' : 'offline'"></span>
            <span class="item-main">
              <strong>{{ friend.friend.nickname }}</strong>
              <small>{{ friend.online ? '在线' : '离线' }} · {{ friend.friend.account }}</small>
            </span>
            <el-badge v-if="friend.unread" :value="friend.unread" />
          </button>
        </el-scrollbar>
      </aside>

      <aside class="sidebar conversation-sidebar">
        <div class="panel-toolbar conversation-toolbar">
          <span>对话</span>
          <small>{{ store.connected ? '连接正常' : '连接中断' }}</small>
        </div>

        <el-scrollbar class="conversation-list-scroll">
          <button
            v-for="conversation in store.conversations"
            :key="conversation.id"
            class="list-item conversation-item"
            :class="{ active: peerAccount(conversation) === peerAccount(store.currentConversation) }"
            @click="selectConversation(conversation)"
          >
            <span class="item-avatar">{{ avatarText(conversation) }}</span>
            <span class="item-main">
              <strong>{{ peerName(conversation) }}</strong>
              <small>{{ conversationSubtitle(conversation) }}</small>
            </span>
            <time>{{ shortTime(conversation.update_time) }}</time>
          </button>
          <div v-if="!store.conversations.length" class="empty-list">暂无对话</div>
        </el-scrollbar>
      </aside>

      <section class="conversation">
        <header class="conversation-header">
          <div v-if="store.currentConversation">
            <h2>{{ peerName(store.currentConversation) }}</h2>
            <span>{{ store.currentPeerOnline ? '在线' : '离线' }} · {{ peerAccount(store.currentConversation) }}</span>
          </div>
          <div v-else>
            <h2>选择对话开始聊天</h2>
            <span>{{ store.connected ? '连接正常' : '连接中断' }}</span>
          </div>
          <el-button :disabled="!store.currentConversation" @click="loadHistory">刷新记录</el-button>
        </header>

        <el-scrollbar ref="messageScroll" class="message-scroll">
          <div class="message-list">
            <div
              v-for="(message, index) in store.currentMessages"
              :key="message.id || index"
              class="message-row"
              :class="{ mine: message.from === store.account }"
            >
              <div class="message-bubble">
                <time>{{ formatTime(message.send_time) }}</time>
                <p>{{ message.content }}</p>
                <small v-if="message.from === store.account" class="message-status">{{ message.status_text }}</small>
                <el-button
                  v-if="message.from === store.account && message.status !== 'recalled'"
                  link
                  type="danger"
                  size="small"
                  class="message-recall"
                  @click="handleRecallMessage(message.id)"
                >
                  撤回
                </el-button>
              </div>
            </div>
          </div>
        </el-scrollbar>

        <footer class="composer">
          <el-input v-model="draft" :disabled="!store.currentConversation" placeholder="输入消息" @keyup.enter="sendMessage" />
          <el-button type="primary" :disabled="!store.currentConversation" @click="sendMessage">发送</el-button>
        </footer>
      </section>
    </section>

    <el-dialog v-model="addFriendVisible" title="申请好友" width="360px">
      <el-input v-model="friendAccount" placeholder="好友账号，例如 10000002" />
      <template #footer>
        <el-button @click="addFriendVisible = false">取消</el-button>
        <el-button type="primary" @click="handleAddFriend">发送申请</el-button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { nextTick, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { conversationPeerAccount, conversationPeerNickname, useChatStore } from './stores/chat'

const store = useChatStore()
const authMode = ref('login')
const loading = ref(false)
const addFriendVisible = ref(false)
const friendAccount = ref('')
const draft = ref('')
const messageScroll = ref(null)

const loginForm = reactive({
  account: store.rememberedAccount,
  password: '',
  remember: Boolean(store.rememberedAccount)
})
const registerForm = reactive({ nickname: '', password: '', confirmPassword: '' })

if (store.isLoggedIn) {
  Promise.all([store.refreshFriends(), store.refreshFriendRequests(), store.refreshConversations()]).then(() => store.connectSocket())
}

async function handleRegister() {
  if (!registerForm.nickname || !registerForm.password || !registerForm.confirmPassword) {
    ElMessage.warning('请填写完整注册信息')
    return
  }
  loading.value = true
  try {
    const data = await store.registerUser(registerForm)
    authMode.value = 'login'
    loginForm.account = data.account
    ElMessageBox.alert(`账号：${data.account}`, '注册成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}

async function handleLogin() {
  if (!loginForm.account || !loginForm.password) {
    ElMessage.warning('请输入账号和密码')
    return
  }
  loading.value = true
  try {
    await store.loginUser(loginForm)
    ElMessage.success('登录成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}

async function startConversation(friend) {
  await store.startConversation(friend)
  scrollToBottom()
}

async function selectConversation(conversation) {
  await store.selectConversation(conversation)
  scrollToBottom()
}

async function loadHistory() {
  if (!store.currentConversation) return
  await store.loadHistory()
  scrollToBottom()
}

function sendMessage() {
  if (store.sendMessage(draft.value)) {
    draft.value = ''
    scrollToBottom()
  }
}

function handleRecallMessage(messageID) {
  if (!messageID) return
  store.recallMessage(messageID)
    .then(() => scrollToBottom())
    .catch((error) => {
      ElMessage.error(error.response?.data?.error || '撤回失败')
    })
}

async function handleAddFriend() {
  if (!friendAccount.value.trim()) return
  try {
    await store.addFriendByAccount(friendAccount.value.trim())
    addFriendVisible.value = false
    friendAccount.value = ''
    ElMessage.success('申请已发送，等待对方同意')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '申请失败')
  }
}

async function handleAcceptRequest(id) {
  try {
    await store.acceptFriendRequest(id)
    ElMessage.success('已添加为好友')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

function scrollToBottom() {
  nextTick(() => {
    const wrap = messageScroll.value?.wrapRef
    if (wrap) wrap.scrollTop = wrap.scrollHeight
  })
}

function formatTime(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', { hour12: false })
}

function peerAccount(conversation) {
  return conversationPeerAccount(conversation)
}

function peerName(conversation) {
  return conversationPeerNickname(conversation)
}

function avatarText(conversation) {
  return peerName(conversation).slice(0, 1).toUpperCase()
}

function conversationSubtitle(conversation) {
  return conversation.last_message?.content || '暂无消息'
}

function shortTime(value) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', hour12: false })
}
</script>
