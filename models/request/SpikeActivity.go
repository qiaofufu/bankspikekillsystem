package request

type ReleaseSpikeActivity struct {
	ActivityName      string `json:"activity_name" binding:"required"`      // 秒杀活动名称
	StartTime         string `json:"start_time" binding:"required"`         // 秒杀开始时间 YYYY-MM-DD hh:mm:ss
	EndingTime        string `json:"ending_time" binding:"required"`        // 秒杀结束时间 YYYY-MM-DD hh:mm:ss
	ProductID         uint   `json:"product_id" binding:"required"`         // 产品ID
	Description       string `json:"description" binding:"required"`        // 活动描述
	RolesIntroduction string `json:"roles_introduction" binding:"required"` // 活动规则介绍
}

type UpdateActivity struct {
	ActivityID        uint   `json:"activity_id" binding:"required"` // 活动ID
	ActivityName      string `json:"activity_name"`                  // 活动名称
	StartTime         string `json:"start_time" `                    // 秒杀开始时间 YYYY-MM-DD hh:mm:ss
	EndingTime        string `json:"ending_time" `                   // 秒杀结束时间 YYYY-MM-DD hh:mm:ss
	Description       string `json:"description" `                   // 活动描述
	ProductID         uint   `json:"product_id"`                     // 活动id
	RolesIntroduction string `json:"roles_introduction" `            // 活动规则介绍
}

type DeleteActivity struct {
	ActivityID uint `json:"activity_id"` // 活动id
}

type FilterRuleConfiguration struct {
	RoleText   string `json:"role_text"`   // 规则文本
	ActivityID uint   `json:"activity_id"` // 活动ID
}

type DeleteActivityProduct struct {
	ActivityID uint   `json:"activity_id" binding:"required"` // 活动id
	ProductID  uint   `json:"product_id" binding:"required"`  // 产品id
	UUID       string `json:"uuid" binding:"required,uuid4"`  // 用户uuid
}

type AddActivityProduct struct {
	ActivityID uint `json:"activity_id" binding:"required"` // 活动id
	ProductID  uint `json:"product_id" binding:"required"`  // 产品id
}

type GetActivityGuidanceInfo struct {
	ActivityID uint `json:"activity_id" binding:"required"` // 活动id
}

type Attention struct {
	ActivityID uint `json:"activity_id"` // 活动id
}

type GetActivityStatus struct {
	ActivityID uint `json:"activity_id" binding:"required"` // 活动id
}

type GetActivityList struct {
	PageInfo
	ActivityStatus int    `json:"activity_status"` // 活动状态, 为0不筛选 1-筛选配置阶段 2-资源配置 4-等待开始
	AuditStatus    int    `json:"audit_status"`    // 活动审核状态， 为0不进行筛选 1-等待审核 2-审核成功 -1审核失败
	StartTime      string `json:"start_time"`      // 开始时间 为空不筛选
	EndingTime     string `json:"ending_time"`     // 结束时间 为空不进行筛选
}

type ParticipateActivity struct {
	ActivityID uint `json:"activity_id"`
	ProductID  uint `json:"product_id"`
}

type ActivityID struct {
	ActivityID uint `json:"activity_id"`
}

type QuerySpikeResult struct {
	ActivityID uint `json:"activity_id"`
	ProductID  uint `json:"product_id"`
}

type CacheSpike struct {
	ActivityID uint `json:"activity_id"`
}

type GetSpikeVerifyCode struct {
	ActivityID uint `json:"activity_id"`
}

type ParticipateActivityTest struct {
	ActivityID uint `json:"activity_id"`
	ProductID  uint `json:"product_id"`
	UserID     uint `json:"user_id"`
}
