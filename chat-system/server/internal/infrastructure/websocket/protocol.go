package websocket

import (
	"encoding/json"

	presencedomain "chat-system/server/internal/presence/domain"
)

const (
	TypeChat         = "chat"
	TypeRecall       = "revoke"
	TypeSystem       = "system"
	TypeOnline       = "online"
	TypeOffline      = "offline"
	TypeLoginSuccess = "login_success"
	TypeHeartbeat    = "heartbeat"
)

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

type ChatData struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Content  string `json:"content,omitempty"`
	SendTime string `json:"send_time,omitempty"`
}

type RecallData struct {
	MessageID uint   `json:"message_id"`
	Content   string `json:"content,omitempty"`
	Status    string `json:"status,omitempty"`
}

func Encode(message Message) []byte {
	data, _ := json.Marshal(message)
	return data
}

type SystemData = presencedomain.SystemData

type PresenceData = presencedomain.PresenceData

func NewMessage(messageType string, data interface{}) Message {
	encoded, _ := json.Marshal(data)
	return Message{Type: messageType, Data: encoded}
}
