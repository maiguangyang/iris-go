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

type IdpAdminsGroup struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Aid int64 `json:"aid"`
  State int64 `json:"state"`
  DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}


// 用户组列表
func GroupList (ctx context.Context) {
  // 获取分页、总数、limit
  page, count, limit, _ := DB.Limit(ctx)
  list := make([]IdpAdminsGroup, 0)


  // 获取统计总数
  var table IdpAdminsGroup
  data := context.Map{}

  total, err := DB.Engine.Count(&table)

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else {
    // 获取列表
    err = DB.Engine.Desc("id").Limit(count, limit).Find(&list)

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
  var table IdpAdminsGroup

  var rules Utils.Rules

  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      return Utils.NewResData(1, err.Error(), ctx)
    }

    reqData := decData.(map[string]interface{})

    table.Id    = int64(reqData["id"].(float64))
    table.Name  = reqData["name"].(string)
    table.State = int64(reqData["state"].(float64))

  } else {
    ctx.ReadJSON(&table)
  }

  // 验证参数
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
    return Utils.NewResData(1, tipsText + "失败", ctx)
  }

  return Utils.NewResData(0, tipsText + "成功", ctx)
}

// 删除
func GroupDel (ctx context.Context) {
  var table IdpAdminsGroup

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
  var exist IdpAdminsGroup
  value := []interface{}{table.Id}
  has, err := DB.Engine.Where("id=?", value...).Exist(&exist)

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


