package router

import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"

  Public "../public"
  Auth "../authorization"
  AppTest "../app/test"
  Admin "../app/admin"
)

func Init() {
  appNew := iris.New()

  // options
  appNew.Options("*", func(ctx context.Context) {
    header(ctx)
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

    if Public.NODE_ENV {
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

      // 验证hash
      hash := Public.CheckHash(key)

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
    }

    ctx.Next()
  })


  // 检测是否设置数据库
  app.Get("/sys/check/database", AppTest.CheckDataBase)
  app.Post("/sys/check/database", AppTest.CheckDataBasePost)


  // admin
  app.Put("/admin/login", Admin.Login)      // 登陆
  admin := app.Party("/admin", Auth.CheckAuthAdmin)
  {

    admin.Get("/detail", Admin.Detail)                // 账户详情

    // group
    admin.Get("/group", Admin.GroupList)              // 部门列表
    admin.Get("/group/{id:int64}", Admin.GroupDetail) // 部门详情
    admin.Post("/group", Admin.GroupAdd)              // 添加部门
    admin.Put("/group", Admin.GroupPut)               // 修改部门
    admin.Delete("/group", Admin.GroupDel)            // 删除部门

    // role
    admin.Get("/role", Admin.RoleList)              // 角色列表
    admin.Get("/role/{id:int64}", Admin.RoleDetail) // 角色详情
    admin.Post("/role", Admin.RoleAdd)              // 添加角色
    admin.Put("/role", Admin.RolePut)               // 修改角色
    admin.Delete("/role", Admin.RoleDel)            // 删除角色

    // user
    admin.Get("/user", Admin.UserList)              // 员工列表
    admin.Get("/user/{id:int64}", Admin.UserDetail) // 员工详情
    admin.Post("/user", Admin.UserAdd)              // 添加员工
    admin.Put("/user", Admin.UserPut)               // 修改员工
    admin.Delete("/user", Admin.UserDel)            // 删除员工

    // set auth
    admin.Get("/config/table", Admin.CongifTable)                    // 获取数据库表
    admin.Get("/config/routes", Admin.CongifRoutesList)              // 权限路由列表
    admin.Get("/config/routes/{id:int64}", Admin.CongifRoutesDetail) // 权限路由详情
    admin.Post("/config/routes", Admin.CongifRoutesAdd)              // 添加权限路由
    admin.Put("/config/routes", Admin.CongifRoutesPut)               // 修改权限路由
    admin.Delete("/config/routes", Admin.CongifRoutesDel)            // 删除权限路由

    // 角色权限设置
    admin.Get("/roleAuth", Admin.AdminAuthList)              // 权限路由列表
    admin.Get("/roleAuth/{id:int64}", Admin.AdminAuthDetail) // 权限路由详情
    admin.Post("/roleAuth", Admin.AdminAuthAdd)              // 添加权限路由
    admin.Put("/roleAuth", Admin.AdminAuthPut)               // 修改权限路由
    admin.Delete("/roleAuth", Admin.AdminAuthDel)            // 删除权限路由

    // sys.Post("/test/sqlopen", AppSyßs.TestOpen)        // 测试数据库连接
    // sys.Get("/database", AppSys.GetDatabase)          // 获取数据库库列表
    // sys.Put("/database", AppSys.AddDatabase)          // 添加库
  }

  appNew.Run(iris.Addr(":1874"))
}

// Header设置跨域
func header(ctx context.Context) context.Context {
  ctx.Header("Access-Control-Allow-Origin", "*")
  ctx.Header("Access-Control-Allow-Credentials", "true")
  ctx.Header("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,Secret-Key,Hash");
  ctx.Header("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS")
  ctx.Header("Access-Control-Expose-Headers", "*")
  return ctx
}
