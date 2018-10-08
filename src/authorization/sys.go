package authorization

import (
  "github.com/kataras/iris/context"
)


/**
 *  系统配置token
 */
const (
  SysSecretKey = "85d34afdffc5d9ddfbd19e2bb31018cf"
)

/**
 * 生成json web token
 */
func SysSetToken(info interface{}, content, ip string) string {
  return SetToken(info, []byte(SysSecretKey), int64(0), int(0), content, ip)
}

/**
 * 获取token内容
 */
func GetSysTokenContent(authHeader string) (interface{}, error) {
  info, err := GetTokenContent(authHeader, []byte(SysSecretKey))
  return info, err
}

/**
 * 后台授权验证
 * Bearer
 */
func SysVerify(authHeader, ip string, ctx context.Context) interface{} {
  return Verify(authHeader, ip, []byte(SysSecretKey), ctx)
}
