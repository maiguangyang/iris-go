package utils

var Rule = map[string]map[string]interface{}{
  "empty"    : { "rgx": "^\\S", "msg": "不能为空", "bool": true, },
  "int"      : { "rgx": "^[0-9]\\d*$", "msg": "必须是0-9的整数", },
  "code"     : { "rgx": "^([0-9]){6}$", "msg": "验证码必须是6位整数", },
  "url"      : { "rgx": "^https?:\\/\\/.+$", "msg": "网址格式不正确", },
  "email"    : { "rgx": "^([a-z0-9\\+\\_\\-]+)(\\.[a-z0-9\\+\\_\\-]+)*@([a-z0-9\\-]+\\.)+[a-z]{2,6}$", "msg": "邮箱格式不正确", },
  "identity" : { "rgx": "^\\d{6}(18|19|20)?\\d{2}(0[1-9]|1[012])(0[1-9]|[12]\\d|3[01])\\d{3}(\\d|X|x)$", "msg": "身份证号码格式不正确",
  },
  "phone"    : { "rgx": "^(((13[0-9]{1})|(15[0-9]{1})|(18[0-9]{1})|(17[0-9]{1})|(14[0-9]{1}))+\\d{8})$", "msg": "必须是11位手机号码",
  },
}