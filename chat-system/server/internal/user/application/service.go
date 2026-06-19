package application

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	userdomain "chat-system/server/internal/user/domain"
	userrepository "chat-system/server/internal/user/repository"

	"golang.org/x/crypto/bcrypt"
)

const sessionTTL = 24 * time.Hour

type Service struct {
	repository   userrepository.UserRepository
	sessionStore userrepository.SessionStore
}

type Session struct {
	Token         string     `json:"token"`
	Account       string     `json:"account"`
	UserID        uint       `json:"user_id"`
	Nickname      string     `json:"nickname"`
	LastLoginTime *time.Time `json:"last_login_time"`
}

func NewService(repository userrepository.UserRepository, sessionStore userrepository.SessionStore) *Service {
	return &Service{repository: repository, sessionStore: sessionStore}
}

func (s *Service) Register(nickname string, password string) (*userdomain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &userdomain.User{
		Nickname: nickname,
		Password: string(hash),
	}
	return user, s.repository.Create(user)
}

func (s *Service) Login(account string, password string) (*Session, error) {
	user, err := s.repository.FindByAccount(account)
	if err != nil {
		if errors.Is(err, userrepository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user by account %s: %w", account, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := s.repository.UpdateLastLogin(user.ID); err != nil {
		return nil, fmt.Errorf("update last login for account %s: %w", account, err)
	}
	token, err := newToken()
	if err != nil {
		return nil, fmt.Errorf("generate session token for account %s: %w", account, err)
	}
	session := &Session{
		Token:         token,
		Account:       user.Account,
		UserID:        user.ID,
		Nickname:      user.Nickname,
		LastLoginTime: user.LastLoginTime,
	}
	value := fmt.Sprintf("%d|%s|%s", user.ID, user.Account, user.Nickname)
	if err := s.sessionStore.Set(token, value, sessionTTL); err != nil {
		return nil, fmt.Errorf("store session for account %s: %w", account, err)
	}
	return session, nil
}

func (s *Service) ParseSession(token string) (*Session, error) {
	value, err := s.sessionStore.Get(token)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(value, "|", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid session")
	}
	userID64, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}
	return &Session{UserID: uint(userID64), Account: parts[1], Nickname: parts[2], Token: token}, nil
}

func (s *Service) FindByAccount(account string) (*userdomain.User, error) {
	return s.repository.FindByAccount(account)
}

func (s *Service) FindByID(userID uint) (*userdomain.User, error) {
	return s.repository.FindByID(userID)
}

func newToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
