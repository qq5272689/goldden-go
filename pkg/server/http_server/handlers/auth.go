package handlers

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/qq5272689/goldden-go/pkg/service"
	"github.com/qq5272689/goldden-go/pkg/utils/auth"
	"github.com/qq5272689/goldden-go/pkg/utils/captcha"
	ghttp "github.com/qq5272689/goldden-go/pkg/utils/http"
	"github.com/qq5272689/goldden-go/pkg/utils/jwt"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"go.uber.org/zap"
)

// @Tags user API
// ShowAccount godoc
// @Summary 获取验证码
// @Description 获取验证码
// @Produce  json
// @Router /v1/verify [get]
// @Success 200 {object} response.HttpResult
func Verify(ctx *gin.Context) {
	id, bs, err := captcha.GetCaptcha(ctx).Generate()
	if err != nil {
		logger.Error("Generate captcha fail!!!", zap.Error(err))
		ghttp.CommonErrorCodeResponse(ctx, 50000, err)
		return
	}
	ctx.SetCookie("captchaid", id, 60, "", "", false, false)
	ghttp.CommonSuccessResponse(ctx, bs)
}

// @Tags 用户相关接口
// ShowAccount godoc
// @Summary 本地用户登录
// @Description 本地用户登录
// @Produce  json
// @Param data body auth.LoginData  true "登录信息"
// @Router /login/local [post]
// @Success 200 {object} response.HttpResult
func LoginLocal(ctx *gin.Context) {
	ld := &auth.LoginData{}
	if err := ghttp.GetBody(ctx, ld); err != nil {
		logger.Warn("调用服务 GetBody 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonErrorCodeResponse(ctx, 50000, err)
		return
	}
	captchaid, err := ctx.Cookie("captchaid")
	if err != nil {
		logger.Warn("验证码已过期!!!")
		ghttp.CommonFailCodeResponse(ctx, 50001, "验证码已过期!!!")
		return
	}
	cok := captcha.GetCaptcha(ctx).Verify(captchaid, ld.Verify, true)
	if !cok {
		logger.Warn("验证码验证失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50002, "验证码验证失败!!!")
		return
	}

	ok, _ := service.GetUserServiceDBWithContext(ctx).CheckPassword(ld.Name, ld.Password)
	if !ok {
		logger.Warn("用户名密码验证失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50003, "用户名密码验证失败!!!")
		return
	}
	u, err := service.GetUserServiceDBWithContext(ctx).GetUserWithName(ld.Name)
	if err != nil {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50004, "获取用户信息失败!!!")
		return
	}
	u.Password = ""
	goldden_jwt, ok := ctx.MustGet("goldden_jwt").(*jwt.GolddenJwt)
	if !ok {
		logger.Warn("获取JWT失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50005, "获取JWT失败!!!")
		return
	}
	claims := jwtgo.MapClaims{}
	mapstructure.Decode(u, &claims)
	tokenStr, _ := goldden_jwt.CreateToken(claims)
	ghttp.CommonSuccessResponse(ctx, tokenStr)
}

// @Tags user API
// ShowAccount godoc
// @Summary 获取用户信息
// @Description 获取用户信息
// @Produce  json
// @Router /v1/userinfo [get]
// @Success 200 {object} response.HttpResult
func UserInfo(ctx *gin.Context) {
	goldden_claims, ok := ctx.MustGet("goldden_claims").(jwtgo.MapClaims)
	if !ok {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50000, "获取用户信息失败!!!")
		return
	}
	ghttp.CommonSuccessResponse(ctx, goldden_claims)
}

// @Tags user API
// ShowAccount godoc
// @Summary 登出
// @Description 登出
// @Produce  json
// @Router /v1/logout [get]
// @Success 200 {object} response.HttpResult
func LogOut(ctx *gin.Context) {
	ctx.SetCookie("goldden_key", "", 0, "", "", false, false)
	ghttp.CommonSuccessResponse(ctx, nil)
}

//
//func UserInfo(rw http.ResponseWriter, r *http.Request) {
//	re := new(HttpResult)
//	defer func() {
//		res, _ := json.Marshal(re)
//		rw.Write(res)
//	}()
//
//	ss := sessions.GetSession(r)
//
//	//u := &models.User{SuperAdmin: true}
//	//u.GetUser()
//	//u.Password = ""
//	//ss.Set("userinfo", *u)
//	uim := ss.Get("userinfo")
//	if uim == nil {
//		re.Info = "Get userinfo fail!!!"
//		rw.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//	bd, _ := json.Marshal(uim)
//	ui := &models.User{}
//	json.Unmarshal(bd, ui)
//
//	switch r.Method {
//	case "GET":
//		params := mux.Vars(r)
//		if params["list"] == "list" {
//			us, err := user.ListUser()
//			if err != nil {
//				re.Info = "ListUser err"
//				return
//			}
//			re.Data = us
//			re.Result = true
//		} else {
//			re.Data = ui
//			re.Result = true
//		}
//
//	case "POST":
//		if !ui.SuperAdmin {
//			re.Info = "you are not admin!!!"
//			rw.WriteHeader(http.StatusForbidden)
//			return
//		}
//		u, err := load_rbody_data(r)
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		err = user.CreateUser(u)
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		re.Result = true
//	case "PUT":
//		u, err := load_rbody_data(r)
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		var userd *models.User
//		userd, err = user.ToUser(u)
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		uid := ui
//		if uid.ID != userd.ID && !uid.SuperAdmin {
//			re.Info = "you are not self or admin!!!"
//			rw.WriteHeader(http.StatusForbidden)
//			return
//		}
//		if userd.OldPassword != "" {
//			ok, _ := user.CheckPassword(userd.Name, userd.OldPassword)
//			if !ok {
//				re.Info = "CheckPassword fail!!!"
//				return
//			}
//		}
//		err = user.UpdateUser(u)
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		re.Result = true
//	case "DELETE":
//		if !ui.SuperAdmin {
//			re.Info = "you are not admin!!!"
//			rw.WriteHeader(http.StatusForbidden)
//			return
//		}
//		r.ParseForm()
//		id := r.Form.Get("id")
//		id_int, err := strconv.Atoi(id)
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		err = user.DelUser(int64(id_int))
//		if err != nil {
//			re.Info = err.Error()
//			return
//		}
//		re.Result = true
//
//	}
//}
//
//func UserGroup(rw http.ResponseWriter, r *http.Request) {
//	re := new(HttpResult)
//	defer func() {
//		res, _ := json.Marshal(re)
//		rw.Write(res)
//	}()
//
//	ss := sessions.GetSession(r)
//	uim := ss.Get("userinfo")
//	if uim == nil {
//		re.Info = "Get userinfo fail!!!"
//		rw.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//	bd, _ := json.Marshal(uim)
//	ui := &models.User{}
//	json.Unmarshal(bd, ui)
//	bd, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		re.Info = "read body err:" + err.Error()
//		return
//	}
//	type data struct {
//		Group []string `json:"group"`
//		Id    []string `json:"id"`
//	}
//
//	d := &data{}
//	err = json.Unmarshal(bd, d)
//
//	if err != nil {
//		logger.Error("json body err", zap.Error(err), zap.String("body", string(bd)))
//		re.Info = "json body err:" + err.Error()
//		return
//	}
//	ids := []interface{}{}
//	for _, dd := range d.Id {
//		id, err := strconv.Atoi(dd)
//		if err != nil {
//			logger.Error("to int err", zap.Error(err), zap.String("id", dd))
//			continue
//		}
//		ids = append(ids, id)
//	}
//	ui.Groups = d.Group
//	err = user.UpdateUserGroups(ui, ids)
//	if err != nil {
//		logger.Error("UpdateUserGroups err", zap.Error(err))
//		re.Info = "json body err:" + err.Error()
//		return
//	}
//	re.Result = true
//}
