package api

import (
	conversationhttp "chat-system/server/internal/conversation/interfaces/http"
	friendhttp "chat-system/server/internal/friend/interfaces/http"
	messagehttp "chat-system/server/internal/message/interfaces/http"
	userhttp "chat-system/server/internal/user/interfaces/http"
)

// Handlers 汇总各模块的 HTTP 处理器，供路由层统一注入。
type Handlers struct {
	Users         *userhttp.Handler
	Friends       *friendhttp.Handler
	Conversations *conversationhttp.Handler
	Messages      *messagehttp.Handler
}
