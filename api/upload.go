package api

import (
	"whu-campus-auth/service"
	"whu-campus-auth/utils"
	"github.com/gin-gonic/gin"
)

type UploadAPI struct {
	uploadService *service.UploadService
}

func NewUploadAPI(uploadService *service.UploadService) *UploadAPI {
	return &UploadAPI{uploadService: uploadService}
}

func (api *UploadAPI) UploadFile(c *gin.Context) {
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

	utils.SuccessWithData(c, gin.H{
		"file_name": fileName,
		"file_path": filePath,
		"url":       api.uploadService.GetFileURL(fileName),
	})
}

func (api *UploadAPI) DeleteFile(c *gin.Context) {
	fileName := c.Param("file_name")
	if err := api.uploadService.DeleteFile(fileName); err != nil {
		utils.ErrorWithMessage(c, "删除失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "删除成功")
}
