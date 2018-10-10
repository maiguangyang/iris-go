// 1、生成token -> 记录token
// 2、验证长度 -> 检查token记录 -> 校验token -> 检查ip地址

package authorization

import (
  "time"
  "github.com/kataras/iris/context"
  jwt "github.com/dgrijalva/jwt-go"

  Utils "../utils"
)

// role group
var SecretKey = context.Map{
  "user": "85d34afdffc5d9ddfbd19e2bb31018cf",
  "admin": "85d34afdffc5d9ddfbd19e2bb31018cf",
}

// 检查user的Token
func CheckAuthUser(ctx context.Context) {
  Verify(ctx.GetHeader("Authorization"), []byte(SecretKey["user"].(string)), ctx)
}

// 检查admin的Token
func CheckAuthAdmin(ctx context.Context) {
  Verify(ctx.GetHeader("Authorization"), []byte(SecretKey["admin"].(string)), ctx)
}

/**
 * 统一使用
 * 生成json web token
 */
func SetToken(str interface{}, role string) string {
  timeNow := time.Now().Unix()
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "content"   : str,
    "nbf"       : int64(timeNow),
    "exp"       : int64(timeNow + 60 * 60 * 24),
    "timestamp" : int64(timeNow),
  })

  ss, _ := token.SignedString([]byte(SecretKey[role].(string)))
  return ss
}

/**
 * 校验token是否有效
 */
func ParseToken(tokenStr string, key []byte) (jwt.MapClaims, error) {
  token, err := jwt.Parse(tokenStr, func (token *jwt.Token) (interface{}, error) {
    return key, nil
  })

  if err != nil {
    return nil, err
  }

  claims := token.Claims.(jwt.MapClaims)

  return claims, nil
}

/**
 * token解密
 */
func DecryptToken(tokenStr, role string) (interface{}, error) {
  claims, err := GetTokenContent(tokenStr, []byte(SecretKey[role].(string)))
  return claims, err
}


/**
 * 获取token内容
 */
func GetTokenContent(authHeader string, key []byte) (interface{}, error) {
  claims, err := ParseToken(authHeader[7:], key)
  return claims["content"], err
}

/**
 * 授权验证
 * 验证长度 -> 检查token记录 -> 校验token -> 检查ip地址
 */
func Verify(authHeader string, key []byte, ctx context.Context) {
  if authHeader == "" || len(authHeader) <= 7 {
    ctx.JSON(Utils.NewResData(401, "未登录", ctx))
    return
  }

  token := authHeader[7:]

  // 校验token
  _, err := ParseToken(token, key)
  if err != nil {
    ctx.JSON(Utils.NewResData(401, "登录授权已失效", ctx))
    return
  }

  ctx.Next()
}
