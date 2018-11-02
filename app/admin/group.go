package admin

import (
  "fmt"
  // "reflect"
  "github.com/kataras/iris/context"

  // Auth "../../authorization"
  Utils "../../utils"
  DB "../../database"
)

type IdpAdminsGroup struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  State int64 `json:"state"`
  DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

// 用户组列表
func GroupList (ctx context.Context) {
  // userinfo, _ := Auth.DecryptToken(ctx.GetHeader("Authorization"), "admin")
  // reqData     := userinfo.(map[string]interface{})

  // 获取分页
  page := Utils.StrToInt64(ctx.URLParam("page"))

  // 获取统计总数
  var table IdpAdminsGroup
  total := DB.Count(&table)

  // 获取列表
  list := make([]IdpAdminsGroup, 0)
  err := DB.Find(&list, page)

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
func GroupDetail (ctx context.Context) {

  var table IdpAdminsGroup
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
func GroupAdd (ctx context.Context) {

  var table IdpAdminsGroup
  ctx.ReadJSON(&table)

  // 验证参数
  rules := Utils.Rules{
    "Name": {
      "required": true,
      // "rgx": "identity",
    },
  }

  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    ctx.JSON(Utils.NewResData(1, errMsgs, ctx))
    return
  }


  // 写入数据库
  err := DB.Post(&table)

  data := context.Map{}
  if err == nil {
    data = Utils.NewResData(0, "添加成功", ctx)
  } else {
    data = Utils.NewResData(1, "添加失败", ctx)
  }

  ctx.JSON(data)
}

// 修改
func GroupPut (ctx context.Context) {

  var table IdpAdminsGroup
  ctx.ReadJSON(&table)


  // 验证参数
  rules := Utils.Rules{
    "Name": {
      "required": true,
      // "rgx": "identity",
    },
  }

  errMsgs := rules.Validate(Utils.StructToMap(table))
  if errMsgs != nil {
    ctx.JSON(Utils.NewResData(1, errMsgs, ctx))
    return
  }

  // 判断数据库里面是否已经存在
  var exist IdpAdminsGroup
  has := DB.Exist(&exist, "id<>? and name=?", []interface{}{table.Id, table.Name})
  fmt.Println(has)

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
func GroupDel (ctx context.Context) {
  var table IdpAdminsGroup
  ctx.ReadJSON(&table)

  err := DB.Delete(&table)

  data := context.Map{}
  if err == nil {
    data = Utils.NewResData(0, "删除成功", ctx)
  } else {
    data = Utils.NewResData(1, "删除失败", ctx)
  }

  ctx.JSON(data)
}


