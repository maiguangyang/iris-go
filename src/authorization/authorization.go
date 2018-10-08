// 1、生成token -> 记录token
// 2、验证长度 -> 检查token记录 -> 校验token -> 检查ip地址

package authorization

import (
  "time"
  "strconv"
  "github.com/kataras/iris/context"
  jwt "github.com/dgrijalva/jwt-go"

  Database "../database"
  Utils "../utils"
)

/**
 * 生成json web token
 * 生成token -> 记录token
 */
func SetToken(str interface{}, key []byte, uid int64, loginType int, content, ip string) string {
  timeNow := time.Now().Unix()
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "content"   : str,
    "nbf"       : int64(timeNow),
    "exp"       : int64(timeNow + 60 * 60 * 24),
    "timestamp" : int64(timeNow),
  })

  ss, err := token.SignedString(key)

  errText := "登录类型：" + Utils.Filter("login_type", loginType) + "；uid：" + strconv.FormatInt(uid, 10)
  if err != nil {
    Database.ErrorLogAdd(1, ss, "插入登录记录出错：" + err.Error() + "；" + errText)
    return ""
  }

  // 添加登录记录
  err = Database.LoginLogAdd(uid, loginType, ss, content, ip)
  if err != nil {
    Database.ErrorLogAdd(1, ss, "插入登录记录出错：" + err.Error() + "；" + errText)
    return ""
  }

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
func Verify(authHeader, ip string, key []byte, ctx context.Context) interface{} {
  if authHeader == "" || len(authHeader) <= 7 {
    return Utils.NewResData(401, "未登录", ctx)
  }

  token := authHeader[7:]

  // 检查记录
  if Database.LoginLogHas(token, ip) == false {
    return Utils.NewResData(401, "登录授权已失效", ctx)
  }

  // 校验token
  _, err := ParseToken(token, key)
  if err != nil {
    return Utils.NewResData(401, "登录授权已失效", ctx)
  }

  return nil
}
