package public

// 判断系统环境
var (
  NODE_ENV bool
)

func IsNodeEnv(str string) bool {
  NODE_ENV = str == "master"
  // NODE_ENV = true
  return NODE_ENV
}