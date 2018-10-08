package utils

var (
  FILTER_DATA = map[string]map[int]string{
    "login_type": {
      0: "系统配置登录",
      1: "后台管理登录",
      2: "用户登录",
    },
  }
)

func Filter(name string, value int) string {
  data := FILTER_DATA[name]
  if data == nil {
    return ""
  }

  return data[value]
}
