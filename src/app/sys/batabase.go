package sys

import (
  "regexp"
  "github.com/kataras/iris/context"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  Utils "../../utils"
  Database "../../database"
  Rsa "../../rsa_key"
)

// 获取数据库库列表
func GetDatabase(ctx context.Context) {
  // 用户base64公key
  pubBase64 := ctx.URLParam("pktb")
  if pubBase64 == "" {
    ctx.JSON(Utils.NewResData(400, "键不能为空", ctx))
    return
  }

  // 连接数据库参数
  content := ctx.URLParam("content")
  if content == "" {
    ctx.JSON(Utils.NewResData(400, "数据不能为空", ctx))
    return
  }

  // 解密
  cData, err := Rsa.DecryptJson(content)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  // 验证
  if r := ruleDatabase(cData, nil, ctx); r != nil {
    ctx.JSON(r)
    return
  }

  // 连接数据库
  db, err := sql.Open(
    cData["driver_name"].(string),
    cData["data_user"].(string) + ":" +
    cData["data_password"].(string) + "@tcp(" +
    cData["data_ip"].(string) + ":" +
    cData["data_port"].(string) + ")/" +
    "?charset=" + cData["charset"].(string))

  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }
  defer db.Close()

  // 获取库表
  sql := "SHOW DATABASES;"
  rows, err := db.Query(sql)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }
  defer rows.Close()

  list, err := Database.RowsScan(rows)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  var rList []string
  // 过虑自带库与其它软件库
  dbrgx := regexp.MustCompile("^(information_schema|mysql|performance_schema|test|walle|confluence_data|sys)$")
  for _, row := range list {
    database := row["Database"].(string)
    if dbrgx.MatchString(database) == false {
      rList = append(rList,database)
    }
  }

  // 加密数据
  content, err = Rsa.EncryptJosn(rList, pubBase64)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  ctx.JSON(Utils.NewResData(200, content, ctx))
}

// 添加库
func AddDatabase(ctx context.Context) {
  // 连接数据库参数
  type data struct {
    Content     string
  }
  var cData data
  ctx.ReadJSON(&cData)
  if cData.Content == "" {
    ctx.JSON(Utils.NewResData(400, "提交数据不能为空", ctx))
    return
  }

  // 解密
  c, err := Rsa.DecryptJson(cData.Content)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  // 验证
  rules := Utils.Rules{
    "database": {
      { "required": true, "msg": "库名值不能为空" },
      { "type": "string", "rgx_s": "^\\w+$", "min": 1, "max": 15, "msg": "库名值为1~15个字符（只能包括数字、英文字母和_）" },
    },
  }
  if r := ruleDatabase(c, rules, ctx); r != nil {
    ctx.JSON(r)
    return
  }

  // 连接数据库
  db, err := sql.Open(
    c["driver_name"].(string),
    c["data_user"].(string) + ":" +
    c["data_password"].(string) + "@tcp(" +
    c["data_ip"].(string) + ":" +
    c["data_port"].(string) + ")/" +
    "?charset=" + c["charset"].(string))

  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }
  defer db.Close()

  // 添加库
  sql := "CREATE DATABASE " + c["database"].(string)
  _, err = db.Exec(sql)
  if err != nil {
    ctx.JSON(Utils.NewResData(400, err.Error(), ctx))
    return
  }

  ctx.JSON(Utils.NewResData(200, "添加成功", ctx))
}

// 验证数据库配置参数
func ruleDatabase(d Utils.Map, rs Utils.Rules, ctx context.Context) interface{} {
  rules := Utils.Rules{
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
    "charset": {
      { "required": true, "msg": "字符集值不能为空" },
      { "type": "string", "msg": "字符集值类型为字符串" },
    },
  }

  for k, v := range rs {
    rules[k] = v
  }

  errMsgs := rules.Validate(d)
  if errMsgs != nil {
    return Utils.NewResData(400, errMsgs, ctx)
  }

  return nil
}