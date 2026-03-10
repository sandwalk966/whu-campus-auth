package api

import (
	"time"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

type AuthAPI struct {
	redisService *service.RedisService
}

func NewAuthAPI() *AuthAPI {
	return &AuthAPI{
		redisService: service.NewRedisService(),
	}
}

// Logout 用户登出
func (api *AuthAPI) Logout(c *gin.Context) {
	token := utils.GetTokenFromRequest(c)

	if token == "" {
		utils.ErrorWithMessage(c, "No token found")
		return
	}

	j := utils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		utils.SuccessWithMessage(c, "Logged out successfully")
		return
	}

	expiresAt := claims.ExpiresAt.Unix() - time.Now().Unix()
	if expiresAt < 0 {
		expiresAt = 0
	}

	err = api.redisService.AddTokenToBlacklist(token, time.Duration(expiresAt)*time.Second)
	if err != nil {
		utils.LogErrorf("Logout failed: %v", err)
		utils.ErrorWithMessage(c, "Logout failed, please try again later")
		return
	}

	api.redisService.DeleteUserToken(claims.ID)

	utils.SuccessWithMessage(c, "Logged out successfully")
}
