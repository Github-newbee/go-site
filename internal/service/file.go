package service

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileService interface {
	UploadFile(ctx *gin.Context) error
	DownloadFile(ctx *gin.Context) error
}

func NewFileService(
	service *Service,
) FileService {
	return &fileService{
		Service: service,
	}
}

type fileService struct {
	*Service
}

func (s *fileService) UploadFile(ctx *gin.Context) error {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}
	// 确保上传目录存在
	uploadDir := filepath.Join("storage", "upload")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	// 保存文件到指定目录
	filePath := filepath.Join(uploadDir, file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return err
	}

	return nil
}

func (s *fileService) DownloadFile(ctx *gin.Context) error {
	return nil
}
