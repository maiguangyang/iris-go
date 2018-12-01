package utils

import (
  // "fmt"
  "math"
  // "strings"
  "github.com/kataras/iris/context"
  Public "../public"
)

func NewResData(code int, data interface{}, ctx context.Context) context.Map {
  var msg interface{}
  var resData interface{}
  // var err error

  if code == 0 {
    msg = "success"
  } else {
    msg = "error"
  }

  resData = data
  // if Public.NODE_ENV && strings.ToUpper(ctx.Method()) == "GET" {
  if Public.NODE_ENV {
    secretKey := ctx.GetHeader("Secret-Key")
    headHash  := ctx.GetHeader("Hash")

    if secretKey == "" || headHash == "" {
      msg     = "error"
      resData = "非法数据请求"
    } else if hash := Public.CheckHash(secretKey); hash != headHash {
        msg     = "error"
        resData = "非法数据请求"

      // hash := Public.CheckHash(secretKey)
      // if headHash != hash {
      //   msg     = "error"
      //   resData = "非法数据请求1"
      // } else {
      //   resData = data
      // }

      // 暂时注释，前端解密卡顿
      // else {
      //   resData, err = Public.EncryptJosn(data, secretKey)
      //   if err != nil {
      //     msg     = "error"
      //     resData = "非法数据请求2"
      //   }
      // }
    }
  }
  // else {
  //   resData = data
  // }

  text := context.Map{
    "code"        : code,
    "data"        : resData,
    "msg"         : msg,
    "status_code" : 200,
  }

  return text

}

// 列表、当前页、总数量、每页数量
func TotalData(list interface{}, page, total int64, count int) context.Map {
  var per_page int64 = 20

  if count > 0 {
    per_page = int64(count)
  }

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