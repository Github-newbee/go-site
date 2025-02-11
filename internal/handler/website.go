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

func (h *WebsiteHandler) GetWebsite(ctx *gin.Context) {

}
