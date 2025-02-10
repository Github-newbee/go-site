package v1

import "go-my-demo/pkg/request"

type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// 查询参数标识 column 数据库列名 operate 操作符 $contains 表示包含 $in 表示in $gte 表示大于等于 $lte 表示小于等于
type GetAllUsersRequest struct {
	request.BaseFindRequest
	Username string `form:"username" column:"username" operate:"$contains"`
}

type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

// min 3 max 20 表示长度限制
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"alan" binding:"min=3,max=20"`
}

type GetProfileResponse struct {
	Response
	Data interface{}
}
