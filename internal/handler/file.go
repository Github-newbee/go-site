package handler

import (
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	*Handler
	fileService service.FileService
}

func NewFileHandler(
	handler *Handler,
	fileService service.FileService,
) *FileHandler {
	return &FileHandler{
		Handler:     handler,
		fileService: fileService,
	}
}

func (h *FileHandler) UploadFile(ctx *gin.Context) {
	newFileName, err := h.fileService.UploadFile(ctx)
	if err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, newFileName)
}
