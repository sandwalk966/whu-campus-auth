package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"whu-campus-auth/model/db"

	"github.com/go-redis/redis/v8"
)

// RedisService Redis 服务
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService 创建 Redis 服务实例
func NewRedisService() *RedisService {
	// 使用全局的 RedisClient（在 initializer 中初始化）
	// 为了避免循环导入，我们直接从 initializer 包获取
	return &RedisService{
		client: getRedisClient(),
		ctx:    context.Background(),
	}
}

// getRedisClient 获取 Redis 客户端（从全局变量）
func getRedisClient() *redis.Client {
	// 这里需要导入 initializer 包，但为了避免循环导入
	// 我们使用一个变通方法：在 initializer 中设置全局变量
	return redisClient
}

// SetRedisClient 设置 Redis 客户端（在初始化时调用）
var redisClient *redis.Client

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

// ==================== Token 黑名单管理 ====================

// StoreUserToken 存储用户 token（JWT + Redis 方案）
// userId: 用户 ID
// token: JWT token 字符串
// duration: token 有效期
func (s *RedisService) StoreUserToken(userId uint, token string, duration time.Duration) error {
	key := fmt.Sprintf("user:token:%d", userId)
	return s.client.Set(s.ctx, key, token, duration).Err()
}

// GetUserToken 获取用户存储的 token
func (s *RedisService) GetUserToken(userId uint) (string, error) {
	key := fmt.Sprintf("user:token:%d", userId)
	return s.client.Get(s.ctx, key).Result()
}

// DeleteUserToken 删除用户存储的 token
func (s *RedisService) DeleteUserToken(userId uint) error {
	key := fmt.Sprintf("user:token:%d", userId)
	return s.client.Del(s.ctx, key).Err()
}

// AddTokenToBlacklist 将 token 加入黑名单
// token: JWT token 字符串
// expiresAt: token 过期时间（用于设置 Redis 过期时间）
func (s *RedisService) AddTokenToBlacklist(token string, expiresAt time.Duration) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return s.client.Set(s.ctx, key, "1", expiresAt).Err()
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func (s *RedisService) IsTokenBlacklisted(token string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", token)
	_, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// RemoveTokenFromBlacklist 从黑名单中移除 token（可选功能）
func (s *RedisService) RemoveTokenFromBlacklist(token string) error {
	key := fmt.Sprintf("blacklist:token:%s", token)
	return s.client.Del(s.ctx, key).Err()
}

// ==================== 用户信息缓存管理 ====================

// CacheUserInfo 缓存用户信息
// userId: 用户 ID
// userInfo: 用户信息对象
// duration: 缓存时长
func (s *RedisService) CacheUserInfo(userId uint, userInfo *db.User, duration time.Duration) error {
	key := fmt.Sprintf("user:info:%d", userId)
	data, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, key, data, duration).Err()
}

// GetCachedUserInfo 获取缓存的用户信息
func (s *RedisService) GetCachedUserInfo(userId uint) (*db.User, error) {
	key := fmt.Sprintf("user:info:%d", userId)
	data, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // 缓存不存在
	}
	if err != nil {
		return nil, err
	}

	var userInfo *db.User
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// DeleteUserCache 删除用户缓存
func (s *RedisService) DeleteUserCache(userId uint) error {
	key := fmt.Sprintf("user:info:%d", userId)
	return s.client.Del(s.ctx, key).Err()
}

// CacheUserByUsername 根据用户名缓存用户信息
func (s *RedisService) CacheUserByUsername(username string, userInfo *db.User, duration time.Duration) error {
	key := fmt.Sprintf("user:username:%s", username)
	data, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}
	return s.client.Set(s.ctx, key, data, duration).Err()
}

// GetCachedUserByUsername 根据用户名获取缓存的用户信息
func (s *RedisService) GetCachedUserByUsername(username string) (*db.User, error) {
	key := fmt.Sprintf("user:username:%s", username)
	data, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var userInfo *db.User
	err = json.Unmarshal([]byte(data), &userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// DeleteUserCacheByUsername 根据用户名删除用户缓存
func (s *RedisService) DeleteUserCacheByUsername(username string) error {
	key := fmt.Sprintf("user:username:%s", username)
	return s.client.Del(s.ctx, key).Err()
}

// GetUserCacheKey 获取用户缓存的 key（用于 GetUserByID）
func (s *RedisService) GetUserCacheKey(userId uint) string {
	return fmt.Sprintf("user:info:%d", userId)
}

// GetUserFromCache 从缓存获取用户（简化版本）
func (s *RedisService) GetUserFromCache(userId uint) (*db.User, error) {
	return s.GetCachedUserInfo(userId)
}

// SetUserCache 设置用户缓存（简化版本）
func (s *RedisService) SetUserCache(user *db.User, key string) error {
	// 从 key 中提取 userId
	var userId uint
	fmt.Sscanf(key, "user:info:%d", &userId)
	if userId == 0 {
		return nil
	}
	return s.CacheUserInfo(userId, user, 24*time.Hour)
}

// ==================== API 限流 ====================

// RateLimit API 限流
// key: 限流的唯一标识（如：IP 地址、用户 ID 等）
// limit: 限制次数
// window: 时间窗口
// 返回：是否超过限制，当前计数，错误
func (s *RedisService) RateLimit(key string, limit int64, window time.Duration) (bool, int64, error) {
	rateKey := fmt.Sprintf("rate:limit:%s", key)

	// 使用 Redis 的 INCR 和 EXPIRE 实现滑动窗口限流
	count, err := s.client.Incr(s.ctx, rateKey).Result()
	if err != nil {
		return false, 0, err
	}

	// 如果是第一次访问，设置过期时间
	if count == 1 {
		err = s.client.Expire(s.ctx, rateKey, window).Err()
		if err != nil {
			return false, 0, err
		}
	}

	// 检查是否超过限制
	if count > limit {
		return true, count, nil // 超过限制
	}

	return false, count, nil // 未超过限制
}

// GetRateLimitCount 获取当前限流计数
func (s *RedisService) GetRateLimitCount(key string) (int64, error) {
	rateKey := fmt.Sprintf("rate:limit:%s", key)
	return s.client.Get(s.ctx, rateKey).Int64()
}

// ResetRateLimit 重置限流计数
func (s *RedisService) ResetRateLimit(key string) error {
	rateKey := fmt.Sprintf("rate:limit:%s", key)
	return s.client.Del(s.ctx, rateKey).Err()
}

// ==================== 用户状态标记管理 ====================

// SetUserDisabled 设置用户禁用状态
// userId: 用户 ID
// duration: 禁用时长（通常设置为 Token 的过期时间）
func (s *RedisService) SetUserDisabled(userId uint, duration time.Duration) error {
	key := fmt.Sprintf("user:disabled:%d", userId)
	return s.client.Set(s.ctx, key, "1", duration).Err()
}

// IsUserDisabled 检查用户是否被禁用
func (s *RedisService) IsUserDisabled(userId uint) (bool, error) {
	key := fmt.Sprintf("user:disabled:%d", userId)
	_, err := s.client.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return false, nil // 用户未被禁用
	}
	if err != nil {
		return false, err
	}
	return true, nil // 用户已被禁用
}

// RemoveUserDisabled 移除用户禁用状态（用户重新登录时调用）
func (s *RedisService) RemoveUserDisabled(userId uint) error {
	key := fmt.Sprintf("user:disabled:%d", userId)
	return s.client.Del(s.ctx, key).Err()
}

// ==================== 辅助方法 ====================

// Ping 测试 Redis 连接
func (s *RedisService) Ping() error {
	return s.client.Ping(s.ctx).Err()
}
