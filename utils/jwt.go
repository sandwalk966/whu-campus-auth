package utils

import (
	"errors"
	"strings"
	"time"
	"whu-campus-auth/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	TokenExpired     = errors.New("token 已过期")
	TokenNotValidYet = errors.New("token 尚未生效")
	TokenMalformed   = errors.New("token 格式不正确")
	TokenInvalid     = errors.New("token 无效")
)

type JWT struct {
	SigningKey  []byte
	ExpiresTime time.Duration
	BufferTime  time.Duration
}

type Claims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWT() *JWT {
	cfg := config.GlobalConfig
	if cfg == nil {
		return &JWT{
			SigningKey:  []byte("whu-campus-auth-secret-key-2026"),
			ExpiresTime: time.Hour * 24 * 7,
			BufferTime:  time.Hour * 24,
		}
	}

	expiresTime, err := time.ParseDuration(cfg.JWT.ExpiresTime)
	if err != nil {
		expiresTime = time.Hour * 24 * 7
	}

	bufferTime, err := time.ParseDuration(cfg.JWT.BufferTime)
	if err != nil {
		bufferTime = time.Hour * 24
	}

	return &JWT{
		SigningKey:  []byte(cfg.JWT.Secret),
		ExpiresTime: expiresTime,
		BufferTime:  bufferTime,
	}
}

func (j *JWT) CreateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

func (j *JWT) CreateTokenByOldToken(oldToken string, claims Claims) (string, error) {
	return j.CreateToken(claims)
}

// GetTokenFromRequest 从请求中提取 token
func GetTokenFromRequest(c *gin.Context) string {
	token := c.GetHeader("x-token")
	if token == "" {
		token = c.GetHeader("Authorization")
		if len(token) > 7 && strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}
	}
	return token
}
