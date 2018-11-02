package utils

import (
  // "fmt"
  "strconv"
)

// String转Int
func StrToInt(v string) int {
  s, _ := strconv.Atoi(v)
  return s
}

// String转Int64
func StrToInt64(v string) int64 {
  s, _ := strconv.ParseInt(v, 10, 64)
  return s
}


// Int转String
func IntToStr(v int) string {
  s := strconv.Itoa(v)
  return s
}

// Int64转String
func Int64ToStr(v int64) string {
  s := strconv.FormatInt(v, 10)
  return s
}