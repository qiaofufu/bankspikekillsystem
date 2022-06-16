package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models/response"
	"bank-product-spike-system/utils"
	"bank-product-spike-system/utils/captcha"
	"context"
	"errors"
	"time"
)

type CaptchaService struct{}

// Captcha
// 生成验证码
func (receiver CaptchaService) Captcha() (*response.CaptchaResponse, error) {
	id, base64, err := captcha.GenerateCaptcha()
	if err != nil {
		return nil, err
	}
	return &response.CaptchaResponse{
		CaptchaType:   "string",
		CaptchaBase64: base64,
		CaptchaID:     id,
	}, nil
}

// CaptchaPhone
// 生成手机验证码
func (receiver CaptchaService) CaptchaPhone(phone string) error {
	code := utils.RandomCode(6)
	err := utils.SendSMS(phone, code)
	ctx := context.Background()
	if err != nil {
		return errors.New("发送手机验证码失败！ " + err.Error())
	}
	global.REDIS.Set(ctx, "phone:"+phone, code, time.Second*300)
	return nil
}
