package middleware

import (
	"errors"
	"time"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
)

const DBContextKey = "db"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := utils.GetTokenFromRequest(c)

		if token == "" {
			utils.Unauthorized(c, "No token found, please login")
			c.Abort()
			return
		}

		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				utils.Unauthorized(c, "Token expired, please login again")
				c.Abort()
				return
			}
			utils.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 检查 token 是否在黑名单中
		redisService := service.NewRedisService()
		isBlacklisted, err := redisService.IsTokenBlacklisted(token)
		if err != nil {
			utils.LogErrorf("检查 token 黑名单失败：%v", err)
		}
		if isBlacklisted {
			utils.Unauthorized(c, "Token is invalid, please login again")
			c.Abort()
			return
		}

		// 验证 Redis 中存储的 token 是否与当前 token 一致（JWT + Redis 方案）
		storedToken, err := redisService.GetUserToken(claims.ID)
		if err == nil && storedToken != "" {
			// Redis 中有存储的 token，检查是否一致
			if storedToken != token {
				utils.Unauthorized(c, "Token expired, please login again")
				c.Abort()
				return
			}
		} else if err != nil && !errors.Is(err, redis.Nil) {
			// 非 Redis Nil 错误，记录日志但不影响验证
			utils.LogErrorf("获取 Redis 中存储的 token 失败：%v", err)
		}

		// 检查用户是否被禁用/删除
		isDisabled, err := redisService.IsUserDisabled(claims.ID)
		if err != nil {
			utils.LogErrorf("检查用户状态失败：%v", err)
		}
		if isDisabled {
			utils.Unauthorized(c, "User has been disabled or deleted, please login again")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Set("username", claims.Username)

		if claims.ExpiresAt.Unix()-time.Now().Unix() < int64(j.BufferTime.Seconds()) {
			newClaims := claims
			newClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(j.ExpiresTime))
			newToken, _ := j.CreateTokenByOldToken(token, *newClaims)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", time.Now().Add(j.ExpiresTime).Format(time.RFC3339))
		}

		c.Next()
	}
}
