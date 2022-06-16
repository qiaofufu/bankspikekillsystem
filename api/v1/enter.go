package v1

import (
	"bank-product-spike-system/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Group struct {
	SystemAPI APIGroup
}

var APIGroupAPP = new(Group)

type APIGroup struct {
	AdminAPI
	CaptchaAPI
	UserAPI
	ProductAPI
	SpikeActivityAPI
	FilterAPI
	SlideShowAPI
	ArticleAPI
	OrderAPI
	AuditAPI
}

var (
	AdminService         = service.ServiceGroupAPP.SystemService.AdminService
	CaptchaService       = service.ServiceGroupAPP.SystemService.CaptchaService
	UserService          = service.ServiceGroupAPP.SystemService.UserService
	ProductService       = service.ServiceGroupAPP.SystemService.ProductService
	AuthorityService     = service.ServiceGroupAPP.SystemService.AuthorityService
	SpikeActivityService = service.ServiceGroupAPP.SystemService.SpikeActivityService
	FilterService        = service.ServiceGroupAPP.SystemService.FilterService
	SlideShowService     = service.ServiceGroupAPP.SystemService.SlideShowService
	ArticleService       = service.ServiceGroupAPP.SystemService.ArticleService
	OrderService         = service.ServiceGroupAPP.SystemService.OrderService
	RedisService         = service.ServiceGroupAPP.SystemService.RedisServer
	JWTService           = service.ServiceGroupAPP.SystemService.JWTService
)

func getFormData(ctx *gin.Context, key string) (value string, err error) {
	value = ctx.PostForm(key)
	if value == "" {
		err = errors.New(fmt.Sprintf("获取参数[%s]失败", key))
	}
	return
}
