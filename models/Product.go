package models

import (
	"bank-product-spike-system/global"
)

const (
	AnnualInterestRate  = 1
	MonthlyInterestRate = 2
	DailyInterestRate   = 3
)

const (
	LoanProduct    = 1 // 贷款
	DepositProduct = 2 // 存款
)

const (
	YearRate = 1 // 年利率
	MoonRate = 2 // 月利率
	DayRate  = 3 // 日利率
)

type Product struct {
	global.Model
	ProductName         string  `gorm:"not null;" json:"product_name"`              // 产品名称
	ProductDescription  string  `gorm:"not null;" json:"product_description"`       // 产品描述
	ProductNumber       int64   `json:"product_number"`                             // 产品数量
	SoldNumber          int64   `gorm:"default:0" json:"sold_number"`               // 售出数量
	ProductType         int     `json:"product_type"`                               // 产品类型
	ProductTags         []Tag   `gorm:"many2many:product_tags" json:"product_tags"` // 产品标签
	ProductInterestRate float64 `json:"product_interest_rate"`                      // 产品利率
	InterestRateType    int     `json:"interest_rate_type"`                         // 产品利率类型 1-年利率 2-月利率 3-日利率
	MinHoldTime         string  `json:"min_hold_time"`                              // 最小持有时间
	ProductPrice        float64 `json:"product_price"`                              // 产品单笔价格
	Audit
	IsVisible bool  `gorm:"default:true" json:"is_visible"` // 是否可见
	Weights   int64 `gorm:"default:0" json:"weights"`       // 产品权重
}

type ProductTags struct {
	ProductID uint `json:"product_id"`
	TagsID    uint `json:"tags_id"`
}

func (receiver *ProductTags) TableName() string {
	return "product_tags"
}

func (p Product) GetMargin() (margin int64) {
	margin = p.ProductNumber - p.SoldNumber
	return
}
