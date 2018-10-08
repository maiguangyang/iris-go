package database

import (
  "errors"
  "strings"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  config "../config"
)

var DB *sql.DB
type Map map[string]interface{}

// 连接
func OpenSql() error {
  configJson := config.ConfigJson
  if configJson.Database == "" {
    return errors.New("未设置数据库连接")
  }

  var dataSourceName string
  // 登录环境
  if config.IsNodeDev() {
    if configJson.DevDataUser == "" || configJson.DevDataPassword == "" || configJson.DevDataIp == "" || configJson.DevDataPort == "" {
      return errors.New("未设置数据库连接")
    }
    dataSourceName = configJson.DevDataUser + ":" + configJson.DevDataPassword + "@tcp(" + configJson.DevDataIp + ":" + configJson.DevDataPort + ")/"
  } else {
    if configJson.DataUser == "" || configJson.DataPassword == "" || configJson.DataIp == "" || configJson.DataPort == "" {
      return errors.New("未设置数据库连接")
    }
    dataSourceName = configJson.DataUser + ":" + configJson.DataPassword + "@tcp(" + configJson.DataIp + ":" + configJson.DataPort + ")/"
  }

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
    if d == configJson.Database {
      // 重新连接，不能用use语句，因为rows.Close()时use会被清空
      return SetDB(dataSourceName + configJson.Database + "?charset=" + configJson.Charset)
    }
  }

  // 找不到库，创建库
  _, err = DB.Exec("create DATABASE " + configJson.Database)
  if err != nil {
    defer DB.Close()
    return err
  }

  // 重新连接，不能用use语句，因为rows.Close()时use会被清空
  err = SetDB(dataSourceName + configJson.Database + "?charset=" + configJson.Charset)
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
  DB, err = sql.Open(config.ConfigJson.DriverName, dataSourceName)

  err = DB.Ping()
  if err != nil {
    defer DB.Close()
    return err
  }

  // 检查必需表
  err = HasTable()
  return err
}

// 读取数据库数据列表
func GetListData(sql string, qAs QueryArgs) ([]Map, error) {
  stmt, err := DB.Prepare(sql)
  if err != nil {
    return []Map{}, err
  }
  defer stmt.Close()

  rows, err := stmt.Query(qAs...)
  if err != nil {
    return []Map{}, err
  }
  defer rows.Close()

  return RowsScan(rows)
}

// 读取数据
func RowsScan(rows *sql.Rows) ([]Map, error) {
  ListData := []Map{}

  columns, err := rows.Columns()
  if err != nil {
    return ListData, err
  }

  dest := make([]interface{}, len(columns))
  values := make([]interface{}, len(columns))
  for i := range values {
    dest[i] = &values[i]
  }

  for rows.Next() {
    err = rows.Scan(dest...)
    if err != nil {
      return ListData, err
    }

    d := Map{}
    dl := map[string]Map{}
    for i, v := range values {
      switch v.(type) {
      case []byte:
        v = string(v.([]byte))
      }

      key := columns[i]
      keys := strings.Split(key, ".")
      if len(keys) > 1 {
        if dl[keys[0]] == nil {
          dl[keys[0]] = Map{}
        }
        dl[keys[0]][keys[1]] = v
      } else {
        d[key] = v
      }
    }
    for k, v := range dl {
      d[k] = v
    }
    ListData = append(ListData, d)
  }

  return ListData, nil
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
