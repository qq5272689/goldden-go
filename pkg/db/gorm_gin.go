package db

import (
	"context"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GormMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		goldden_claims, ok := c.MustGet("goldden_claims").(jwtgo.MapClaims)
		if !ok {
			return
		}
		if username, exists := goldden_claims["username"].(string); exists {
			ctx = context.WithValue(ctx, "username", username)
		} else {
			ctx = context.WithValue(ctx, "username", "nobody")
		}
		if userid, exists := goldden_claims["userid"].(string); exists {
			ctx = context.WithValue(ctx, "userid", fmt.Sprintf("%v", userid))
		} else {
			ctx = context.WithValue(ctx, "userid", "nobody")
		}
		c.Set("DB", DB.WithContext(ctx))
	}

}
