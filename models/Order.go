package models

import (
	"bank-product-spike-system/global"
)

const (
	PaymentWaiting    = 1  // 等待支付
	PaymentSuccessful = 2  // 支付成功
	PaymentFailed     = -1 // 支付失败
)

const (
	OrderWait    = 1
	OrderSuccess = 2
	OrderFailed  = -1
)

type Order struct {
	global.Model
	ActivityID    uint `json:"activity_id"`                     // 活动id
	ProductID     uint `json:"product_id"`                      // 产品ID
	UserID        uint `json:"user_id"`                         // 用户ID
	OrderStatus   int  `json:"order_status" gorm:"default:1"`   // 订单状态 1- 等待处理 2-处理成功 -1订单失败
	PaymentStatus int  `json:"payment_status" gorm:"default:1"` // 支付状态 1- 未支付 2-支付成功 -1 -支付失败
}
