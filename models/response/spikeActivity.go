package response

import (
	"bank-product-spike-system/models"
	"time"
)

type SpikeActivityBaseInfo struct {
	ActivityID            uint      `json:"activity_id"`                             // 活动ID
	ActivityName          string    `json:"activity_name"`                           // 秒杀活动名称
	DescriptionText       string    `json:"description"`                             // 活动描述文本
	ActivityStatus        int       `json:"activity_status" gorm:"default:1"`        // 活动状态 1-筛选配置 2-资源配置 3-等待开始
	RolesIntroductionText string    `json:"roles_introduction"`                      // 活动规则介绍文本
	StartTime             time.Time `json:"start_time"`                              // 秒杀开始时间
	EndingTime            time.Time `json:"ending_time"`                             // 秒杀结束时间
	NumberOfParticipants  int       `json:"number_of_participants" gorm:"default:0"` // 秒杀活动参与人数
	NumberOfFavorites     int       `json:"number_of_favorites" gorm:"default:0"`    // 秒杀活动关注人数
	AuditStatus           int       `json:"audit_status"`                            // 审核状态
	AuditMessage          string    `json:"audit_message"`                           // 审核消息
	AdminID               uint      `json:"admin_id"`                                // 创建者ID
	AdminName             string    `json:"admin_name"`                              // 创建者名称
}

type SpikeActivityBaseInfoList struct {
	Total int64                   `json:"total"` // 列表数据总数
	List  []SpikeActivityBaseInfo `json:"list"`  // 活动列表
}

type AttentionList struct {
	List []ActivityBaseInfo `json:"list"` // 活动列表
}

type IsAttention struct {
	IsAttention bool `json:"is_attention"` // 是否关注
}

type ActivityBaseInfo struct {
	ActivityName         string    `json:"activity_name"`          // 活动名称
	ActivityID           uint      `json:"activity_id"`            // 活动ID
	StartTime            time.Time `json:"start_time"`             // 活动开始时间
	NumberOfParticipants int       `json:"number_of_participants"` // 秒杀活动参与人数
	NumberOfFavorites    int       `json:"number_of_favorites"`    // 秒杀活动关注人数
	ProductRate          float64   `json:"product_rate"`           // 产品利率
	ProductTotal         int64     `json:"product_total"`          // 产品总数
	ProductNowNumber     int64     `json:"product_now_number"`     // 产品当前数量
}

type GetActivityList struct {
	RunningActivity  []ActivityBaseInfo `json:"running_activity"`  // 进行中活动
	UpcomingActivity []ActivityBaseInfo `json:"upcoming_activity"` // 即将开放的活动
}

type ParticipateSpike struct {
	QueueID uint `json:"queue_id"` // 查询id
}

type QuerySpikeResult struct {
	Result  int    `json:"result"` // 结果状态 -1 秒杀失败 0-等待秒杀 1-秒杀成功 2-等待支付
	Message string `json:"message"`
}

type GetActivityGuidanceInformation struct {
	ActivityID           uint                                 `json:"activity_id"`                             // 活动ID
	ActivityName         string                               `json:"activity_name"`                           // 秒杀活动名称
	DescriptionText      string                               `json:"description"`                             // 活动描述文本
	ActivityStatus       int                                  `json:"activity_status" gorm:"default:1"`        // 活动状态 1-筛选配置 2-资源配置 3-等待审核 4-等待开始
	StartTime            time.Time                            `json:"start_time"`                              // 秒杀开始时间
	EndingTime           time.Time                            `json:"ending_time"`                             // 秒杀结束时间
	NumberOfParticipants int                                  `json:"number_of_participants" gorm:"default:0"` // 秒杀活动参与人数
	NumberOfFavorites    int                                  `json:"number_of_favorites" gorm:"default:0"`    // 秒杀活动关注人数
	AdminID              uint                                 `json:"admin_id"`                                // 创建者ID
	AdminName            string                               `json:"admin_name"`                              // 创建者名称
	Product              models.Product                       `json:"product"`                                 // 产品
	GuidanceInformation  []models.ActivityGuidanceInformation `json:"guidance_information"`                    // 引导信息
}

type GetSpikeVerifyCode struct {
	VerifyCode string `json:"verify_code"`
}
