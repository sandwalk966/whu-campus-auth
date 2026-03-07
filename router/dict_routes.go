package router

import (
	"whu-campus-auth/api"
	"whu-campus-auth/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterDictRoutes 注册字典路由（包含公开和受保护）
func RegisterDictRoutes(r *gin.Engine, protected *gin.RouterGroup, dictAPI *api.DictAPI) {
	// 公开接口：根据编码查询字典
	r.GET("/api/dict/code/:code", dictAPI.GetDictByCode)

	// 需要认证的接口
	dict := protected.Group("/dict")
	{
		dict.GET("/:id", dictAPI.GetDictByID)
		dict.GET("/list", dictAPI.GetDictList)
		dict.POST("", middleware.IsAdmin(), dictAPI.CreateDict)
		dict.PUT("", middleware.IsAdmin(), dictAPI.UpdateDict)
		dict.DELETE("/:id", middleware.IsAdmin(), dictAPI.DeleteDict)
	}
}
