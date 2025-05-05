package config

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("fail to connect to redis: %v", err)
	}

	return nil
}

func GetRedisContext() context.Context{
	return context.Background()
}