package api

import (
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
	// 获取 token
	token := utils.GetTokenFromRequest(c)

	if token == "" {
		utils.ErrorWithMessage(c, "No token found")
		return
	}

	// 解析 token 获取过期时间
	j := utils.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		// token 无效，直接返回成功（反正已经不能用了）
		utils.SuccessWithMessage(c, "Logged out successfully")
		return
	}

	// 计算 token 剩余有效时间
	expiresAt := claims.ExpiresAt.Unix() - int64(utils.NewJWT().BufferTime.Seconds())
	if expiresAt < 0 {
		expiresAt = 0
	}

	// 将 token 加入黑名单（使用 token 的剩余有效期）
	err = api.redisService.AddTokenToBlacklist(token, utils.NewJWT().ExpiresTime)
	if err != nil {
		utils.LogErrorf("Logout failed: %v", err)
		utils.ErrorWithMessage(c, "Logout failed, please try again later")
		return
	}

	utils.SuccessWithMessage(c, "Logged out successfully")
}
