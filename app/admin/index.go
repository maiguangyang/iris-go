package admin

import (
	"encoding/json"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/kataras/iris.v6"
)

const (
	SecretKey = "Hv0B2wCIKkNL8KdRpkHDXY8DXBniz2Ft"
)

/**
 * 生成json web token
 */
func GenToken() string {
	// claims := &jwt.StandardClaims{
	// 	NotBefore: int64(time.Now().Unix()),
	// 	ExpiresAt: int64(time.Now().Unix() + 1000),
	// 	Issuer:    "Token",
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":    "123456",
		"timestamp": int64(time.Now().Unix()),
	})
	ss, err := token.SignedString([]byte(SecretKey))

	if err != nil {
		return ""
	}

	return ss
}

/**
 * 校验token是否有效
 */
func keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(SecretKey), nil
}

func parseToken(unparseToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(unparseToken, keyFunc)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

/**
 * 返回类型
 */
func responseJsonStr(text string) string {
	jsonStr := `{
    "code": 401,
    "data": "Authorization ` + text + `",
    "msg" : "error"
  }`

	return jsonStr
}

/**
 * 授权验证
 */
func Authorization(ctx *iris.Context) {

	authHeader := ctx.RequestHeader("Authorization")
	var responseData interface{}

	if authHeader == "" || len(authHeader) <= 7 {
		text := responseJsonStr("未授权")
		err := json.Unmarshal([]byte(text), &responseData)

		if err == nil {
			ctx.JSON(iris.StatusOK, responseData)
		}

		return
	}

	claims, err := parseToken(authHeader[7:])

	if err != nil {
		text := responseJsonStr("授权已失效")
		err := json.Unmarshal([]byte(text), &responseData)

		if err == nil {
			ctx.JSON(iris.StatusOK, responseData)
		}
		return
	}

	println(claims["userid"].(string))

	ctx.Next()

}

/**
 * 控制器
 */

func SetSession(ctx *iris.Context) {
	token := GenToken()
	ctx.Writef(token)
}

func GetSession(ctx *iris.Context) {
	ctx.Writef("Get：%s", ctx.Session().GetString("name"))
}

func Index(ctx *iris.Context) {
	ctx.Writef("/")
}

func Profile(ctx *iris.Context) {
	ctx.Writef("profile")
}
