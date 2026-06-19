package repository

import frienddomain "chat-system/server/internal/friend/domain"

type FriendRepository interface {
	Request(userID uint, friendID uint) error
	AcceptRequest(requestID uint, userID uint) error
	ExistsAccepted(userID uint, friendID uint) (bool, error)
	ListAccepted(userID uint) ([]frienddomain.Friend, error)
	ListPendingRequests(userID uint) ([]frienddomain.Friend, error)
}
