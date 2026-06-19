package application

import (
	notificationdomain "chat-system/server/internal/notification/domain"
	notificationrepository "chat-system/server/internal/notification/repository"
)

type Service struct {
	presenceStore notificationrepository.PresenceStore
}

func NewService(presenceStore notificationrepository.PresenceStore) *Service {
	return &Service{presenceStore: presenceStore}
}

func (s *Service) MarkOnline(account string, userID uint) notificationdomain.Event {
	_ = s.presenceStore.SetOnline(account, userID)
	return notificationdomain.Event{Type: "online", Data: notificationdomain.PresenceData{Account: account}}
}

func (s *Service) MarkOffline(account string) notificationdomain.Event {
	_ = s.presenceStore.SetOffline(account)
	return notificationdomain.Event{Type: "offline", Data: notificationdomain.PresenceData{Account: account}}
}

func (s *Service) IsOnline(account string) bool {
	return s.presenceStore != nil && s.presenceStore.IsOnline(account)
}

func (s *Service) MessageDelivered(from string, to string, content string) notificationdomain.Event {
	return notificationdomain.Event{Type: "message_delivered", Data: map[string]string{"from": from, "to": to, "content": content}}
}

func System(content string) notificationdomain.Event {
	return notificationdomain.Event{Type: "system", Data: notificationdomain.SystemData{Content: content}}
}

func Heartbeat(now string) notificationdomain.Event {
	return notificationdomain.Event{Type: "heartbeat", Data: map[string]string{"time": now}}
}
