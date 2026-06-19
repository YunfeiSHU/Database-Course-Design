package repository

import (
	"errors"
	"time"

	userdomain "chat-system/server/internal/user/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user *userdomain.User) error
	FindByID(userID uint) (*userdomain.User, error)
	FindByAccount(account string) (*userdomain.User, error)
	UpdateLastLogin(userID uint) error
}

type SessionStore interface {
	Set(token string, value string, ttl time.Duration) error
	Get(token string) (string, error)
}
