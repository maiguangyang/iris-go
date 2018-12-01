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
  DB.Model
  Phone string `json:"phone"`
  Username string `json:"username"`
  Sex int64 `json:"sex"`
  Super int64 `json:"super" gorm:"default:1"`
  Gid string `json:"gid"`
  Rid string `json:"rid"`
  Aid int64 `json:"aid"`
  Money int64 `json:"money"`
  State int64 `json:"state"`
  JobState int64 `json:"job_state"`
  LoginCount int64 `json:"login_count"`
  LoginTime int64 `json:"login_time" gorm:"default:null"`
  LastTime int64 `json:"last_time" gorm:"default:null"`
  LoginIp string `json:"login_ip"`
  LastIp string `json:"last_ip"`
  EntryTime int64 `json:"entry_time" gorm:"default:null"`
  QuitTime int64 `json:"quit_time" gorm:"default:null"`
  TrialTime int64 `json:"trial_time" gorm:"default:null"`
  ContractTime int64 `json:"contract_time" gorm:"default:null"`

  // Groups []IdpAdminGroups `json:"groups" gorm:"foreignkey:Id;association_foreignkey:Gid"`
  // Roles []IdpAdminRoles `json:"roles" gorm:"FOREIGNKEY:Id"`
}

type IdpAdminsPass struct {
  IdpAdmins
  Password string `json:"password"`
}

func (IdpAdminsPass) TableName() string {
  return "idp_admins"
}

// 登陆
func Login(ctx context.Context) {

  var table IdpAdminsPass
  timestamp := time.Now().Unix()

  // 根据不同环境返回数据
  err := Utils.ResNodeEnvData(&table, ctx)
  if err != nil {
    ctx.JSON(Utils.NewResData(1, err.Error(), ctx))
    return
  }

  table.Password  = Public.EncryptPassword(table.Password)

  data := context.Map{}


  result := DB.EngineBak.Model(&table).First(&table, 1)

  if result.Error != nil {
    data = Utils.NewResData(1, "请检查账号密码输入是否正确", ctx)
  } else {
    if table.State == 2 {
      data = Utils.NewResData(1, "您的账户已被禁用", ctx)
    } else {
      // 返回前端的Token
      ip := ctx.RemoteAddr()
      token := Auth.SetToken(context.Map{
        "id": table.Id,
        "gid": table.Gid,
        "rid": table.Rid,
        "super": table.Super,
      }, Public.EncryptMd5(ip + Auth.SecretKey["admin"].(string)), "admin")

      table.LoginCount = table.LoginCount + 1
      table.LastTime   = table.LoginTime

      table.LoginTime  = timestamp
      table.LastIp     = table.LoginIp
      table.LoginIp    = ip

      // 更新用户登陆信息
      result = DB.EngineBak.Model(&table).UpdateColumns(&table)
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

  var table IdpAdmins
  // table.Id = uid

  if err := DB.EngineBak.Model(&table).Where("id=?", uid).First(&table).Error; err != nil {
    return Utils.NewResData(1, err, ctx)
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

  data := Utils.NewResData(1, context.Map{}, ctx)
  if !Utils.IsEmpty(reqData["rid"]) {
    type dataJson struct {
      Id int64 `json:"id" gorm:"primary_key"`
      Rid int64 `json:"rid"`
      Sids string `json:"sids"`
      Content string `json:"content"`
      TableName string `json:"table_name"`
      Sid string `json:"sid"`
      Name string `json:"name"`
      Routes string `json:"routes"`
      SubId string `json:"sub_id"`
      Bid string `json:"bid"`
      SubRoutes string `json:"sub_routes"`
    }

    var list dataJson
    sql := `select auth.id, auth.rid, auth.sid as sids, auth.content, a.table_name, a.id as sid, a.name, a.routes, a.sub_id, b.id as bid, b.routes as sub_routes from idp_admin_auth as auth left join idp_auth_set as a ON FIND_IN_SET(a.id, auth.sid) left join idp_auth_set as b ON FIND_IN_SET(b.id, a.sub_id) where auth.rid = ?`
    if err := DB.EngineBak.Raw(sql, reqData["rid"]).Scan(&list).Error; err != nil {
      data = Utils.NewResData(1, err, ctx)
    } else {
      data = Utils.NewResData(0, list, ctx)
    }
  }

  ctx.JSON(data)
}
