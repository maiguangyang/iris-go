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

// type IdpAdmins struct {
//   Id int64 `json:"id"`
//   Name string `json:"name"`
//   Gid int64 `json:"gid"`
//   Aid int64 `json:"aid"`
//   State int64 `json:"state"`
//   DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
//   UpdatedAt int64 `json:"updated_at" xorm:"updated"`
//   CreatedAt int64 `json:"created_at" xorm:"created"`
// }

type AdminsGroup struct {
  IdpAdmins `xorm:"extends"`
  Group IdpAdminsGroup `json:"group" xorm:"extends"`
  Role IdpAdminsRole `json:"role" xorm:"extends"`
}

func (AdminsGroup) TableName() string {
  return "idp_admins"
}

// 用户组列表
func UserList (ctx context.Context) {
  // 获取分页
  page  := Utils.StrToInt64(ctx.URLParam("page"))
  count := Utils.StrToInt64(ctx.URLParam("count"))
  list := make([]AdminsGroup, 0)

  // 获取统计总数
  var table AdminsGroup
  // total := DB.Count(&table)

  total := DB.Count(context.Map{
    "type"  : 1,
    "table" : &table,
    "sql"   : "select count(*) from idp_admins as a LEFT JOIN idp_admins_group as g ON a.gid = g.id LEFT JOIN idp_admins_role as r ON a.rid = r.id",
  })

  // 获取列表
  err := DB.Find(context.Map{
    "type"  : 1,
    "table" : &list,
    "page"  : page,
    "count" : count,
    "where": "",
    "value": []interface{}{},
    "sql"   : "select * from idp_admins as a LEFT JOIN idp_admins_group as g ON a.gid = g.id LEFT JOIN idp_admins_role as r ON a.rid = r.id order by a.id desc",
  })

  // 返回数据
  data := context.Map{}

  if err != nil {
    data = Utils.NewResData(404, err.Error(), ctx)
  } else {

    resData := Utils.TotalData(list, page, total, count)

    data = Utils.NewResData(0, resData, ctx)
  }

  ctx.JSON(data)

}

// 详情
func UserDetail (ctx context.Context) {

  var table AdminsGroup
  ctx.ReadJSON(&table)

  id, _ := ctx.Params().GetInt64("id")
  // has := DB.Get(&table, "id=?", []interface{}{id})

  has := DB.Get(context.Map{
    "type": 1,
    "table": &table,
    "where": "id=?",
    "value": []interface{}{id},
    "sql": `select * from idp_admins_role as r LEFT JOIN idp_admins_group as g ON r.gid = g.id where r.id=` + Utils.Int64ToStr(id),
  })

  data := context.Map{}
  if has == true {
    data = Utils.NewResData(0, table, ctx)
  } else {
    data = Utils.NewResData(1, "记录不存在", ctx)
  }

  ctx.JSON(data)

}

// 新增
// func UserAdd (ctx context.Context) {
//   var table IdpAdmins

//   var rules Utils.Rules

//   // 线上环境
//   if Public.NODE_ENV {
//     decData, err := Public.DecryptReqData(ctx)

//     if err != nil {
//       ctx.JSON(Utils.NewResData(1, err, ctx))
//       return
//     }

//     reqData := decData.(map[string]interface{})

//     table.Name  = reqData["name"].(string)
//     table.State = int64(reqData["state"].(float64))

//   } else {
//     ctx.ReadJSON(&table)
//   }

//   // 验证参数
//   rules = Utils.Rules{
//     "Name": {
//       "required": true,
//     },
//     "Gid": {
//       "required": true,
//       "rgx": "int",
//     },
//   }


//   errMsgs := rules.Validate(Utils.StructToMap(table))
//   if errMsgs != nil {
//     ctx.JSON(Utils.NewResData(1, errMsgs, ctx))
//     return
//   }

//   // 判断数据库里面是否已经存在
//   var exist IdpAdmins
//   // has := DB.Exist(&exist, "id<>? and gid=? and name=?", []interface{}{table.Id, table.Gid, table.Name})
//   has := DB.Exist(context.Map{
//     "type": 0,
//     "table": &exist,
//     "where": "id<>? and gid=? and name=?",
//     "value": []interface{}{table.Id, table.Gid, table.Name},
//     "sql": "",
//   })

//   data := context.Map{}
//   if has == true {
//     data = Utils.NewResData(1, table.Name + "已存在", ctx)
//     ctx.JSON(data)
//     return
//   }


//   // 写入数据库
//   err := DB.Post(&table)

//   if err == nil {
//     data = Utils.NewResData(0, "添加成功", ctx)
//   } else {
//     data = Utils.NewResData(1, "添加失败", ctx)
//   }

//   ctx.JSON(data)
// }

// // 修改
// func UserPut (ctx context.Context) {
//   var table IdpAdmins

//   var rules Utils.Rules

//   // 线上环境
//   if Public.NODE_ENV {
//     decData, err := Public.DecryptReqData(ctx)

//     if err != nil {
//       ctx.JSON(Utils.NewResData(1, err, ctx))
//       return
//     }

//     reqData := decData.(map[string]interface{})

//     table.Id    = int64(reqData["id"].(float64))
//     table.Name  = reqData["name"].(string)
//     table.Gid   = int64(reqData["gid"].(float64))
//     table.State = int64(reqData["state"].(float64))

//   } else {
//     ctx.ReadJSON(&table)
//   }

//   // 验证参数
//   rules = Utils.Rules{
//     "Name": {
//       "required": true,
//     },
//     "Gid": {
//       "required": true,
//       "rgx": "int",
//     },
//   }


//   errMsgs := rules.Validate(Utils.StructToMap(table))
//   if errMsgs != nil {
//     ctx.JSON(Utils.NewResData(1, errMsgs, ctx))
//     return
//   }

//   // 判断数据库里面是否已经存在
//   var exist IdpAdmins
//   // has := DB.Exist(&exist, "id<>? and gid=? and name=?", []interface{}{table.Id, table.Gid, table.Name})
//   has := DB.Exist(context.Map{
//     "type": 0,
//     "table": &exist,
//     "where": "id<>? and gid=? and name=?",
//     "value": []interface{}{table.Id, table.Gid, table.Name},
//     "sql": "",
//   })

//   data := context.Map{}
//   if has == true {
//     data = Utils.NewResData(1, table.Name + "已存在", ctx)
//     ctx.JSON(data)
//     return
//   }

//   // 写入数据库
//   err := DB.Put(table.Id, &table)

//   if err == nil {
//     data = Utils.NewResData(0, "修改成功", ctx)
//   } else {
//     data = Utils.NewResData(1, "修改失败", ctx)
//   }

//   ctx.JSON(data)
// }

// 删除
func UserDel (ctx context.Context) {
  var table IdpAdmins

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

  err := DB.Delete(table.Id, &table)

  data := context.Map{}
  if err == nil {
    data = Utils.NewResData(0, "删除成功", ctx)
  } else {
    data = Utils.NewResData(1, err.Error(), ctx)
  }


  ctx.JSON(data)
}


