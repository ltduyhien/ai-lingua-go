// Package cache provides a Redis-backed cache for translation results.
// It is used at runtime by the gRPC handler to avoid calling the LLM when the same translation was already done.
// The caller only creates a cache when config.RedisAddr is non-empty; otherwise cache is skipped.
package cache

import (
	// context is used to pass cancellation and timeouts into Redis calls so we can abort if needed.
	"context"
	// fmt is used to build the cache key string from the translation inputs.
	"fmt"
	// time is used to set the TTL (time-to-live) on cached values so they expire after RedisTTLSeconds.
	"time"

	// redis is the go-redis client; we use it to get and set key-value pairs in Redis.
	"github.com/redis/go-redis/v9"
)

// Cache holds the Redis client and the TTL duration so we can set expiry on every Set call.
// We store TTL as a duration so we don't convert from seconds on every Set.
type Cache struct {
	// client is the go-redis client; it manages the connection pool and sends commands to Redis.
	client *redis.Client
	// ttl is how long each cached translation stays in Redis; after this duration the key expires and Redis deletes it.
	ttl time.Duration
}

// New creates a Cache that talks to Redis at addr (e.g. "localhost:6379") and expires keys after ttlSeconds.
// It pings Redis to verify the connection; if Redis is down or addr is wrong we return an error so the server can fail fast.
// The caller passes config.RedisAddr and config.RedisTTLSeconds from config; cache does not import config.
func New(ctx context.Context, addr string, ttlSeconds int) (*Cache, error) {
	// redis.NewClient creates a client with default options; we only set Addr so it connects to our Redis instance.
	client := redis.NewClient(&redis.Options{Addr: addr})
	// Ping verifies that we can reach Redis; we use a short timeout so we don't block startup for long if Redis is down.
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		// Return the error so the caller (e.g. main) can log and exit or skip using cache.
		return nil, err
	}
	// Convert seconds to time.Duration so we can pass it to client.Set(..., ttl) later.
	ttl := time.Duration(ttlSeconds) * time.Second
	return &Cache{client: client, ttl: ttl}, nil
}

// Key builds a cache key for a translation request so the same text + language pair always maps to the same key.
// We use a simple format so we can reproduce it in Get and Set; we don't hash because the key is for lookup only.
func (c *Cache) Key(text, sourceLang, targetLang string) string {
	// Sprintf produces a string like "trans:en:es:Hello world"; the prefix "trans:" avoids collisions with other keys.
	return fmt.Sprintf("trans:%s:%s:%s", sourceLang, targetLang, text)
}

// Get returns the cached translation for key if it exists.
// The first return is the value; the second is true only when the key was found (false when missing or on error).
// redis.Nil is the error Redis returns when the key does not exist; we treat it as "not found" and return false, nil.
func (c *Cache) Get(ctx context.Context, key string) (string, bool, error) {
	// Get sends a GET command to Redis; it returns redis.Nil when the key is missing.
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key not in cache; caller will translate and then call Set.
		return "", false, nil
	}
	if err != nil {
		// Network or Redis error; we return it so the caller can decide to retry or fall back to uncached translation.
		return "", false, err
	}
	return val, true, nil
}

// Set stores the translated value at key and sets the TTL so Redis automatically deletes it after c.ttl.
// If Redis is down or the write fails we return the error; the caller already has the translation so they can still respond.
func (c *Cache) Set(ctx context.Context, key, value string) error {
	// Set(key, value, ttl) stores the value and sets the expiration; 0 ttl would mean no expiry, we always pass c.ttl.
	return c.client.Set(ctx, key, value, c.ttl).Err()
}
