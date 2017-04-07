package admin

import (
	"gopkg.in/kataras/iris.v6"
)

func Login(ctx *iris.Context) {
	println(ctx.RequestHeader("Cookie"))
	ctx.Writef("test：%s", ctx.Header())
}
