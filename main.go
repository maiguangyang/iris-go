package main

import (
  "fmt"
  "os"
  "os/signal"

  DB "./database"
  Public "./public"
  Router "./router"
)

var (
  NODE_ENV string
)

func main() {
  // 初始化系统环境变量
  Public.IsNodeEnv(NODE_ENV)

  // 连接数据库
  err := DB.OpenSql()
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
  if DB.Engine != nil && DB.Engine.Ping() == nil {
    DB.Engine.Close()
  }
}

// go build -ldflags "-X 'main.NODE_ENV=master'" ./main.go
