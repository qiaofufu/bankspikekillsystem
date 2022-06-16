package response

import "bank-product-spike-system/models"

type OrderInfo struct {
	models.Order
	ProductName         string  `json:"product_name"`
	ProductDescription  string  `json:"product_description"`
	ProductType         int     `json:"product_type"`
	ProductInterestRate float64 `json:"product_interest_rate"`
	InterestRateType    int     `json:"interest_rate_type"`
	MinHoldTime         string  `json:"min_hold_time"`
	ProductPrice        float64 `json:"product_price"`
}

type OrderList struct {
	Total int64
	List  []OrderInfo `json:"list"`
}

type OrderBaseInfo struct {
	models.Order
	RealName string `json:"real_name"`
}

type GetOrderListByActivityID struct {
	Total int64           `json:"total"` // 总数
	List  []OrderBaseInfo `json:"list"`  // 订单信息
}
