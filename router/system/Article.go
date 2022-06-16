package system

import (
	v1 "bank-product-spike-system/api/v1"
	"bank-product-spike-system/middleware"
	"github.com/gin-gonic/gin"
)

type ArticleRouter struct {
}

func (a ArticleRouter) InitArticleRouter(group *gin.RouterGroup) {
	articleAPI := v1.APIGroupAPP.SystemAPI.ArticleAPI
	router := group.Group("article")
	{
		router.POST("list", articleAPI.GetArticleList)
		router.POST("getArticleContent", articleAPI.GetArticleContent)
	}
	routerAuth := group.Group("article")
	routerAuth.Use(middleware.JWTAuth())
	{
		// 后台
		routerAuth.POST("publish", articleAPI.PublishArticle)
		routerAuth.PUT("update", articleAPI.UpdateArticle)
		routerAuth.DELETE("delete", articleAPI.DeleteArticle)
		routerAuth.POST("getList", articleAPI.GetArticleListAdmin)
		router.POST("getArticleContentAdmin", articleAPI.GetArticleContentAdmin)
		// 前台
		routerAuth.POST("like", articleAPI.LikeArticle)
	}
}
