package utils

import (
  "fmt"
  "reflect"
  "regexp"
)

type Rules map[string][]map[string]interface{}
type ErrMsgs []string

var rule = map[string]map[string]string{
  "ip": {
    "rgx": "^([1-9]|[1-9]\\d|1\\d\\d|2[0-4]\\d|25[0-5])(\\.(\\d|[1-9]\\d|1\\d\\d|2[0-4]\\d|25[0-5])){3}$",
    "msg": "ip格式错误",
  },
  "node_env": {
    "rgx": "^(0|1)$",
    "msg": "环境值为0或1（0为测试环境；1为正则环境）",
  },
  "phone": {
    "rgx": "^(((13[0-9]{1})|(15[0-9]{1})|(18[0-9]{1})|(17[0-9]{1})|(14[0-9]{1}))+\\d{8})$",
    "msg": "手机号码格式错误",
  },
}

func (rs Rules) Validate(d Map) ErrMsgs {
  var errMsgs ErrMsgs

  for f, v := range rs {
    for _, r := range v {
      if IsEmpty(d[f]) == true {
        // 不能为空
        if r["required"] != nil && r["required"].(bool) == true {
          errMsgs = append(errMsgs, r["msg"].(string))
        }
      } else {
        typ := reflect.TypeOf(d[f]).String()

        if r["type"] != nil {
          if r["type"].(string) == typ {
            // string 、int 长度大小判断
            if r["min"] != nil {
              if typ == "string" {
                if len([]rune(d[f].(string))) < r["min"].(int) {
                  errMsgs = append(errMsgs, r["msg"].(string))
                  goto ENDTYPE  // 小于最小值，下面的判断跳过
                }
              } else if typ == "float64" {
                if d[f].(float64) < r["min"].(float64) {
                  errMsgs = append(errMsgs, r["msg"].(string))
                  goto ENDTYPE  // 小于最小值，下面的判断跳过
                }
              }
            }
            if r["max"] != nil {
              if typ == "string" {
                if len([]rune(d[f].(string))) > r["max"].(int) {
                  errMsgs = append(errMsgs, r["msg"].(string))
                  goto ENDTYPE
                }
              } else if typ == "float64" {
                if d[f].(float64) > r["min"].(float64) {
                  errMsgs = append(errMsgs, r["msg"].(string))
                  goto ENDTYPE
                }
              }
            }
          } else {
            errMsgs = append(errMsgs, r["msg"].(string))
            goto ENDTYPE
          }
        }

        // 正则类型
        if r["rgx"] != nil {
          rl := rule[r["rgx"].(string)]
          if rl != nil {
            if regexp.MustCompile(rl["rgx"]).MatchString(fmt.Sprint(d[f])) == false {
              errMsgs = append(errMsgs, rl["msg"])
              goto ENDTYPE
            }
          }
        }

        // 正则语句
        if r["rgx_s"] != nil {
          if regexp.MustCompile(r["rgx_s"].(string)).MatchString(fmt.Sprint(d[f])) == false {
            errMsgs = append(errMsgs, r["msg"].(string))
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