package router

import "github.com/gin-gonic/gin"

// RegisterStaticRoutes 注册静态文件路由
func RegisterStaticRoutes(r *gin.Engine) {
	r.Static("/uploads", "./uploads")
}
