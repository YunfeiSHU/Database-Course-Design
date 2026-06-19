package domain

import "time"

type User struct {
	ID            uint       `json:"id"`
	Account       string     `json:"account"`
	Nickname      string     `json:"nickname"`
	Password      string     `json:"-"`
	CreateTime    time.Time  `json:"create_time" gorm:"autoCreateTime"`
	LastLoginTime *time.Time `json:"last_login_time"`
}

const TableName = "user"
