package database

import (
  "fmt"
  "errors"
  // "reflect"
  // "database/sql"
  "encoding/json"
  "github.com/kataras/iris/context"
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-xorm/xorm"
  Utils "../utils"
  Public "../public"
)

var Engine *xorm.Engine

// 连接
func OpenSql() error {
  dataSourceName := "root:123456@tcp(192.168.1.235:3306)/idongpin?charset=utf8mb4"

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

  HasInitTable()
  return nil

}

// 判断初始化表是否已经存在，不存在则创建
func HasInitTable() {

  // 用户组
  type IdpAdminsGroup struct {
    Id int64
    Name string
    Aid int64
    CreatedAt int64 `xorm:"created"`
  }

  var group IdpAdminsGroup
  group.Name = "超级管理员"
  group.Aid  = 1

  has, _   := Engine.IsTableExist("idp_admins_group")
  empty, _ := Engine.IsTableEmpty(&group)

  if has == false {
    CreateTables(IDP_ADMIN_GROUP, group)
  } else if empty == true {
    Post(group)
  }

  // 角色表
  type IdpAdminsRole struct {
    Id int64
    Name string
    Gid int64
    Aid int64
    CreatedAt int64 `xorm:"created"`
  }

  var role IdpAdminsRole
  role.Name = "超级管理员"
  role.Gid  = 1
  role.Aid  = 1

  has, _   = Engine.IsTableExist("idp_admins_role")
  empty, _ = Engine.IsTableEmpty(&role)

  if has == false {
    CreateTables(IDP_ADMIN_ROLE, role)
  } else if empty == true {
    Post(role)
  }

  // 人员表
  type IdpAdmins struct {
    Id int64
    Phone string
    Password string
    Username string
    Gid int64
    Rid int64
    Aid int64
    CreatedAt int64 `xorm:"created"`
  }

  var admin IdpAdmins
  admin.Phone    = "13800138000"
  admin.Password = Public.EncryptPassword("123456")
  admin.Username = "admin"
  admin.Gid      = 1
  admin.Rid      = 1
  admin.Aid      = 1

  has, _ = Engine.IsTableExist("idp_admins")
  empty, _  = Engine.IsTableEmpty(&admin)

  if has == false {
    CreateTables(IDP_ADMIN, admin)
  } else if empty == true {
    Post(admin)
  }

  // 员工资料
  has, _ = Engine.IsTableExist("idp_admin_archive")

  if has == false {
    CreateTables(IDP_ADMIN_ARCHIVE, nil)
  }

  // 权限表
  type IdpAuth struct {
    Id int64
    Content string
    CreatedAt int64 `xorm:"created"`
  }

  var auth IdpAuth
  auth.Content = ""

  has, _   = Engine.IsTableExist("idp_auth")
  empty, _ = Engine.IsTableEmpty(&auth)

  if has == false {
    CreateTables(IDP_AUTH, auth)
  } else if empty == true {
    Post(auth)
  }

  // fmt.Println(Public.DecryptPassword("admin", "c02f8a145384b65099b17b9d64dd03e1"))

}

// 创建表
func CreateTables(tableName string, table interface{}) {
  _, err := Engine.Exec(tableName)

  if err != nil {
    CheckErr(err)
  } else if table != nil {
    err = Post(table)
    CheckErr(err)
  }
}

// 新增数据
func Post(table interface{}) error {
  // TODU 用户权限验证
  _, err := Engine.Insert(table)
  return err
}

// 更新数据
func Put(id int64, table interface{}) error {
  // TODU 用户权限验证
  _, err := Engine.Id(id).Update(table)
  return err
}

// 删除记录
func Delete(id int64, table interface{}) error {

  has, err := Exist(context.Map{
    "type": 0,
    "table": table,
    "where": "id=?",
    "value": []interface{}{id},
    "sql": "",
  })

  if err != nil {
    return err
  }

  if has == true {
    _, err = Engine.Delete(table)
  } else {
    err = errors.New("删除失败，该记录不存在")
  }

  return err
}

// 查询记录
// func Get(table, where interface{}, value []interface{}) bool {
func Get(object context.Map) bool {
  // 单条记录
  table := object["table"]
  where := object["where"]
  value := object["value"].([]interface{})

  var has bool

  if object["type"] == 1 {
    sqlTxt := object["sql"].(string)
    has, _ = Engine.Where(where, value...).Sql(sqlTxt).Get(table)
  } else {
    has, _ = Engine.Where(where, value...).Get(table)
  }

  return has
}

// Exist查询记录
// func Exist(table, where interface{}, value []interface{}) bool {
func Exist(object context.Map) (bool, error) {
  // 单条记录
  table := object["table"]
  where := object["where"]
  value := object["value"].([]interface{})

  var has bool
  var err error = nil

  if object["type"] == 1 {
    sqlTxt := object["sql"].(string)
    has, err = Engine.Sql(sqlTxt).Where(where, value...).Exist(table)
  } else {
    has, err = Engine.Where(where, value...).Exist(table)
  }

  return has, err
}

// 列表 join["type"]为1的时候使用SQL语句连表
func Find(object context.Map) error {
  table := object["table"]
  var count int64 = 20

  if object["count"].(int64) > 0 {
    count = object["count"].(int64)
  }

  limit := (int(object["page"].(int64)) - 1) * int(count)

  var err error = nil

  if object["type"] == 1 {
    sqlTxt := object["sql"].(string)
    err = Engine.Sql(sqlTxt).Desc("id").Limit(int(count), limit).Find(table)
  } else {
    where := object["where"]
    value := object["value"].([]interface{})
    err = Engine.Desc("id").Where(where, value...).Limit(int(count), limit).Find(table)
  }

  return err
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

// 统计
// func Count(table interface{}) int64 {
func Count(object context.Map) int64 {
  // total, _ := Engine.Count(table)
  // return total
  table := object["table"]
  var total int64

  if object["type"] == 1 {
    sqlTxt := object["sql"].(string)
    res, _ := Engine.Query(sqlTxt)
    if len(res) >= 0 {
      data := string(res[0]["count"])
      total = Utils.StrToInt64(data)
    }
  } else {
    where := object["where"]
    var value []interface{}

    if object["value"] != nil {
      value = object["value"].([]interface{})
    }

    total, _ = Engine.Where(where, value...).Count(table)
  }

  return total

}

func CheckErr(err error) {
  if err != nil {
    fmt.Println(err.Error())
  }
}


