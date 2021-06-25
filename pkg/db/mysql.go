//+build mysql

package db

import (
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

func OpenDB(serviceName, dsn string) (err error) {

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: strings.ToLower(serviceName) + "_", // 表名前缀，`User` 的表名应该是 `t_users`
			//SingularTable: true,                              // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		logger.Error("Database connection failed.", zap.Error(err))
		return err
	}
	return nil
}

func SetupDatabase(db *gorm.DB) error {
	//db.Exec("create extension IF NOT EXISTS hstore;")
	//db.AutoMigrate(ModelNoHistory...)
	err := db.AutoMigrate(ModelWithHistory...)
	if err != nil {
		logger.Error("setup database failed.", zap.Error(err))
		return err
	}
	return nil
}
