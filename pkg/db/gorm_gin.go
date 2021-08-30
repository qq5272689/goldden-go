package db

import (
	"context"
	"fmt"

	"gitee.com/golden-go/golden-go/pkg/utils/logger"
	"github.com/gin-gonic/gin"
	jwtgo "github.com/golang-jwt/jwt"
)

func GormMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		defer func() {
			c.Set("DB", DB.WithContext(ctx))
		}()
		logger.Debug("设置数据库接口成功！！！")
		golden_claims_I, exists := c.Get("golden_claims")
		if !exists {
			return
		}
		golden_claims, ok := golden_claims_I.(jwtgo.MapClaims)
		if !ok {
			logger.Error("转换golden_claims失败")
			return
		}
		if golden_claims["name"] != nil {
			ctx = context.WithValue(ctx, "userid", fmt.Sprintf("%v", golden_claims["name"]))
			ctx = context.WithValue(ctx, "username", fmt.Sprintf("%v", golden_claims["name"]))
		} else {
			ctx = context.WithValue(ctx, "userid", "nobody")
			ctx = context.WithValue(ctx, "username", "nobody")
		}
		if golden_claims["display_name"] != nil {
			ctx = context.WithValue(ctx, "username", fmt.Sprintf("%v", golden_claims["display_name"]))
		}

	}

}
