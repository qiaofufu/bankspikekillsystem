package system

import (
	v1 "bank-product-spike-system/api/v1"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (receiver BaseRouter) InitBaseRouter(group *gin.RouterGroup) {
	captchaAPI := v1.APIGroupAPP.SystemAPI.CaptchaAPI
	baseRouter := group.Group("base")
	{
		baseRouter.GET("captcha", captchaAPI.Captcha)
		baseRouter.POST("captchaPhone", captchaAPI.CaptchaPhone)
	}
}
