package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type SlideShowRouter struct{}

func (s SlideShowRouter) InitSlideShowRouter(group *gin.RouterGroup) {
	router := group.Group("SlideShow")
	slideShowAPI := v1.APIGroupAPP.SystemAPI.SlideShowAPI
	{
		router.GET("list", slideShowAPI.GetSlideShow)
	}
	routerAuth := group.Group("SlideShow")
	routerAuth.Use(middleware.JWTAuth())
	{
		// 后台
		routerAuth.POST("add", slideShowAPI.AddSlideShow)
		routerAuth.PUT("update", slideShowAPI.UpdateSlideShow)
		routerAuth.DELETE("delete", slideShowAPI.DeleteSlideShow)
		routerAuth.GET("listAdmin", slideShowAPI.GetSlideShowListByAdmin)
	}
}
