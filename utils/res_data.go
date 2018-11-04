package utils

import (
  // "fmt"
  "math"
  "strings"
  "github.com/kataras/iris/context"
  Public "../public"
)

func NewResData(code int, data interface{}, ctx context.Context) context.Map {
  var msg interface{}
  var resData interface{}
  var err error

  if code == 0 {
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

// 列表分页、总数量
func TotalData(list interface{}, page, total int64) context.Map {
  per_page     := 2
  total_page   := int64(math.Ceil(float64(total) / float64(per_page)))

  if page > total_page {
    list = []context.Map{}
  }

  pageData := context.Map{
    "total"        : total,
    "current_page" : page,
    "per_page"     : per_page,
    "total_page"   : total_page,
    "data"         : list,
  }

  return pageData
}