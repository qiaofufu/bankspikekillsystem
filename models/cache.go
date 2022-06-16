package models

import "time"

type SpikeActivityCache struct {
	ActivityBaseInfoCache
	ProductNumberCache
}

type ActivityBaseInfoCache struct {
	ActivityName string    `json:"activity_name"` // 活动名称
	ActivityID   uint      `json:"activity_id"`   // 活动ID
	StartTime    time.Time `json:"start_time"`    // 活动开始时间
	EndingTime   time.Time `json:"ending_time"`   // 活动结束时间
	ProductID    uint      `json:"product_id"`    // 产品ID
	ProductTotal int64     `json:"product_total"` // 产品总数
}

type ProductNumberCache struct {
	ProductNumber int64 `json:"number"` // 产品库存
}
