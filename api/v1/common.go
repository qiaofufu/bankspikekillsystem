package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"

	"errors"
	"github.com/gin-gonic/gin"
)

// getToken 获取Token, 并进行相关处理
func getToken(ctx *gin.Context) (token string, err error) {
	token = ctx.Request.Header.Get("Authorization")
	if token == "" {
		err = errors.New("获取Token失败，请重新登录")
	}
	return
}

// getTags 获取Tags
func getTags(tagsArr []request.ProductTag) (tags []models.Tag) {
	for _, v := range tagsArr {
		tags = append(tags, models.Tag{
			Model:   global.Model{ID: v.TagID},
			TagName: v.TagName,
		})
	}
	return
}

// getAuthority 从authID数组中获取authority struct
func getAuthority(authID []uint) (auth []models.Authority) {
	for _, v := range authID {
		auth = append(auth, models.Authority{
			Model: global.Model{ID: v},
		})
	}
	return
}
