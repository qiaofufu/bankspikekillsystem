package request

type PageInfo struct {
	Page     int `json:"page"`      // 页数
	PageSize int `json:"page_size"` // 页面数据个数
}

type ShotMessageCode struct {
	Number string `json:"number" binding:"required"` // 手机号
}
