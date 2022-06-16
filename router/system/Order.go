package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct {
}

func (o OrderRouter) InitOrderRouter(group *gin.RouterGroup) {
	orderAPI := v1.APIGroupAPP.SystemAPI.OrderAPI

	routerAuth := group.Group("order")
	routerAuth.Use(middleware.JWTAuth())
	{
		// 后台
		routerAuth.POST("getOrderListByActivityID", orderAPI.GetOrderListByActivityID)
		// 用户
		routerAuth.GET("getSelfOrder", orderAPI.GetSelfOrder)
		//routerAuth.POST("generateOrder", orderAPI.GenerateOrder)
	}
}
