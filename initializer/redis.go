package initializer

import (
	"context"
	"fmt"
	"whu-campus-auth/config"
	"whu-campus-auth/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}
	
	utils.LogInfo("Redis 连接成功")
	return nil
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// InitLogger 初始化日志系统
func InitLogger(cfg *config.Config) error {
	return utils.InitLogger(cfg)
}

// SyncLogger 同步并关闭日志
func SyncLogger() {
	utils.Sync()
}

// LogInfo 记录信息日志
func LogInfo(msg string) {
	zap.L().Info(msg)
}

// LogInfof 记录格式化信息日志
func LogInfof(format string, args ...interface{}) {
	zap.L().Info(fmt.Sprintf(format, args...))
}

// LogErrorf 记录错误日志
func LogErrorf(format string, args ...interface{}) {
	zap.L().Error(fmt.Sprintf(format, args...))
}

// LogFatal 记录致命错误日志并退出
func LogFatalf(format string, args ...interface{}) {
	zap.L().Fatal(fmt.Sprintf(format, args...))
}
