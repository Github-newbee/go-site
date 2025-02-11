package handler

import (
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/service"
	"go-my-demo/pkg/request"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	*Handler
	categoryService service.CategoryService
}

func NewCategoryHandler(
	handler *Handler,
	categoryService service.CategoryService,
) *CategoryHandler {
	return &CategoryHandler{
		Handler:         handler,
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) GetAllCategory(ctx *gin.Context) {
	req := v1.GetCategoryRequest{}
	request.Assign(ctx, &req)
	res, err := h.categoryService.GetAllCategory(req, ctx)
	if err != nil {
		h.logger.WithContext(ctx).Error("categoryService.GetAllCategory", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, res)
}

func (h *CategoryHandler) GetCategory(ctx *gin.Context) {

}

func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	req := v1.CategoryRequest{}
	request.Assign(ctx, &req)
	res, err := h.categoryService.CreateCategory(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("categoryService.CreateCategory", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, res)
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	req := v1.CategoryRequest{}
	request.Assign(ctx, &req)
	err := h.categoryService.UpdateCategory(ctx, id, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("categoryService.UpdateCategory", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}
