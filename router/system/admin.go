package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type AdminRouter struct{}

func (receiver AdminRouter) InitAdminRouter(router *gin.RouterGroup) {
	adminAPI := v1.APIGroupAPP.SystemAPI.AdminAPI
	adminRouter := router.Group("admin")
	{
		adminRouter.POST("login", adminAPI.Login)
	}

	adminRouterAuth := router.Group("admin")
	adminRouterAuth.Use(middleware.JWTAuth())
	{
		// 管理员基本操作
		adminRouterAuth.PUT("changePassword", adminAPI.ChangePassword)
		adminRouterAuth.PUT("changeSelfPassword", adminAPI.ChangeSelfPassword)
		adminRouterAuth.POST("getAdminInfoList", adminAPI.GetAdminInfoList)
		adminRouterAuth.POST("generateAdminAccount", adminAPI.GenerateAdminAccount)
		adminRouterAuth.PUT("setAdminAuthorities", adminAPI.SetAdminAuthorities)
		adminRouterAuth.PUT("setAdminInfo", adminAPI.SetAdminInfo)
		adminRouterAuth.PUT("setSelfInfo", adminAPI.SetSelfInfo)
		adminRouterAuth.GET("getInformation", adminAPI.GetInformation)
		// 权限
		authorityRouterAuth := adminRouterAuth.Group("authority")
		{
			authorityRouterAuth.POST("all", adminAPI.GetAllAuthority)
		}
	}
}
