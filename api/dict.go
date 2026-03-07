package api

import (
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

	utils.SuccessWithMessage(c, "创建成功")
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

	utils.SuccessWithMessage(c, "更新成功")
}

func (api *DictAPI) GetDictByID(c *gin.Context) {
	id := c.Param("id")
	dict, err := api.dictService.GetDictByID(getIDFromParam(id))
	if err != nil {
		utils.ErrorWithMessage(c, "字典不存在")
		return
	}

	utils.SuccessWithData(c, dict)
}

func (api *DictAPI) GetDictList(c *gin.Context) {
	var listReq req.DictListRequest
	if err := c.ShouldBindJSON(&listReq); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	dicts, total, err := api.dictService.GetDictList(listReq.Page, listReq.PageSize, listReq.Name, listReq.Status)
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
	id := c.Param("id")
	if err := api.dictService.DeleteDict(getIDFromParam(id)); err != nil {
		utils.ErrorWithMessage(c, err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功")
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
