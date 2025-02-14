package handler

import (
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/service"
	"go-my-demo/pkg/request"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebsiteHandler struct {
	*Handler
	websiteService service.WebsiteService
}

func NewWebsiteHandler(
	handler *Handler,
	websiteService service.WebsiteService,
) *WebsiteHandler {
	return &WebsiteHandler{
		Handler:        handler,
		websiteService: websiteService,
	}
}

func (h *WebsiteHandler) CreateWebsite(ctx *gin.Context) {

	req := v1.WebsiteRequest{}
	request.Assign(ctx, &req)
	res, err := h.websiteService.CreateWebsite(ctx, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("websiteService.CreateWebsite", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, res)
}

func (h *WebsiteHandler) GetAllWebsite(ctx *gin.Context) {
	req := v1.GetWebsiteRequest{}
	queryErr := request.Assign(ctx, &req)
	if queryErr != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, queryErr.Error())
		return
	}

	res, err := h.websiteService.Get(ctx, req)

	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrInternalServerError, err.Error())
		return
	}

	v1.HandleSuccess(ctx, res)
}

func (h *WebsiteHandler) UpdateWebsite(ctx *gin.Context) {
	id := ctx.Param("id")
	req := v1.WebsiteRequest{}
	queryErr := request.Assign(ctx, &req)
	if queryErr != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, queryErr.Error())
		return
	}
	err := h.websiteService.Update(ctx, id, &req)
	if err != nil {
		h.logger.WithContext(ctx).Error("websiteService.Update", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}
	v1.HandleSuccess(ctx, nil)
}
