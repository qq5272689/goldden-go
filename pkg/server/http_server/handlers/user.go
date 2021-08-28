package handlers

import (
	"strconv"

	"gitee.com/goldden-go/goldden-go/pkg/models"
	"gitee.com/goldden-go/goldden-go/pkg/service"
	ghttp "gitee.com/goldden-go/goldden-go/pkg/utils/http"
	"gitee.com/goldden-go/goldden-go/pkg/utils/logger"
	"gitee.com/goldden-go/goldden-go/pkg/utils/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags 用户相关接口
// ShowAccount godoc
// @Summary 搜索用户
// @Description 搜索用户
// @Produce  json
// @Param filter query string  false "过滤关键词"
// @Param pageNo query []int  false "多个ID 每个ID之间用,分隔，例：123,233 注：跟 name 参数只有一个会生效，hostids参数优先级"
// @Param pageSize query int  false "单页条数"
// @Router /v1/user [get]
// @Success 200 {object} ghttp.HttpResult
func SearchUser(ctx *gin.Context) {
	filter := ctx.Query("filter")
	keyword := ctx.Query("keyword")
	if keyword != "" && filter == "" {
		filter = keyword
	}
	pageNo, err := strconv.Atoi(ctx.Query("pageNo"))
	if err != nil {
		pageNo = 1
	}
	if pageNo < 1 {
		pageNo = 1
	}
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil {
		pageSize = 1000
	}
	if pageSize < 1 {
		pageSize = 1
	}

	if d, err := service.GetUserServiceDBWithContext(ctx).SearchUser(filter, pageNo, pageSize); err != nil {
		logger.Warn("调用服务 SearchUser 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
	} else {
		ghttp.CommonSuccessResponse(ctx, d)
	}
}

// @Tags 用户相关接口
// ShowAccount godoc
// @Summary 获取用户
// @Description 获取用户
// @Produce  json
// @Param userid path int  false "用户ID"
// @Router /v1/user/{userid} [get]
// @Success 200 {object} ghttp.HttpResult
func GetUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("userid"))
	if err != nil {
		logger.Warn("get服务 id 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
		return
	}
	if d, err := service.GetUserServiceDBWithContext(ctx).GetUser(id); err != nil {
		logger.Warn("调用服务 SearchUser 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
	} else {
		ghttp.CommonSuccessResponse(ctx, d)
	}
}

func GetUserWithGroup(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("groupid"))
	if err != nil {
		logger.Warn("get group id 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
		return
	}
	if d, err := service.GetUserServiceDBWithContext(ctx).GetUserWithGroup(id); err != nil {
		logger.Warn("调用服务 GetUserWithGroup 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
	} else {
		ghttp.CommonSuccessResponse(ctx, d)
	}
}

// @Tags 用户相关接口
// ShowAccount godoc
// @Summary 创建用户
// @Description 创建用户
// @Produce  json
// @Param data body models.User  true "用户"
// @Router /v1/user [post]
// @Success 200 {object} ghttp.HttpResult
func CreateUser(ctx *gin.Context) {
	args := &models.User{}
	if err := ghttp.GetBody(ctx, args); err != nil {
		logger.Warn("调用服务 GetBody 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
		return
	}
	if err := service.GetUserServiceDBWithContext(ctx).CreateUser(args); err != nil {
		logger.Warn("调用服务 SearchUser 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
	} else {
		if d, err := service.GetUserServiceDBWithContext(ctx).SearchUser("", 1, 1000); err != nil {
			logger.Warn("调用服务 SearchUser 错误!!!错误信息：", zap.Error(err))
			ghttp.CommonFailResponse(ctx, err.Error())
		} else {
			ghttp.CommonSuccessResponse(ctx, d)
		}
	}
}

// @Tags 用户相关接口
// ShowAccount godoc
// @Summary 更新用户
// @Description 更新用户
// @Produce  json
// @Param data body models.User  true "用户"
// @Router /v1/user [put]
// @Success 200 {object} ghttp.HttpResult
func UpdateUser(ctx *gin.Context) {
	args := &models.User{}
	if err := ghttp.GetBody(ctx, args); err != nil {
		return
	}
	if err := service.GetUserServiceDBWithContext(ctx).UpdateUser(args); err != nil {
		logger.Warn("调用服务 SearchUser 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
	} else {
		if d, err := service.GetUserServiceDBWithContext(ctx).SearchUser("", 1, 1000); err != nil {
			logger.Warn("调用服务 SearchUser 错误!!!错误信息：", zap.Error(err))
			ghttp.CommonFailResponse(ctx, err.Error())
		} else {
			ghttp.CommonSuccessResponse(ctx, d)

		}
	}
}

// @Tags 用户相关接口
// ShowAccount godoc
// @Summary 删除user
// @Description 删除user
// @Produce  json
// @Param ids query []int  false "多个ID 每个ID之间用,分隔，例：123,233"
// @Router /v1/user [delete]
// @Success 200 {object} ghttp.HttpResult
func DeleteUser(ctx *gin.Context) {
	id_str := ctx.QueryArray("ids")
	ids, err := types.SliceStringToInt(id_str)
	if err != nil {
		logger.Warn("id，无法转化！！！", zap.Any("ids", id_str), zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
		return
	}
	if err := service.GetUserServiceDBWithContext(ctx).DelUser(ids); err != nil {
		logger.Warn("调用服务 DelUser 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonFailResponse(ctx, err.Error())
	} else {
		ghttp.CommonSuccessResponse(ctx, nil)
	}
}
