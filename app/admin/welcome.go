package admin

import(
  // "fmt"
  // "encoding/json"
  "github.com/kataras/iris/context"

  Auth "../../authorization"
  Public "../../public"
  Utils "../../utils"
  DB "../../database"
)

type IdpAdmins struct {
  Id int64 `json:"id"`
  Phone string `json:"phone"`
  Password string `json:"password"`
  Username string `json:"username"`
  Sex int64 `json:"sex"`
  Gid string `json:"gid"`
  Rid string `json:"rid"`
  Aid int64 `json:"aid"`
  Money int64 `json:"money" xorm:"default(0)"`
  State int64 `json:"state"`
  JobState int64 `json:"job_state"`
  LoginCount int64 `json:"login_count" xorm:"version"`
  LoginTime int64 `json:"login_time" xorm:"updated"`
  LastTime int64 `json:"last_time"`
  LoginIp string `json:"login_ip"`
  LastIp string `json:"last_ip"`
  EntryTime int64 `json:"entry_time"`
  QuitTime int64 `json:"quit_time"`
  TrialTime int64 `json:"trial_time"`
  ContractTime int64 `json:"contract_time"`
  DeletedAt int64 `json:"deleted_at" xorm:"deleted"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
  // LoginTime int64 `xorm:"created"`
}



type UserDetailGroup struct {
  // Role IdpAdminsGroup `json:"role" xorm:"extends"`
  // Group IdpAdminsRole `json:"group" xorm:"extends"`
  Group string `json:"group"`
  IdpAdmins `xorm:"extends"`
}

func (UserDetailGroup) TableName() string {
  return "idp_admins"
}

// 登陆
func Login(ctx context.Context) {

  var table IdpAdmins

  // 根据不同环境返回数据
  err := Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }


  table.Password = Public.EncryptPassword(table.Password)

  has, err := DB.Engine.Get(&table)

  data := context.Map{}

  if err != nil {
    data = Utils.NewResData(1, err.Error(), ctx)
  } else if has == true {
    if table.State == 2 {
      data = Utils.NewResData(1, "您的账户已被禁用", ctx)
    } else {
      // 返回前端的Token
      ip := ctx.RemoteAddr()
      token := Auth.SetToken(context.Map{
        "id": table.Id,
        "gid": table.Gid,
        "rid": table.Rid,
      }, Public.EncryptMd5(ip + Auth.SecretKey["admin"].(string)), "admin")

      table.LastTime = table.LoginTime
      table.LastTime = table.LoginTime
      table.LastIp   = table.LoginIp
      table.LoginIp  = ip

      // 更新用户登陆信息
      // UpdataUserLoginInfo(table)

      _, err := DB.Engine.Id(table.Id).Update(&table)
      if err != nil {
        ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
        return
      }

      data = Utils.NewResData(0, token, ctx)
    }
  } else {
    data = Utils.NewResData(1, "请检查账号密码输入是否正确", ctx)
  }

  ctx.JSON(data)

}

// 用户详情
func Detail (ctx context.Context) {
  author      := ctx.GetHeader("Authorization")
  userinfo, _ := Auth.DecryptToken(author, "admin")
  reqData     := userinfo.(map[string]interface{})

  if len(reqData) <= 0 {
    ctx.JSON(context.Map{})
    return
  }

  id := int64(reqData["id"].(float64))

  res := GetUserDetail(id, ctx)

  ctx.JSON(res)
}

// 获取用户信息方法
func GetUserDetail(uid int64, ctx context.Context) context.Map {

  var table IdpAdmins
  table.Id = uid

  has, err := DB.Engine.Omit("password").Get(&table)
  if err != nil {
    return Utils.NewResData(1, err.Error(), ctx)
  }

  data := context.Map{}
  if has == true {
    data = Utils.NewResData(0, table, ctx)
  } else {
    data = Utils.NewResData(1, "该账户不存在", ctx)
  }

  return data

}
