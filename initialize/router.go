package initialize

import (
	"bank-product-spike-system/middleware"
	"bank-product-spike-system/router"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	engine.GET("", func(context *gin.Context) {
		context.String(200, "success")
	})
	Router := engine.Group("api")

	// Tls
	Router.Use(middleware.Tls())

	// CORS 跨域
	Router.Use(middleware.Cors())

	// Log 日志
	Router.Use(middleware.Log())

	systemRouter := router.RouterGroupAPP.SystemRouter

	// 路由初始化
	systemRouter.AdminRouter.InitAdminRouter(Router)
	systemRouter.BaseRouter.InitBaseRouter(Router)
	systemRouter.UserRouter.InitUserRouter(Router)
	systemRouter.ProductRouter.InitProductRouter(Router)
	systemRouter.SpikeActivityRouter.InitSpikeActivityRouter(Router)
	systemRouter.FilterRouter.InitFilterRouter(Router)
	systemRouter.SlideShowRouter.InitSlideShowRouter(Router)
	systemRouter.ArticleRouter.InitArticleRouter(Router)
	systemRouter.AuditRouter.InitAuditRouter(Router)
	systemRouter.OrderRouter.InitOrderRouter(Router)
}
