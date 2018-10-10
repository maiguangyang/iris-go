package public

import(
  "fmt"
  "crypto/md5"
  "encoding/base64"
)

// 密码加密
func EncryptPassword(text string) string {
  password := EncryptMd5(text)
  password = base64.StdEncoding.EncodeToString([]byte(password))
  password = fmt.Sprintf("%x", md5.Sum([]byte(password)));
  return password
}

// MD5验证
func DecryptPassword(text, md5Text string) bool {
  password := EncryptMd5(text)
  password = base64.StdEncoding.EncodeToString([]byte(password))
  password = EncryptMd5(password);
  return password == md5Text
}

// MD5加密
func EncryptMd5(text string) string {
  md5Data := fmt.Sprintf("%x", md5.Sum([]byte(text)));
  return md5Data
}