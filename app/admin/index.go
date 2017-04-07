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
	claims := &jwt.StandardClaims{
		NotBefore: int64(time.Now().Unix()),
		ExpiresAt: int64(time.Now().Unix() + 1000),
		Issuer:    "hzwy23",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
	println(token.Header["Issuer"])
	return []byte(SecretKey), nil
}

func CheckToken(token string) bool {
	_, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return false
	}
	return true
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

	claims := CheckToken(authHeader[7:])

	println(claims)

	// if err != nil {
	// 	text := responseJsonStr("授权已失效")
	// 	err := json.Unmarshal([]byte(text), &responseData)

	// 	if err == nil {
	// 		ctx.JSON(iris.StatusOK, responseData)
	// 	}
	// 	return
	// }

	ctx.Next()

}

/**
 * 控制器
 */

func SetSession(ctx *iris.Context) {
	// ctx.Session().Set("name", "iris")
	// ctx.Writef("Set：%s", ctx.Session().GetString("name"))

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"foo":       "bar",
	// 	"timestamp": time.Now().Unix(),
	// })

	// tokenString, _ := token.SignedString([]byte(SecretKey))

	token := GenToken()
	println(token)

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
