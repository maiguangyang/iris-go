package public

import(
  "fmt"
  "strings"
  "crypto/md5"
  "encoding/base64"
  "github.com/kataras/iris/context"
  "errors"
)

// 密码加密
func EncryptPassword(text string) string {
  password := EncryptMd5(text)
  password = base64.StdEncoding.EncodeToString([]byte(password))
  password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
  return password
}

// MD5验证
func DecryptPassword(text, md5Text string) bool {
  password := EncryptMd5(text)
  password = base64.StdEncoding.EncodeToString([]byte(password))
  password = EncryptMd5(password)
  return password == md5Text
}

// MD5加密
func EncryptMd5(text string) string {
  md5Data := fmt.Sprintf("%x", md5.Sum([]byte(text)))
  return md5Data
}

// 组装返回header hash加密后的值
func CheckHash(text string) string {
  data := fmt.Sprintf("%x", md5.Sum([]byte(text[8:108] + "EQUOYpl72tsjwzJnnY")))
  return data
}

// 解密前端传过来的数据
func DecryptReqData(ctx context.Context) (context.Map, error) {

  if NODE_ENV && strings.ToUpper(ctx.Method()) != "GET" {
    type ReqData struct {
      Content string `json:"content"`
    }

    var reqData ReqData
    ctx.ReadJSON(&reqData)
    content := reqData.Content

    secretKey := ctx.GetHeader("Secret-Key")
    headHash  := ctx.GetHeader("Hash")

    if secretKey == "" || headHash == "" {
      return context.Map{"content": "非法数据请求"}, errors.New("failed to parse PEM block containing the private key")
    } else {
      hash := CheckHash(secretKey)
      if headHash != hash {
        return context.Map{"content": "非法数据请求"}, errors.New("failed to parse PEM block containing the private key")
      } else {
        res, err := DecryptJson(content, secretKey)
        return res, err
      }
    }
  } else {
    var reqData context.Map
    ctx.ReadJSON(&reqData)
    return reqData, nil
  }
}