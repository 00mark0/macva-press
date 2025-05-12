// db/redis/redis.go

package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"), // should be "mn-redis:6379"
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // default
	})
}

// CacheService provides methods for working with Redis cache
type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new cache service
func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{
		client: client,
	}
}

// Get retrieves data from cache and decodes it into the provided destination
// Returns true if found in cache, false if not found
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Key doesn't exist
		return false, nil
	} else if err != nil {
		// Some other Redis error
		return false, fmt.Errorf("redis get error: %w", err)
	}

	// Decode the data into the destination
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(dest)
	if err != nil {
		return false, fmt.Errorf("decode error: %w", err)
	}

	return true, nil
}

// Set stores data in cache with the specified expiration
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	err := c.client.Set(ctx, key, buf.Bytes(), expiration).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// Delete removes a key from the cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	return nil
}

// DeleteByPattern removes all keys matching the pattern
func (c *CacheService) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("redis keys error: %w", err)
	}

	if len(keys) > 0 {
		err = c.client.Del(ctx, keys...).Err()
		if err != nil {
			return fmt.Errorf("redis delete error: %w", err)
		}
	}

	return nil
}

// GenerateKey creates a standardized cache key
func GenerateKey(prefix string, parts ...interface{}) string {
	key := prefix
	for _, part := range parts {
		key += ":" + fmt.Sprint(part)
	}
	return key
}
