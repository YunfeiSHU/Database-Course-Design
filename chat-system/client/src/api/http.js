import axios from 'axios'

const http = axios.create({
  baseURL: '',
  timeout: 10000
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('chat_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export function register(payload) {
  return http.post('/api/register', payload).then((res) => res.data)
}

export function login(payload) {
  return http.post('/api/login', payload).then((res) => res.data)
}

export function listFriends() {
  return http.get('/api/friends').then((res) => res.data)
}

export function listConversations() {
  return http.get('/api/conversations').then((res) => res.data)
}

export function addFriend(account) {
  return http.post('/api/friends', { account }).then((res) => res.data)
}

export function listFriendRequests() {
  return http.get('/api/friend-requests').then((res) => res.data)
}

export function acceptFriendRequest(id) {
  return http.post(`/api/friend-requests/${id}/accept`).then((res) => res.data)
}

export function getHistory(params) {
  return http.get('/api/history', { params }).then((res) => res.data)
}
