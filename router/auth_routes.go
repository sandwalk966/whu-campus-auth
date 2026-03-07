package router

import (
	"whu-campus-auth/api"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(r *gin.Engine, userAPI *api.UserAPI) {
	r.POST("/api/auth/login", userAPI.Login)
	r.POST("/api/auth/register", userAPI.Register)
}
