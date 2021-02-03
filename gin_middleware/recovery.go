package gin_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/qq5272689/goutils/logger"
	"go.uber.org/zap"
	"net/http"
)

type JsonData interface {
	SetErr(err error)
}

func GinZapRecovery(log *zap.Logger, jd JsonData) gin.HandlerFunc {
	logger.SetLogger(log)
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
