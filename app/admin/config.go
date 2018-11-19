package admin

import(
  "github.com/kataras/iris/context"
  // Public "../../public"
  DB "../../database"
  Utils "../../utils"
)

type IdpAuthSet struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  TableName string `json:"table_name"`
  Routes string `json:"routes"`
  Path string `json:"path"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

// 获取数据库所有表
func CongifTable(ctx context.Context) {
  allTable, _ := DB.Engine.DBMetas()

  array :=  []string{}
  for _,v := range allTable{
    array = append(array, v.Name)
    // tabMap[v.Name] = v.ColumnsSeq()
  }

  data := Utils.NewResData(0, array, ctx)
  ctx.JSON(data)
}

// 列表
func CongifRoutesList (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 获取分页、总数、limit
  page, count, limit, _ := DB.Limit(ctx)
  list := make([]IdpAuthSet, 0)

  // 获取统计总数
  var table IdpAuthSet
  data := context.Map{}

  total, err := DB.Engine.Count(&table)

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    // 获取列表
    err = DB.Engine.Desc("id").Limit(count, limit).Find(&list)

    // 返回数据
    if err != nil {
      data = Utils.NewResData(1, err.Error(), ctx)
    } else {
      resData := Utils.TotalData(list, page, total, count)
      data = Utils.NewResData(0, resData, ctx)
    }
  }

  ctx.JSON(data)

}

// 详情
func CongifRoutesDetail (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  var table IdpAuthSet
  ctx.ReadJSON(&table)

  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  has, err := DB.Engine.Get(&table)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  data := context.Map{}
  if has == true {
    data = Utils.NewResData(0, table, ctx)
  } else {
    data = Utils.NewResData(1, "记录不存在", ctx)
  }

  ctx.JSON(data)

}

// 新增
func CongifRoutesAdd (ctx context.Context) {
  data := sumbitCongifRoutesData(0, ctx)
  ctx.JSON(data)
}

// 修改
func CongifRoutesPut (ctx context.Context) {
  data := sumbitCongifRoutesData(1, ctx)
  ctx.JSON(data)
}


// 提交数据 0新增、1修改
func sumbitCongifRoutesData(tye int, ctx context.Context) context.Map {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  var table IdpAuthSet

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  // 验证参数
  var rules Utils.Rules
  rules = Utils.Rules{
    "Name": {
      "required": true,
    },
  }


  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    return Utils.NewResData(1, errMsgs, ctx)
  }

  // 判断数据库里面是否已经存在
  var exist IdpAuthSet
  value := []interface{}{table.Id, table.Name, table.TableName}
  has, err := DB.Engine.Where("id<>? and name=? and table_name=?", value...).Exist(&exist)

  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  if has == true {
    return Utils.NewResData(1, table.Name + "已存在", ctx)
  }

  // 写入数据库
  tipsText := "添加"
  if tye == 1 {
    tipsText = "修改"
    // 修改
    _, err = DB.Engine.Id(table.Id).Update(&table)
  } else {
    // 新增
    _, err = DB.Engine.Insert(&table)
  }


  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  return Utils.NewResData(0, tipsText + "成功", ctx)
}

// 删除
func CongifRoutesDel (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  var table IdpAuthSet

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAuthSet
  has, err := DB.Engine.Where("id=?", table.Id).Exist(&exist)

  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  if has != true {
    ctx.JSON(Utils.NewResData(1, "该信息不存在", ctx))
    return
  }

  // 开始删除
  _, err = DB.Engine.Id(table.Id).Delete(&table)

  data := context.Map{}
  if err == nil {
    data = Utils.NewResData(0, "删除成功", ctx)
  } else {
    data = Utils.NewResData(1, err.Error(), ctx)
  }

  ctx.JSON(data)
}
