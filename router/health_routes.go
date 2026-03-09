package router

import "github.com/gin-gonic/gin"

// RegisterHealthRoutes 注册健康检查路由
func RegisterHealthRoutes(r *gin.Engine) {
	// 处理根路径的 HEAD 请求
	r.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})

	// 处理根路径的 GET 请求
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "service is running"})
	})
}
