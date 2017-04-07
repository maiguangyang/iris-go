package router

import (
	"encoding/json"
	"time"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/gorillamux"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"

	Admin "web/app/admin"
	Home "web/app/index"
	User "web/app/user"
)

func Routers() {
	app := iris.New()
	// app.Adapt(iris.DevLogger())
	app.Adapt(gorillamux.New())

	/**
	 * 定义错误路由
	 */
	app.OnError(iris.StatusNotFound, func(ctx *iris.Context) {

		jsonStr := `{
      "code": 404,
      "data": "404 not found",
      "msg" : "error"
    }`

		var responseData interface{}
		err := json.Unmarshal([]byte(jsonStr), &responseData)
		if err == nil {
			ctx.JSON(iris.StatusOK, responseData)
		}
	})

	/**
	 * sessions配置
	 */
	mySessions := sessions.New(sessions.Config{
		Cookie:                      "sessionId",
		Expires:                     time.Hour * 2,
		CookieLength:                64,
		DisableSubdomainPersistence: false,
	})

	app.Adapt(mySessions)

	/**
	 * 定义路由
	 */
	app.Get("/", Home.Hello)
	app.Get("/user", User.UserSay)

	admin := app.Party("/admin", Admin.Authorization)
	{
		admin.Get("/", Admin.Index)
		admin.Get("/profile", Admin.Profile)
		admin.Get("/set", Admin.SetSession)
		admin.Get("/get", Admin.GetSession)
		admin.Get("/login", Admin.Login)
	}

	app.Listen(":8080")
}
