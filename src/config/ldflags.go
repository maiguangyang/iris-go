// 在编译时用ldflags设置变量的值
package config

type ldflags struct {
  NODE_ENV string
}

var Ldflags ldflags

const (
  NODE_DEV = "dev"
  NODE_MASTER = "master"
)

func IsNodeDev () bool {
  return Ldflags.NODE_ENV != NODE_MASTER
}
