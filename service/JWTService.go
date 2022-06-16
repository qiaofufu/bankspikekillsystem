package service

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type JWTService struct{}

var Secret = []byte("0C018992A72AAD2FDA3F3FE2D5D9E598")

var SignatureAlgorithm = jwa.HS256

var issuer = `qiaowiwi`

// VerifyToken 验证Token是否有效
func (receiver JWTService) VerifyToken(tokenStr string) error {
	token, err := receiver.ParseToken(tokenStr)
	if err != nil {
		return errors.New("token 验证内部错误")
	}
	if len(token.Audience()) > 0 && token.Audience()[0] != "User" {
		return errors.New("非法接收者")
	}
	if token.Expiration().Before(time.Now()) {
		return errors.New("登陆信息已过期，请重新登陆")
	}
	return nil
}

// ParseToken 解析TokenStr并获取其中的对象
func (receiver JWTService) ParseToken(tokenStr string) (token jwt.Token, err error) {
	token, err = jwt.ParseString(tokenStr, jwt.WithValidate(true), jwt.WithVerify(SignatureAlgorithm, Secret))
	return
}

// GenerateToken 签发用户 token
func (receiver JWTService) GenerateToken(payload map[string]interface{}) (tokenStr []byte, err error) {
	token := jwt.New()
	token.Set(jwt.IssuerKey, issuer)
	token.Set(jwt.AudienceKey, "User")
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().AddDate(0, 0, 15).Unix())

	for key, value := range payload {
		token.Set(key, value)
	}

	tokenStr, err = jwt.Sign(token, SignatureAlgorithm, Secret)
	return
}
