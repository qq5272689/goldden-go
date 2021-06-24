package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	HandleUserCode string         `json:"handle_user_code" gorm:"column:handle_user_code" swaggerignore:"true"`
	HandleUserName string         `json:"handle_user_name" gorm:"column:handle_user_name"` //上次操作用户
	CreatedAt      time.Time      `json:"create_time" gorm:"column:create_time"`           //创建时间
	UpdatedAt      time.Time      `json:"update_time" gorm:"column:update_time"`           //更新时间
	DeletedAt      gorm.DeletedAt `json:"deleted_at"  gorm:"index" swaggertype:"string" swaggerignore:"true"`
}

func (b *BaseModel) BeforeSave(tx *gorm.DB) error {
	return b.setUser(tx)
}

func (b *BaseModel) setUser(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	if userid := ctx.Value("userid"); userid != nil {
		b.HandleUserCode = fmt.Sprintf("%v", userid)
	}
	if username := ctx.Value("username"); username != nil {
		b.HandleUserName = fmt.Sprintf("%v", username)
	}
	return nil
}
