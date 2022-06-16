package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type AuditRouter struct {
}

func (a AuditRouter) InitAuditRouter(group *gin.RouterGroup) {
	auditAPI := v1.APIGroupAPP.SystemAPI.AuditAPI
	routerAuth := group.Group("audit")
	routerAuth.Use(middleware.JWTAuth())
	{
		routerAuth.POST("product", auditAPI.AuditProduct)
		routerAuth.POST("article", auditAPI.AuditArticle)
		routerAuth.POST("activity", auditAPI.AuditActivity)
	}
}
