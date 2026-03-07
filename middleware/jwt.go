package middleware

import (
	"errors"
	"time"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const DBContextKey = "db"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-token")
		if token == "" {
			token = c.GetHeader("Authorization")
			if len(token) > 7 && token[:7] == "Bearer " {
				token = token[7:]
			}
		}

		if token == "" {
			utils.Unauthorized(c, "未登录或非法访问，请登录")
			c.Abort()
			return
		}

		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				utils.Unauthorized(c, "登录已过期，请重新登录")
				c.Abort()
				return
			}
			utils.Unauthorized(c, err.Error())
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
