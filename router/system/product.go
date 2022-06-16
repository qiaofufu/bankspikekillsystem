package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type ProductRouter struct{}

func (receiver ProductRouter) InitProductRouter(router *gin.RouterGroup) {
	productAPI := v1.APIGroupAPP.SystemAPI.ProductAPI
	productRouter := router.Group("product")
	{
		productRouter.POST("featuredProduct", productAPI.GetFeaturedProduct)
		productRouter.POST("get", productAPI.GetVisibleProduct)
	}
	productRouterAuth := router.Group("product")
	productRouterAuth.Use(middleware.JWTAuth())
	{
		productRouterAuth.POST("release", productAPI.ReleaseProduct)
		productRouterAuth.PUT("update", productAPI.UpdateProduct)
		productRouterAuth.DELETE("delete", productAPI.DeleteProduct)
		productRouterAuth.POST("all", productAPI.GetAllProduct)
		productTagRouterAuth := productRouterAuth.Group("tag")
		{
			productTagRouterAuth.POST("list", productAPI.GetProductTagsList)
		}

	}
}
