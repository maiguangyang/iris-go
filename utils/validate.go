package utils

import (
  "fmt"
  "reflect"
  "regexp"
  "github.com/kataras/iris/context"
)

// type Rules map[string][]map[string]interface{}

type Rules map[string]map[string]interface{}

func StructToMap(obj interface{}) map[string]interface{} {
  t := reflect.TypeOf(obj)
  v := reflect.ValueOf(obj)

  var data = make(map[string]interface{})
  for i := 0; i < t.NumField(); i++ {
    data[t.Field(i).Name] = v.Field(i).Interface()
  }
  return data
}

func (rs Rules) Validate(d context.Map) interface{} {

  var errMsgs []context.Map

  for f, v := range rs {
    if v["required"] != nil && v["required"].(bool) == true && IsEmpty(d[f]) == true {
      errMsgs = append(errMsgs, context.Map{
        "value": d[f],
        "msg": f + "不能为空",
      })
    } else {
      typ := reflect.TypeOf(d[f]).String()

      if v["type"] != nil {
        if v["type"].(string) == typ {
          // string 、int 长度大小判断
          if v["min"] != nil {
            if typ == "string" {
              if len([]rune(d[f].(string))) < v["min"].(int) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE  // 小于最小值，下面的判断跳过
              }
            } else if typ == "float64" {
              if d[f].(float64) < float64(v["min"].(int)) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE  // 小于最小值，下面的判断跳过
              }
            } else if typ == "[]interface {}" {
              if len(d[f].([]interface{})) < v["min"].(int) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE
              }
            } else if typ == "[]context.Map" {
              if len(d[f].([]context.Map)) < v["min"].(int) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE  // 小于最小值，下面的判断跳过
              }
            }
          }
          if v["max"] != nil {
            if typ == "string" {
              if len([]rune(d[f].(string))) > v["max"].(int) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE
              }
            } else if typ == "float64" {
              if d[f].(float64) > float64(v["max"].(int)) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE
              }
            } else if typ == "[]interface {}" {
              if len(d[f].([]interface{})) > v["max"].(int) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE
              }
            } else if typ == "[]context.Map" {
              if len(d[f].([]context.Map)) > v["max"].(int) {
                errMsgs = append(errMsgs, context.Map{
                  "value": d[f],
                  "msg": v["msg"],
                })
                goto ENDTYPE  // 小于最小值，下面的判断跳过
              }
            }
          }
        } else {
          errMsgs = append(errMsgs, context.Map{
            "value": d[f],
            "msg": v["msg"],
          })
          goto ENDTYPE
        }
      }

      if v["len"] != nil {
        if (typ != "string" && typ != "[]interface {}" && typ != "[]context.Map") ||
        (typ == "string" && v["len"].(int) != len([]rune(d[f].(string)))) ||
        (typ == "[]interface {}" && v["len"].(int) != len(d[f].([]interface{}))) ||
        (typ == "[]context.Map" && v["len"].(int) != len(d[f].([]context.Map))) {
          errMsgs = append(errMsgs, context.Map{
            "value": d[f],
            "msg": v["msg"],
          })
          goto ENDTYPE  // 不等于长度，下面的判断跳过
        }
      }


      // 正则类型
      if v["rgx"] != nil {
        rl := Rule[v["rgx"].(string)]
        if rl != nil {
          if regexp.MustCompile(rl["rgx"].(string)).MatchString(fmt.Sprint(d[f])) == false {
            str := rl["msg"].(string)

            if rl["bool"] == true {
              str = v["rgx"].(string) + rl["msg"].(string)
            }

            errMsgs = append(errMsgs, context.Map{
              "value": d[f],
              "msg": str,
            })
            goto ENDTYPE
          }
        }
      }

      // 正则语句
      if v["rgx_s"] != nil {
        var str string
        if typ == "float64" {
          str = fmt.Sprint(int64(d[f].(float64)))
        } else {
          str = fmt.Sprint(d[f])
        }
        if regexp.MustCompile(v["rgx_s"].(string)).MatchString(str) == false {
          errMsgs = append(errMsgs, context.Map{
            "value": d[f],
            "msg": v["msg"],
          })
        }
      }

      ENDTYPE:
    }
  }

  if len(errMsgs) > 0 {
    return errMsgs
  }
  return nil
}