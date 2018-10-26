package admin

import (
  // "fmt"
  // "reflect"
  "github.com/kataras/iris/context"

  // Auth "../../authorization"
  Utils "../../utils"
  DB "../../database"
)

type IdpAdminsGroup struct {
  Id int64 `json:"id"`
  Name string `json:"name"`
  Value int64 `json:"value"`
  State int64 `json:"state"`
  DeletedAt int64 `json:"deleted_at"`
  UpdatedAt int64 `json:"updated_at"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

// 用户组列表
func GroupList (ctx context.Context) {
  // userinfo, _ := Auth.DecryptToken(ctx.GetHeader("Authorization"), "admin")
  // reqData     := userinfo.(map[string]interface{})

  list := make([]IdpAdminsGroup, 0)
  err := DB.Find(&list)

  data := context.Map{}

  if err != nil {
    data = Utils.NewResData(404, err.Error(), ctx)
  } else {
    data = Utils.NewResData(0, list, ctx)
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

  err := DB.Put(table.Id, &table)

  data := context.Map{}
  if err == nil {
    data = Utils.NewResData(0, "修改成功", ctx)
  } else {
    data = Utils.NewResData(1, "修改失败", ctx)
  }

  ctx.JSON(data)
}
