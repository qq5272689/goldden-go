package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/qq5272689/goldden-go/pkg/models"
	"github.com/qq5272689/goldden-go/pkg/utils/captcha"
	ghttp "github.com/qq5272689/goldden-go/pkg/utils/http"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
)

type LoginData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Verify   string `json:"verify"`
}

func Verify() {
	id, bs, err := captcha.MemCaptcha.Generate()
	if err != nil {
		fmt.Println("Generate mem captcha fail!!! err:", err)
		re.Info = err.Error()
		//rw.WriteHeader(502)
	}
	ss.Set("captchaid", id)
	//rw.Write([]byte(bs))
	re.Data = bs
	re.Result = true
}

func LoginLocal(ctx *gin.Context) {
	token, err := ctx.Cookie("goldden_key")
	//没有登录
	if err != nil || token == "" {
		ld := &LoginData{}
		if err := ghttp.GetBody(ctx, ld); err != nil {
			logger.Warn("调用服务 GetBody 错误!!!错误信息：", zap.Error(err))
			ghttp.CommonFailResponse(ctx, err.Error())
			return
		}
	}

	ss := sessions.GetSession(r)
	captchaid, str := ss.Get("captchaid").(string)
	if !str {
		re.Info = "get captcha id fail!!!"
		return
	}
	cok := captcha.MemCaptcha.Verify(captchaid, ld.Verify, true)
	if !cok {
		re.Info = "Check Verify fail!!!"
		return
	}

	ok, _ := user.CheckPassword(ld.Name, ld.Password)
	if !ok {
		re.Info = "CheckPassword fail!!!"
		return
	}
	u, err := user.GetUser(ld.Name)
	if err != nil {
		re.Info = "GetUser fail!!! err:" + err.Error()
		return
	}
	//ss := sessions.GetSession(r)
	u.Password = ""
	ss.Set("userinfo", *u)
	re.Result = true
}

func LogOut(ctx *gin.Context) {
	ctx.SetCookie("goldden_key", "", 0, "", "", false, false)
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
