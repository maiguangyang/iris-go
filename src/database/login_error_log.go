package database

import (
  "time"
)

// 登录错误记录表
func LoginErrorLogAdd(loginType int, ip string) error {
  timeNow := time.Now().Unix()
  sql := `INSERT login_error_log SET type=?, ip=?, created_time=?`
  _, err := DB.Exec(sql, loginType, ip, timeNow)
  return err
}

// 是否限制登录
func IsRestrictLogin(loginType int, ip string) bool {
  timeOhb := time.Now().Add(- time.Hour * 1).Unix()
  sql := `SELECT count(*) FROM login_error_log WHERE type=? AND ip=? AND created_time>?`
  var errNum int
  err := DB.QueryRow(sql, loginType, ip, timeOhb).Scan(&errNum)
  if err != nil {
    return true
  }
  if errNum >= 5 {
    return true
  }
  return false
}
