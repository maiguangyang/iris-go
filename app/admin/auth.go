package admin

import(
  "fmt"
  "github.com/kataras/iris/context"
  // Utils "../../utils"
  // Auth "../../authorization"
)

// 新增部门
func GroupAdd(ctx context.Context) {
  fmt.Println(ctx)
  // fmt.Println(Auth.SetToken(context.Map{"name":1}, "user"))

  // data := context.Map{ "has": false }
  // ctx.JSON(Utils.NewResData(200, data, ctx))
}