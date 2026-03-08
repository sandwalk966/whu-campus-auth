package router

import (
	"whu-campus-auth/api"
	"whu-campus-auth/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(protected *gin.RouterGroup, userAPI *api.UserAPI, uploadAPI *api.UploadAPI) {
	user := protected.Group("/user")
	{
		user.GET("/info", userAPI.GetUserInfo)
		user.POST("", middleware.IsAdmin(), userAPI.CreateUser)
		user.PUT("", middleware.IsAdmin(), userAPI.UpdateUser)
		user.PUT("/password", userAPI.ChangePassword)
		user.GET("/list", userAPI.GetUserList)
		user.DELETE("/:id", middleware.IsAdmin(), userAPI.DeleteUser)
		user.POST("/assign-roles", middleware.IsAdmin(), userAPI.AssignRoles)
		user.POST("/avatar", userAPI.UploadAvatar)
	}

	// 上传路由
	upload := protected.Group("/upload")
	{
		upload.POST("", uploadAPI.UploadFile)
		upload.DELETE("/:file_name", uploadAPI.DeleteFile)
	}
}
