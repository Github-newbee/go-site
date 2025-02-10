package model

import (
	"go-my-demo/pkg/sid"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        uint64    `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
	// 如果模型有DeletedAt字段，将自动获得软删除的功能！
	// 在调用Delete时不会从数据库永久删除，而是只将字段DeletedAt的值设置为当前时间
	DeletedAt gorm.DeletedAt `json:"-"`
}

// 生成唯一int类型的id作为主键
func (obj *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	sidGen := sid.NewSid()
	obj.Id, err = sidGen.GenUint64()
	if err != nil {
		return err
	}
	return nil
}

func (obj *BaseModel) BeforeUpdate(db *gorm.DB) (err error) {
	obj.UpdatedAt = time.Now()
	return nil
}
