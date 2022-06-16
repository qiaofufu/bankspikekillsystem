package request

type GenerateOrder struct {
	ActivityID uint `json:"activity_id"` // 活动id
	ProductID  uint `json:"product_id"`  // 产品id
	Number     uint `json:"number"`      // 购买数量
}

type GetOrderListByActivityID struct {
	PageInfo
	ActivityID  uint `json:"activity_id"`  // 活动ID
	OrderStatus int  `json:"order_status"` // 订单状态  1- 等待处理 2-处理成功 -1订单失败
}
