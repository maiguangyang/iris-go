package database

import (
  "fmt"
  // "reflect"
  // "database/sql"
  // "strings"
  // "time"
  "errors"
  "encoding/json"
  "github.com/kataras/iris/context"
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-xorm/xorm"

  "github.com/jinzhu/gorm"

  Auth "../authorization"
  Utils "../utils"
  Public "../public"
)

var Engine *xorm.Engine
var EngineBak *gorm.DB

type Model gorm.Model

// 连接
func OpenSql() error {
  dataSourceName := "root:123456@tcp(192.168.1.235:3306)/idongpin?charset=utf8mb4&parseTime=True&loc=Local"

  var err error

  // 接连数据库，已经连接上，要手动断开连接
  if Engine != nil && Engine.Ping() == nil {
    Engine.Close()
  }

  Engine, _ = xorm.NewEngine("mysql", dataSourceName)

  err = Engine.Ping()
  if err != nil {
    defer Engine.Close()
    return err
  }



  // gorm
  // 接连数据库，已经连接上，要手动断开连接
  if EngineBak != nil && EngineBak.DB().Ping() == nil {
    EngineBak.Close()
  }

  EngineBak, err = gorm.Open("mysql", dataSourceName)
  if err != nil {
    return err
  }

  err = EngineBak.DB().Ping()
  if err != nil {
    defer EngineBak.Close()
    return err
  }

  gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
    return "idp_" + defaultTableName;
  }

  EngineBak.SingularTable(true)
  EngineBak.DB().SetMaxIdleConns(2000)
  EngineBak.DB().SetMaxOpenConns(10000)
  HasInitTable()

  return nil

}

// 判断初始化表是否已经存在，不存在则创建
func HasInitTable() {
  // 用户组
  type AdminGroup struct {
    gorm.Model
    Name string
    Aid int64
  }

  var group AdminGroup
  group.Name      = "超级管理员"
  group.Aid       = 1

  has := EngineBak.HasTable(&AdminGroup{})
  if (has == false) {
    EngineBak.Exec(IDP_ADMIN_GROUP)
    EngineBak.Create(&group)
  }

  // 角色表
  type AdminRole struct {
    gorm.Model
    Name string
    Gid int64
    Aid int64
  }

  var role AdminRole
  role.Name = "超级管理员"
  role.Gid  = 1
  role.Aid  = 1

  has = EngineBak.HasTable(&AdminRole{})
  if (has == false) {
    EngineBak.Exec(IDP_ADMIN_ROLE)
    EngineBak.Create(&role)
  }

  // 人员表
  type Admins struct {
    gorm.Model
    Phone string
    Password string
    Username string
    Gid int64
    Rid int64
    Aid int64
    Super int64
  }

  var admin Admins
  admin.Phone    = "13800138000"
  admin.Password = Public.EncryptPassword("123456")
  admin.Username = "admin"
  admin.Gid      = 1
  admin.Rid      = 1
  admin.Aid      = 1
  admin.Super    = 2

  has = EngineBak.HasTable(&Admins{})
  if (has == false) {
    EngineBak.Exec(IDP_ADMIN)
    EngineBak.Create(&admin)
  }


  // 权限配置表
  has = EngineBak.HasTable("idp_auth_set")
  if has == false {
    EngineBak.Exec(IDP_AUTH_SET)
  }

  // 员工资料
  has = EngineBak.HasTable("idp_admin_archive")
  if has == false {
    EngineBak.Exec(IDP_ADMIN_ARCHIVE)
  }

  // 权限表
  type AdminAuth struct {
    gorm.Model
    Rid int64
    Sid string
    Content string
    Auth int64
  }

  var auth AdminAuth
  auth.Rid = 1
  auth.Sid = "*"
  auth.Content = "*"
  auth.Auth = 1

  has = EngineBak.HasTable(&AdminAuth{})
  if (has == false) {
    EngineBak.Exec(IDP_ADMIN_AUTH)
    EngineBak.Create(&auth)
  }

  // fmt.Println(Public.DecryptPassword("admin", "c02f8a145384b65099b17b9d64dd03e1"))

}

func CheckAdminAuth(ctx context.Context, table string) (bool, bool, int, error) {
  type IdpAdminAuth struct {
    Id int64 `json:"id"`
    Rid int64 `json:"rid"`
    Sid string `json:"sid"`
    Content string `json:"content"`
    Auth int64 `json:"auth"`
    UpdatedAt int64 `json:"updated_at" xorm:"updated"`
    CreatedAt int64 `json:"created_at" xorm:"created"`
  }

  // 获取服务端用户信息
  reqData, err := Auth.HandleUserJWTToken(ctx, "admin")
  if err != nil {
    return false, false, 407, err
  }

  rid   := reqData["rid"].(string)
  super := int64(reqData["super"].(float64))

  // 如果是超级账户的话，直接返回所有权限
  if super == 2 {
    return true, true, 0, nil
  }

  list := make([]IdpAdminAuth, 0)
  has, auth, err := AuthData(ctx, &list, rid, table)

  if err != nil {
    return false, false, 407, err
  }

  return has, auth, 0, nil
}

// 返回用户的权限
func AuthData(ctx context.Context, str interface{}, rid, table string) (bool, bool, error) {
  err := Engine.Desc("id").Where("rid in(" + rid + ")").Limit(10000, 0).Find(str)

  if err != nil {
    return false, false, err
  }

  data,_ := json.Marshal(str)

  method := ctx.Method()
  methodType := map[string]string {
    "GET"    : "info",
    "POST"   : "add",
    "PUT"    : "edit",
    "DELETE" : "del",
  }

  id, err := ctx.Params().GetInt64("id")

  // 如果是详情
  if id == -1 {
    methodType["GET"] = "list"
    err = errors.New("操作权限不足")
  }

  // if err == nil {
  //   methodType["GET"] = "info"
  // }

  var list []map[string]interface{}
  _ = json.Unmarshal([]byte(string(data)), &list)

  var has bool     = false
  var hasAuth bool = false

  if len(list) > 0 {
    var dat []map[string]interface{}
    for _, v := range list{
      content := v["content"].(string)
      _ = json.Unmarshal([]byte(content), &dat)

      if int64(v["auth"].(float64)) == 1 {
        hasAuth = true
      } else {
        hasAuth = false
      }

      if len(dat) > 0 {
        for _, v := range dat{
          tye := v["type"].(map[string]interface{})
          for k, _ := range tye{
            if methodType[method] == k && v["name"].(string) == table {
                return true, hasAuth, nil
              } else {
                has = false
                err = errors.New("操作权限不足")
              }
          }

        }
      }
    }
  }

  return has, hasAuth, err
}


// 返回权限
func checkAuthValue(ctx context.Context, data map[string][]interface{}, table string) (bool, error) {
  method := ctx.Method()
  var index int = -1


  if method == "GET" {
    _, err := ctx.Params().GetInt64("id")
    // 如果是详情
    if err == nil {
      index = Utils.IndexOf(data["INFO"], table)
    } else {
      // arr := strings.Split(data["GET"]), ",")
      index = Utils.IndexOf(data["GET"], table)
    }

  }
  if method == "POST" {
    index = Utils.IndexOf(data["POST"], table)
  }

  if method == "PUT" {
    index = Utils.IndexOf(data["PUT"], table)
  }

  if method == "DELETE" {
    index = Utils.IndexOf(data["DELETE"], table)
  }

  if index == -1 {
    return false, errors.New("操作权限不足")
  }

  return true, nil
}


// 返回分页、总数、limit
func Limit(ctx context.Context) (int64, int, int, map[string]interface{}) {
  page    := Utils.StrToInt64(ctx.URLParam("page"))
  count   := Utils.StrToInt(ctx.URLParam("count"))
  filters := ctx.URLParam("filters")

  var dat map[string]interface{}
  _ = json.Unmarshal([]byte(filters), &dat)

  if page <= 0 {
    page = 1
  }

  if count <= 0 {
    count = 20
  }

  limit := (int(page) - 1) * int(count)

  return page, count, limit, dat
}

func IsWhereEmpty(data, str interface{}) string {

  if data.(string) != "" {
    return data.(string) + " and " + str.(string)
  }
  return " " + str.(string)
}


func CheckErr(err error) {
  if err != nil {
    fmt.Println(err.Error())
  }
}


