package db

import (
	"context"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
)

func GormMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Set("DB", DB.WithContext(ctx))
		goldden_claims_I, exists := c.Get("goldden_claims")
		if !exists {
			return
		}
		goldden_claims, ok := goldden_claims_I.(jwtgo.MapClaims)
		if !ok {
			logger.Error("转换goldden_claims失败")
			return
		}
		if goldden_claims["username"] != nil {
			if username, exists := goldden_claims["username"].(string); exists {
				ctx = context.WithValue(ctx, "username", username)
			} else {
				ctx = context.WithValue(ctx, "username", "nobody")
			}
		} else {
			ctx = context.WithValue(ctx, "username", "nobody")
		}
		if goldden_claims["userid"] != nil {
			if userid, exists := goldden_claims["userid"].(string); exists {
				ctx = context.WithValue(ctx, "userid", fmt.Sprintf("%v", userid))
			} else {
				ctx = context.WithValue(ctx, "userid", "nobody")
			}
		} else {
			ctx = context.WithValue(ctx, "userid", "nobody")
		}

	}

}
