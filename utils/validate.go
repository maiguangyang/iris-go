package utils

import (
  "fmt"
  "reflect"
  "regexp"
  "github.com/kataras/iris/context"
)

type Rules map[string][]map[string]interface{}

var rule = map[string]map[string]string{
  "url": {
    "rgx": "^https?:\\/\\/.+$",
    "msg": "网址格式不正确",
  },
  "identity": {
    "rgx": "^\\d{6}(18|19|20)?\\d{2}(0[1-9]|1[012])(0[1-9]|[12]\\d|3[01])\\d{3}(\\d|X|x)$",
    "msg": "身份证号码格式不正确",
  },
  "phone": {
    "rgx": "^(((13[0-9]{1})|(15[0-9]{1})|(18[0-9]{1})|(17[0-9]{1})|(14[0-9]{1}))+\\d{8})$",
    "msg": "必须是11位手机号码",
  },
}


func IsEmpty(v interface{}) bool {
  if v == nil {
    return true
  }

  switch v.(type) {
  case string:
    if v == "" {
      return true
    }
  default:
    return false
  }
  return false
}

func (rs Rules) Validate(d context.Map) interface{} {
  var errMsgs []context.Map

  for f, v := range rs {
    for _, r := range v {
      if IsEmpty(d[f]) == true {
        // 不能为空
        if r["required"] != nil && r["required"].(bool) == true {
          errMsgs = append(errMsgs, context.Map{
            "value": d[f],
            "msg": r["msg"],
          })
        }
      } else {
        typ := reflect.TypeOf(d[f]).String()

        if r["type"] != nil {
          if r["type"].(string) == typ {
            // string 、int 长度大小判断
            if r["min"] != nil {
              if typ == "string" {
                if len([]rune(d[f].(string))) < r["min"].(int) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE  // 小于最小值，下面的判断跳过
                }
              } else if typ == "float64" {
                if d[f].(float64) < float64(r["min"].(int)) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE  // 小于最小值，下面的判断跳过
                }
              } else if typ == "[]interface {}" {
                if len(d[f].([]interface{})) < r["min"].(int) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE
                }
              } else if typ == "[]context.Map" {
                if len(d[f].([]context.Map)) < r["min"].(int) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE  // 小于最小值，下面的判断跳过
                }
              }
            }
            if r["max"] != nil {
              if typ == "string" {
                if len([]rune(d[f].(string))) > r["max"].(int) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE
                }
              } else if typ == "float64" {
                if d[f].(float64) > float64(r["max"].(int)) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE
                }
              } else if typ == "[]interface {}" {
                if len(d[f].([]interface{})) > r["max"].(int) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE
                }
              } else if typ == "[]context.Map" {
                if len(d[f].([]context.Map)) > r["max"].(int) {
                  errMsgs = append(errMsgs, context.Map{
                    "value": d[f],
                    "msg": r["msg"],
                  })
                  goto ENDTYPE  // 小于最小值，下面的判断跳过
                }
              }
            }
          } else {
            errMsgs = append(errMsgs, context.Map{
              "value": d[f],
              "msg": r["msg"],
            })
            goto ENDTYPE
          }
        }

        if r["len"] != nil {
          if (typ != "string" && typ != "[]interface {}" && typ != "[]context.Map") ||
          (typ == "string" && r["len"].(int) != len([]rune(d[f].(string)))) ||
          (typ == "[]interface {}" && r["len"].(int) != len(d[f].([]interface{}))) ||
          (typ == "[]context.Map" && r["len"].(int) != len(d[f].([]context.Map))) {
            errMsgs = append(errMsgs, context.Map{
              "value": d[f],
              "msg": r["msg"],
            })
            goto ENDTYPE  // 不等于长度，下面的判断跳过
          }
        }


        // 正则类型
        if r["rgx"] != nil {
          rl := rule[r["rgx"].(string)]
          if rl != nil {
            if regexp.MustCompile(rl["rgx"]).MatchString(fmt.Sprint(d[f])) == false {
              errMsgs = append(errMsgs, context.Map{
                "value": d[f],
                "msg": rl["msg"],
              })
              goto ENDTYPE
            }
          }
        }

        // 正则语句
        if r["rgx_s"] != nil {
          var str string
          if typ == "float64" {
            str = fmt.Sprint(int64(d[f].(float64)))
          } else {
            str = fmt.Sprint(d[f])
          }
          if regexp.MustCompile(r["rgx_s"].(string)).MatchString(str) == false {
            errMsgs = append(errMsgs, context.Map{
              "value": d[f],
              "msg": r["msg"],
            })
          }
        }

        ENDTYPE:
      }
    }
  }

  if len(errMsgs) > 0 {
    return errMsgs
  }
  return nil
}