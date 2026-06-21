package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	messageapplication "chat-system/server/internal/message/application"
	presencedomain "chat-system/server/internal/presence/domain"
	userapplication "chat-system/server/internal/user/application"
)

type ClientMessage struct {
	Client  *Client
	Message Message
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Inbound    chan ClientMessage
	sessions   SessionParser
	messages   MessageSender
	notices    NoticeService
}

type SessionParser interface {
	ParseSession(token string) (*userapplication.Session, error)
}

type MessageSender interface {
	Send(senderID uint, senderAccount string, receiverAccount string, content string) (*messageapplication.DeliveredMessage, error)
}

type NoticeService interface {
	MarkOnline(account string, userID uint) presencedomain.Event
	MarkOffline(account string) presencedomain.Event
}

func NewHub(sessions SessionParser, messages MessageSender, notices NoticeService) *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Inbound:    make(chan ClientMessage, 256),
		sessions:   sessions,
		messages:   messages,
		notices:    notices,
	}
}

func (h *Hub) ParseSession(token string) (*userapplication.Session, error) {
	return h.sessions.ParseSession(token)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.Account] = client
			h.broadcastEvent(h.notices.MarkOnline(client.Account, client.UserID))
		case client := <-h.Unregister:
			if current, ok := h.Clients[client.Account]; ok && current == client {
				delete(h.Clients, client.Account)
				close(client.Send)
				h.broadcastEvent(h.notices.MarkOffline(client.Account))
			}
		case inbound := <-h.Inbound:
			h.handleMessage(inbound.Client, inbound.Message)
		}
	}
}

func (h *Hub) handleMessage(client *Client, msg Message) {
	switch msg.Type {
	case TypeChat:
		var data ChatData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("websocket message decode failed for account %q: %v", client.Account, err)
			client.Send <- Encode(NewMessage(TypeSystem, SystemData{Content: "消息格式错误"}))
			return
		}
		if data.To == "" || data.Content == "" {
			return
		}
		delivered, err := h.messages.Send(client.UserID, client.Account, data.To, data.Content)
		if err != nil {
			log.Printf("websocket send message failed from %q to %q: %v", client.Account, data.To, err)
			client.Send <- Encode(NewMessage(TypeSystem, SystemData{Content: err.Error()}))
			return
		}
		out := NewMessage(TypeChat, delivered)
		if target, ok := h.Clients[data.To]; ok {
			target.Send <- Encode(out)
		}
		client.Send <- Encode(out)
	case TypeSystem:
		var data SystemData
		_ = json.Unmarshal(msg.Data, &data)
		h.broadcast(NewMessage(TypeSystem, data))
	case TypeHeartbeat:
		client.Send <- Encode(NewMessage(TypeHeartbeat, map[string]string{"time": time.Now().Format(time.RFC3339)}))
	default:
		client.Send <- Encode(NewMessage(TypeSystem, SystemData{Content: fmt.Sprintf("unsupported message type: %s", msg.Type)}))
	}
}

func (h *Hub) broadcastEvent(event presencedomain.Event) {
	h.broadcast(NewMessage(event.Type, event.Data))
}

func (h *Hub) broadcast(msg Message) {
	data := Encode(msg)
	for _, client := range h.Clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(h.Clients, client.Account)
		}
	}
}