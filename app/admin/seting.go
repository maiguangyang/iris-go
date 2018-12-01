package admin

import(
  // "fmt"
  "time"
  "github.com/kataras/iris/context"
  // Public "../../public"
  DB "../../database"
  Utils "../../utils"
)

type IdpAuthSet struct {
  DB.Model
  Name string `json:"name"`
  TableName string `json:"table_name"`
  Routes string `json:"routes"`
  SubId string `json:"sub_id"`
  SubName string `json:"sub_name"`
}

// 获取数据库所有表
func AuthSetTable(ctx context.Context) {

  type handleTableName struct {
    TableName string `json:"table_name"`
    CreateTime time.Time `json:"create_time"`
  }

  list := make([]handleTableName, 0)
  sql := "select table_name, create_time from information_schema.tables where table_schema = (select database())"
  DB.Engine.Raw(sql).Scan(&list)

  data := Utils.NewResData(0, list, ctx)
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
  // page, count, _, _ := DB.Limit(ctx)
  list := make([]IdpAuthSet, 0)
  data := context.Map{}

  sql := "SELECT s.*, (SELECT GROUP_CONCAT(cast(`name` as char(10)) SEPARATOR ',') FROM idp_auth_set where FIND_IN_SET(id, s.sub_id)) as sub_name from idp_auth_set as s order by id desc"
  if err := DB.Engine.Raw(sql).Scan(&list).Error; err != nil {
    data = Utils.NewResData(1, err, ctx)
  } else {
    resData := Utils.TotalData(list, 1, int64(len(list)), len(list))
    data = Utils.NewResData(0, resData, ctx)
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

  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  if err := DB.Engine.Where("id =?", table.Id).First(&table).Error; err != nil {
    ctx.JSON(Utils.NewResData(1, "return data is empty.", ctx))
    return
  }

  ctx.JSON(Utils.NewResData(0, table, ctx))
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
  IdpAuthSet
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
  // var table1 IdpAuthSet
  var exist IdpAdminRoles
  value := []interface{}{table.Id, table.Name}
  if err := DB.Engine.Where("id<>? and name = ?", value...).First(&exist).Error; err == nil {
    return Utils.NewResData(1, table.Name + "已存在", ctx)
  }

  // 修改
  if tye == 1 {
    if err := DB.Engine.Model(&table).Omit("last_table", "sub_name").Where("id =?", table.IdpAuthSet.Id).Updates(&table).Error; err != nil {
      return Utils.NewResData(1, "修改失败", ctx)
    }

    // 判断数据表是否改变
    if table.Last_table != table.IdpAuthSet.TableName {
      type Result struct {
        Id int64
      }
      var result Result
      sql := "SELECT GROUP_CONCAT(cast(`id` as char(10)) SEPARATOR ',') as id  FROM idp_admin_auth WHERE find_in_set(?,sid)"
      DB.Engine.Raw(sql, table.Id).Scan(&result)

      if result.Id > 0 {
        sql = "update idp_admin_auth set content=replace(content, ?, ?) where id in(?)"
        DB.Engine.Exec(sql, table.Last_table, table.IdpAuthSet.TableName, result.Id)
      }
    }

    return Utils.NewResData(0, "修改成功", ctx)
  }

  // 新增
  if err := DB.Engine.Omit("last_table", "sub_name").Create(&table).Error; err != nil {
    return Utils.NewResData(1, "添加失败", ctx)
  }

  return Utils.NewResData(0, "添加成功", ctx)

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
  if err := DB.Engine.Where("id=?", table.Id).First(&table).Error; err != nil {
    ctx.JSON(Utils.NewResData(1, "该信息不存在", ctx))
    return
  }

  // 判断是否被使用，如果存在的话，不予删除
  var roleExist IdpAdminAuth
  if err := DB.Engine.Where("sid=?", table.Id).First(&roleExist).Error; err == nil {
    ctx.JSON(Utils.NewResData(1, "无法删除，用户权限中使用了该值", ctx))
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
