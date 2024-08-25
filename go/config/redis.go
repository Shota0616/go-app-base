package config

import (
	"github.com/go-redis/redis/v8"
	"context"
)

var RDB *redis.Client

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		panic("Failed to connect to Redis!")
	}
}