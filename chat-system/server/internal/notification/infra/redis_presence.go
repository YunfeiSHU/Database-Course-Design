package infra

import (
	"fmt"

	dao "chat-system/server/internal/infrastructure"

	"github.com/go-redis/redis/v8"
)

type RedisPresenceStore struct {
	client *redis.Client
}

func NewRedisPresenceStore(client *redis.Client) RedisPresenceStore {
	return RedisPresenceStore{client: client}
}

func (s RedisPresenceStore) SetOnline(account string, userID uint) error {
	return s.client.Set(dao.Ctx, "online:"+account, fmt.Sprint(userID), 0).Err()
}

func (s RedisPresenceStore) SetOffline(account string) error {
	return s.client.Del(dao.Ctx, "online:"+account).Err()
}

func (s RedisPresenceStore) IsOnline(account string) bool {
	if s.client == nil {
		return false
	}
	return s.client.Exists(dao.Ctx, "online:"+account).Val() > 0
}
