package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(ctx context.Context, addr string, ttlSeconds int) (*Cache, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		return nil, err
	}
	ttl := time.Duration(ttlSeconds) * time.Second
	return &Cache{client: client, ttl: ttl}, nil
}

func (c *Cache) Key(text, sourceLang, targetLang, customPrompt string) string {
	return fmt.Sprintf("trans:%s:%s:%s:%s", sourceLang, targetLang, customPrompt, text)
}

func (c *Cache) Get(ctx context.Context, key string) (string, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Cache) Set(ctx context.Context, key, value string) error {
	return c.client.Set(ctx, key, value, c.ttl).Err()
}
