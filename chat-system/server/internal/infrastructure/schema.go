package infrastructure

import "time"

type userRecord struct {
	ID            uint      `gorm:"primaryKey"`
	Account       string    `gorm:"size:32;uniqueIndex;not null"`
	Nickname      string    `gorm:"size:64;not null"`
	Password      string    `gorm:"size:255;not null"`
	CreateTime    time.Time `gorm:"autoCreateTime"`
	LastLoginTime *time.Time
}

func (userRecord) TableName() string {
	return "user"
}

type friendRecord struct {
	ID       uint       `gorm:"primaryKey"`
	UserID   uint       `gorm:"uniqueIndex:idx_user_friend;not null"`
	FriendID uint       `gorm:"uniqueIndex:idx_user_friend;not null"`
	Status   string     `gorm:"size:16;not null;default:accepted;index"`
	User     userRecord `gorm:"foreignKey:UserID"`
	Friend   userRecord `gorm:"foreignKey:FriendID"`
}

func (friendRecord) TableName() string {
	return "friend"
}

type messageRecord struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   uint      `gorm:"index;not null"`
	ReceiverID uint      `gorm:"index;not null"`
	Content    string    `gorm:"type:text;not null"`
	Status     string    `gorm:"size:16;not null;default:created;index"`
	SendTime   time.Time `gorm:"autoCreateTime"`
}

func (messageRecord) TableName() string {
	return "message"
}

type conversationRecord struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"uniqueIndex:idx_user_peer;not null"`
	PeerID        uint      `gorm:"uniqueIndex:idx_user_peer;not null"`
	LastMessageID uint      `gorm:"index;not null"`
	Status        string    `gorm:"size:16;not null;default:normal;index"`
	UpdateTime    time.Time `gorm:"autoUpdateTime"`
}

func (conversationRecord) TableName() string {
	return "conversation"
}
