package models

import "bank-product-spike-system/global"

type ActivityGuidanceInformation struct {
	global.Model
	ActivityID uint `json:"activity_id"` // 活动id
	Url
}
