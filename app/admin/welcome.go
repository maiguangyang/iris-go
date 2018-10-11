package admin

import(
  "fmt"

  "github.com/kataras/iris/context"

  Auth "../../authorization"
  Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

// 登陆
func Login(ctx context.Context) {

  type ReqData struct {
    Phone string `json:"phone"`
    Password string `json:"password"`
  }
  var reqData ReqData

  type IdpAdmins struct {
    Id int64 `json:"id"`
    // LoginTime int64 `xorm:"created"`
  }

  var table IdpAdmins


  ctx.ReadJSON(&reqData)

  phone    := reqData.Phone
  Password := Public.EncryptPassword(reqData.Password)

  has, err := DB.Engine.Where("state = 0 and phone = ? and password = ?", phone, Password).Get(&table)

  data := context.Map{}

  if err != nil {
    data = Utils.NewResData(401, err.Error(), ctx)
  } else if has == true {
    data = Utils.NewResData(200, Auth.SetToken(table, "admin"), ctx)
  } else {
    data = Utils.NewResData(401, "请检查账号密码输入是否正确", ctx)
  }

  ctx.JSON(data)

}

// 用户详情
func Detail (ctx context.Context) {
  res := GetUserDetail(ctx.GetHeader("Authorization"), ctx)
  ctx.JSON(res)
}

// 获取用户信息方法
func GetUserDetail(author string, ctx context.Context) context.Map {
  userinfo, _ := Auth.DecryptToken(author, "admin")
  reqData := userinfo.(map[string]interface{})

  if len(reqData) <= 0 {
    return context.Map{}
  }

  type IdpAdmins struct {
    Id int64 `json:"id"`
    Phone string `json:"phone"`
    Realname string `json:"realname"`
    Nickname string `json:"nickname"`
    Avatar string `json:"avatar"`
    Sex int64 `json:"sex"`
    Identity string `json:"identity"`
    Groups int64 `json:"groups"`
    Roles int64 `json:"roles"`
    LoginCount int64 `json:"login_count"`
    LastTime int64 `json:"last_time"`
    CreatedAt int64 `json:"created_at"`
    // LoginTime int64 `xorm:"created"`
  }

  var table IdpAdmins

  table.Id = int64(reqData["id"].(float64))

  has, err := DB.Engine.Get(&table)
  if has == true {
    return Utils.NewResData(200, table, ctx)
  }

  fmt.Println(has, err, table)
  return Utils.NewResData(200, "该用户不存在", ctx)
}

