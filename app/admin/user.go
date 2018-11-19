package admin

import (
  // "fmt"
  // "reflect"
  "github.com/kataras/iris/context"

  Auth "../../authorization"
  Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

// 列表
func UserList (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 获取分页、总数、limit
  page, count, limit, filters := DB.Limit(ctx)
  list := make([]IdpAdmins, 0)

  // 下面开始是查询条件 where
  whereData  := ""
  whereValue :=  []interface{}{}

  start_time := filters["start_time"]
  end_time   := filters["end_time"]
  username   := filters["username"]
  phone      := filters["phone"]
  state      := filters["state"]
  gid        := filters["gid"]
  rid        := filters["rid"]
  job_state  := filters["job_state"]
  sex        := filters["sex"]

  if !Utils.IsEmpty(start_time) && !Utils.IsEmpty(end_time) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.entry_time >= ? and idp_admins.entry_time <= ?`)
    whereValue = append(whereValue, start_time, end_time)
  }

  if !Utils.IsEmpty(username) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.username like ?`)
    whereValue = append(whereValue, `%` + username.(string) + `%`)
  }

  if !Utils.IsEmpty(phone) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.phone = ?`)
    whereValue = append(whereValue, phone)
  }

  if !Utils.IsEmpty(state) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.state = ?`)
    whereValue = append(whereValue, state)
  }

  if !Utils.IsEmpty(gid) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.gid like ?`)
    whereValue = append(whereValue, `%` + Utils.Float64ToStr(gid.(float64)) + `%`)
  }

  if !Utils.IsEmpty(rid) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.rid like ?`)
    whereValue = append(whereValue, `%` + Utils.Float64ToStr(rid.(float64)) + `%`)
  }

  if !Utils.IsEmpty(job_state) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.job_state = ?`)
    whereValue = append(whereValue, job_state)
  }

  if !Utils.IsEmpty(sex) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.sex = ?`)
    whereValue = append(whereValue, sex)
  }

  // 获取服务端用户信息
  reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }
  if !Utils.IsEmpty(reqData["gid"]) {
    whereData = DB.IsWhereEmpty(whereData, "idp_admins.gid =?")
    whereValue = append(whereValue, reqData["gid"])
  }
  // 查询条件结束

  // 获取统计总数
  var table IdpAdmins
  data := context.Map{}
  total, err := DB.Engine.Desc("id").Where(whereData, whereValue...).Count(&table)

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    // 获取列表
    err = DB.Engine.Omit("password").Desc("idp_admins.id").Where(whereData, whereValue...).Limit(count, limit).Find(&list)
    // // err = DB.Engine.Sql("SELECT GROUP_CONCAT(cast(`name` as char(10)) SEPARATOR ',') as `group` from idp_admins_group where idp_admins_group.id = idp_admins.id ").Where(whereData, whereValue...).Limit(count, limit).Find(&list)
    // err = DB.Engine.Omit("password").Sql("select idp_admins.*, (SELECT GROUP_CONCAT(cast(`name` as char(10)) SEPARATOR ',') from idp_admins_group where FIND_IN_SET(id, idp_admins.gid)) as `group` from idp_admins").Limit(count, limit).Find(&list)

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

// 用户详情
func UserDetail(ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  uid, err := ctx.Params().GetInt64("id")
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  res := GetUserDetail(uid, ctx)
  ctx.JSON(res)
}

// 新增
func UserAdd (ctx context.Context) {
  data := sumbitUserData(0, ctx)
  ctx.JSON(data)
}

// 修改
func UserPut (ctx context.Context) {
  data := sumbitUserData(1, ctx)
  ctx.JSON(data)
}

// 提交数据 0新增、1修改
func sumbitUserData(tye int, ctx context.Context) context.Map {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  var table IdpAdmins


  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  // 验证参数
  var rules Utils.Rules
  rules = Utils.Rules{
    "Phone": {
      "required": true, "rgx": "phone",
    },
    "Username": {
      "required": true,
    },
    "Sex": {
      "required": true, "rgx": "int",
    },
  }

  // 新增的时候，必须验证密码
  if !Utils.IsEmpty(table.Password) {
    rules["Password"] = map[string]interface{}{
      "required": true, "rgx": "password",
    }
  }

  // 验证参数
  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    return Utils.NewResData(1, errMsgs, ctx)
  }


  // 获取服务端用户信息
  author      := ctx.GetHeader("Authorization")
  userinfo, _ := Auth.DecryptToken(author, "admin")
  reqData     := userinfo.(map[string]interface{})

  if len(reqData) <= 0 {
    return Utils.NewResData(1, "获取服务端用户信息失败", ctx)
  }

  id := int64(reqData["id"].(float64))
  res := GetUserDetail(id, ctx)
  userData := res["data"].(IdpAdmins)
  // 获取服务端用户信息 END


  // 看看是否修改密码
  var isUser IdpAdmins
  has, err := DB.Engine.Cols("password").Where("id=?", userData.Id).Get(&isUser)

  if Utils.IsEmpty(table.Password) {
    table.Password = isUser.Password
  } else {
    table.Password = Public.EncryptPassword(table.Password)
  }


  // 判断数据库里面是否已经存在
  var exist IdpAdmins
  value := []interface{}{userData.Id, userData.Phone}
  has, err = DB.Engine.Where("id<>? and phone=?", value...).Exist(&exist)

  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  if has == true {
    return Utils.NewResData(1, table.Phone + "已存在", ctx)
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
      "uid": table.Id,
    }, ctx)
  }

  return Utils.NewResData(0, tipsText + "成功", ctx)
}

// 删除
func UserDel (ctx context.Context) {
  // 判断权限
  hasAuth, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  var table IdpAdmins

  // 根据不同环境返回数据
  err = Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdmins
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


