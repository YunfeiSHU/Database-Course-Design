package infra

import (
	"time"

	dao "chat-system/server/internal/infrastructure"

	"github.com/go-redis/redis/v8"
)

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(client *redis.Client) RedisSessionStore {
	return RedisSessionStore{client: client}
}

func (s RedisSessionStore) Set(token string, value string, ttl time.Duration) error {
	return s.client.Set(dao.Ctx, "session:"+token, value, ttl).Err()
}

func (s RedisSessionStore) Get(token string) (string, error) {
	return s.client.Get(dao.Ctx, "session:"+token).Result()
}
