package admin

import (
  // "fmt"
  // "reflect"
  "github.com/kataras/iris/context"

  Auth "../../authorization"
  Utils "../../utils"
  DB "../../database"
)

type IdpAdminsGroup struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Aid int64 `json:"aid"`
  State int64 `json:"state"`
  DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}


type roleList = IdpAdminsRole

type GroupAndRole struct {
  IdpAdminsGroup `xorm:"extends"`
  // Count int64 `json:"count"`
}

func (GroupAndRole) TableName() string {
  return "idp_admins_group"
}

// 列表
func GroupList (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins_group")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 获取分页、总数、limit
  page, count, limit, _ := DB.Limit(ctx)
  list := make([]GroupAndRole, 0)

  // 下面开始是查询条件 where
  whereData  := ""
  whereValue :=  []interface{}{}

  // 获取服务端用户信息
  reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }
  if !Utils.IsEmpty(reqData["gid"]) {
    whereData = DB.IsWhereEmpty(whereData, "id =?")
    whereValue = append(whereValue, reqData["gid"])
  }
  // 查询条件结束

  // 获取统计总数
  var table IdpAdminsRole
  data := context.Map{}

  total, err := DB.Engine.Where(whereData, whereValue...).Count(&table)

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    // 获取列表
    // err = DB.Engine.Desc("g.id").Sql("select g.*, r.* from idp_admins_group as g, idp_admins_role as r where r.gid = g.id").Limit(count, limit).Find(&list)
    // err = DB.Engine.Desc("g.id").Sql("select g.*, (select count(id) from idp_admins_role as r where r.gid = g.id) as count from idp_admins_group as g").Where(whereData, whereValue...).Limit(count, limit).Find(&list)
    err = DB.Engine.Desc("id").Where(whereData, whereValue...).Limit(count, limit).Find(&list)

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
func GroupDetail (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins_group")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  var table IdpAdminsGroup
  ctx.ReadJSON(&table)

  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  data := context.Map{}

  has, err := DB.Engine.Get(&table)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  if has == true {
    data = Utils.NewResData(0, table, ctx)
  } else {
    data = Utils.NewResData(1, "记录不存在", ctx)
  }

  ctx.JSON(data)

}

// 新增
func GroupAdd (ctx context.Context) {
  data := sumbitGroupData(0, ctx)
  ctx.JSON(data)
}

// 修改
func GroupPut (ctx context.Context) {
  data := sumbitGroupData(1, ctx)
  ctx.JSON(data)
}

// 提交数据 0新增、1修改
func sumbitGroupData(tye int, ctx context.Context) context.Map {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins_group")
  if hasAuth != true {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  var table IdpAdminsGroup

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
      // "rgx": "identity",
    },
  }


  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    return Utils.NewResData(1, errMsgs, ctx)
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsGroup
  value := []interface{}{table.Id, table.Name}
  has, err := DB.Engine.Where("id<>? and name=?", value...).Exist(&exist)

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
func GroupDel (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins_group")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  var table IdpAdminsGroup

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsGroup
  has, err := DB.Engine.Where("id=?", table.Id).Exist(&exist)

  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  if has != true {
    ctx.JSON(Utils.NewResData(1, "该信息不存在", ctx))
    return
  }

  // 判断角色管理表是否存在，如果存在的话，不予删除
  var roleExist IdpAdminsRole
  has, err = DB.Engine.Where("gid=?", table.Id).Exist(&roleExist)

  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  if has == true {
    ctx.JSON(Utils.NewResData(1, "无法删除，角色管理中使用了该值", ctx))
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


