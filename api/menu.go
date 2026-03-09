package api

import (
	"strconv"
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

	utils.SuccessWithMessage(c, "Created successfully")
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

	utils.SuccessWithMessage(c, "Updated successfully")
}

func (api *MenuAPI) GetMenuByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "Menu ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "Menu ID must be a positive integer")
		return
	}

	menu, err := api.menuService.GetMenuByID(uint(id))
	if err != nil {
		utils.ErrorWithMessage(c, "Menu not found")
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
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "Menu ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "Menu ID must be a positive integer")
		return
	}

	if err := api.menuService.DeleteMenu(uint(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Deleted successfully")
}

func (api *MenuAPI) GetMenusByRoleID(c *gin.Context) {
	roleIDStr := c.Param("role_id")
	if roleIDStr == "" {
		utils.ErrorWithMessage(c, "Role ID cannot be empty")
		return
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil || roleID < 1 {
		utils.ErrorWithMessage(c, "Role ID must be a positive integer")
		return
	}

	menus, err := api.menuService.GetMenusByRoleID(uint(roleID))
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, menus)
}
