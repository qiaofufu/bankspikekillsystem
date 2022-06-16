package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 将string进行md5加密
func MD5(str string) string {
	MD5 := md5.New()
	MD5.Write([]byte(str))
	return hex.EncodeToString(MD5.Sum(nil))
}
