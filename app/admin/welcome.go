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
type IdpAdmins struct {
  Id int64 `json:"id"`
  State int64 `json:"state"`
  // LoginTime int64 `xorm:"created"`
}


func Login(ctx context.Context) {

  type ReqData struct {
    Phone string `json:"phone"`
    Password string `json:"password"`
  }
  var reqData ReqData

  var table IdpAdmins


  ctx.ReadJSON(&reqData)

  phone    := reqData.Phone
  Password := Public.EncryptPassword(reqData.Password)


  has, err := DB.Engine.Where("state = 0 and phone = ? and password = ?", phone, Password).Get(&table)

  data := context.Map{}

  fmt.Println(has, err)
  if err != nil {
    data = Utils.NewResData(401, err.Error(), ctx)
  } else if has == true {
    data = Utils.NewResData(200, table, ctx)
  } else {
    data = Utils.NewResData(401, "Authorization 未授权", ctx)
  }

  ctx.JSON(data)

}

// 用户详情
func Detail (ctx context.Context) {
  GetUserDetail(ctx.GetHeader("Authorization"))
  // ctx.JSON(res)
}

// 获取用户信息方法

func GetUserDetail(author string) {
  userinfo, _ := Auth.DecryptToken(author, "admin")
  fmt.Println(userinfo)
}