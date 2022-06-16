package utils

import (
	"bank-product-spike-system/global"
	"context"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// VerifyEmailCode 验证email验证码
func VerifyEmailCode(emailCode string, email string) bool {
	redis := global.REDIS
	ctx := context.Background()
	n, err := redis.Exists(ctx, email).Result()
	if err != nil {
		panic(err)
	}
	if n > 0 {
		res, err := redis.Get(ctx, email).Result()
		redis.Del(ctx, email)
		if err != nil {
			panic(err)
		} else {
			return res == emailCode
		}
	}
	return false
}

// SendEmail
// 发送邮件
func SendEmail(email string, code string) bool {
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("email.address"))
	m.SetHeader("To", email)
	m.SetAddressHeader("Aa", viper.GetString("email.address"), "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello you register code is <b>"+code+"</b>!")
	d := gomail.NewDialer(viper.GetString("email.host"), viper.GetInt("email.port"), viper.GetString("email.address"), viper.GetString("email.password"))
	if err := d.DialAndSend(m); err != nil {
		return false
	}
	return true
}
