package public

import (
  "fmt"
  "time"
  "strconv"
  "crypto/md5"
  "github.com/kataras/iris/context"

  Ras "../../rsa_key"
  Utils "../../utils"
)

func GetRsaPubKey(ctx context.Context) {
  if Ras.PubPemEnc == "" {
    err := Ras.Gen()
    if err != nil {
      ctx.JSON(Utils.NewResData(400, Utils.Map{
        "msg": err.Error(),
      }, ctx))
      return
    }

    ctx.JSON(Utils.NewResData(400, Utils.Map{
      "msg": "服务器公钥生成出错",
    }, ctx))
    return
  }

  b := []byte(strconv.FormatInt(time.Now().Unix(), 10))
  ctx.JSON(Utils.NewResData(200, Utils.Map{
    "p": Ras.PubPemEnc,
    "h": fmt.Sprintf("%x", md5.Sum(b)),
  }, ctx))
}
