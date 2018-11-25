package admin

import (
  // "fmt"
  // "reflect"
  // "encoding/json"
  "github.com/kataras/iris/context"

  Auth "../../authorization"
  // Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

type IdpAdminsRole struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Gid int64 `json:"gid"`
  Aid int64 `json:"aid"`
  State int64 `json:"state"`
  DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

type RoleAndGroup struct {
  IdpAdminsRole `xorm:"extends"`
  Group IdpAdminsGroup `json:"group" xorm:"extends"`
  Auth IdpAdminAuth `json:"auth" xorm:"extends"`
}

func (RoleAndGroup) TableName() string {
  return "idp_admins_role"
}


// 角色列表
func RoleList (ctx context.Context) {
  // 判断权限
  hasAuth, stride, code, err := DB.CheckAdminAuth(ctx, "idp_admins_role")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  // 获取分页、总数、limit
  page, count, limit, filters := DB.Limit(ctx)
  list := make([]RoleAndGroup, 0)

  // 下面开始是查询条件 where
  whereData  := ""
  whereValue :=  []interface{}{}

  group := filters["group"]
  name  := filters["name"]
  state := filters["state"]
  gid   := filters["gid"]

  if !Utils.IsEmpty(group) {
    whereData = DB.IsWhereEmpty(whereData, "idp_admins_role.gid in(" + Utils.ArrayInt64ToString(group) + ")")
  }

  if !Utils.IsEmpty(name) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins_role.name like ?`)
    whereValue = append(whereValue, `%` + name.(string) + `%`)
  }

  if !Utils.IsEmpty(state) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins_role.state = ?`)
    whereValue = append(whereValue, state)
  }

  if !Utils.IsEmpty(gid) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins_role.gid = ?`)
    whereValue = append(whereValue, gid)
  }

  // 是否跨部门
  if stride != true {
    // 获取服务端用户信息
    reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
    if err != nil {
      ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
      return
    }

    if !Utils.IsEmpty(reqData["gid"]) {
      whereData = DB.IsWhereEmpty(whereData, "idp_admins_role.gid =?")
      whereValue = append(whereValue, reqData["gid"])
    }
  }
  // 查询条件结束


  // 连表查询，下面进行了1个连表
  joinTable  := make(map[int]map[string]string)
  joinTable[0] = map[string]string {
    "type"  : "LEFT",
    "table" : "idp_admins_group",
    "where" : "idp_admins_role.gid = idp_admins_group.id",
  }

  // 获取统计总数
  var table RoleAndGroup
  data := context.Map{}

  total, err := DB.Engine.Join(joinTable[0]["type"], joinTable[0]["table"], joinTable[0]["where"]).Where(whereData, whereValue...).Count(&table)

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    // 获取列表
    err = DB.Engine.Desc("idp_admins_role.id").Join(joinTable[0]["type"], joinTable[0]["table"], joinTable[0]["where"]).Where(whereData, whereValue...).Limit(count, limit).Find(&list)

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
func RoleDetail (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admins_role")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  var table RoleAndGroup
  ctx.ReadJSON(&table)

  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  has, err := DB.Engine.Join("LEFT", "idp_admins_group", "idp_admins_role.gid = idp_admins_group.id").Join("LEFT", "idp_admin_auth", "idp_admin_auth.rid = idp_admins_role.id").Get(&table)
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
func RoleAdd (ctx context.Context) {
  data := sumbitRoleData(0, ctx)
  ctx.JSON(data)
}

// 修改
func RolePut (ctx context.Context) {
  data := sumbitRoleData(1, ctx)
  ctx.JSON(data)
}


// 提交数据 0新增、1修改
func sumbitRoleData(tye int, ctx context.Context) context.Map {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admins_role")
  if hasAuth != true {
    return Utils.NewResData(code, err.Error(), ctx)
  }

  var table IdpAdminsRole

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
    "Gid": {
      "required": true,
      "rgx": "int",
    },
  }


  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    return Utils.NewResData(1, errMsgs, ctx)
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsRole
  value := []interface{}{table.Id, table.Gid, table.Name}
  has, err := DB.Engine.Where("id<>? and gid=? and name=?", value...).Exist(&exist)

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

  // 新增返回并
  if tye == 0 {
    return Utils.NewResData(0, context.Map{
      "rid": table.Id,
    }, ctx)
  }

  return Utils.NewResData(0, tipsText + "成功", ctx)
}

// 删除
func RoleDel (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admins_role")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  var table IdpAdminsRole

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsRole
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
  var adminExist IdpAdmins
  has, err = DB.Engine.Where("rid=?", table.Id).Exist(&adminExist)

  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  if has == true {
    ctx.JSON(Utils.NewResData(1, "无法删除，员工管理中使用了该值", ctx))
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


