package admin

import (
  "fmt"
  // "reflect"
  "encoding/json"
  "github.com/kataras/iris/context"

  // Auth "../../authorization"
  Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

type AdminsGroup struct {
  IdpAdmins `xorm:"extends"`
  Group IdpAdminsGroup `json:"group" xorm:"extends"`
  Role IdpAdminsRole `json:"role" xorm:"extends"`
}

func (AdminsGroup) TableName() string {
  return "idp_admins"
}


func formFilter(ctx context.Context) {

  data := ctx.URLParam("filters")
  var dat [][3]string

  _ = json.Unmarshal([]byte(data), &dat)
  fmt.Println(dat)

}


// 用户组列表
func UserList (ctx context.Context) {
  // 获取分页、总数、limit
  page, count, limit, filters := DB.Limit(ctx)
  list := make([]AdminsGroup, 0)



  // 获取统计总数
  var table AdminsGroup
  data := context.Map{}

  // 连表查询，下面进行了2个连表
  joinTable  := make(map[int]map[string]string)


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

  joinTable[0] = map[string]string {
    "type"      : "LEFT",
    "table" : "idp_admins_group",
    "where"     : "idp_admins.gid = idp_admins_group.id",
  }

  joinTable[1] = map[string]string {
    "type"      : "LEFT",
    "table" : "idp_admins_role",
    "where"     : "idp_admins.gid = idp_admins_role.id",
  }

  total, err := DB.Engine.Table("idp_admins").Join(joinTable[0]["type"], joinTable[0]["table"], joinTable[0]["where"]).Join(joinTable[1]["type"], joinTable[1]["table"], joinTable[1]["where"]).Where(whereData, whereValue...).Count(&table)

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    // 获取列表
    err = DB.Engine.Table("idp_admins").Join(joinTable[0]["type"], joinTable[0]["table"], joinTable[0]["where"]).Join(joinTable[1]["type"], joinTable[1]["table"], joinTable[1]["where"]).Where(whereData, whereValue...).Limit(count, limit).Find(&list)

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


