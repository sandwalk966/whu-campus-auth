package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"whu-campus-auth/config"
)

var RedisClient *redis.Client

func InitRedis() error {
	cfg := config.GlobalConfig.Redis
	
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}
	
	return nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
