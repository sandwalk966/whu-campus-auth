package middleware

import (
	"time"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware API 限流中间件
// limit: 限制次数
// window: 时间窗口
func RateLimitMiddleware(limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisService := service.NewRedisService()

		// 使用 IP 地址作为限流 key
		key := c.ClientIP()

		// 检查是否超过限制
		exceeded, count, err := redisService.RateLimit(key, limit, window)
		if err != nil {
			utils.LogErrorf("限流检查失败：%v", err)
			c.Next() // 限流失败不影响正常请求
			return
		}

		if exceeded {
			utils.LogWarnf("IP %s 触发限流：%d/%d", key, count, limit)
			utils.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser 基于用户的限流中间件
// limit: 限制次数
// window: 时间窗口
func RateLimitByUser(limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		redisService := service.NewRedisService()

		// 使用用户 ID 作为限流 key
		userId, exists := c.Get("userId")
		if !exists {
			// 未登录，使用 IP 限流
			RateLimitMiddleware(limit, window)(c)
			return
		}

		key := userId.(string)

		// 检查是否超过限制
		exceeded, count, err := redisService.RateLimit(key, limit, window)
		if err != nil {
			utils.LogErrorf("用户限流检查失败：%v", err)
			c.Next()
			return
		}

		if exceeded {
			utils.LogWarnf("用户 %s 触发限流：%d/%d", key, count, limit)
			utils.TooManyRequests(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
