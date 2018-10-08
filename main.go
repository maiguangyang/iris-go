package main

import (
  "fmt"
  "os"
  "os/signal"

  Config "./src/config"
  RsaKey "./src/rsa_key"
  Database "./src/database"
  Router "./src/router"
)

var (
  NODE_ENV string
)

func main() {
  // 运行环境
  Config.Ldflags.NODE_ENV = NODE_ENV

  // 读取“./config.json”配置文件
  err := Config.ConfigJson.Read()
  if err != nil {
    fmt.Println("读取“./config.json”配置文件出错：", err.Error())
    return
  }

  // rsa加密
  err = RsaKey.Gen()
  if err != nil {
    fmt.Println("生成rsa key出错：", err.Error())
    return
  }


  // 连接数据库
  err = Database.OpenSql()
  if err != nil {
    fmt.Println("连接数据库失败：", err.Error())
  } else {
    err = Database.LoginLogDelContent()
    if err != nil {
      fmt.Println("清空登录提交数据出错：", err.Error())
    }
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
