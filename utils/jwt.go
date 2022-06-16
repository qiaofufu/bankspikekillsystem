package utils

import (
	"bank-product-spike-system/global"
	"errors"
	"github.com/lestrrat-go/jwx/jwt"
)

func GetTokenPayload(tokenStr string) (jwt.Token, error) {
	return global.ParseToken(tokenStr, global.SignatureAlgorithm, global.Secret)
}

// GetUUIDByToken
// 从TOKEN中获取UUID
func GetUUIDByToken(tokenStr string) (interface{}, error) {
	token, err := GetTokenPayload(tokenStr)
	if err != nil {
		return "", err
	}
	uuid, ok := token.Get("uuid")
	if !ok {
		return "", errors.New("token中信息不全，请重新登录")
	}
	return uuid, nil
}
