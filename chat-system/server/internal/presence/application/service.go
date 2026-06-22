package application

import (
	presencedomain "chat-system/server/internal/presence/domain"
	presencerepository "chat-system/server/internal/presence/repository"
)

type Service struct {
	presenceStore presencerepository.PresenceStore
}

func NewService(presenceStore presencerepository.PresenceStore) *Service {
	return &Service{presenceStore: presenceStore}
}

func (s *Service) MarkOnline(account string, userID uint) presencedomain.Event {
	_ = s.presenceStore.SetOnline(account, userID)
	return presencedomain.Event{Type: "online", Data: presencedomain.PresenceData{Account: account}}
}

func (s *Service) MarkOffline(account string) presencedomain.Event {
	_ = s.presenceStore.SetOffline(account)
	return presencedomain.Event{Type: "offline", Data: presencedomain.PresenceData{Account: account}}
}

func (s *Service) IsOnline(account string) bool {
	return s.presenceStore != nil && s.presenceStore.IsOnline(account)
}

func (s *Service) MessageDelivered(from string, to string, content string) presencedomain.Event {
	return presencedomain.Event{Type: "message_delivered", Data: map[string]string{"from": from, "to": to, "content": content}}
}

func (s *Service) MessageRecalled(from string, to string, messageID uint) presencedomain.Event {
	return presencedomain.Event{Type: "revoke", Data: map[string]interface{}{"from": from, "to": to, "message_id": messageID}}
}

func System(content string) presencedomain.Event {
	return presencedomain.Event{Type: "system", Data: presencedomain.SystemData{Content: content}}
}

func Heartbeat(now string) presencedomain.Event {
	return presencedomain.Event{Type: "heartbeat", Data: map[string]string{"time": now}}
}
