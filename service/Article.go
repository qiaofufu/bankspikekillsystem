package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"bank-product-spike-system/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ArticleService struct {
}

func (a ArticleService) GetArticleBaseList(page int, pageSize int, status int) (list []models.Article, cnt int64, err error) {
	result := global.DB.Model(&list)
	if status != 0 {
		result = result.Where("audit_status = ?", status)
	}
	result.Model(&models.Article{}).Count(&cnt)
	result = result.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	err = result.Error
	return
}

func (a ArticleService) GetArticleContent(articleID uint, status int) (content models.Article, err error) {
	result := global.DB.Model(&content)
	if status != 0 {
		result = result.Where("audit_status = ?", status)
	}
	result = result.Where("id = ?", articleID).First(&content)
	err = result.Error
	return
}

func (a ArticleService) PublishArticle(ctx *gin.Context, article models.Article) (err error) {
	url, err := utils.UploadFile(ctx, "front_cover")
	if err != nil {
		err = errors.New(fmt.Sprintf("上传封面失败，[%v]", err))
		return
	}
	article.FrontCover = models.Url{
		Host:         url.Host,
		RelativePATH: url.RelativePath,
	}
	result := global.DB.Create(&article)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("发表文章失败， err:%v", result.Error))
	}
	return
}

func (a ArticleService) UpdateArticle(article models.Article) (err error) {

	result := global.DB.Updates(&article)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("更新文章失败， err:%v", result.Error))
	}
	return
}

func (a ArticleService) DeleteArticle(articleID uint) (err error) {
	result := global.DB.Delete(&models.Article{}, articleID)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("删除文章失败， err：%v", result.Error))
	}
	return
}

func (a ArticleService) LikeArticle(articleID uint) (err error) {
	sql := fmt.Sprintf("UPDATE `articles` SET `like_number`=like_number+1,`updated_at`='2022-03-21 21:47:44.65' WHERE id = %d AND `articles`.`deleted_at` IS NULL", articleID)
	result := global.DB.Exec(sql)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("操作失败， err:%v", result.Error))
	}
	return
}
