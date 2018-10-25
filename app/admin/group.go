package admin

import (
  // "fmt"
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
    data = Utils.NewResData(200, list, ctx)
  }

  ctx.JSON(data)

}

func GroupAdd (ctx context.Context) {
  var table IdpAdminsGroup
  ctx.ReadJSON(&table)

  _ = DB.Post(&table)
}