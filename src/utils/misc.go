package utils

import (
  "reflect"
  "crypto/md5"
  "encoding/hex"
)

// ########################################### UTILITY FUNCTIONS
func StringInSlice(a string, list []string) bool {
  for _, b := range list {
    if b == a {
        return true
    }
  }
  return false
}

// panic if s is not a slice
func ReverseSlice(s interface{}) {
  size := reflect.ValueOf(s).Len()
  swap := reflect.Swapper(s)
  for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
    swap(i, j)
  }
}

func GetMD5Hash(text string) string {
   hash := md5.Sum([]byte(text))
   return hex.EncodeToString(hash[:])
}
