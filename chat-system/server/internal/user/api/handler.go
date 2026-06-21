package api

import (
	userapplication "chat-system/server/internal/user/application"
	userhttp "chat-system/server/internal/user/interfaces/http"
)

type Handler = userhttp.Handler

func NewHandler(users *userapplication.Service) *Handler {
	return userhttp.NewHandler(users)
}
