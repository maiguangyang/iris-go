package sys

import (
  "github.com/kataras/iris/context"

  Config "../../config"
  Utils "../../utils"
  Database "../../database"
  Rsa "../../rsa_key"
)

// 获取配置
func GetConfig(ctx context.Context) {
  // 用户base64公key
  pubBase64 := ctx.URLParam("pktb")

  if pubBase64 == "" {
    ctx.JSON(Utils.NewResData(400, "键不能为空", ctx))
    return
  }

  // 加密数据
  content, err := Rsa.EncryptJosn(Config.ConfigJson, pubBase64)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }
  ctx.JSON(Utils.NewResData(200, content, ctx))
}

// 修改配置
func EditConfig(ctx context.Context) {
  var cData Utils.Map
  ctx.ReadJSON(&cData)

  if cData["content"] == nil || cData["content"].(string) == "" {
    ctx.JSON(Utils.NewResData(400, "提交数据不能为空", ctx))
    return
  }

  // 解密数据
  c, err := Rsa.DecryptJson(cData["content"].(string))
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  // 验证数据
  if r := ruleConfig(c, ctx); r != nil {
    ctx.JSON(r)
    return
  }

  if c["node_env"].(float64) != 1 {
    // 开发环境配置
    Config.ConfigJson.DriverName      = c["driver_name"].(string)
    Config.ConfigJson.DevDataUser     = c["data_user"].(string)
    Config.ConfigJson.DevDataPassword = c["data_password"].(string)
    Config.ConfigJson.DevDataIp       = c["data_ip"].(string)
    Config.ConfigJson.DevDataPort     = c["data_port"].(string)
    Config.ConfigJson.Database        = c["database"].(string)
    Config.ConfigJson.Charset         = c["charset"].(string)
  } else {
    // 正式环境配置
    Config.ConfigJson.DriverName      = c["driver_name"].(string)
    Config.ConfigJson.DataUser        = c["data_user"].(string)
    Config.ConfigJson.DataPassword    = c["data_password"].(string)
    Config.ConfigJson.DataIp          = c["data_ip"].(string)
    Config.ConfigJson.DataPort        = c["data_port"].(string)
    Config.ConfigJson.Database        = c["database"].(string)
    Config.ConfigJson.Charset         = c["charset"].(string)
  }

  // 写入配置
  err = Config.ConfigJson.Write()
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  // 重新连数据库
  err = Database.OpenSql()
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  ctx.JSON(Utils.NewResData(200, "修改成功", ctx))
}

// 连接测试
func TestOpen(ctx context.Context) {
  var cData Utils.Map
  ctx.ReadJSON(&cData)

  if cData["content"] == nil || cData["content"].(string) == "" {
    ctx.JSON(Utils.NewResData(400, "提交数据不能为空", ctx))
    return
  }

  // 解密数据
  c, err := Rsa.DecryptJson(cData["content"].(string))
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }
  // 验证数据
  if r := ruleConfig(c, ctx); r != nil {
    ctx.JSON(r)
    return
  }

  // 测试连接数据库
  err = Database.TestOpen(
    c["driver_name"].(string),
    c["data_user"].(string) + ":" +
    c["data_password"].(string) + "@tcp(" +
    c["data_ip"].(string) + ":" +
    c["data_port"].(string) + ")/" +
    c["database"].(string) + "?charset=" +
    c["charset"].(string))
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  ctx.JSON(Utils.NewResData(200, "测试成功", ctx))
}

// 验证数据库配置参数
func ruleConfig(d Utils.Map, ctx context.Context) interface{} {
  rules := Utils.Rules{
    "node_env": {
      { "required": true, "msg": "环境值不能为空" },
      { "type": "float64", "msg": "环境值类型为数字" },
      { "rgx": "node_env" },
    },
    "driver_name": {
      { "required": true, "msg": "数据库驱动名值不能为空" },
      { "type": "string", "msg": "数据库驱动名值类型为字符串" },
    },
    "data_user": {
      { "required": true, "msg": "账号值不能为空" },
      { "type": "string", "msg": "账号值类型为字符串" },
    },
    "data_password": {
      { "required": true, "msg": "密码值不能为空" },
      { "type": "string", "msg": "密码值类型为字符串" },
    },
    "data_ip": {
      { "required": true, "msg": "IP值不能为空" },
      { "rgx": "ip" },
    },
    "data_port": {
      { "required": true, "msg": "端口值不能为空" },
      { "type": "string", "msg": "端口值类型为字符串" },
    },
    "database": {
      { "required": true, "msg": "库名值不能为空" },
      { "type": "string", "msg": "库名值类型为字符串" },
    },
    "charset": {
      { "required": true, "msg": "字符集值不能为空" },
      { "type": "string", "msg": "字符集值类型为字符串" },
    },
  }

  errMsgs := rules.Validate(d)
  if errMsgs != nil {
    return Utils.NewResData(400, errMsgs, ctx)
  }

  return nil
}
