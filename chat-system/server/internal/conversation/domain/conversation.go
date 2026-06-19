package domain

import "time"

const TableName = "conversation"

type Conversation struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	PeerID        uint      `json:"peer_id"`
	LastMessageID uint      `json:"last_message_id"`
	Status        string    `json:"status"`
	UpdateTime    time.Time `json:"update_time"`
}

func NewConversation(userID uint, peerID uint, messageID uint, status string) Conversation {
	return Conversation{UserID: userID, PeerID: peerID, LastMessageID: messageID, Status: status}
}

func (c *Conversation) Archive(status string) {
	c.Status = status
}

func (c *Conversation) Delete(status string) {
	c.Status = status
}

func (c *Conversation) Restore(status string) {
	c.Status = status
}
