package models

import (
	"bank-product-spike-system/global"
	"time"
)

const (
	AccessDenied  = 0 // 拒绝进入
	AccessAllowed = 1 // 允许进入
)

type ActivityUser struct {
	global.Model
	SpikeActivityID uint      `json:"spike_activity_id"` // 活动id
	UserID          uint      `json:"user_id"`           // 用户id
	Status          int       `json:"status"`            // 状态 0-拒绝进入 1-允许进入
	ApplyTime       time.Time `json:"apply_time"`        // 申请时间
	Cause           string    `json:"cause"`             // 原因
}

func (a ActivityUser) TableName() string {
	return "activity_user"
}
