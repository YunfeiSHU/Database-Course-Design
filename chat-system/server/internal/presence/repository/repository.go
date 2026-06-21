package repository

type PresenceStore interface {
	SetOnline(account string, userID uint) error
	SetOffline(account string) error
	IsOnline(account string) bool
}
