package admin

import(
  // "fmt"
  "time"
  // "reflect"
  // "encoding/json"
  "github.com/kataras/iris/context"
  // Public "../../public"
  Auth "../../authorization"
  DB "../../database"
  Utils "../../utils"
)

type IdpAdminAuth struct {
  Id int64 `json:"id" gorm:"primary_key;"`
  Rid int64 `json:"rid"`
  Sid string `json:"sid"`
  Content string `json:"content"`
  Auth int64 `json:"auth"`
  DeletedAt *time.Time `json:"deleted_at"`
  UpdatedAt time.Time `json:"updated_at"`
  CreatedAt time.Time `json:"created_at"`
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
  result := DB.EngineBak.Model(&list).Order("id desc").Limit(count).Offset(offset).Find(&list).Count(&total)

  if result.Error != nil {
    data = Utils.NewResData(1, "return data is empty.", ctx)
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

  result := DB.EngineBak.Where("id =?", table.Id).First(&table)

  if result.Error != nil {
    ctx.JSON(Utils.NewResData(1, "return data is empty.", ctx))
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
  if err := DB.EngineBak.Where("rid =?", value...).First(&exist).Error; err == nil {
    // return Utils.NewResData(1, "已分配权限", ctx)
    if err := DB.EngineBak.Model(&table).Where("rid =?", table.Rid).Updates(&table).Error; err != nil {
      return Utils.NewResData(1, "修改失败", ctx)
    }
    return Utils.NewResData(0, "修改成功", ctx)
  }

  // 新增
  if err := DB.EngineBak.Create(&table).Error; err != nil {
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
  if err := DB.EngineBak.Where("id=?", table.Id).First(&table).Error; err != nil {
    ctx.JSON(Utils.NewResData(1, "该信息不存在", ctx))
    return
  }

  // 开始删除
  data := context.Map{}
  if err := DB.EngineBak.Where("id =?", table.Id).Delete(&table).Error; err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    data = Utils.NewResData(0, "删除成功", ctx)
  }

  ctx.JSON(data)
}
