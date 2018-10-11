package utils

import (
  "strings"
  "github.com/kataras/iris/context"
  Public "../public"
)

func NewResData(code int, data interface{}, ctx context.Context) context.Map {
  var msg interface{}
  var resData interface{}
  var err error

  if code == 200 {
    msg = "success"
  } else {
    msg = "error"
  }

  if Public.NODE_ENV && strings.ToUpper(ctx.Method()) == "GET" {
    secretKey := ctx.GetHeader("Secret-Key")
    headHash  := ctx.GetHeader("Hash")

    if secretKey == "" || headHash == "" {
      msg     = "error"
      resData = "非法数据请求"
    } else {
      hash := Public.CheckHash(secretKey)
      if headHash != hash {
        msg     = "error"
        resData = "非法数据请求"
      } else {
        resData, err = Public.EncryptJosn(data, secretKey)
        if err != nil {
          msg     = "error"
          resData = "非法数据请求"
        }
      }
    }
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