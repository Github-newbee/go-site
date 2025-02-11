package model

type Website struct {
	BaseModel

	Name        string   `json:"name" gorm:"type:varchar(200);not null"`
	URL         string   `json:"link" gorm:"type:varchar(255);not null"`
	Icon        string   `json:"icon" gorm:"type:varchar(255)"`
	CategoryID  string   `json:"category_id" gorm:"index"`
	Description string   `json:"description" gorm:"type:text"`
	Category    Category `json:"-" gorm:"foreignKey:CategoryID"`
	Status      int      `json:"status" gorm:"not null; type:tinyint(1); default:1; comment:状态 1:正常 2:禁用"`
}
