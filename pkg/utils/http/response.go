package http

import (
	"net/http"

	"gitee.com/golden-go/golden-go/pkg/utils/types"
	"github.com/gin-gonic/gin"
)

type HttpResult struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func CommonSuccessResult(data interface{}) HttpResult {
	return HttpResult{
		Code:    20000,
		Data:    data,
		Message: "OK",
	}
}

func CommonSuccessPageResult(total int, items []interface{}) HttpResult {
	type data struct {
		Items []interface{} `json:"items"`
		Total int           `json:"total"`
	}
	return CommonSuccessResult(data{items, total})
}

func CommonFailResult(err string) HttpResult {
	return HttpResult{
		Code:    50000,
		Message: "err:" + err,
	}
}

func CommonErrResult(err error) HttpResult {
	return HttpResult{
		Code:    50000,
		Message: "err:" + err.Error(),
	}
}

func CommonSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, CommonSuccessResult(data))
}

func CommonSuccessPageResponse(c *gin.Context, total int, items []interface{}) {
	c.JSON(http.StatusOK, CommonSuccessPageResult(total, items))
}

func CommonFailResponse(c *gin.Context, err string) {
	c.JSON(http.StatusOK, CommonFailResult(err))
}

func CommonErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusOK, CommonErrResult(err))
}

func CommonFailCodeResponse(c *gin.Context, code int, err string) {
	r := CommonFailResult(err)
	r.Code = code
	c.JSON(http.StatusOK, CommonFailResult(err))
}

func CommonErrorCodeResponse(c *gin.Context, code int, err error) {
	r := CommonErrResult(err)
	r.Code = code
	c.JSON(http.StatusOK, CommonErrResult(err))
}

func NewTableData(data interface{}, pageNo, pageSize, count int) (td *types.TableData) {
	td = &types.TableData{
		Data:       data,
		PageNo:     pageNo,
		PageSize:   pageSize,
		TotalCount: count,
		TotalPage:  count / pageSize,
	}
	if count%pageSize != 0 {
		td.TotalPage += 1
	}
	return td
}
