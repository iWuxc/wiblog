package internal

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

// RedisSet redis string set
func RedisSet(ctx context.Context, key string, value interface{}, time time.Duration) error {
	err := redisClient.Set(ctx, key, value, time).Err()
	if err != nil {
		return err
	}
	return nil
}

// RedisGet redis string get
func RedisGet(ctx context.Context, key string) (string, error) {
	val, err := redisClient.Get(ctx, "key").Result()
	if err != nil {
		return "", err
	}
	return val, nil
}


