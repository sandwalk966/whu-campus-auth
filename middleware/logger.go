package middleware

import (
	"time"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		costTime := time.Since(startTime)

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		utils.LogInfof("[%s] %d | %13v | %15s | %-7s %s %s",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			costTime,
			clientIP,
			c.Request.Method,
			path,
			query,
		)
	}
}
