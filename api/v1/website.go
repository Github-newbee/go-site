package v1

import (
	"go-my-demo/internal/model"
	"go-my-demo/pkg/request"
)

type WebsiteRequest struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type GetWebsiteRequest struct {
	request.BaseFindRequest
	// Name string `form:"name" column:"name" operate:"$contains"`
}

type WebsiteResponse struct {
	model.Website
	CategoryName string `json:"category_name"`
	CategoryDesc string `json:"category_desc"`
}
