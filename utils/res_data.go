package utils

import (
  "github.com/kataras/iris/context"
  Public "../public"
)

func NewResData(code int, data interface{}, ctx context.Context) context.Map {
  var msg interface{}
  resData := ""

  if code == 200 {
    msg = "success"
    switch data.(type) {
    case context.Map:
      resData, _ = Public.EncryptJosn(data, ctx.GetHeader("Secret-Key"))
    }
  } else {
    msg = "error"
    switch data.(type) {
    case string:
      msg = data
    case context.Map:
      d := data.(context.Map)
      if d["msg"] != nil {
        msg = d["msg"]
      }
      resData, _ = Public.EncryptJosn(d, ctx.GetHeader("Secret-Key"))
    }
  }



  return context.Map{
    "code"        : code,
    "data"        : resData,
    "msg"         : msg,
    "status_code" : 200,
  }
}