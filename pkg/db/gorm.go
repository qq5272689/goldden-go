package db

import (
	"github.com/qq5272689/goldden-go/pkg/models"

	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	ModelWithHistory = []interface{}{&models.User{}}
)
