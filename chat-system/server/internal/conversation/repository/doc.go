package repository

import conversationdomain "chat-system/server/internal/conversation/domain"

type ConversationRepository interface {
	Upsert(conversation conversationdomain.Conversation) error
	ListByUser(userID uint) ([]conversationdomain.Conversation, error)
}
