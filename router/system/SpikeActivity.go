package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type SpikeActivityRouter struct{}

func (receiver SpikeActivityRouter) InitSpikeActivityRouter(group *gin.RouterGroup) {
	spikeAPI := v1.APIGroupAPP.SystemAPI.SpikeActivityAPI
	spikeRouter := group.Group("spikeActivity")
	{
		spikeRouter.GET("getActivityList", spikeAPI.GetActivity)
	}
	spikeRouterAuth := group.Group("spikeActivity")
	spikeRouterAuth.Use(middleware.JWTAuth())
	{
		// 秒杀活动后台用户
		spikeRouterAuth.POST("release", spikeAPI.ReleaseSpikeActivity)
		spikeRouterAuth.PUT("update", spikeAPI.UpdateSpikeActivity)
		spikeRouterAuth.DELETE("delete", spikeAPI.DeleteActivity)
		spikeRouterAuth.POST("list", spikeAPI.GetActivityList)
		spikeRouterAuth.POST("addGuidanceInformation", spikeAPI.AddActivityGuidanceInformation)
		spikeRouterAuth.POST("updateGuidanceInformation", spikeAPI.UpdateActivityGuidanceInformation)
		spikeRouterAuth.POST("getActivityUserRecord", spikeAPI.GetActivityUserRecord)
		spikeRouterAuth.POST("getActivityInfoByID", spikeAPI.GetActivityInfoByID)
		spikeRouterAuth.POST("cache", spikeAPI.CacheSpike)
		spikeRouterAuth.POST("getCacheSpike", spikeAPI.GetCacheSpike)

		// 秒杀活动前台用户
		spikeRouterAuth.POST("getActivityGuidanceInformation", spikeAPI.GetActivityGuidanceInformation)
		spikeRouterAuth.POST("attention", spikeAPI.Attention)
		spikeRouterAuth.GET("attention", spikeAPI.GetAttentionList)
		spikeRouterAuth.DELETE("attention", spikeAPI.UnAttention)
		spikeRouterAuth.POST("isAttention", spikeAPI.IsAttention)
		spikeRouterAuth.POST("participateTest", spikeAPI.ParticipateActivityTest)
		spikeRouterAuth.POST("participate/:verifyCode", spikeAPI.ParticipateActivity)
		spikeRouterAuth.POST("getSpikeResult", spikeAPI.GetSpikeResult)
		spikeRouterAuth.POST("getSpikeVerifyCode", spikeAPI.GetSpikeVerifyCode)

	}
}
