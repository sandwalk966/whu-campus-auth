package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetUploadDir() string {
	return "./uploads"
}

func EnsureUploadDir() error {
	dir := GetUploadDir()
	return os.MkdirAll(dir, os.ModePerm)
}

func GenerateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Format("20060102150405")
	random := fmt.Sprintf("%d", time.Now().UnixNano())
	return fmt.Sprintf("%s_%s%s", timestamp, random, ext)
}

func GetFilePath(fileName string) string {
	return filepath.Join(GetUploadDir(), fileName)
}

func IsValidFileType(fileName string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".zip", ".rar"}
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
