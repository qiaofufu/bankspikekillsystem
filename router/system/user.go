package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (receiver UserRouter) InitUserRouter(group *gin.RouterGroup) {
	userRouter := group.Group("user")
	userAPI := v1.APIGroupAPP.SystemAPI.UserAPI
	//userRouter.Use(middleware.IpLimit())
	{
		userRouter.POST("register", userAPI.Register)
		userRouter.POST("login", userAPI.Login)
	}
	userRouterAuth := group.Group("user")
	userRouterAuth.Use(middleware.JWTAuth())
	{
		// 后台
		userRouterAuth.PUT("update", userAPI.UpdateUserInfo, userAPI.UpdateUserInfo)
		userRouterAuth.POST("list", userAPI.GetUserList)
		// 前台
		userRouterAuth.PUT("changeSelfPassword", userAPI.ChangeSelfPassword)
		userRouterAuth.GET("getInformation", userAPI.GetInformation)
		userRouterAuth.GET("changePhone", userAPI.ChangePhone)
		userRouterAuth.PUT("updateSelfInfo", userAPI.UpdateUserSelfInfo)
	}
}
