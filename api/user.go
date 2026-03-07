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
	var listReq req.UserListRequest
	if err := c.ShouldBindJSON(&listReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	users, total, err := api.userService.GetUserList(listReq.Page, listReq.PageSize, listReq.Username, listReq.Status)
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

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
