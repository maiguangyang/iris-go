package database

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

)

var DB *sql.DB
type Map map[string]interface{}

// 连接
func OpenSql() error {
  database := "idongpin"
  dataSourceName := "root:123456@tcp(192.168.0.234:3306)/"

  var err error
  // 接连数据库
  err = SetDB(dataSourceName)
  if err != nil {
    return err
  }

  // 查询库
  rows, err := DB.Query("SHOW DATABASES")
  if err != nil {
    defer DB.Close()
    return err
  }

  defer rows.Close()
  for rows.Next() {
    var d string
    err = rows.Scan(&d)
    if err != nil {
      defer DB.Close()
      return err
    }
    // 找到库
    if d == database {
      // 重新连接，不能用use语句，因为rows.Close()时use会被清空
      return SetDB(dataSourceName + database + "?charset=utf8mb4")
    }
  }

  // 找不到库，创建库
  _, err = DB.Exec("create DATABASE " + database)
  if err != nil {
    defer DB.Close()
    return err
  }

  // 重新连接，不能用use语句，因为rows.Close()时use会被清空
  err = SetDB(dataSourceName + database + "?charset=utf8mb4")
  if err != nil {
    return err
  }
  return nil
}

// 设置DB
func SetDB(dataSourceName string) error {
  // 已经连接上，要手动断开连接
  if DB != nil && DB.Ping() == nil {
    DB.Close()
  }

  var err error
  DB, err = sql.Open("mysql", dataSourceName)

  err = DB.Ping()
  if err != nil {
    defer DB.Close()
    return err
  }

  // 检查必需表
  // err = HasTable()
  return nil
}

// 测试连接
func TestOpen(driverName, dataSourceName string) error {
  db, err := sql.Open(driverName, dataSourceName)
  if err != nil {
    return err
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    return err
  }

  return nil
}
