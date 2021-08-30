package http

import (
	"encoding/json"
	"io/ioutil"

	"gitee.com/golden-go/golden-go/pkg/utils/logger"
	"github.com/gin-gonic/gin"
)

func GetBody(ctx *gin.Context, v interface{}) error {
	req_data, _ := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body.Close()
	if err := json.Unmarshal(req_data, v); err != nil {
		logger.Warn("json.Unmarshal Fail！！！data:" + string(req_data))
		// CommonFailResponse(ctx, err.Error())
		return err
	}
	return nil
}
