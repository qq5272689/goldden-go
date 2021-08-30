package jwt

import (
	"crypto/rsa"
	"errors"
	"time"

	"gitee.com/golden-go/golden-go/pkg/utils/logger"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

type GoldenJwt struct {
	Exp        int
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

//func init() {
//	publicKeyByte, err := ioutil.ReadFile("公钥的路径/public.key")
//	if err != nil {
//		log.Println(err.Error())
//	}
//	publicKey, err = jwtgo.ParseRSAPublicKeyFromPEM(publicKeyByte)
//	privateKeyByte, err := ioutil.ReadFile("私钥的路径/private.key")
//	if err != nil {
//		log.Println(err.Error())
//	}
//	privateKey, _ = jwtgo.ParseRSAPrivateKeyFromPEM(privateKeyByte)
//}

func NewGoldenJwt(exp int, puk, prk string) (gj *GoldenJwt, err error) {
	gj = &GoldenJwt{Exp: exp}
	gj.publicKey, err = jwtgo.ParseRSAPublicKeyFromPEM([]byte(puk))
	if err != nil {
		return nil, err
	}
	gj.privateKey, err = jwtgo.ParseRSAPrivateKeyFromPEM([]byte(prk))
	if err != nil {
		return nil, err
	}
	return gj, nil
}

func (gj *GoldenJwt) GinJwtMiddleware(ctx *gin.Context) {
	ctx.Set("golden_jwt", gj)
	claims := jwtgo.MapClaims{}
	token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor, gj.keyFunc, request.WithClaims(&claims))
	if err == nil && token.Valid {
		ctx.Set("golden_claims", claims)
		return
	}
	golden_key, _ := ctx.Cookie("golden_key")
	claims, err = gj.GetClaimsFromToken(golden_key)
	if err == nil {
		ctx.Set("golden_claims", claims)
		return
	}
	logger.Info("token不存在")
}

// createToken 生成一个RS256验证的Token
// Token里面包括的值，可以自己根据情况添加，
func (gj *GoldenJwt) CreateToken(claims jwtgo.MapClaims) (tokenStr string, err error) {
	//	jwtgo.MapClaims{
	//	"iat":      time.Now().Unix(), // Token颁发时间
	//	//"nbf":      time.Now().Unix(), // Token生效时间
	//	"exp":      time.Now().Add(time.Hour * time.Duration(gj.Exp)).Unix(), // Token过期时间，目前是24小时
	//	//"iss":      "liwenbin.com", // 颁发者
	//	//"sub":      "AuthToken", // 主题
	//	//  "role":     "guest", // 角色（附加）
	//}
	now := time.Now()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(time.Minute * time.Duration(gj.Exp)).Unix()
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS512, claims)
	return token.SignedString(gj.privateKey)
}

// createToken 生成一个RS256验证的Token
// Token里面包括的值，可以自己根据情况添加，
func (gj *GoldenJwt) CreateTokenAndSetCookie(claims jwtgo.MapClaims, ctx *gin.Context) (tokenStr string, err error) {
	tokenStr, err = gj.CreateToken(claims)
	if err != nil {
		return
	}
	ctx.SetCookie("golden_key", tokenStr, gj.Exp*60, "", "", false, true)
	return
}

func (gj *GoldenJwt) keyFunc(token *jwtgo.Token) (interface{}, error) {
	// 基于JWT的第一部分中的alg字段值进行一次验证
	if _, ok := token.Method.(*jwtgo.SigningMethodRSA); !ok {
		return nil, errors.New("验证Token的加密类型错误")
	}
	return gj.publicKey, nil
}

// getSubFromToken 获取Token的主题（也可以更改获取其他值）
// 参数tokenStr指的是 从客户端传来的待验证Token
// 验证Token过程中，如果Token生成过程中，指定了iat与exp参数值，将会自动根据时间戳进行时间验证
func (gj *GoldenJwt) GetClaimsFromToken(tokenStr string) (claims jwtgo.MapClaims, err error) {
	// 基于公钥验证Token合法性
	token, err := jwtgo.Parse(tokenStr, gj.keyFunc)
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwtgo.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Token无效或者无对应值")
}
