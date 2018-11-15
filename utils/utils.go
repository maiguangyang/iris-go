package utils

import(
  // "fmt"
  "reflect"
  "github.com/kataras/iris/context"
  // "errors"
  Public "../public"
)

func SetField(obj interface{}, name string, value interface{}) error {
  structValue := reflect.ValueOf(obj).Elem()
  structFieldValue := structValue.FieldByName(name)


  // if !structFieldValue.IsValid() {
  //   return fmt.Errorf("No such field: %s in obj", name)
  // }

  // if !structFieldValue.CanSet() {
  //   return fmt.Errorf("Cannot set %s field value", name)
  // }
  if structFieldValue.IsValid() && structFieldValue.CanSet() {
    structFieldType := structFieldValue.Type()
    val := reflect.ValueOf(value)

    if structFieldType == val.Type() {
      structFieldValue.Set(val)

      // return errors.New("Provided value type didn't match obj field type")
      return nil
    }

    return nil

  }

  return nil
}

// map 映射 struct
func FillStruct(s interface{}, m map[string]interface{}) error {
  for k, v := range m {
    if reflect.TypeOf(v).String() == "float64" {
      v = int64(v.(float64))
    }

    err := SetField(s, CamelString(k), v)
    if err != nil {
      return err
    }
  }
  return nil
}

// 转驼峰命名
func CamelString(s string) string {
  data := make([]byte, 0, len(s))
  j := false
  k := false
  num := len(s) - 1

  for i := 0; i <= num; i++ {
    d := s[i]
    if k == false && d >= 'A' && d <= 'Z' {
      k = true
    }
    if d >= 'a' && d <= 'z' && (j || k == false) {
      d = d - 32
      j = false
      k = true
    }

    if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
      j = true
      continue
    }

    data = append(data, d)
  }
  return string(data[:])
}


// 返回环境数据
func ResNodeEnvData(table interface{}, ctx context.Context) error {
  // 线上环境
  if Public.NODE_ENV {
    decData, err := Public.DecryptReqData(ctx)

    if err != nil {
      return err
    }

    reqData := decData.(map[string]interface{})

    // map 映射 struct
    err = FillStruct(table, reqData)
    if err != nil {
      return err
    }

    return nil

  }
  ctx.ReadJSON(table)
  return nil
}



