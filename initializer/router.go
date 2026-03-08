package initializer

import (
	"whu-campus-auth/middleware"
	"whu-campus-auth/router"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter(deps *Dependencies) *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.Logger())

	// 注册认证路由（公开）
	router.RegisterAuthRoutes(r, deps.UserAPI)

	// 注册需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// 用户相关路由
		router.RegisterUserRoutes(protected, deps.UserAPI, deps.UploadAPI)

		// 系统管理路由（角色、菜单）
		router.RegisterSystemRoutes(protected, deps.RoleAPI, deps.MenuAPI)
	}

	// 注册字典路由（包含公开和受保护）
	router.RegisterDictRoutes(r, protected, deps.DictAPI)

	return r
}
