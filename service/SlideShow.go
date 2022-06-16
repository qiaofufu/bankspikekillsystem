package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
)

type SlideShowService struct {
}

// GetSlideShowBaseList 获取轮播图基本列表
func (s SlideShowService) GetSlideShowBaseList() (list []models.SlideShow, err error) {
	result := global.DB.Order("weight desc").Find(&list)
	err = result.Error
	return
}

// AddSlideShow 添加轮播图
func (s SlideShowService) AddSlideShow(slideShow models.SlideShow) (err error) {
	result := global.DB.Create(&slideShow)
	err = result.Error
	return
}

// UpdateSlideShow 更新轮播图
func (s SlideShowService) UpdateSlideShow(slideShow models.SlideShow) (err error) {
	result := global.DB.Updates(&slideShow)
	err = result.Error
	return
}

// DeleteSlideShow 删除轮播图
func (s SlideShowService) DeleteSlideShow(slideShowID uint) (err error) {
	result := global.DB.Delete(&models.SlideShow{}, slideShowID)
	err = result.Error
	return
}

// GetSlideShowList 获取轮播图列表
func (s SlideShowService) GetSlideShowList() (list []models.SlideShow, err error) {
	result := global.DB.Preload("Admin").Order("weight desc").Find(&list)
	err = result.Error
	return
}
