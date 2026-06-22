package application

import (
	"errors"
	"time"

	"chat-system/server/internal/common"
	messagedomain "chat-system/server/internal/message/domain"
	messagerepository "chat-system/server/internal/message/repository"
	presencedomain "chat-system/server/internal/presence/domain"
	userdomain "chat-system/server/internal/user/domain"
)

var ErrNotFriends = errors.New("receiver is not your friend")

type UserFinder interface {
	FindByAccount(account string) (*userdomain.User, error)
}

type FriendChecker interface {
	AreFriends(userID uint, friendID uint) (bool, error)
}

type ConversationUpdater interface {
	MarkConversationUpdated(senderID uint, receiverID uint, messageID uint) error
}

type Notifier interface {
	MessageDelivered(from string, to string, content string) presencedomain.Event
	MessageRecalled(from string, to string, messageID uint) presencedomain.Event
}

type Service struct {
	repository    messagerepository.MessageRepository
	users         UserFinder
	friends       FriendChecker
	conversations ConversationUpdater
	notifier      Notifier
}

type DeliveredMessage struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Content  string `json:"content,omitempty"`
	SendTime string `json:"send_time,omitempty"`
	Status   string `json:"status,omitempty"`
}

type RecalledMessage struct {
	ID      uint   `json:"id"`
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Status  string `json:"status,omitempty"`
	Content string `json:"content,omitempty"`
}

var ErrCannotRecall = errors.New("message cannot be recalled")

func NewService(repository messagerepository.MessageRepository, users UserFinder, friends FriendChecker, conversations ConversationUpdater, notifier Notifier) *Service {
	return &Service{repository: repository, users: users, friends: friends, conversations: conversations, notifier: notifier}
}

func (s *Service) ListHistory(userID uint, friendID uint, limit int) ([]messagedomain.Message, error) {
	return s.repository.List(userID, friendID, limit)
}

func (s *Service) FindByID(messageID uint) (*messagedomain.Message, error) {
	return s.repository.FindByID(messageID)
}

func (s *Service) ListHistoryByAccount(userID uint, account string, limit int) ([]messagedomain.Message, error) {
	friend, err := s.users.FindByAccount(account)
	if err != nil {
		return nil, err
	}
	return s.ListHistory(userID, friend.ID, limit)
}

func (s *Service) Send(senderID uint, senderAccount string, receiverAccount string, content string) (*DeliveredMessage, error) {
	receiver, err := s.users.FindByAccount(receiverAccount)
	if err != nil {
		return nil, err
	}
	accepted, err := s.friends.AreFriends(senderID, receiver.ID)
	if err != nil {
		return nil, err
	}
	if !accepted {
		return nil, ErrNotFriends
	}
	message := messagedomain.NewMessage(senderID, receiver.ID, content, common.MessageStatusCreated)
	message.MarkSending(common.MessageStatusSending)
	if err := s.repository.Save(&message); err != nil {
		return nil, err
	}
	message.MarkDelivered(common.MessageStatusDelivered)
	if err := s.repository.UpdateStatus(message.ID, message.Status); err != nil {
		return nil, err
	}
	if err := s.conversations.MarkConversationUpdated(senderID, receiver.ID, message.ID); err != nil {
		return nil, err
	}
	sendTime := message.SendTime
	if sendTime.IsZero() {
		sendTime = time.Now()
	}
	_ = s.notifier.MessageDelivered(senderAccount, receiverAccount, content)
	return &DeliveredMessage{From: senderAccount, To: receiverAccount, Content: content, SendTime: sendTime.Format(time.RFC3339), Status: message.Status}, nil
}

func (s *Service) Recall(senderID uint, senderAccount string, messageID uint) (*RecalledMessage, error) {
	message, err := s.repository.FindByID(messageID)
	if err != nil {
		return nil, err
	}
	if message.SenderID != senderID {
		return nil, ErrCannotRecall
	}
	if message.Status == common.MessageStatusRecalled {
		return nil, ErrCannotRecall
	}
	if message.Status != common.MessageStatusCreated && message.Status != common.MessageStatusSending && message.Status != common.MessageStatusDelivered {
		return nil, ErrCannotRecall
	}
	recalledContent := "[消息已撤回]"
	message.Recall(common.MessageStatusRecalled)
	message.Content = recalledContent
	if err := s.repository.UpdateContentAndStatus(message.ID, message.Content, message.Status); err != nil {
		return nil, err
	}
	if err := s.conversations.MarkConversationUpdated(message.SenderID, message.ReceiverID, message.ID); err != nil {
		return nil, err
	}
	_ = s.notifier.MessageRecalled(senderAccount, "", message.ID)
	return &RecalledMessage{ID: message.ID, From: senderAccount, To: "", Status: message.Status, Content: recalledContent}, nil
}
