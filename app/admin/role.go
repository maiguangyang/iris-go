package admin

import (
  // "fmt"
  // "time"
  // "reflect"
  // "encoding/json"
  "github.com/kataras/iris/context"

  Auth "../../authorization"
  // Public "../../public"
  Utils "../../utils"
  DB "../../database"
)


type IdpAdminRoles struct {
  DB.Model
  Name string `json:"name"`
  Gid int64 `json:"gid"`
  Aid int64 `json:"aid"`
  State int64 `json:"state"`
}

func (IdpAdminRoles) TableName() string {
  return "idp_admin_role"
}

type IdpAdminRole struct {
  IdpAdminRoles
  Group IdpAdminGroups `json:"group" gorm:"foreignkey:Gid"`
  Auth IdpAdminAuth `json:"auth" gorm:"foreignkey:Rid"`
}

// 角色列表
func RoleList (ctx context.Context) {
  // 判断权限
  hasAuth, stride, code, err := DB.CheckAdminAuth(ctx, "idp_admin_role")
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
  // page, count, limit, filters := DB.Limit(ctx)
  page, count, offset, filters := DB.Limit(ctx)
  lists := make([]IdpAdminRole, 0)

  // 下面开始是查询条件 where
  whereData  := ""
  whereValue :=  []interface{}{}

  group := filters["group"]
  name  := filters["name"]
  state := filters["state"]
  gid   := filters["gid"]

  if !Utils.IsEmpty(group) {
    whereData = DB.IsWhereEmpty(whereData, "gid in(" + Utils.ArrayInt64ToString(group) + ")")
  }

  if !Utils.IsEmpty(name) {
    whereData = DB.IsWhereEmpty(whereData, `name like ?`)
    whereValue = append(whereValue, `%` + name.(string) + `%`)
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
    whereData = DB.IsWhereEmpty(whereData, `gid = ?`)
    whereValue = append(whereValue, gid)
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
  // result := DB.Engine.Model(&lists).Order("id desc").Where(whereData, whereValue...).Limit(count).Offset(offset).Preload("Group").Preload("Auth").Find(&lists).Count(&total)
  if err := DB.Engine.Model(&lists).Order("id desc").Where(whereData, whereValue...).Count(&total).Limit(count).Offset(offset).Find(&lists).Error; err != nil {
    data = Utils.NewResData(1, "return data is empty.", ctx)
  } else {
    // 然后循环列表，关联查询roles表
    for key, list := range lists {
      if err := DB.Engine.Model(&list).Related(&list.Group, "Group").Error; err == nil {
        lists[key] = list
      }

      if err := DB.Engine.Model(&list).Related(&list.Auth, "Auth").Error; err == nil {
        lists[key] = list
      }
    }

    resData := Utils.TotalData(lists, page, total, count)
    data = Utils.NewResData(0, resData, ctx)
  }


  // if result.Error != nil {
  //   data = Utils.NewResData(1, "return data is empty.", ctx)
  // } else {
  //   resData := Utils.TotalData(lists, page, total, count)
  //   data = Utils.NewResData(0, resData, ctx)
  // }

  ctx.JSON(data)

}

// 详情
func RoleDetail (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_role")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  data := context.Map{}

  var table IdpAdminRole
  id, _ := ctx.Params().GetInt64("id")
  table.Id = id

  if err := DB.Engine.First(&table).Error; err != nil {
    data = Utils.NewResData(1, err, ctx)
  } else {
    if err := DB.Engine.Model(&table).Order("id desc").Related(&table.Group, "Group").Related(&table.Auth, "Auth").Error; err != nil {
      data = Utils.NewResData(1, err, ctx)
    } else {
      data = Utils.NewResData(0, table, ctx)
    }
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
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_role")
  if hasAuth != true {
    return Utils.NewResData(code, err.Error(), ctx)
  }

  var table IdpAdminRoles

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
  var exist IdpAdminRoles
  value := []interface{}{table.Id, table.Name}
  if err := DB.Engine.Where("id<>? and name=?", value...).First(&exist).Error; err == nil {
    return Utils.NewResData(1, table.Name + "已存在", ctx)
  }

  // tipsText := "添加"
  if tye == 1 {

    if table.State == 2 {
      // 判断角色管理表是否存在，如果存在的话，不予删除
      var adminsExist IdpAdmins
      if err := DB.Engine.Where("gid=?", table.Id).First(&adminsExist).Error; err == nil {
        return Utils.NewResData(1, "状态禁用失败，员工管理中使用了该值", ctx)
      }
    }

    if err := DB.Engine.Model(&table).Where("id =?", table.Id).Updates(&table).Error; err != nil {
      return Utils.NewResData(1, "修改失败", ctx)
    }
    return Utils.NewResData(0, "修改成功", ctx)
  }
    // 新增
  if err := DB.Engine.Create(&table).Error; err != nil {
    return Utils.NewResData(1, "添加失败", ctx)
  }

  return Utils.NewResData(0, context.Map{
    "rid": table.Id,
  }, ctx)

}

// 删除
func RoleDel (ctx context.Context) {
  // 判断权限
  hasAuth, _, code, err := DB.CheckAdminAuth(ctx, "idp_admin_role")
  if hasAuth != true {
    ctx.JSON(Utils.NewResData(code, err.Error(), ctx))
    return
  }

  var table IdpAdminRoles

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

  // 判断角色管理表是否存在，如果存在的话，不予删除
  var adminsExist IdpAdmins
  if err := DB.Engine.Where("gid=?", table.Id).First(&adminsExist).Error; err == nil {
    ctx.JSON(Utils.NewResData(1, "无法删除，员工管理中使用了该值", ctx))
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


