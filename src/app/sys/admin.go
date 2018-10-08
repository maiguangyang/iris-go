package sys

import (
  "github.com/kataras/iris/context"

  Utils "../../utils"
)

func AdminAdd(ctx context.Context) {
  var cData Utils.Map
  ctx.ReadJSON(&cData)

  rule := Utils.Rules{
    "phone": {
      { "required": true, "msg": "手机号码不能为空" },
      { "rgx": "phone" },
    },
    "password": {
      { "required": true, "msg": "密码不能为空" },
    },
    "name": {
      { "required": true, "msg": "姓名不能为空" },
    },
    "role": {
      { "required": true, "msg": "角色不能为空" },
      { "rgx_s": "^\\d+$", "msg": "角色只能为正整数"},
    },
    "state": {
      { "required": true, "msg": "状态不能为空" },
      { "rgx_s": "^\\d+$", "msg": "状态只能为正整数"},
    },
  }
  errMsgs := rule.Validate(cData)
  if errMsgs != nil {
    ctx.JSON(Utils.NewResData(400, errMsgs, ctx))
    return
  }
}