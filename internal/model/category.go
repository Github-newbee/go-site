package model

type Category struct {
	BaseModel

	CategoryName string `json:"category_name" gorm:"type:varchar(100);unique;not null" `
	Description  string `json:"description" gorm:"type:text" `
	Status       int    `json:"status" gorm:"not null; type:tinyint(1); default:1; comment:状态 1:正常 2:禁用"`
}
