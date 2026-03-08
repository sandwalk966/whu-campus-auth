package middleware

import (
	dbModel "whu-campus-auth/model/db"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			utils.Forbidden(c, "无法获取用户信息")
			c.Abort()
			return
		}

		db := GetDB()
		var user dbModel.User
		if err := db.Preload("Roles").First(&user, userId).Error; err != nil {
			utils.ErrorWithMessage(c, "用户不存在")
			c.Abort()
			return
		}

		hasAdminRole := false
		for _, role := range user.Roles {
			if role.Code == "admin" {
				hasAdminRole = true
				break
			}
		}

		if !hasAdminRole {
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

func HasRole(roleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			utils.Forbidden(c, "无法获取用户信息")
			c.Abort()
			return
		}

		db := GetDB()
		var user dbModel.User
		if err := db.Preload("Roles").First(&user, userId).Error; err != nil {
			utils.ErrorWithMessage(c, "用户不存在")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range user.Roles {
			if role.Code == roleCode {
				hasRole = true
				break
			}
		}

		if !hasRole {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}

func HasAnyRole(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			utils.Forbidden(c, "无法获取用户信息")
			c.Abort()
			return
		}

		db := GetDB()
		var user dbModel.User
		if err := db.Preload("Roles").First(&user, userId).Error; err != nil {
			utils.ErrorWithMessage(c, "用户不存在")
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range user.Roles {
			for _, rc := range roleCodes {
				if rc == role.Code {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}
		if !hasRole {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}
		c.Next()
	}
}
