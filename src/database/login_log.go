package database

// 插入登录记录
func LoginLogAdd(uid int64, loginType int, token, content, ip string) error {
  sql := `INSERT login_log SET uid=?, type=?, token=?, content=?, ip=?`
  _, err := DB.Exec(sql, uid, loginType, token, content, ip)
  return err
}

// 检查记录
func LoginLogHas(token, ip string) bool {
  sql := `SELECT id FROM login_log WHERE token=? AND ip=?`
  var id int64
  err := DB.QueryRow(sql, token, ip).Scan(&id)
  if err != nil || id == 0 {
    return false
  }
  return true
}

// 检查登录提交数据
func LoginLogHasContent(content string) bool {
  sql := `SELECT id FROM login_log WHERE content=?`
  var id int64
  err := DB.QueryRow(sql, content).Scan(&id)
  if err != nil || id == 0 {
    return false
  }
  return true
}

// 清空登录提交数据
func LoginLogDelContent() error {
  sql := `UPDATE login_log SET content = '' WHERE content != ''`
  _, err := DB.Exec(sql)
  return err
}