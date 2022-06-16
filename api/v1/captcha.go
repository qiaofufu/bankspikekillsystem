package v1

import (
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CaptchaAPI struct{}

// Captcha
// @Summary 获取验证码
// @Tags 验证
// @Accept json
// @Product json
// @Success 200 {object} response.DTO{data=response.CaptchaResponse} "获取成功"
// @Router /base/captcha [get]
func (receiver *CaptchaAPI) Captcha(ctx *gin.Context) {
	captcha, err := CaptchaService.Captcha()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取验证信息失败,err[%v]", err), ctx)
		return
	}
	response.Success(captcha, "获取成功", ctx)
}

// CaptchaPhone
// @Summary 获取手机验证码
// @Tags 验证
// @Accept json
// @Product json
// @Param data body request.ShotMessageCode true "传入参数"
// @Success 200 {object} response.DTO{} "获取成功"
// @Router /base/captchaPhone [post]
func (receiver CaptchaAPI) CaptchaPhone(ctx *gin.Context) {
	var phoneCaptcha request.ShotMessageCode
	if err := ctx.ShouldBindJSON(&phoneCaptcha); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	err := CaptchaService.CaptchaPhone(phoneCaptcha.Number)
	if err != nil {
		response.FailWithMessage("获取失败!"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("获取成功", ctx)
}
