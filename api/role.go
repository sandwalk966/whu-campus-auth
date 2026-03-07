package api

import (
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

	utils.SuccessWithMessage(c, "创建成功")
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

	utils.SuccessWithMessage(c, "更新成功")
}

func (api *RoleAPI) GetRoleByID(c *gin.Context) {
	id := c.Param("id")
	role, err := api.roleService.GetRoleByID(getIDFromParam(id))
	if err != nil {
		utils.ErrorWithMessage(c, "角色不存在")
		return
	}

	utils.SuccessWithData(c, role)
}

func (api *RoleAPI) GetRoleList(c *gin.Context) {
	var listReq req.RoleListRequest
	if err := c.ShouldBindJSON(&listReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	roles, total, err := api.roleService.GetRoleList(listReq.Page, listReq.PageSize, listReq.Name, listReq.Status)
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
	id := c.Param("id")
	if err := api.roleService.DeleteRole(getIDFromParam(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功")
}

func (api *RoleAPI) GetAllRoles(c *gin.Context) {
	roles, err := api.roleService.GetAllRoles()
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, roles)
}
