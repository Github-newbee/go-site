package model

type User struct {
	BaseModel

	Username string `json:"username" gorm:"unique;not null"`
	Nickname string `json:"nickname" gorm:"not null"`
	Password string `json:"-" gorm:"not null"`
	Status   int    `json:"status" gorm:"not null; type:tinyint(1); default:1; comment:状态 1:正常 2:禁用"`
}
