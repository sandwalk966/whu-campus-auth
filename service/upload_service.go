package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"whu-campus-auth/utils"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) UploadFile(file *multipart.FileHeader) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	if !utils.IsValidFileType(file.Filename) {
		return "", "", fmt.Errorf("不支持的文件类型")
	}

	fileName := utils.GenerateFileName(file.Filename)
	filePath := utils.GetFilePath(fileName)

	if err := utils.EnsureUploadDir(); err != nil {
		return "", "", err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", "", err
	}

	return fileName, filePath, nil
}

func (s *UploadService) DeleteFile(fileName string) error {
	filePath := utils.GetFilePath(fileName)
	return os.Remove(filePath)
}

func (s *UploadService) GetFileURL(fileName string) string {
	return fmt.Sprintf("/uploads/%s", fileName)
}
