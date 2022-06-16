package captcha

import (
	"github.com/mojocn/base64Captcha"
)

var store = RedisStore{}

func GenerateCaptcha() (id string, base64 string, err error) {
	driver := base64Captcha.DefaultDriverDigit

	captcha := base64Captcha.NewCaptcha(driver, store)
	return captcha.Generate()
}

func VerifyCaptcha(id string, value string) bool {
	return store.Verify(id, value, true)
}
