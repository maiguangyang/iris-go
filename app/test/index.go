package test

import(
  "fmt"
  "github.com/kataras/iris/context"
  Utils "../../utils"
  Auth "../../authorization"
)

// 检测是否设置数据库
func CheckDataBase(ctx context.Context) {
  fmt.Println(Auth.SetToken(context.Map{"name":1}, "user"))

  data := context.Map{ "has": false }
  ctx.JSON(Utils.NewResData(200, data, ctx))
}