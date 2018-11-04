package admin

import (
  // "fmt"
  // "reflect"
  "github.com/kataras/iris/context"

  // Auth "../../authorization"
  Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

type IdpAdminsRole struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Gid int64 `json:"gid"`
  State int64 `json:"state"`
  DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

type GroupRoleGroup struct {
  IdpAdminsRole `xorm:"extends"`
  Group IdpAdminsGroup `json:"group" xorm:"extends"`
}

func (GroupRoleGroup) TableName() string {
  return "idp_admins_role"
}


// 用户组列表
func RoleList (ctx context.Context) {
  // 获取分页
  page := Utils.StrToInt64(ctx.URLParam("page"))

  // 获取统计总数
  var table GroupRoleGroup
  total := DB.Count(&table)

  // 获取列表
  list := make([]GroupRoleGroup, 0)
  err := DB.Find(context.Map{
    "type": 1,
    "table": &list,
    "page": page,
    "sql": "select * from idp_admins_role, idp_admins_group where idp_admins_role.gid = idp_admins_group.id",
  })


  // 返回数据
  data := context.Map{}

  if err != nil {
    data = Utils.NewResData(404, err.Error(), ctx)
  } else {

    resData := Utils.TotalData(list, page, total)

    data = Utils.NewResData(0, resData, ctx)
  }

  ctx.JSON(data)

}

// 详情
func RoleDetail (ctx context.Context) {

  var table IdpAdminsRole
  ctx.ReadJSON(&table)

  id, _ := ctx.Params().GetInt64("id")

  has := DB.Get(&table, "id=?", []interface{}{id})


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
  var table IdpAdminsRole

  var rules Utils.Rules

  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      ctx.JSON(Utils.NewResData(1, err, ctx))
      return
    }

    reqData := decData.(map[string]interface{})

    table.Name  = reqData["name"].(string)
    table.State = int64(reqData["state"].(float64))

  } else {
    ctx.ReadJSON(&table)
  }

  // 验证参数
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
    ctx.JSON(Utils.NewResData(1, errMsgs, ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsRole
  has := DB.Exist(&exist, "id<>? and gid=? and name=?", []interface{}{table.Id, table.Gid, table.Name})

  data := context.Map{}
  if has == true {
    data = Utils.NewResData(1, table.Name + "已存在", ctx)
    ctx.JSON(data)
    return
  }


  // 写入数据库
  err := DB.Post(&table)

  if err == nil {
    data = Utils.NewResData(0, "添加成功", ctx)
  } else {
    data = Utils.NewResData(1, "添加失败", ctx)
  }

  ctx.JSON(data)
}

// 修改
func RolePut (ctx context.Context) {
  var table IdpAdminsRole

  var rules Utils.Rules

  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      ctx.JSON(Utils.NewResData(1, err, ctx))
      return
    }

    reqData := decData.(map[string]interface{})

    table.Id    = int64(reqData["id"].(float64))
    table.Name  = reqData["name"].(string)
    table.Gid   = int64(reqData["gid"].(float64))
    table.State = int64(reqData["state"].(float64))

  } else {
    ctx.ReadJSON(&table)
  }

  // 验证参数
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
    ctx.JSON(Utils.NewResData(1, errMsgs, ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsRole
  has := DB.Exist(&exist, "id<>? and gid=? and name=?", []interface{}{table.Id, table.Gid, table.Name})

  data := context.Map{}
  if has == true {
    data = Utils.NewResData(1, table.Name + "已存在", ctx)
    ctx.JSON(data)
    return
  }

  // 写入数据库
  err := DB.Put(table.Id, &table)

  if err == nil {
    data = Utils.NewResData(0, "修改成功", ctx)
  } else {
    data = Utils.NewResData(1, "修改失败", ctx)
  }

  ctx.JSON(data)
}

// 删除
func RoleDel (ctx context.Context) {
  var table IdpAdminsRole

  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      ctx.JSON(Utils.NewResData(1, err, ctx))
      return
    }

    reqData  := decData.(map[string]interface{})
    table.Id = int64(reqData["id"].(float64))

  } else {
    ctx.ReadJSON(&table)
  }

  err := DB.Delete(&table)

  data := context.Map{}
  if err == nil {
    data = Utils.NewResData(0, "删除成功", ctx)
  } else {
    data = Utils.NewResData(1, "删除失败", ctx)
  }

  ctx.JSON(data)
}


