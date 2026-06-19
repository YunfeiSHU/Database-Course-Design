package domain

import userdomain "chat-system/server/internal/user/domain"

type Friend struct {
	ID       uint            `json:"id"`
	UserID   uint            `json:"user_id"`
	FriendID uint            `json:"friend_id"`
	Status   string          `json:"status"`
	User     userdomain.User `json:"user"`
	Friend   userdomain.User `json:"friend"`
}

const TableName = "friend"
