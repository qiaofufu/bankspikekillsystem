package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type FilterRouter struct{}

func (receiver FilterRouter) InitFilterRouter(group *gin.RouterGroup) {
	filterAPI := v1.APIGroupAPP.SystemAPI.FilterAPI
	filterRouterAuth := group.Group("filter")
	filterRouterAuth.Use(middleware.JWTAuth())
	{
		// 配置准入规则
		filterRouterAuth.GET("tableList", filterAPI.GetTableNameList)
		filterRouterAuth.POST("fieldList", filterAPI.GetFieldNameList)
		filterRouterAuth.POST("setValueRange", filterAPI.SetValueRange)
		filterRouterAuth.DELETE("delete", filterAPI.DeleteFilterRoles)
		filterRouterAuth.DELETE("deleteByID", filterAPI.DeleteFilterRolesByID)
		filterRouterAuth.GET("list", filterAPI.GetRolesList)
		filterRouterAuth.POST("getFilterRolesByID", filterAPI.GetFilterRolesByID)
		filterRouterAuth.POST("finishFilterConfiguration", filterAPI.FinishFilterConfiguration)
		filterRouterAuth.POST("finishFilterConfiguration2", filterAPI.FinishFilterConfiguration2)
		filterRouterAuth.POST("getFilterTree", filterAPI.GetFilterTree)
		// 准入规则校验
		filterRouterAuth.POST("check", filterAPI.Check)
		filterRouterAuth.POST("check2", filterAPI.Check2)
	}
}
