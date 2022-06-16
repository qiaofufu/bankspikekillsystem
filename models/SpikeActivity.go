package models

import (
	"bank-product-spike-system/global"
	"time"
)

const (
	ActivityFilterRoleConfiguration = 1 // 筛选配置
	ActivityResourceConfiguration   = 2 // 资源配置
	ActivityWaitAudit               = 3 // 等待审核
	ActivityWaitStart               = 4 // 等待开始

)

type SpikeActivity struct {
	global.Model
	ActivityName   string `json:"activity_name"`                    // 秒杀活动名称
	ActivityStatus int    `json:"activity_status" gorm:"default:1"` // 活动状态
	Audit
	BaseURL               string                        `json:"base_url"`                                // 本次活动的随机URL
	DescriptionText       string                        `json:"description"`                             // 活动描述文本
	RolesIntroductionText string                        `json:"roles_introduction"`                      // 活动规则介绍文本
	StartTime             time.Time                     `json:"start_time"`                              // 秒杀开始时间
	EndingTime            time.Time                     `json:"ending_time"`                             // 秒杀结束时间
	NumberOfParticipants  int                           `json:"number_of_participants" gorm:"default:0"` // 秒杀活动参与人数
	NumberOfFavorites     int                           `json:"number_of_favorites" gorm:"default:0"`    // 秒杀活动关注人数
	ProductID             uint                          `json:"product_id"`
	Product               Product                       `json:"product" gorm:""` // 参与本次秒杀活动的产品
	AdminID               uint                          `json:"admin_id"`
	Admin                 Admin                         `json:"admin"`                                             // 创建秒杀活动管理员信息
	User                  []User                        `json:"user" gorm:"many2many:activity_user"`               // 参与本次活动的用户信息
	GuidanceInformation   []ActivityGuidanceInformation `json:"guidance_information" gorm:"foreignKey:ActivityID"` // 本次活动的指示信息
}
