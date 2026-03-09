package middleware

import (
	dbModel "whu-campus-auth/model/db"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

// getUserWithRoles 从数据库获取用户及其角色
func getUserWithRoles(c *gin.Context) (*dbModel.User, bool) {
	userId, exists := c.Get("userId")
	if !exists {
		utils.Forbidden(c, "无法获取用户信息")
		c.Abort()
		return nil, false
	}

	db := GetDB()
	var user dbModel.User
	if err := db.Preload("Roles").First(&user, userId).Error; err != nil {
		utils.ErrorWithMessage(c, "用户不存在")
		c.Abort()
		return nil, false
	}

	return &user, true
}

// hasRole 检查用户是否拥有指定角色
func hasRole(user *dbModel.User, roleCode string) bool {
	for _, role := range user.Roles {
		if role.Code == roleCode {
			return true
		}
	}
	return false
}

// hasAnyRole 检查用户是否拥有指定角色中的任意一个
func hasAnyRole(user *dbModel.User, roleCodes []string) bool {
	for _, role := range user.Roles {
		for _, rc := range roleCodes {
			if rc == role.Code {
				return true
			}
		}
	}
	return false
}

// IsAdmin 管理员权限检查
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := getUserWithRoles(c)
		if !ok {
			return
		}

		if !hasRole(user, "admin") {
			utils.LogWarnf("IsAdmin 中间件：用户 %s 不是管理员", user.Username)
			utils.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		utils.LogInfof("IsAdmin 中间件：用户 %s 通过管理员验证", user.Username)
		c.Next()
	}
}

// HasRole 指定角色权限检查
func HasRole(roleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := getUserWithRoles(c)
		if !ok {
			return
		}

		if !hasRole(user, roleCode) {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}

// HasAnyRole 任意角色权限检查
func HasAnyRole(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := getUserWithRoles(c)
		if !ok {
			return
		}

		if !hasAnyRole(user, roleCodes) {
			utils.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}
