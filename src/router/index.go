package router

import (
  "fmt"
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"
  "crypto/md5"

  AppPublic "../app/public"
  // Admin "../app/admin"
  AppSys "../app/sys"
  Config "../config"
)

func Init() {
  appNew := iris.New()

  // options
  appNew.Options("*", func(ctx context.Context) {
    if Config.IsNodeDev() {
      header(ctx)
    }
  })

  // 定义404错误路由
  appNew.OnErrorCode(iris.StatusNotFound, func(ctx context.Context) {
    header(ctx)
    peter := context.Map{
      "code"       : iris.StatusNotFound,
      "data"       : "",
      "msg"        : "404 not found",
      "status_code": iris.StatusOK,
    }
    ctx.JSON(peter)
  })

  // 定义500错误路由
  appNew.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
    header(ctx)
    peter := context.Map{
      "code"       : iris.StatusInternalServerError,
      "data"       : "",
      "msg"        : "Internal server error",
      "status_code": iris.StatusOK,
    }
    ctx.JSON(peter)
  })

  app := appNew.Party("/", func(ctx context.Context) {
    header(ctx)

    key := ctx.GetHeader("Secret-Key")
    headHash := ctx.GetHeader("Hash")

    if key == "" || headHash == "" {
      peter := context.Map{
        "code"       : 403,
        "data"       : "",
        "msg"        : "key or hash not found",
        "status_code": 200,
      }
      ctx.JSON(peter)
      return
    }

    hash := fmt.Sprintf("%x", md5.Sum([]byte(key[8:108] + "EQUOYpl72tsjwzJnnY")));

    if headHash != hash {
      peter := context.Map{
        "code"       : 403,
        "data"       : "",
        "msg"        : "非法请求",
        "status_code": 200,
      }
      ctx.JSON(peter)
      return
    }
    ctx.Next()
  })

  // 获取公Key
  app.Get("/rsa", AppPublic.GetRsaPubKey)
  // 是否设置配置管理员
  app.Get("/sys/has/user", AppSys.HasUser)
  // 登录系统配置
  app.Post("/sys/login", AppSys.Login)

  // 检测是否设置数据库
  app.Get("/sys/check/database", AppSys.CheckDataBase)


  // 系统配置
  sys := app.Party("/sys", AppSys.Authorization)
  {
    sys.Post("/password", AppSys.ModifyPassword)      // 修改密码
    sys.Get("/config", AppSys.GetConfig)              // 获取配置
    sys.Post("/config", AppSys.EditConfig)            // 修改配置
    sys.Post("/test/sqlopen", AppSys.TestOpen)        // 测试数据库连接
    sys.Get("/database", AppSys.GetDatabase)          // 获取数据库库列表
    sys.Put("/database", AppSys.AddDatabase)          // 添加库
  }

  appNew.Run(iris.Addr(":1874"))
}

// Header设置跨域
func header(ctx context.Context) context.Context {
  if Config.IsNodeDev() {
    ctx.Header("Access-Control-Allow-Origin", "*")
    ctx.Header("Access-Control-Allow-Credentials", "true")
    ctx.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,Secret-Key,Hash");
    ctx.Header("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS");
    ctx.Header("Access-Control-Expose-Headers", "*");
  }
  return ctx
}
