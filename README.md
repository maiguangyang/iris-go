# 项目说明

## 本地环境安装

```
# 安装 Golang
$ http://www.golangtc.com/download
$ wget xxx
$ n 1.8.x 以上


# 安装 Package依赖包
$ gopkg.in/kataras/iris.v6
$ github.com/dgrijalva/jwt-go

```

## 代码编译与发布

```
# 编译生成代码
$ go build ./app/main.go

```

#### Nginx 配置文件导入（Todo以下内容待定）

```
# 生成vhost


# 导入配置文件
$ echo "include /path/to/project/vhosts/nginx.conf;" >> /path/to/nginx/nginx.conf

# 重启nginx
# Linux
$ sudo service nginx restart
# OSX
$ sudo brew services restart nginx

```

