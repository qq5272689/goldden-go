package gin_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qq5272689/goutils/logger"
	"net/http"
)

type JsonData interface {
	SetErr(err error)
}

func GinZapRecovery(jd JsonData) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err.(error).Error())
				jd.SetErr(err.(error))
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, jd)
			}
		}()
		c.Next()
	}
}
