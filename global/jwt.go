package global

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

var Secret = []byte("0C018992A72AAD2FDA3F3FE2D5D9E598")
var SignatureAlgorithm = jwa.HS256
var issuer = `qiaowiwi`

// GenerateUserToken 生成user token
func GenerateUserToken(payload map[string]interface{}) (string, error) {
	token := jwt.New()
	token.Set(jwt.IssuerKey, issuer)
	token.Set(jwt.AudienceKey, `User`)
	token.Set(jwt.IssuedAtKey, time.Now().Unix())
	token.Set(jwt.ExpirationKey, time.Now().AddDate(0, 0, 15).Unix())

	for key, value := range payload {
		token.Set(key, value)
	}

	tokenString, err := jwt.Sign(token, SignatureAlgorithm, Secret)
	return string(tokenString), err
}

// VerifyToken 验证Token在有效期内是否正确
func VerifyToken(tokenString string, algorithm jwa.SignatureAlgorithm, secret []byte) bool {
	token, err := ParseToken(tokenString, algorithm, secret)
	if err != nil {
		return false
	}
	if token.Issuer() != issuer {
		return false
	}
	if len(token.Audience()) > 0 && token.Audience()[0] != "User" {
		return false
	}
	if token.Expiration().Before(time.Now()) {
		return false
	}
	return true
}

// ParseToken ParseToken 解析Token并获取其中的对象
func ParseToken(tokenString string, algorithm jwa.SignatureAlgorithm, secret []byte) (jwt.Token, error) {
	return jwt.ParseString(tokenString, jwt.WithValidate(true), jwt.WithVerify(algorithm, secret))
}
