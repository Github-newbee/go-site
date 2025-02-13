package service

import (
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
)

type FileService interface {
	UploadFile(ctx *gin.Context) (string, error)
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

func (s *fileService) UploadFile(ctx *gin.Context) (string, error) {
	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		return "", err
	}
	// 确保上传目录存在
	uploadDir := filepath.Join("storage", "upload")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	newFileName, err := random.UUIdV4()
	if err != nil {
		return "", err
	}
	newFileName += ext
	// 保存文件到指定目录
	filePath := filepath.Join(uploadDir, newFileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return newFileName, nil
}

func (s *fileService) DownloadFile(ctx *gin.Context) error {
	return nil
}
