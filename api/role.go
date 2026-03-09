package api

import (
	"strconv"
	"whu-campus-auth/model/req"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

type RoleAPI struct {
	roleService *service.RoleService
}

func NewRoleAPI(roleService *service.RoleService) *RoleAPI {
	return &RoleAPI{roleService: roleService}
}

func (api *RoleAPI) CreateRole(c *gin.Context) {
	var createReq req.CreateRoleRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.roleService.CreateRole(createReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Created successfully")
}

func (api *RoleAPI) UpdateRole(c *gin.Context) {
	var updateReq req.UpdateRoleRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.roleService.UpdateRole(updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Updated successfully")
}

func (api *RoleAPI) GetRoleByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "Role ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "Role ID must be a positive integer")
		return
	}

	role, err := api.roleService.GetRoleByID(uint(id))
	if err != nil {
		utils.ErrorWithMessage(c, "Role not found")
		return
	}

	utils.SuccessWithData(c, role)
}

func (api *RoleAPI) GetRoleList(c *gin.Context) {
	// 从查询参数中读取数据
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	name := c.Query("name")
	statusStr := c.Query("status")

	// 参数转换和校验
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.ErrorWithMessage(c, "Page number must be a positive integer")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		utils.ErrorWithMessage(c, "Page size must be between 1 and 100")
		return
	}

	var status int
	if statusStr != "" {
		status, err = strconv.Atoi(statusStr)
		if err != nil || (status != 0 && status != 1) {
			utils.ErrorWithMessage(c, "Role status must be 0 or 1")
			return
		}
	}

	roles, total, err := api.roleService.GetRoleList(page, pageSize, name, status)
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, gin.H{
		"list":  roles,
		"total": total,
	})
}

func (api *RoleAPI) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "Role ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "Role ID must be a positive integer")
		return
	}

	if err := api.roleService.DeleteRole(uint(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Deleted successfully")
}

func (api *RoleAPI) GetAllRoles(c *gin.Context) {
	roles, err := api.roleService.GetAllRoles()
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, roles)
}
