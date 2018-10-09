package main

import (
  "fmt"
  "os"
  "os/signal"

  Database "./database"
  Router "./router"
)

var (
  NODE_ENV string
)

func main() {

  // 连接数据库
  err := Database.OpenSql()
  if err != nil {
    fmt.Println("连接数据库失败：", err.Error())
  }

  // 开启路由
  Router.Init()

  // 安全退出
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt, os.Kill)
  <-c
  // 关闭数据库连接
  if Database.DB != nil && Database.DB.Ping() == nil {
    Database.DB.Close()
  }
}
