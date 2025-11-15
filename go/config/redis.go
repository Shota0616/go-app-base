package config

import (
	"context"
	"time" // Added for time.Duration

	"github.com/go-redis/redis/v8"
)

// RedisClientInterface defines the methods used from redis.Client that we want to mock.
type RedisClientInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Ping(ctx context.Context) *redis.StatusCmd // Added Ping method
}

var RDB RedisClientInterface

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis!")
	}
}