package api

import (
	messageapplication "chat-system/server/internal/message/application"
	messagehttp "chat-system/server/internal/message/interfaces/http"
)

type Handler = messagehttp.Handler

func NewHandler(messages *messageapplication.Service) *Handler {
	return messagehttp.NewHandler(messages)
}
