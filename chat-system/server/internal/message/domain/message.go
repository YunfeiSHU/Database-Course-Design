package domain

import "time"

type Message struct {
	ID         uint      `json:"id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	Status     string    `json:"status"`
	SendTime   time.Time `json:"send_time"`
}

const TableName = "message"

func NewMessage(senderID uint, receiverID uint, content string, status string) Message {
	return Message{SenderID: senderID, ReceiverID: receiverID, Content: content, Status: status, SendTime: time.Now()}
}

func (m *Message) MarkSending(status string) {
	m.Status = status
}

func (m *Message) MarkDelivered(status string) {
	m.Status = status
}

func (m *Message) MarkRead(status string) {
	m.Status = status
}

func (m *Message) Recall(status string) {
	m.Status = status
}
