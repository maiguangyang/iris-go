package admin

import (
  "fmt"
  // "reflect"
  "github.com/kataras/iris/context"

  // Auth "../../authorization"
  Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

// 用户组列表
func UserList (ctx context.Context) {
  // 获取分页、总数、limit
  page, count, limit, filters := DB.Limit(ctx)
  list := make([]IdpAdmins, 0)


  // 下面开始是查询条件 where
  whereData  := ""
  whereValue :=  []interface{}{}

  start_time := filters["start_time"]
  end_time   := filters["end_time"]
  phone      := filters["phone"]

  if !Utils.IsEmpty(start_time) && !Utils.IsEmpty(end_time) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.entry_time >= ? and idp_admins.entry_time <= ?`)
    whereValue = append(whereValue, start_time, end_time)
  }


  if !Utils.IsEmpty(phone) {
    whereData = DB.IsWhereEmpty(whereData, `idp_admins.phone = ?`)
    whereValue = append(whereValue, phone)
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
    err = DB.Engine.Omit("password").Desc("id").Where(whereData, whereValue...).Limit(count, limit).Find(&list)

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
  var table IdpAdmins

  var rules Utils.Rules

  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      return Utils.NewResData(1, err.Error(), ctx)
    }

    reqData := decData.(map[string]interface{})
    // map 映射 struct
    err = Utils.FillStruct(&table, reqData)
    if err != nil {
      return Utils.NewResData(1, err.Error(), ctx)
    }

  } else {
    ctx.ReadJSON(&table)
  }

  fmt.Println(table)

  // 验证参数
  rules = Utils.Rules{
    "Phone": {
      "required": true, "rgx": "phone",
    },
    "Password": {
      "required": true, "rgx": "password",
    },
    "Username": {
      "required": true,
    },
    "Sex": {
      "required": true, "rgx": "int",
    },
  }

  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    return Utils.NewResData(1, errMsgs, ctx)
  }


  // 判断数据库里面是否已经存在
  var exist IdpAdmins
  value := []interface{}{table.Id, table.Phone}
  has, err := DB.Engine.Where("id<>? and phone=?", value...).Exist(&exist)

  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  if has == true {
    return Utils.NewResData(1, table.Phone + "已存在", ctx)
  }

  table.Password = Public.EncryptPassword(table.Password)

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

  if tye == 0 {
    return Utils.NewResData(0, context.Map{
      "uid": table.Id,
    }, ctx)
  }

  return Utils.NewResData(0, tipsText + "成功", ctx)
}


// 删除
func UserDel (ctx context.Context) {
  var table IdpAdmins

  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
      return
    }

    reqData  := decData.(map[string]interface{})
    table.Id = int64(reqData["id"].(float64))

  } else {
    ctx.ReadJSON(&table)
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


