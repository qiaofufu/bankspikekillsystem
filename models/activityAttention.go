package models

import "bank-product-spike-system/global"

type ActivityAttention struct {
	global.Model
	ActivityID uint `json:"activity_id"` // 活动id
	UserID     uint `json:"user_id"`     // 用户id
}
