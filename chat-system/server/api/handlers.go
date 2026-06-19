package api

import (
	conversationapi "chat-system/server/internal/conversation/api"
	friendapi "chat-system/server/internal/friend/api"
	messageapi "chat-system/server/internal/message/api"
	userapi "chat-system/server/internal/user/api"
)

type Handlers struct {
	Users         *userapi.Handler
	Friends       *friendapi.Handler
	Conversations *conversationapi.Handler
	Messages      *messageapi.Handler
}
