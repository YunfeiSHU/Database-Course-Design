export function createChatSocket(token, handlers = {}) {
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const socket = new WebSocket(`${protocol}://${window.location.host}/ws?token=${encodeURIComponent(token)}`)
  let heartbeatTimer = null

  socket.addEventListener('open', () => {
    handlers.onOpen?.()
    heartbeatTimer = window.setInterval(() => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'heartbeat', data: {} }))
      }
    }, 30000)
  })

  socket.addEventListener('message', (event) => {
    try {
      handlers.onMessage?.(JSON.parse(event.data))
    } catch (error) {
      handlers.onError?.(error)
    }
  })

  socket.addEventListener('close', () => {
    if (heartbeatTimer) window.clearInterval(heartbeatTimer)
    handlers.onClose?.()
  })

  socket.addEventListener('error', (event) => handlers.onError?.(event))

  return {
    send(type, data) {
      if (socket.readyState !== WebSocket.OPEN) return false
      socket.send(JSON.stringify({ type, data }))
      return true
    },
    close() {
      socket.close()
    }
  }
}
