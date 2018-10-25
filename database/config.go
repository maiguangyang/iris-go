package database

import (
  "fmt"
  // "reflect"
  // "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-xorm/xorm"
  Public "../public"
)

var Engine *xorm.Engine

// 连接
func OpenSql() error {
  dataSourceName := "root:123456@tcp(192.168.0.234:3306)/idongpin?charset=utf8mb4"

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
    Value int64
    CreatedAt int64 `xorm:"created"`
  }
  var group IdpAdminsGroup
  group.Name  = "超级管理员"
  group.Value = 0

  has, _ := Engine.IsTableExist("idp_admins_group")
  empty, _  := Engine.IsTableEmpty(&group)

  if has == false {
    CreateTables(IDP_ADMIN_GROUP, group)
  } else if empty == true {
    Post(group)
  }

  // 角色表
  type IdpAdminsRole struct {
    Id int64
    Name string
    Value int64
    CreatedAt int64 `xorm:"created"`
  }
  var role IdpAdminsRole
  role.Name  = "超级管理员"
  role.Value = 0

  has, _ = Engine.IsTableExist("idp_admins_role")
  empty, _  = Engine.IsTableEmpty(&role)

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
    Nickname string
    Groups int64
    Roles int64
    CreatedAt int64 `xorm:"created"`
  }
  var admin IdpAdmins
  admin.Phone    = "13800138000"
  admin.Password = Public.EncryptPassword("123456")
  admin.Nickname = "admin"
  admin.Groups   = 0
  admin.Roles    = 0

  has, _ = Engine.IsTableExist("idp_admins")
  empty, _  = Engine.IsTableEmpty(&admin)

  if has == false {
    CreateTables(IDP_ADMIN, admin)
  } else if empty == true {
    Post(admin)
  }

  // 权限表
  type IdpAuth struct {
    Id int64
    Content string
    CreatedAt int64 `xorm:"created"`
  }

  var auth IdpAuth
  auth.Content = ""

  has, _ = Engine.IsTableExist("idp_auth")
  empty, _  = Engine.IsTableEmpty(&auth)

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
  } else {
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
func Delete(table interface{}) error {
  // TODU 用户权限验证
  _, err := Engine.Delete(table)
  return err
}

// 查询记录
func Get(table, where interface{}, value []interface{}) (interface{}, bool) {
  // 单条记录
  has, _ := Engine.Where(where, value...).Get(table)
  return table, has
}

func Find(table interface{}) error {
  err := Engine.Find(table)
  return err
}

func CheckErr(err error) {
  if err != nil {
    fmt.Println(err.Error())
  }
}


