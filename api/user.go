package api

import (
	"strconv"
	"strings"
	"whu-campus-auth/model/req"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

type UserAPI struct {
	userService   *service.UserService
	uploadService *service.UploadService
	redisService  *service.RedisService
}

func NewUserAPI(userService *service.UserService, uploadService *service.UploadService) *UserAPI {
	return &UserAPI{
		userService:   userService,
		uploadService: uploadService,
		redisService:  service.NewRedisService(),
	}
}

func (api *UserAPI) Login(c *gin.Context) {
	var loginReq req.LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		utils.LogErrorf("登录 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, "请求参数错误")
		return
	}

	// 参数校验
	loginReq.Username = strings.TrimSpace(loginReq.Username)
	if len(loginReq.Username) < 3 || len(loginReq.Username) > 50 {
		utils.ErrorWithMessage(c, "Username must be between 3 and 50 characters")
		return
	}

	if len(loginReq.Password) < 6 || len(loginReq.Password) > 50 {
		utils.ErrorWithMessage(c, "Password must be between 6 and 50 characters")
		return
	}

	// 登录并获取 token
	token, user, err := api.userService.Login(loginReq)
	if err != nil {
		utils.LogErrorf("登录失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	// 将 token 存入 Redis（支持多点登录控制）
	err = api.redisService.StoreUserToken(user.ID, token, utils.NewJWT().ExpiresTime)
	if err != nil {
		utils.LogErrorf("存储 token 到 Redis 失败：%v", err)
		// 存储失败不影响登录，只记录日志
	}

	utils.LogInfof("用户登录成功：%s (ID: %d)", loginReq.Username, user.ID)
	utils.SuccessWithData(c, gin.H{
		"token":      token,
		"expires_in": int64(utils.NewJWT().ExpiresTime.Seconds()),
	})
}

func (api *UserAPI) CreateUser(c *gin.Context) {
	var createReq req.CreateUserRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		utils.LogErrorf("创建用户 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, "请求参数错误")
		return
	}

	// 参数校验
	createReq.Username = strings.TrimSpace(createReq.Username)
	if len(createReq.Username) < 3 || len(createReq.Username) > 50 {
		utils.ErrorWithMessage(c, "Username must be between 3 and 50 characters")
		return
	}

	if len(createReq.Password) < 6 || len(createReq.Password) > 50 {
		utils.ErrorWithMessage(c, "Password must be between 6 and 50 characters")
		return
	}

	// 邮箱格式校验（如果提供了邮箱）
	if createReq.Email != "" && !utils.IsValidEmail(createReq.Email) {
		utils.ErrorWithMessage(c, "Invalid email format")
		return
	}

	// 手机号格式校验（如果提供了手机号）
	if createReq.Phone != "" && !utils.IsValidPhone(createReq.Phone) {
		utils.ErrorWithMessage(c, "Invalid phone number format")
		return
	}

	// 状态校验
	if createReq.Status != 0 && createReq.Status != 1 {
		utils.ErrorWithMessage(c, "User status must be 0 or 1")
		return
	}

	utils.LogInfof("创建用户 - 请求参数：%+v", createReq)

	if err := api.userService.CreateUser(createReq); err != nil {
		utils.LogErrorf("创建用户失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("创建用户成功：%s", createReq.Username)
	utils.SuccessWithMessage(c, "Created successfully")
}

func (api *UserAPI) Register(c *gin.Context) {
	var registerReq req.RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		utils.LogErrorf("注册 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, "请求参数错误")
		return
	}

	// 参数校验
	registerReq.Username = strings.TrimSpace(registerReq.Username)
	if len(registerReq.Username) < 3 || len(registerReq.Username) > 50 {
		utils.ErrorWithMessage(c, "Username must be between 3 and 50 characters")
		return
	}

	if len(registerReq.Password) < 6 || len(registerReq.Password) > 50 {
		utils.ErrorWithMessage(c, "Password must be between 6 and 50 characters")
		return
	}

	// 邮箱格式校验（如果提供了邮箱）
	if registerReq.Email != "" && !utils.IsValidEmail(registerReq.Email) {
		utils.ErrorWithMessage(c, "Invalid email format")
		return
	}

	// 手机号格式校验（如果提供了手机号）
	if registerReq.Phone != "" && !utils.IsValidPhone(registerReq.Phone) {
		utils.ErrorWithMessage(c, "Invalid phone number format")
		return
	}

	if err := api.userService.Register(registerReq); err != nil {
		utils.LogErrorf("注册失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("用户注册成功：%s", registerReq.Username)
	utils.SuccessWithMessage(c, "Registered successfully")
}

func (api *UserAPI) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userId")
	// Fix: userID is uint, not string
	var userId uint
	switch v := userID.(type) {
	case uint:
		userId = v
	case float64: // JSON 解析后可能是 float64
		userId = uint(v)
	default:
		utils.ErrorWithMessage(c, "Invalid user ID type")
		return
	}

	// 直接调用 Service 层的带缓存方法（不关心缓存细节）
	user, err := api.userService.GetUserByIDWithCache(userId)
	if err != nil {
		utils.ErrorWithMessage(c, "Failed to get user info")
		return
	}

	utils.SuccessWithData(c, user)
}

func (api *UserAPI) UpdateUser(c *gin.Context) {
	var updateReq req.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		utils.LogErrorf("更新用户 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, "请求参数错误")
		return
	}

	// 邮箱格式校验（如果提供了邮箱）
	if updateReq.Email != "" && !utils.IsValidEmail(updateReq.Email) {
		utils.ErrorWithMessage(c, "Invalid email format")
		return
	}

	// 手机号格式校验（如果提供了手机号）
	if updateReq.Phone != "" && !utils.IsValidPhone(updateReq.Phone) {
		utils.ErrorWithMessage(c, "Invalid phone number format")
		return
	}

	// 性别校验
	if updateReq.Gender < 0 || updateReq.Gender > 2 {
		utils.ErrorWithMessage(c, "Gender must be 0(unknown), 1(male) or 2(female)")
		return
	}

	// 状态校验
	if updateReq.Status != 0 && updateReq.Status != 1 {
		utils.ErrorWithMessage(c, "User status must be 0 or 1")
		return
	}

	if err := api.userService.UpdateUser(updateReq); err != nil {
		utils.LogErrorf("更新用户失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("更新用户成功：%d", updateReq.ID)
	utils.SuccessWithMessage(c, "Updated successfully")
}

func (api *UserAPI) ChangePassword(c *gin.Context) {
	var changeReq req.ChangePasswordRequest
	if err := c.ShouldBindJSON(&changeReq); err != nil {
		utils.LogErrorf("修改密码 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, "请求参数错误")
		return
	}

	// 参数校验
	if len(changeReq.OldPassword) < 6 || len(changeReq.OldPassword) > 50 {
		utils.ErrorWithMessage(c, "Old password must be between 6 and 50 characters")
		return
	}

	if len(changeReq.NewPassword) < 6 || len(changeReq.NewPassword) > 50 {
		utils.ErrorWithMessage(c, "New password must be between 6 and 50 characters")
		return
	}

	// 检查新旧密码是否相同
	if changeReq.OldPassword == changeReq.NewPassword {
		utils.ErrorWithMessage(c, "New password cannot be the same as old password")
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		utils.LogErrorf("修改密码 - 未找到用户 ID")
		utils.ErrorWithMessage(c, "Unauthorized")
		return
	}

	if err := api.userService.ChangePassword(userID.(uint), changeReq); err != nil {
		utils.LogErrorf("修改密码失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("用户修改密码成功：%d", userID)
	utils.SuccessWithMessage(c, "Password changed successfully")
}

func (api *UserAPI) GetUserList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	username := c.Query("username")
	status := c.Query("status")

	// 参数转换和校验
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		utils.ErrorWithMessage(c, "Page number must be a positive integer")
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
		utils.ErrorWithMessage(c, "Page size must be between 1 and 100")
		return
	}

	var statusInt int
	if status != "" {
		statusInt, err = strconv.Atoi(status)
		if err != nil || (statusInt != 0 && statusInt != 1) {
			utils.ErrorWithMessage(c, "User status must be 0 or 1")
			return
		}
	}

	utils.LogInfof("获取用户列表 - page: %d, page_size: %d, username: %s, status: %s",
		pageInt, pageSizeInt, username, status)

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
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "User ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "User ID must be a positive integer")
		return
	}

	// 不允许删除 ID 为 1 的管理员账户（可选的安全措施）
	if id == 1 {
		utils.ErrorWithMessage(c, "System administrator account cannot be deleted")
		return
	}

	// 不允许删除自己
	currentUserID, exists := c.Get("userId")
	if !exists {
		utils.ErrorWithMessage(c, "Current user information not found")
		return
	}

	if uint(id) == currentUserID.(uint) {
		utils.ErrorWithMessage(c, "You cannot delete your own account")
		return
	}

	if err := api.userService.DeleteUser(uint(id)); err != nil {
		utils.LogErrorf("删除用户失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("删除用户成功：%d", id)
	utils.SuccessWithMessage(c, "Deleted successfully")
}

func (api *UserAPI) AssignRoles(c *gin.Context) {
	var assignReq req.AssignRoleRequest
	if err := c.ShouldBindJSON(&assignReq); err != nil {
		utils.LogErrorf("分配角色 - 参数绑定失败：%v", err)
		utils.ErrorWithMessage(c, "请求参数错误")
		return
	}

	// 参数校验
	if assignReq.UserID < 1 {
		utils.ErrorWithMessage(c, "User ID must be a positive integer")
		return
	}

	if len(assignReq.RoleIDs) == 0 {
		utils.ErrorWithMessage(c, "At least one role must be assigned")
		return
	}

	// 校验所有角色 ID
	for _, roleID := range assignReq.RoleIDs {
		if roleID < 1 {
			utils.ErrorWithMessage(c, "Role ID must be a positive integer")
			return
		}
	}

	if err := api.userService.AssignRoles(assignReq.UserID, assignReq.RoleIDs); err != nil {
		utils.LogErrorf("分配角色失败：%v", err)
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.LogInfof("分配角色成功 - 用户 ID: %d, 角色数：%d", assignReq.UserID, len(assignReq.RoleIDs))
	utils.SuccessWithMessage(c, "Roles assigned successfully")
}

func (api *UserAPI) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.LogErrorf("上传头像 - 获取文件失败：%v", err)
		utils.ErrorWithMessage(c, "Upload failed: Please provide a file")
		return
	}

	// 校验文件大小（限制 5MB）
	if file.Size > 5*1024*1024 {
		utils.ErrorWithMessage(c, "File size cannot exceed 5MB")
		return
	}

	// 校验文件类型
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	if !allowedTypes[ext] {
		utils.ErrorWithMessage(c, "Only JPG, JPEG, PNG, and GIF formats are supported")
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		utils.LogErrorf("上传头像 - 未找到用户 ID")
		utils.ErrorWithMessage(c, "Unauthorized")
		return
	}

	fileName, filePath, err := api.uploadService.UploadFile(file)
	if err != nil {
		utils.LogErrorf("上传头像失败：%v", err)
		utils.ErrorWithMessage(c, "Upload failed: "+err.Error())
		return
	}

	avatarURL := "/uploads/" + fileName

	if err := api.userService.UpdateAvatar(userID.(uint), avatarURL); err != nil {
		api.uploadService.DeleteFile(fileName)
		utils.LogErrorf("更新头像失败：%v", err)
		utils.ErrorWithMessage(c, "Failed to update avatar: "+err.Error())
		return
	}

	utils.LogInfof("用户上传头像成功：%d, %s", userID, fileName)
	utils.SuccessWithData(c, gin.H{
		"avatar_url": avatarURL,
		"file_name":  fileName,
		"file_path":  filePath,
	})
}
