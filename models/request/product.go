package request

type ReleaseProduct struct {
	ProductName         string       `json:"product_name" binding:"required"`          // 产品名称
	ProductDescription  string       `json:"product_description" binding:"required"`   // 产品描述
	ProductNumber       int64        `json:"product_number" binding:"required"`        // 产品数量
	ProductTags         []ProductTag `json:"product_tags" binding:""`                  // 产品标签
	ProductInterestRate float64      `json:"product_interest_rate" binding:"required"` // 产品利率
	InterestRateType    int          `json:"interest_rate_type" binding:""`            // 产品利率类型 1-年利率 2-月利率 3-日利率
	MinHoldTime         string       `json:"min_hold_time"`                            // 产品最小持有时间 YYYY:MM:DD
	ProductPrice        float64      `json:"product_price" binding:"required"`         // 产品单笔价格
}

type ProductTag struct {
	TagName string `json:"tag_name" binding:"required"` // 标签名
	TagID   uint   `json:"tag_id"`                      // 标签id 新建标签填0
}

type UpdateProduct struct {
	ProductID           uint         `json:"product_id" binding:"required"`    // 产品ID
	ProductName         string       `json:"product_name" binding:""`          // 产品名称
	ProductDescription  string       `json:"product_description" binding:""`   // 产品描述
	ProductNumber       int64        `json:"product_number" binding:""`        // 产品数量
	ProductTags         []ProductTag `json:"product_tags" binding:""`          // 产品标签
	ProductInterestRate float64      `json:"product_interest_rate" binding:""` // 产品利率
	MinHoldTime         string       `json:"min_hold_time"`                    // 产品最小持有时间 YYYY:MM:DD
	InterestRateType    int          `json:"interest_rate_type" binding:""`    // 产品利率类型 1-年利率 2-月利率 3-日利率
	ProductPrice        float64      `json:"product_price" binding:""`         // 产品单笔价格
}

type DeleteProduct struct {
	ProductID uint `json:"product_id"` // 被删除产品ID
}

type GetAllProduct struct {
	PageInfo
	ProductName      string `json:"product_name"`       // 产品名称
	InterestRateType int    `json:"interest_rate_type"` // 利率类型 1-年利率 2-月利率 3-日利率
	AuditStatus      int    `json:"audit_status"`       // 产品状态 1-未审核， 2-审核成功， -1-审核失败
}

type GetVisibleProduct struct {
	ProductID uint `json:"product_id"` // 产品ID
}
