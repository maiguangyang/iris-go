package utils

import(
  "strings"

)
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

// 查找数组并返回下标
func IndexOf(str []interface{}, data interface{}) int {
  for k, v := range str{
    if v == data {
      return k
    }
  }

  return - 1
}


// stringToArray
func StrToArr(data, split string) []string {
  if IsEmpty(data) {
    return nil
  }
  return strings.Split(data, split)
}
