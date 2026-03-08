package api

import (
	"fmt"
	"whu-campus-auth/model/req"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	userService   *service.UserService
	uploadService *service.UploadService
}

func NewUserAPI(userService *service.UserService, uploadService *service.UploadService) *UserAPI {
	return &UserAPI{
		userService:   userService,
		uploadService: uploadService,
	}
}

func (api *UserAPI) Login(c *gin.Context) {
	var loginReq req.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	token, err := api.userService.Login(loginReq)
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, gin.H{
		"token": token,
	})
}

func (api *UserAPI) CreateUser(c *gin.Context) {
	var createReq req.CreateUserRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		utils.LogErrorf("创建用户 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("创建用户 - 请求参数：%+v", createReq)

	if err := api.userService.CreateUser(createReq); err != nil {
		utils.LogErrorf("创建用户失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("创建用户成功：%s", createReq.Username)
	utils.SuccessWithMessage(c, "创建成功")
}

func (api *UserAPI) Register(c *gin.Context) {
	var registerReq req.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.userService.Register(registerReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "注册成功")
}

func (api *UserAPI) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userId")
	user, err := api.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.ErrorWithMessage(c, "获取用户信息失败")
		return
	}

	utils.SuccessWithData(c, user)
}

func (api *UserAPI) UpdateUser(c *gin.Context) {
	var updateReq req.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.userService.UpdateUser(updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功")
}

func (api *UserAPI) ChangePassword(c *gin.Context) {
	var changeReq req.ChangePasswordRequest
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	userID, _ := c.Get("userId")
	if err := api.userService.ChangePassword(userID.(uint), changeReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功")
}

func (api *UserAPI) GetUserList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	username := c.Query("username")
	status := c.Query("status")

	utils.LogInfof("获取用户列表 - page: %s, page_size: %s, username: %s, status: %s",
		page, pageSize, username, status)

	var pageInt, pageSizeInt, statusInt int
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(pageSize, "%d", &pageSizeInt)
	if status != "" {
		fmt.Sscanf(status, "%d", &statusInt)
	}

	users, total, err := api.userService.GetUserList(pageInt, pageSizeInt, username, statusInt)
	if err != nil {
		utils.LogErrorf("获取用户列表失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("获取用户列表成功 - total: %d", total)
	utils.SuccessWithData(c, gin.H{
		"list":  users,
		"total": total,
	})
}

func (api *UserAPI) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := api.userService.DeleteUser(getIDFromParam(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功")
}

func (api *UserAPI) AssignRoles(c *gin.Context) {
	var assignReq req.AssignRoleRequest
	if err := c.ShouldBindJSON(&assignReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.userService.AssignRoles(assignReq.UserID, assignReq.RoleIDs); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "分配角色成功")
}

func (api *UserAPI) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorWithMessage(c, "上传失败："+err.Error())
		return
	}

	fileName, filePath, err := api.uploadService.UploadFile(file)
	if err != nil {
		utils.ErrorWithMessage(c, "上传失败："+err.Error())
		return
	}

	userID, _ := c.Get("userId")
	avatarURL := "/uploads/" + fileName

	if err := api.userService.UpdateAvatar(userID.(uint), avatarURL); err != nil {
		api.uploadService.DeleteFile(fileName)
		utils.ErrorWithMessage(c, "更新头像失败："+err.Error())
		return
	}

	utils.SuccessWithData(c, gin.H{
		"avatar_url": avatarURL,
		"file_name":  fileName,
		"file_path":  filePath,
	})
}

func getIDFromParam(id string) uint {
	var result uint
	fmt.Sscanf(id, "%d", &result)
	return result
}
