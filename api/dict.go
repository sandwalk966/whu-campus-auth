package api

import (
	"strconv"
	"whu-campus-auth/model/req"
	"whu-campus-auth/service"
	"whu-campus-auth/utils"

	"github.com/gin-gonic/gin"
)

type DictAPI struct {
	dictService *service.DictService
}

func NewDictAPI(dictService *service.DictService) *DictAPI {
	return &DictAPI{dictService: dictService}
}

func (api *DictAPI) CreateDict(c *gin.Context) {
	var createReq req.CreateDictRequest
	if err := c.ShouldBindJSON(&createReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.dictService.CreateDict(createReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Created successfully")
}

func (api *DictAPI) UpdateDict(c *gin.Context) {
	var updateReq req.UpdateDictRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	if err := api.dictService.UpdateDict(updateReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Updated successfully")
}

func (api *DictAPI) GetDictByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "Dict ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "Dict ID must be a positive integer")
		return
	}

	dict, err := api.dictService.GetDictByID(uint(id))
	if err != nil {
		utils.ErrorWithMessage(c, "Dict not found")
		return
	}

	utils.SuccessWithData(c, dict)
}

func (api *DictAPI) GetDictList(c *gin.Context) {
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
			utils.ErrorWithMessage(c, "Dict status must be 0 or 1")
			return
		}
	}

	dicts, total, err := api.dictService.GetDictList(page, pageSize, name, status)
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, gin.H{
		"list":  dicts,
		"total": total,
	})
}

func (api *DictAPI) DeleteDict(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		utils.ErrorWithMessage(c, "Dict ID cannot be empty")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		utils.ErrorWithMessage(c, "Dict ID must be a positive integer")
		return
	}

	if err := api.dictService.DeleteDict(uint(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "Deleted successfully")
}

func (api *DictAPI) GetDictByCode(c *gin.Context) {
	code := c.Param("code")
	dict, err := api.dictService.GetDictByCode(code)
	if err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithData(c, dict)
}
