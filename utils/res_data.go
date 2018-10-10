package utils

import (
  "github.com/kataras/iris/context"
  Public "../public"
)

func NewResData(code int, data interface{}, ctx context.Context) context.Map {
  var msg interface{}
  var resData interface{}

  if code == 200 {
    msg = "success"
  } else {
    msg = "error"
  }

  if Public.NODE_ENV {
    resData, _ = Public.EncryptJosn(data, ctx.GetHeader("Secret-Key"))
  } else {
    resData = data
  }

  return context.Map{
    "code"        : code,
    "data"        : resData,
    "msg"         : msg,
    "status_code" : 200,
  }
}