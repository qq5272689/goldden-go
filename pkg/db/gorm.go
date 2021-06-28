package db

import (
	"gitee.com/goldden-go/goldden-go/pkg/models"

	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	ModelWithHistory = []interface{}{&models.User{}}
)
