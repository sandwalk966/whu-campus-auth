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

func SetCache(key string, value interface{}, expiration int) error {
	return RedisClient.Set(context.Background(), key, value, 0).Err()
}

func GetCache(key string) (string, error) {
	return RedisClient.Get(context.Background(), key).Result()
}

func DeleteCache(key string) error {
	return RedisClient.Del(context.Background(), key).Err()
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
