package admin

import(
  // "fmt"
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
  SubId string `json:"sub_id"`
  SubName string `json:"sub_name" xorm:"<-"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

// 获取数据库所有表
func AuthSetTable(ctx context.Context) {
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
func AuthSetList (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
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
    // err = DB.Engine.Desc("id").Limit(count, limit).Find(&list)
    err = DB.Engine.Sql("SELECT s.*, (SELECT GROUP_CONCAT(cast(`name` as char(10)) SEPARATOR ',') FROM idp_auth_set where FIND_IN_SET(id, s.sub_id)) as sub_name from idp_auth_set as s").Limit(count, limit).Find(&list)

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
func AuthSetDetail (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  var table IdpAuthSet
  ctx.ReadJSON(&table)

  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  has, err := DB.Engine.Omit("sub_name").Get(&table)
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
func AuthSetAdd (ctx context.Context) {
  data := sumbitAuthSetData(0, ctx)
  ctx.JSON(data)
}

// 修改
func AuthSetPut (ctx context.Context) {


  data := sumbitAuthSetData(1, ctx)
  ctx.JSON(data)
}

type company struct {
  IdpAuthSet `xorm:"extends"`
  Last_table string `json:"last_table"`
}

func (company) TableName() string {
  return "idp_auth_set"
}

// 提交数据 0新增、1修改
func sumbitAuthSetData(tye int, ctx context.Context) context.Map {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    return Utils.NewResData(code, err.Error(), ctx)
  }

  var table company

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


  errMsgs := rules.Validate(Utils.StructToMap(table.IdpAuthSet))
  if errMsgs != nil {
    return Utils.NewResData(1, errMsgs, ctx)
  }


  // 判断数据库里面是否已经存在
  var exist IdpAuthSet
  value := []interface{}{table.IdpAuthSet.Id, table.IdpAuthSet.Name, table.IdpAuthSet.TableName}
  has, err := DB.Engine.Omit("last_table", "sub_name").Where("id<>? and name=? and table_name=?", value...).Exist(&exist)

  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  if has == true {
    return Utils.NewResData(1, table.IdpAuthSet.Name + "已存在", ctx)
  }

  // 写入数据库
  tipsText := "添加"
  if tye == 1 {
    tipsText = "修改"
    // 修改
    _, err = DB.Engine.Omit("last_table", "sub_name").Id(table.IdpAuthSet.Id).Update(&table)

    if err == nil {
      ss, err := DB.Engine.Query("SELECT GROUP_CONCAT(cast(`id` as char(10)) SEPARATOR ',') as id  FROM idp_admin_auth WHERE find_in_set(?,sid)", table.IdpAuthSet.Id)
      if err == nil {
        aid := string(ss[0]["id"])
        sql := "update idp_admin_auth set content=replace(content, ?, ?) where id in(?)"
        _, err = DB.Engine.Exec(sql, table.Last_table, table.IdpAuthSet.TableName, aid)
      }
    }
  } else {
    // 新增
    _, err = DB.Engine.Omit("last_table", "sub_name").Insert(&table)
  }



  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  return Utils.NewResData(0, tipsText + "成功", ctx)
}

// 删除
func AuthSetDel (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_auth_set")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
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
