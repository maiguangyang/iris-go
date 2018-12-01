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
  hasAuth, stride, code, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  // 获取服务端用户信息
  reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  // 获取分页、总数、limit
  page, count, offset, filters := DB.Limit(ctx)
  lists := make([]IdpAdmins, 0)

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
    whereData = DB.IsWhereEmpty(whereData, `entry_time >= ? and entry_time <= ?`)
    whereValue = append(whereValue, start_time, end_time)
  }

  if !Utils.IsEmpty(username) {
    whereData = DB.IsWhereEmpty(whereData, `username like ?`)
    whereValue = append(whereValue, `%` + username.(string) + `%`)
  }

  if !Utils.IsEmpty(phone) {
    whereData = DB.IsWhereEmpty(whereData, `phone = ?`)
    whereValue = append(whereValue, phone)
  }

  // 如果不是超级账户，只显示state状态为1的信息
  super := int64(reqData["super"].(float64))
  if super == 2 {
    if !Utils.IsEmpty(state) {
      whereData = DB.IsWhereEmpty(whereData, `state = ?`)
      whereValue = append(whereValue, state)
    }
  } else {
    whereData = DB.IsWhereEmpty(whereData, "state =?")
    whereValue = append(whereValue, 1)
  }


  if !Utils.IsEmpty(gid) {
    whereData = DB.IsWhereEmpty(whereData, `gid like ?`)
    whereValue = append(whereValue, `%` + Utils.Float64ToStr(gid.(float64)) + `%`)
  }

  if !Utils.IsEmpty(rid) {
    whereData = DB.IsWhereEmpty(whereData, `rid like ?`)
    whereValue = append(whereValue, `%` + Utils.Float64ToStr(rid.(float64)) + `%`)
  }

  if !Utils.IsEmpty(job_state) {
    whereData = DB.IsWhereEmpty(whereData, `job_state = ?`)
    whereValue = append(whereValue, job_state)
  }

  if !Utils.IsEmpty(sex) {
    whereData = DB.IsWhereEmpty(whereData, `sex = ?`)
    whereValue = append(whereValue, sex)
  }

  // 是否跨部门
  if stride != true {
    if !Utils.IsEmpty(reqData["gid"]) {
      whereData = DB.IsWhereEmpty(whereData, "gid =?")
      whereValue = append(whereValue, reqData["gid"])
    }
  }
  // 查询条件结束

  // 查询列表
  data := context.Map{}
  var total int64
  if err := DB.Engine.Model(&lists).Order("id desc").Where(whereData, whereValue...).Count(&total).Limit(count).Offset(offset).Find(&lists).Error; err != nil {
    data = Utils.NewResData(1, err, ctx)
  } else {

    gList := make([]IdpAdminGroups, 0)
    rList := make([]IdpAdminRoles, 0)
    for key, list := range lists {
      if err := DB.Engine.Where("id in(?)", Utils.StrToArr(list.Gid, ",")).Find(&gList).Error; err == nil {
        lists[key].Groups = gList
      }
      if err := DB.Engine.Where("id in(?)", Utils.StrToArr(list.Rid, ",")).Find(&rList).Error; err == nil {
        lists[key].Roles  = rList
      }
    }

    resData := Utils.TotalData(lists, page, total, count)
    data = Utils.NewResData(0, resData, ctx)
  }

  ctx.JSON(data)

}

// 用户详情
func UserDetail(ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
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
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    return Utils.NewResData(code, err.Error(), ctx)
  }

  var table IdpAdminsPass
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
  var isUser IdpAdminsPass
  // has, err := DB.Engine.Cols("password").Where("id=?", userData.Id).Get(&isUser)
  DB.Engine.Where("id=?", userData.Id).First(&isUser)

  if Utils.IsEmpty(table.Password) {
    table.Password = isUser.Password
  } else {
    table.Password = Public.EncryptPassword(table.Password)
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdmins
  if tye == 1 {
    if err := DB.Engine.Where("id<>? and phone=?", table.Id, table.Phone).First(&exist).Error; err == nil {
      return Utils.NewResData(1, table.Phone + "已存在", ctx)
    }

    if err := DB.Engine.Model(&table).Where("id =?", table.Id).Updates(&table).Error; err != nil {
      return Utils.NewResData(1, "修改失败", ctx)
    }
    return Utils.NewResData(0, "修改成功", ctx)
  }

  if err := DB.Engine.Where("phone=?", table.Phone).First(&exist).Error; err == nil {
    return Utils.NewResData(1, table.Phone + "已存在", ctx)
  }
    // 新增
  if err := DB.Engine.Create(&table).Error; err != nil {
    return Utils.NewResData(1, "添加失败", ctx)
  }

  return Utils.NewResData(0, context.Map{
    "uid": table.Id,
  }, ctx)

}

// 删除
func UserDel (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admins")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
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


