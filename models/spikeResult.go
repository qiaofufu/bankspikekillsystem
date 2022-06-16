package models

const (
	SpikeSuccess = 1
	SpikeFailed  = -1
	SpikeWait    = 0
)

type SpikeResult struct {
	ActivityID uint   `json:"activity_id"` // 活动id
	UserID     uint   `json:"user_id"`     // 用户id
	Status     int    `json:"status"`      // 结果状态 -1 秒杀失败 0-等待秒杀 1-秒杀成功
	OrderID    uint   `json:"order_id"`    // 订单id
	Message    string `json:"message"`     // 附加消息
}
