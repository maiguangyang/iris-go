package admin

import(
  // "fmt"
  "time"
  // "reflect"
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
  Super int64 `json:"super"`
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
  DeletedAt int64 `json:"deleted_at"`
  UpdatedAt int64 `json:"updated_at" xorm:"updated"`
  CreatedAt int64 `json:"created_at" xorm:"created"`
}

type Admins struct {
  DB.Model
  Phone string `json:"phone"`
  Password string `json:"password"`
  Username string `json:"username"`
  Sex int64 `json:"sex"`
  Super int64 `json:"super"`
  Gid string `json:"gid"`
  Rid string `json:"rid"`
  Aid int64 `json:"aid"`
  Money int64 `json:"money"`
  State int64 `json:"state"`
  JobState int64 `json:"job_state"`
  LoginCount int64 `json:"login_count"`
  LoginTime time.Time `json:"login_time"`
  LastTime time.Time `json:"last_time"`
  LoginIp string `json:"login_ip"`
  LastIp string `json:"last_ip"`
  EntryTime time.Time `json:"entry_time"`
  QuitTime time.Time `json:"quit_time"`
  TrialTime time.Time `json:"trial_time"`
  ContractTime time.Time `json:"contract_time"`

}

type UserDetailGroup struct {
  Group []*IdpAdminsRole `json:"group" xorm:"extends"`
  IdpAdmins `xorm:"extends"`
}

func (UserDetailGroup) TableName() string {
  return "idp_admins"
}


// 登陆
func Login(ctx context.Context) {

  var table Admins
  timestamp := time.Now()

  // 根据不同环境返回数据
  err := Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  table.Password  = Public.EncryptPassword(table.Password)

  data := context.Map{}

  result := DB.EngineBak.Where("phone =? and password =?", table.Phone, table.Password).First(&table)

  if result.Error != nil {
    data = Utils.NewResData(1, "请检查账号密码输入是否正确", ctx)
  } else {
    if table.State == 2 {
      data = Utils.NewResData(1, "您的账户已被禁用", ctx)
    } else {
      // 返回前端的Token
      ip := ctx.RemoteAddr()
      token := Auth.SetToken(context.Map{
        "id": table.ID,
        "gid": table.Gid,
        "rid": table.Rid,
        "super": table.Super,
      }, Public.EncryptMd5(ip + Auth.SecretKey["admin"].(string)), "admin")

      table.LoginCount = table.LoginCount + 1
      table.LastTime   = table.LoginTime

      table.UpdatedAt  = timestamp
      table.LoginTime  = timestamp
      table.LastIp     = table.LoginIp
      table.LoginIp    = ip

      // 更新用户登陆信息
      result = DB.EngineBak.Model(&table).Updates(&table)
      if result.Error != nil {
        ctx.JSON(Utils.NewResData(1, result.Error, ctx))
        return
      }

      data = Utils.NewResData(0, token, ctx)
    }
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

  var table Admins
  // table.Id = uid

  result := DB.EngineBak.Where("id=?", uid).Omit("password").First(&table)
  if result.Error != nil {
    return Utils.NewResData(1, result.Error, ctx)
  }

  return Utils.NewResData(0, table, ctx)
}

// 获取用户的前端路由
func HandleAdminRoutes(ctx context.Context) {
  // 获取服务端用户信息
  reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  data := context.Map{}
  if !Utils.IsEmpty(reqData["rid"]) {
    sql := `select auth.id, auth.rid, auth.sid, auth.content, a.table_name, a.id as s_id, a.name, a.routes, a.sub_id, b.id as b_id, b.routes as sub_routes from idp_admin_auth as auth left join idp_auth_set as a ON FIND_IN_SET(a.id, auth.sid) left join idp_auth_set as b ON FIND_IN_SET(b.id, a.sub_id) where auth.rid = ?`
    rows, err := DB.Engine.QueryString(sql, reqData["rid"])

    if err == nil {
      data = Utils.NewResData(0, rows, ctx)
    } else {
      data = Utils.NewResData(1, "获取前端路由错误", ctx)
    }
  }

  ctx.JSON(data)
}
