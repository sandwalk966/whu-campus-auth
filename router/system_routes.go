package router

import (
	"whu-campus-auth/api"
	"whu-campus-auth/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统管理路由（角色、菜单）
func RegisterSystemRoutes(protected *gin.RouterGroup, roleAPI *api.RoleAPI, menuAPI *api.MenuAPI) {
	// 角色路由
	role := protected.Group("/role")
	{
		role.GET("/:id", roleAPI.GetRoleByID)
		role.GET("/list", roleAPI.GetRoleList)
		role.POST("", middleware.IsAdmin(), roleAPI.CreateRole)
		role.PUT("", middleware.IsAdmin(), roleAPI.UpdateRole)
		role.DELETE("/:id", middleware.IsAdmin(), roleAPI.DeleteRole)
		role.GET("/all", roleAPI.GetAllRoles)
	}

	// 菜单路由
	menu := protected.Group("/menu")
	{
		menu.GET("/:id", menuAPI.GetMenuByID)
		menu.GET("/list", menuAPI.GetMenuList)
		menu.GET("/tree", menuAPI.GetMenuTree)
		menu.POST("", middleware.IsAdmin(), menuAPI.CreateMenu)
		menu.PUT("", middleware.IsAdmin(), menuAPI.UpdateMenu)
		menu.DELETE("/:id", middleware.IsAdmin(), menuAPI.DeleteMenu)
		menu.GET("/role/:role_id", menuAPI.GetMenusByRoleID)
	}
}
