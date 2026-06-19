package application

import (
	friendrepository "chat-system/server/internal/friend/repository"
	userdomain "chat-system/server/internal/user/domain"
)

type UserFinder interface {
	FindByAccount(account string) (*userdomain.User, error)
}

type OnlineChecker interface {
	IsOnline(account string) bool
}

type Service struct {
	repository    friendrepository.FriendRepository
	users         UserFinder
	onlineChecker OnlineChecker
}

type FriendItem struct {
	ID       uint            `json:"id"`
	UserID   uint            `json:"user_id"`
	FriendID uint            `json:"friend_id"`
	Friend   userdomain.User `json:"friend"`
	Online   bool            `json:"online"`
}

type FriendRequestItem struct {
	ID       uint            `json:"id"`
	UserID   uint            `json:"user_id"`
	FriendID uint            `json:"friend_id"`
	Status   string          `json:"status"`
	User     userdomain.User `json:"user"`
}

func NewService(repository friendrepository.FriendRepository, users UserFinder, onlineChecker OnlineChecker) *Service {
	return &Service{repository: repository, users: users, onlineChecker: onlineChecker}
}

func (s *Service) ListFriends(userID uint) ([]FriendItem, error) {
	friends, err := s.repository.ListAccepted(userID)
	if err != nil {
		return nil, err
	}
	items := make([]FriendItem, 0, len(friends))
	for _, row := range friends {
		items = append(items, FriendItem{
			ID:       row.ID,
			UserID:   row.UserID,
			FriendID: row.FriendID,
			Friend:   row.Friend,
			Online:   s.onlineChecker != nil && s.onlineChecker.IsOnline(row.Friend.Account),
		})
	}
	return items, nil
}

func (s *Service) AddFriend(userID uint, account string) error {
	friend, err := s.users.FindByAccount(account)
	if err != nil {
		return err
	}
	if friend.ID == userID {
		return ErrCannotAddSelf
	}
	return s.repository.Request(userID, friend.ID)
}

func (s *Service) ListFriendRequests(userID uint) ([]FriendRequestItem, error) {
	requests, err := s.repository.ListPendingRequests(userID)
	if err != nil {
		return nil, err
	}
	items := make([]FriendRequestItem, 0, len(requests))
	for _, row := range requests {
		items = append(items, FriendRequestItem{
			ID:       row.ID,
			UserID:   row.UserID,
			FriendID: row.FriendID,
			Status:   row.Status,
			User:     row.User,
		})
	}
	return items, nil
}

func (s *Service) AcceptFriendRequest(requestID uint, userID uint) error {
	return s.repository.AcceptRequest(requestID, userID)
}

func (s *Service) AreFriends(userID uint, friendID uint) (bool, error) {
	return s.repository.ExistsAccepted(userID, friendID)
}
