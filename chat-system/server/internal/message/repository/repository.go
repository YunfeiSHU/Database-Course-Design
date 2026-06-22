package repository

import messagedomain "chat-system/server/internal/message/domain"

type MessageRepository interface {
	Save(message *messagedomain.Message) error
	FindByID(messageID uint) (*messagedomain.Message, error)
	UpdateStatus(messageID uint, status string) error
	UpdateContentAndStatus(messageID uint, content string, status string) error
	List(userID uint, friendID uint, limit int) ([]messagedomain.Message, error)
}
