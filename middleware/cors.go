package middleware

import (
	"os"
	"strings"
	"whu-campus-auth/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		// 开发环境默认值（本地开发）
		if os.Getenv("GIN_MODE") != "release" {
			// 开发环境支持 localhost 和 Docker 内部访问
			allowedOriginsEnv = "http://localhost:3000,http://localhost:5173,https://nginx"
		} else {
			// 生产环境可以为空
			// 如果为空，只有通过 Nginx 反向代理的请求（无 Origin 头）或 HTTPS 请求才能访问
			allowedOriginsEnv = ""
		}
	}

	allowedOrigins := strings.Split(allowedOriginsEnv, ",")
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	// 过滤空字符串
	var filteredOrigins []string
	for _, origin := range allowedOrigins {
		if origin != "" {
			filteredOrigins = append(filteredOrigins, origin)
		}
	}
	allowedOrigins = filteredOrigins

	isRelease := os.Getenv("GIN_MODE") == "release"

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 如果没有 Origin 头，说明不是浏览器请求（如 Nginx 反向代理、curl、Postman 等）
		// 直接放行，不需要 CORS 检查
		if origin == "" {
			c.Next()
			return
		}

		// 如果 ALLOWED_ORIGINS 为空，检查是否是 HTTPS 请求
		// 这是为了方便生产环境部署（不知道具体域名的情况）
		if len(allowedOrigins) == 0 {
			if strings.HasPrefix(origin, "https://") {
				// HTTPS 请求，放行并记录日志
				if isRelease {
					utils.LogInfof("CORS allowed HTTPS origin: %s", origin)
				}
				c.Header("Access-Control-Allow-Origin", origin)
				c.Next()
				return
			} else {
				// HTTP 请求，拒绝（开发环境除外）
				if isRelease {
					utils.LogWarnf("CORS rejected HTTP origin: %s", origin)
				}
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			if isRelease {
				utils.LogWarnf("CORS rejected: %s", origin)
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, x-token")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, new-token, new-expires-at")

		if isRelease {
			c.Header("Access-Control-Max-Age", "86400")
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		} else {
			c.Header("Access-Control-Max-Age", "3600")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
