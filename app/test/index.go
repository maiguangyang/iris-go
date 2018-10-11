package test

import(
  "github.com/kataras/iris/context"
  Public "../../public"
  Utils "../../utils"
)

// 检测是否设置数据库
func CheckDataBase(ctx context.Context) {
  data := context.Map{ "has": false }
  ctx.JSON(Utils.NewResData(200, data, ctx))
}

// 解密前端数据
func CheckDataBasePost(ctx context.Context) {
  data, err := Public.DecryptReqData(ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(403, data, ctx))
    return
  }
  ctx.JSON(Utils.NewResData(200, data, ctx))

}