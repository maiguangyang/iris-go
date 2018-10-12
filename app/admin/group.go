package admin

import (
  "fmt"
  "github.com/kataras/iris/context"
  Auth "../../authorization"
)

// 用户组列表
func GroupList (ctx context.Context) {
  userinfo, _ := Auth.DecryptToken(ctx.GetHeader("Authorization"), "admin")
  reqData     := userinfo.(map[string]interface{})

  fmt.Println(reqData)

}