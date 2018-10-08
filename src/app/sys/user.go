package sys

import (
  "fmt"
  "encoding/json"
  "github.com/kataras/iris/context"

  Config "../../config"
  Utils "../../utils"
  Rsa "../../rsa_key"
  Database "../../database"
  Auth "../../authorization"
)

// 是否设置配置管理员
func HasUser(ctx context.Context) {
  if Config.ConfigJson.User == "" {
    ctx.JSON(Utils.NewResData(200, context.Map{ "has": false }, ctx))
    return
  }
  ctx.JSON(Utils.NewResData(200, context.Map{ "has": true }, ctx))
}


// 检测是否设置数据库
func CheckDataBase(ctx context.Context) {
  var data context.Map
  if Config.ConfigJson.DataUser == "" || Config.ConfigJson.DataPassword == "" || Config.ConfigJson.DataIp == "" || Config.ConfigJson.DataPort == "" || Config.ConfigJson.Database == "" {
    data = context.Map{ "has": false }
  } else {
    data = context.Map{ "has": true }
  }


  ctx.JSON(Utils.NewResData(200, data, ctx))
}

// 检测是否设置数据库
func CheckDataBasePost(ctx context.Context) {
  type data struct {
    Content string
  }

  var cData data
  ctx.ReadJSON(&cData)

  fmt.Println(cData.Content)

  ctx.JSON(Utils.NewResData(200, cData.Content, ctx))
}



// 登录
func Login(ctx context.Context) {
  ip := ctx.RemoteAddr()
  if Database.DB != nil && Database.IsRestrictLogin(0, ip) == true {
    ctx.JSON(Utils.NewResData(400, "短时间内登录失败次数过多，已限制登录", ctx))
    return
  }

  type data struct {
    Content     string
  }

  var cData data
  ctx.ReadJSON(&cData)

  if cData.Content == "" {
    ctx.JSON(Utils.NewResData(400, "提交数据不能为空", ctx))
    return
  }

  if Database.DB != nil && Database.LoginLogHasContent(cData.Content) == true {
    ctx.JSON(Utils.NewResData(400, "不能重复提交数据", ctx))
    return
  }

  b, err := Rsa.Decrypt(cData.Content)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  type content struct {
    User         string     `json:"user"`
    Password     string     `json:"password"`
  }
  var c content
  err = json.Unmarshal(b, &c)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  if c.User != Config.ConfigJson.User || c.Password != Config.ConfigJson.Password {
    if Database.DB != nil { Database.LoginErrorLogAdd(0, ip) }
    ctx.JSON(Utils.NewResData(400, "账号或密码错误", ctx))
    return
  }

  token := Auth.SysSetToken(context.Map{ "user": c.User }, cData.Content, ip)
  ctx.JSON(Utils.NewResData(200, context.Map{
    "token": token,
  }, ctx))
}

// 修改密码
func ModifyPassword(ctx context.Context) {
  type data struct {
    Content     string
  }
  var cData data
  ctx.ReadJSON(&cData)
  if cData.Content == "" {
    ctx.JSON(Utils.NewResData(400, "提交数据不能为空", ctx))
    return
  }

  b, err := Rsa.Decrypt(cData.Content)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  type content struct {
    Password     string     `json:"password"`
  }
  var c content
  err = json.Unmarshal(b, &c)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  if c.Password == "" {
    ctx.JSON(Utils.NewResData(400, "密码不能为空", ctx))
    return
  }

  Config.ConfigJson.Password = c.Password
  err = Config.ConfigJson.Write()
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  ctx.JSON(Utils.NewResData(200, "修改成功", ctx))
}