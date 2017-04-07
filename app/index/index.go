package index

import (
	"gopkg.in/kataras/iris.v6"
)

func Hello(ctx *iris.Context) {
	ctx.Writef("testYang")
}
