package db

import (
	"context"
	"fmt"
	"gitee.com/goldden-go/goldden-go/pkg/utils/logger"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GormMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		defer func() {
			c.Set("DB", DB.WithContext(ctx))
		}()
		logger.Debug("设置数据库接口成功！！！")
		goldden_claims_I, exists := c.Get("goldden_claims")
		if !exists {
			return
		}
		goldden_claims, ok := goldden_claims_I.(jwtgo.MapClaims)
		if !ok {
			logger.Error("转换goldden_claims失败")
			return
		}
		if goldden_claims["name"] != nil {
			ctx = context.WithValue(ctx, "userid", fmt.Sprintf("%v", goldden_claims["name"]))
			ctx = context.WithValue(ctx, "username", fmt.Sprintf("%v", goldden_claims["name"]))
		} else {
			ctx = context.WithValue(ctx, "userid", "nobody")
			ctx = context.WithValue(ctx, "username", "nobody")
		}
		if goldden_claims["display_name"] != nil {
			ctx = context.WithValue(ctx, "username", fmt.Sprintf("%v", goldden_claims["display_name"]))
		}

	}

}
