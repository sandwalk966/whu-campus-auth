package middleware

import (
	"net/http"
	"os"
	"strings"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		if os.Getenv("GIN_MODE") != "release" {
			allowedOriginsEnv = "http://localhost:3000,http://localhost:5173,https://nginx"
		} else {
			allowedOriginsEnv = ""
		}
	}

	allowedOrigins := strings.Split(allowedOriginsEnv, ",")
	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

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

		if origin == "" {
			c.Next()
			return
		}

		if len(allowedOrigins) == 0 {
			if strings.HasPrefix(origin, "https://") {
				if isRelease {
					utils.LogInfof("CORS allowed HTTPS origin: %s", origin)
				}
				c.Header("Access-Control-Allow-Origin", origin)
				c.Next()
				return
			} else {
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
