package models

import (
	"bank-product-spike-system/global"
)

type Tag struct {
	global.Model
	TagName  string    `json:"tag_name"`                               // 产品标签名
	Products []Product `gorm:"many2many:product_tags" json:"products"` // 此标签下的产品
}
