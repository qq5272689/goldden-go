package db

import (
	"gitee.com/golden-go/golden-go/pkg/models"

	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	ModelWithHistory = []interface{}{&models.User{}}
)
