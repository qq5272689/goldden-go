package captcha

import (
	"gitee.com/golden-go/golden-go/pkg/utils/jwt"
	"gitee.com/golden-go/golden-go/pkg/utils/logger"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type CookieStore struct {
	Ctx *gin.Context
}

func (cs *CookieStore) Set(id string, value string) {
	golden_jwt_I, exists := cs.Ctx.Get("golden_jwt")
	if !exists {
		logger.Error("golden_jwt doesn't exist")
		return
	}
	gj, ok := golden_jwt_I.(*jwt.GoldenJwt)
	if !ok {
		logger.Error("golden_jwt doesn't exist")
		return
	}
	tokenStr, err := gj.CreateToken(jwtgo.MapClaims{"captcha_id": id, id: value})
	if err != nil {
		logger.Error("CreateToken fail", zap.Error(err))
		return
	}
	cs.Ctx.SetCookie("golden_captcha", tokenStr, gj.Exp, "", "", false, false)
}

func (cs *CookieStore) Get(id string, clear bool) string {
	tokenStr, err := cs.Ctx.Cookie("golden_captcha")
	if err != nil {
		logger.Error("获取 golden_captcha cookie失败", zap.Error(err))
		return ""
	}
	golden_jwt_I, exists := cs.Ctx.Get("golden_jwt")
	if !exists {
		logger.Error("golden_jwt doesn't exist")
		return ""
	}
	gj, ok := golden_jwt_I.(*jwt.GoldenJwt)
	if !ok {
		logger.Error("golden_jwt doesn't exist")
		return ""
	}
	claims, err := gj.GetClaimsFromToken(tokenStr)
	if err != nil {
		logger.Error("解析token失败", zap.Error(err))
		return ""
	}
	defer func() {
		if clear {
			delete(claims, "captcha_id")
			delete(claims, id)
			tokenStr, err := gj.CreateToken(claims)
			if err != nil {
				logger.Error("CreateToken fail", zap.Error(err))
			}
			cs.Ctx.SetCookie("golden_captcha", tokenStr, gj.Exp, "", "", false, false)
		}
	}()
	if claims["captcha_id"] == nil {
		logger.Error("获取数据失败", zap.String("id", id))
		return ""
	}
	value_I := claims[id]
	if value_I == nil {
		logger.Error("获取值失败", zap.String("id", id), zap.Any("value", claims[id]))
		return ""
	}
	value, ok := value_I.(string)
	if !ok {
		logger.Error("获取值失败", zap.String("id", id), zap.Any("value", claims[id]))
		return ""
	}

	return value
}

func (cs *CookieStore) Verify(id, answer string, clear bool) bool {
	value := cs.Get(id, clear)
	if value != "" && value == answer {
		return true
	}
	return false
}

func GetCaptcha(ctx *gin.Context) *base64Captcha.Captcha {
	var store base64Captcha.Store
	store = base64Captcha.DefaultMemStore
	if ctx != nil {
		logger.Debug("cookiestore")
		store = &CookieStore{Ctx: ctx}
	}
	return base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, store)
}
