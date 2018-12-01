package admin

import(
  // "fmt"
  // "reflect"
  // "encoding/json"
  "github.com/kataras/iris/context"
  // Public "../../public"
  Auth "../../authorization"
  DB "../../database"
  Utils "../../utils"
)

type IdpAdminAuth struct {
  DB.Model
  Rid int64 `json:"rid"`
  Sid string `json:"sid"`
  Content string `json:"content"`
  Auth int64 `json:"auth"`
}

// 列表
func RouteAuthList (ctx context.Context) {

  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_auth")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  // 获取分页、总数、limit
  page, count, offset, _ := DB.Limit(ctx)
  list := make([]IdpAdminAuth, 0)

  // 查询列表
  data := context.Map{}

  var total int64
  if err := DB.Engine.Model(&list).Order("id desc").Count(&total).Limit(count).Offset(offset).Find(&list).Error; err != nil {
    data = Utils.NewResData(1, err, ctx)
  } else {
    resData := Utils.TotalData(list, page, total, count)
    data = Utils.NewResData(0, resData, ctx)
  }

  ctx.JSON(data)
}

// 详情
func RouteAuthDetail (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_auth")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  var table IdpAdminAuth
  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  if err := DB.Engine.Where("id =?", table.Id).First(&table).Error; err != nil {
    ctx.JSON(Utils.NewResData(1, err, ctx))
    return
  }

  ctx.JSON(Utils.NewResData(0, table, ctx))
}

// 新增
func RouteAuthAdd (ctx context.Context) {
  data := sumbitRouteAuthData(0, ctx)
  ctx.JSON(data)
}

// 修改
func RouteAuthPut (ctx context.Context) {
  data := sumbitRouteAuthData(1, ctx)
  ctx.JSON(data)
}


// 提交数据 0新增、1修改
func sumbitRouteAuthData(tye int, ctx context.Context) context.Map {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_auth")
  if hasAuth != true {
    return Utils.NewResData(code, err.Error(), ctx)
  }

  var table IdpAdminAuth

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  // 不能修改自己所在组的权限
  reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  if Utils.StrToInt64(reqData["rid"].(string)) == table.Rid {
    return Utils.NewResData(1, "登陆账户属于该角色，无法修改权限", ctx)
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminAuth
  value := []interface{}{table.Rid}
  if err := DB.Engine.Where("rid =?", value...).First(&exist).Error; err == nil {
    // return Utils.NewResData(1, "已分配权限", ctx)
    if err := DB.Engine.Model(&table).Where("rid =?", table.Rid).Updates(&table).Error; err != nil {
      return Utils.NewResData(1, "修改失败", ctx)
    }
    return Utils.NewResData(0, "修改成功", ctx)
  }

  // 新增
  if err := DB.Engine.Create(&table).Error; err != nil {
    return Utils.NewResData(1, "添加失败", ctx)
  }

  return Utils.NewResData(0, "添加成功", ctx)

}

// 删除
func RouteAuthDel (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_auth")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  var table IdpAdminAuth

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 判断数据库里面是否已经存在
  if err := DB.Engine.Where("id=?", table.Id).First(&table).Error; err != nil {
    ctx.JSON(Utils.NewResData(1, "该信息不存在", ctx))
    return
  }

  // 开始删除
  data := context.Map{}
  if err := DB.Engine.Where("id =?", table.Id).Delete(&table).Error; err != nil {
    data = Utils.NewResData(1, err, ctx)
  } else {
    data = Utils.NewResData(0, "删除成功", ctx)
  }

  ctx.JSON(data)
}
