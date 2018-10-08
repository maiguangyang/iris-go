package config

import (
  "os"
  "io/ioutil"
  "encoding/json"
)

const (
  configFilePath = "./config.json"
)

// 系统配置
type CJStruct struct {
  // 管理账号
  User                  string    `json:"user"`
  // 管理密码
  Password              string    `json:"password"`
  // 数据库驱动名，如：mysql
  DriverName            string    `json:"driver_name"`
  // 数据库账号
  DataUser              string    `json:"data_user"`
  // 数据库密码
  DataPassword          string    `json:"data_password"`
  // 数据库ip地址
  DataIp                string    `json:"data_ip"`
  // 数据库端口
  DataPort              string    `json:"data_port"`
  // 以下4个测试环境数据库
  DevDataUser           string    `json:"dev_data_user"`
  DevDataPassword       string    `json:"dev_data_password"`
  DevDataIp             string    `json:"dev_data_ip"`
  DevDataPort           string    `json:"dev_data_port"`
  // 数据库库名
  Database              string    `json:"database"`
  // 数据库编码
  Charset               string    `json:"charset"`
}

var ConfigJson CJStruct

// 读取配置
func (c *CJStruct) Read() error {
  _, err := os.Stat(configFilePath)
  if err != nil {
    _, err = os.Create(configFilePath)
    if err != nil {
      return err
    }
  }
  data, err := ioutil.ReadFile(configFilePath)
  if err != nil {
    return err
  }

  json.Unmarshal(data, &c)

  // 以后再做选择功能
  c.Charset = "utf8mb4"

  return nil
}

// 写入配置
func (c CJStruct) Write() error {
  data, err := json.MarshalIndent(c, "", "    ")
  if err != nil {
    return err
  }

  return ioutil.WriteFile(configFilePath, data, 0755)
}
