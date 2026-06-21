package api

import (
	conversationapplication "chat-system/server/internal/conversation/application"
	conversationhttp "chat-system/server/internal/conversation/interfaces/http"
)

type Handler = conversationhttp.Handler

func NewHandler(conversations *conversationapplication.Service) *Handler {
	return conversationhttp.NewHandler(conversations)
}
