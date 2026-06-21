package api

import (
	friendapplication "chat-system/server/internal/friend/application"
	friendhttp "chat-system/server/internal/friend/interfaces/http"
)

type Handler = friendhttp.Handler

func NewHandler(friends *friendapplication.Service) *Handler {
	return friendhttp.NewHandler(friends)
}
