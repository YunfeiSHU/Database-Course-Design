package websocket

import "encoding/json"

const (
	TypeChat         = "chat"
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

type SystemData struct {
	Content string `json:"content"`
}

type PresenceData struct {
	Account string `json:"account"`
}

func Encode(message Message) []byte {
	data, _ := json.Marshal(message)
	return data
}

func NewMessage(messageType string, data interface{}) Message {
	encoded, _ := json.Marshal(data)
	return Message{Type: messageType, Data: encoded}
}
