package utils

import (
  "github.com/kataras/iris/context"
  Rsa "../rsa_key"
)

func NewResData(code int, data interface{}, ctx context.Context) Map {
  var msg interface{}
  if code == 200 {
    msg = "success"
  } else {
    msg = "error"

    switch data.(type) {
    case string:
      msg  = data
      data = ""
    case Map:
      d := data.(Map)
      if d["msg"] != nil {
        msg = d["msg"]
      }
    }
  }

  resData, _ := Rsa.EncryptJosn(data, ctx.GetHeader("Secret-Key"))

  return Map{
    "code"        : code,
    "data"        : resData,
    "msg"         : msg,
    "status_code" : 200,
  }
}