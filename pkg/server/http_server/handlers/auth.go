package handlers

import (
	"errors"

	"gitee.com/golden-go/golden-go/pkg/service"
	"gitee.com/golden-go/golden-go/pkg/utils/captcha"
	ghttp "gitee.com/golden-go/golden-go/pkg/utils/http"
	"gitee.com/golden-go/golden-go/pkg/utils/jwt"
	"gitee.com/golden-go/golden-go/pkg/utils/ldap"
	"gitee.com/golden-go/golden-go/pkg/utils/logger"
	"gitee.com/golden-go/golden-go/pkg/utils/types"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @Tags 登录相关接口
// ShowAccount godoc
// @Summary 获取验证码
// @Description 获取验证码
// @Produce  json
// @Router /v1/verify [get]
// @Success 200 {object} ghttp.HttpResult
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

// @Tags 登录相关接口
// ShowAccount godoc
// @Summary 本地用户登录
// @Description 本地用户登录
// @Produce  json
// @Param data body types.LoginData  true "登录信息"
// @Router /v1/login/local [post]
// @Success 200 {object} ghttp.HttpResult
func LoginLocal(ctx *gin.Context) {
	ld, err := loginFirstCheck(ctx)
	if err != nil {
		return
	}
	ok, _ := service.GetUserServiceDBWithContext(ctx).CheckPassword(ld.Name, ld.Password)
	if !ok {
		logger.Warn("用户名密码验证失败!!!")
		if viper.GetBool("auth.ldap.enable") {
			loginLdap(ctx, ld)
		} else {
			ghttp.CommonFailCodeResponse(ctx, 50003, "用户名密码验证失败!!!")
		}

		return
	}
	u, err := service.GetUserServiceDBWithContext(ctx).GetUserWithName(ld.Name)
	if err != nil {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50004, "获取用户信息失败!!!")
		return
	}
	u.Password = ""
	golden_jwt_I, exists := ctx.Get("golden_jwt")
	if !exists {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50005, "获取JWT失败!!!")
		return
	}
	golden_jwt, ok := golden_jwt_I.(*jwt.GoldenJwt)
	if !ok {
		logger.Warn("获取JWT失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50006, "获取JWT失败!!!")
		return
	}
	claims := jwtgo.MapClaims{}
	types.JsonStruct(u, &claims)
	tokenStr, _ := golden_jwt.CreateTokenAndSetCookie(claims, ctx)

	ghttp.CommonSuccessResponse(ctx, tokenStr)
}

func loginFirstCheck(ctx *gin.Context) (*types.LoginData, error) {
	ld := &types.LoginData{}
	if err := ghttp.GetBody(ctx, ld); err != nil {
		logger.Warn("调用服务 GetBody 错误!!!错误信息：", zap.Error(err))
		ghttp.CommonErrorCodeResponse(ctx, 50000, err)
		return nil, err
	}
	captchaid, err := ctx.Cookie("captchaid")
	if err != nil {
		logger.Warn("验证码已过期!!!")
		ghttp.CommonFailCodeResponse(ctx, 50001, "验证码已过期!!!")
		return nil, errors.New("验证码已过期!!!")
	}
	cok := captcha.GetCaptcha(ctx).Verify(captchaid, ld.Verify, true)
	if !cok {
		logger.Warn("验证码验证失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50002, "验证码验证失败!!!")
		return nil, errors.New("验证码验证失败!!!")
	}
	return ld, nil
}

func LoginLdap(ctx *gin.Context) {
	ld, err := loginFirstCheck(ctx)
	if err != nil {
		return
	}
	loginLdap(ctx, ld)
}

func loginLdap(ctx *gin.Context, ld *types.LoginData) {
	imli, ok := ctx.Get("IML")
	if !ok {
		logger.Warn("获取IML失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50006, "获取IML失败!!!")
		return
	}
	iml, ok := imli.(ldap.IMultiLDAP)
	if !ok {
		logger.Warn("转换IML失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50006, "转换IML失败!!!")
		return
	}
	u, err := iml.Login(ld)
	if err != nil {
		logger.Warn("LDAP登录失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50004, "LDAP登录失败!!!")
		return
	}
	golden_jwt_I, exists := ctx.Get("golden_jwt")
	if !exists {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50005, "获取JWT失败!!!")
		return
	}
	golden_jwt, ok := golden_jwt_I.(*jwt.GoldenJwt)
	if !ok {
		logger.Warn("获取JWT失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50006, "获取JWT失败!!!")
		return
	}
	claims := jwtgo.MapClaims{}
	types.JsonStruct(u, &claims)
	tokenStr, _ := golden_jwt.CreateTokenAndSetCookie(claims, ctx)

	ghttp.CommonSuccessResponse(ctx, tokenStr)
}

// @Tags 登录相关接口
// ShowAccount godoc
// @Summary 获取登录用户信息
// @Description 获取登录用户信息
// @Produce  json
// @Router /v1/userinfo [get]
// @Success 200 {object} ghttp.HttpResult
func UserInfo(ctx *gin.Context) {
	golden_claims_I, exists := ctx.Get("golden_claims")
	if !exists {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50000, "获取用户信息失败!!!")
		return
	}
	golden_claims, ok := golden_claims_I.(jwtgo.MapClaims)
	if !ok {
		logger.Warn("获取用户信息失败!!!")
		ghttp.CommonFailCodeResponse(ctx, 50001, "获取用户信息失败!!!")
		return
	}
	ghttp.CommonSuccessResponse(ctx, golden_claims)
}

// @Tags 登录相关接口
// ShowAccount godoc
// @Summary 登出
// @Description 登出
// @Produce  json
// @Router /v1/logout [get]
// @Success 200 {object} ghttp.HttpResult
func LogOut(ctx *gin.Context) {
	ctx.SetCookie("golden_key", "", 0, "", "", false, false)
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
