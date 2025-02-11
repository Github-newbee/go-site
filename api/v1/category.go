package v1

import "go-my-demo/pkg/request"

type CategoryRequest struct {
	CategoryName string `json:"category_name"  `
	Description  string `json:"description" `
	Status       int    `json:"status" `
}

type GetCategoryRequest struct {
	request.BaseFindRequest
	CategoryName string `form:"category_name" column:"category_name" operate:"$contains"`
}

type CreateCategoryResponse struct {
	Response
	Data interface{}
}
