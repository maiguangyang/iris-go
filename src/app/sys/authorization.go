package sys

import (
  "github.com/kataras/iris/context"
  Auth "../../authorization"
)
/**
 * 授权验证
 */
func Authorization(ctx context.Context) {
  r := Auth.SysVerify(ctx.GetHeader("Authorization"), ctx.RemoteAddr(), ctx)
  if r != nil {
    // r 已经过Utils.NewResData加工
    ctx.JSON(r)
    return
  }

  ctx.Next()
}