package api

import (
	"whu-campus-auth/model/req"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"
	"github.com/gin-gonic/gin"
)

type MenuAPI struct {
	menuService *service.MenuService
}

func NewMenuAPI(menuService *service.MenuService) *MenuAPI {
	return &MenuAPI{menuService: menuService}
}

func (api *MenuAPI) CreateMenu(c *gin.Context) {
	var createReq req.CreateMenuRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.menuService.CreateMenu(createReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功")
}

func (api *MenuAPI) UpdateMenu(c *gin.Context) {
	var updateReq req.UpdateMenuRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.menuService.UpdateMenu(updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "更新成功")
}

func (api *MenuAPI) GetMenuByID(c *gin.Context) {
	id := c.Param("id")
	menu, err := api.menuService.GetMenuByID(getIDFromParam(id))
	if err != nil {
		utils.ErrorWithMessage(c, "菜单不存在")
		return
	}

	utils.SuccessWithData(c, menu)
}

func (api *MenuAPI) GetMenuList(c *gin.Context) {
	var listReq req.MenuListRequest
	if err := c.ShouldBindJSON(&listReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	menus, total, err := api.menuService.GetMenuList(listReq.Page, listReq.PageSize)
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, gin.H{
		"list":  menus,
		"total": total,
	})
}

func (api *MenuAPI) GetMenuTree(c *gin.Context) {
	menus, err := api.menuService.GetMenuTree()
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, menus)
}

func (api *MenuAPI) DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	if err := api.menuService.DeleteMenu(getIDFromParam(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功")
}

func (api *MenuAPI) GetMenusByRoleID(c *gin.Context) {
	roleID := c.Param("role_id")
	menus, err := api.menuService.GetMenusByRoleID(getIDFromParam(roleID))
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, menus)
}
