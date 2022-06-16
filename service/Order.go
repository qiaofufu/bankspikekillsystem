package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"errors"
	"fmt"
)

type OrderService struct {
}

// GenerateOrder 生成订单
func (o OrderService) GenerateOrder(order models.Order) (err error) {
	result := global.DB.Create(&order)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("生成订单失败， err:%v", result.Error))
	}
	return
}

func (o OrderService) UpdateOrder(order models.Order) (err error) {
	result := global.DB.Updates(&order)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("更新订单失败， err:%v", result.Error))
	}
	return
}

func (o OrderService) GetOrderList(page int, pageSize int, activityID uint, orderStatus int) (cnt int64, list []models.Order, err error) {
	tx := global.DB.Model(&models.Order{}).Where("activity_id = ?", activityID)
	if orderStatus != 0 {
		tx = tx.Where("order_status = ?", orderStatus)
	}
	tx.Count(&cnt)
	err1 := tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error
	if err1 != nil {
		err = errors.New(fmt.Sprintf("获取订单失败, err:%v", err1))
	}
	return
}

func (o OrderService) GetSelfOrder(userID uint) (list []models.Order, err error) {
	result := global.DB.Where("user_id = ?", userID).Find(&list)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取自身订单失败， err:%v", result.Error))
	}
	return
}
