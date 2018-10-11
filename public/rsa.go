package public

import (
  // "fmt"
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  "encoding/base64"
  "encoding/json"
  "encoding/hex"
  "crypto/cipher"
  "crypto/aes"
  "bytes"
  "errors"
  // Utils "../utils"
)

const (
  aesKey = "!IWCT$C6hSA5iVQ8"
)

func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
  pkey := PaddingLeft(key, '0', 16)
  block, err := aes.NewCipher(pkey) //选择加密算法
  if err != nil {
    return nil, err
  }

  blockModel := cipher.NewCBCDecrypter(block, pkey)
  plantText := make([]byte, len(ciphertext))
  blockModel.CryptBlocks(plantText, []byte(ciphertext))
  plantText = PKCS7UnPadding(plantText, block.BlockSize())
  return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {
  length := len(plantText)
  unpadding := int(plantText[length-1])
  return plantText[:(length - unpadding)]
}

func PaddingLeft(ori []byte, pad byte, length int) []byte {
  if len(ori) >= length {
    return ori[:length]
  }
  pads := bytes.Repeat([]byte{pad}, length-len(ori))
  return append(pads, ori...)
}


// 解密
func Decrypt(text, prvPemEnc string) ([]byte, error) {
  var ciphertext []byte

  // 先base64解密
  prvPem, err := base64.StdEncoding.DecodeString(prvPemEnc)

  // 再使用hex.DecodeString转换
  dexDecode, err := hex.DecodeString(string(prvPem))

  if err != nil {
    return ciphertext, err
  }

  // 最后解密对称加密的内容
  prvPem, err = AesDecrypt(dexDecode, []byte(aesKey))
  if err != nil {
    return ciphertext, err
  }

  // 公私钥加密
  block, _ := pem.Decode(prvPem)
  if block == nil {
    return ciphertext, errors.New("failed to parse PEM block containing the private key")
  }

  prvKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

  if err != nil {
    return ciphertext, err
  }

  limitLen := 172
  textLen := len(text)
  if textLen % limitLen != 0 {
    return ciphertext, errors.New("加密数据非法")
  }
  for i, j := 0, limitLen; i < textLen; i, j = i + limitLen, j + limitLen {
    decode, err := base64.StdEncoding.DecodeString(text[i:j])
    if err != nil {
      return ciphertext, err
    }
    c, err := rsa.DecryptPKCS1v15(rand.Reader, prvKey, decode)
    if err != nil {
      if err.Error() == "crypto/rsa: decryption error" {
        return ciphertext, errors.New("页面已过期，请刷新重试")
      }
      return ciphertext, err
    }
    var buffer bytes.Buffer
    buffer.Write(ciphertext)
    buffer.Write(c)
    ciphertext = buffer.Bytes()
  }
  return ciphertext, nil
}

// 解密后的数据为json
func DecryptJson(text, pubPemEnc string) (map[string]interface{}, error) {
  b, err := Decrypt(text, pubPemEnc)
  if err != nil {
    return map[string]interface{}{}, err
  }
  var cData map[string]interface{}
  err = json.Unmarshal(b, &cData)
  return cData, err
}

// 加密
func Encrypt(text, pubPemEnc string) (string, error) {
  var encrypted string
  // 先base64解密
  pubPem, err := base64.StdEncoding.DecodeString(pubPemEnc)

  // 再使用hex.DecodeString转换
  dexDecode, err := hex.DecodeString(string(pubPem))

  if err != nil {
    return encrypted, err
  }

  // 最后解密对称加密的内容
  pubPem, err = AesDecrypt(dexDecode, []byte(aesKey))
  if err != nil {
    return encrypted, err
  }

  // 公私钥加密
  block, _ := pem.Decode(pubPem)
  if block == nil {
    return encrypted, errors.New("failed to parse PEM block containing the public key")
  }

  pub, err := x509.ParsePKIXPublicKey(block.Bytes)
  if err != nil {
    return encrypted, err
  }

  limitLen := 117
  for i, j, l := 0, limitLen, len(text); i < l; i, j = i + limitLen, j + limitLen {
    if j > l {
      j = l
    }

    enc, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(text[i:j]))
    if err != nil {
      return encrypted, err
    }

    encrypted += base64.StdEncoding.EncodeToString(enc)
  }

  return encrypted, nil
}

// 加密json数据
func EncryptJosn(v interface{}, pubPemEnc string) (string, error) {
  // 配置信息转字符串
  text, err := json.Marshal(v)


  if err != nil {
    return "", err
  }
  // 加密数据
  data, err := Encrypt(string(text), pubPemEnc)
  return data, err
}
