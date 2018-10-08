package database

import (
  "time"
)

// 添加记录
func ErrorLogAdd(typeId int64, token, data string) error {
  timeNow := time.Now().Unix()
  sql := "INSERT error_log SET tid=?, token=?, data=?, created_time=?"
  _, err := DB.Exec(sql, typeId, token, data, timeNow)
  return err
}
