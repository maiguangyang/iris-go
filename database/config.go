package database

import (
  "fmt"
  // "reflect"
  // "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-xorm/xorm"
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

  // 测试数据
  // type IdpAdmins struct {
  //   Id int64
  //   Username string
  //   CreatedAt int64 `xorm:"created"`
  // }

  // var table JdAdmins
  // table.Id = 1

  // err = Post(&table)
  // Delete(&table)

  // Get(&table, `id < ? && username = ?`, []interface{}{5, "123"})

  return nil

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
  fmt.Println(table)
  return table, has
}


