package application

import (
	"fmt"
	"time"

	"chat-system/server/internal/common"
	conversationdomain "chat-system/server/internal/conversation/domain"
	conversationrepository "chat-system/server/internal/conversation/repository"
	messagedomain "chat-system/server/internal/message/domain"
	userdomain "chat-system/server/internal/user/domain"
)

type UserProvider interface {
	FindByID(userID uint) (*userdomain.User, error)
}

type MessageProvider interface {
	ListHistory(userID uint, friendID uint, limit int) ([]messagedomain.Message, error)
	FindByID(messageID uint) (*messagedomain.Message, error)
}

type Service struct {
	repository conversationrepository.ConversationRepository
	users      UserProvider
	messages   MessageProvider
}

type Item struct {
	ID            uint                   `json:"id"`
	UserID        uint                   `json:"user_id"`
	PeerID        uint                   `json:"peer_id"`
	LastMessageID uint                   `json:"last_message_id"`
	Status        string                 `json:"status"`
	UpdateTime    time.Time              `json:"update_time"`
	Peer          userdomain.User        `json:"peer"`
	LastMessage   *messagedomain.Message `json:"last_message,omitempty"`
}

func NewService(repository conversationrepository.ConversationRepository, users UserProvider, messages MessageProvider) *Service {
	return &Service{repository: repository, users: users, messages: messages}
}

func (s *Service) TouchMessage(senderID uint, receiverID uint, messageID uint) error {
	now := time.Now()
	if err := s.repository.Upsert(conversationdomain.Conversation{
		UserID:        senderID,
		PeerID:        receiverID,
		LastMessageID: messageID,
		Status:        common.ConversationStatusNormal,
		UpdateTime:    now,
	}); err != nil {
		return err
	}
	return s.repository.Upsert(conversationdomain.Conversation{
		UserID:        receiverID,
		PeerID:        senderID,
		LastMessageID: messageID,
		Status:        common.ConversationStatusNormal,
		UpdateTime:    now,
	})
}

func (s *Service) List(userID uint) ([]Item, error) {
	conversations, err := s.repository.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	items := make([]Item, 0, len(conversations))
	for _, row := range conversations {
		peer, err := s.users.FindByID(row.PeerID)
		if err != nil {
			return nil, fmt.Errorf("find conversation peer %d: %w", row.PeerID, err)
		}
		item := Item{
			ID:            row.ID,
			UserID:        row.UserID,
			PeerID:        row.PeerID,
			LastMessageID: row.LastMessageID,
			Status:        row.Status,
			UpdateTime:    row.UpdateTime,
			Peer:          *peer,
		}
		if row.LastMessageID != 0 && s.messages != nil {
			lastMessage, err := s.messages.FindByID(row.LastMessageID)
			if err != nil {
				return nil, fmt.Errorf("find conversation last message %d: %w", row.LastMessageID, err)
			}
			item.LastMessage = lastMessage
		}
		items = append(items, item)
	}
	return items, nil
}
