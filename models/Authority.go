package models

import (
	"bank-product-spike-system/global"
)

const (
	Root                 = "root"
	AdminManager         = "adminManager"
	UserManager          = "userManager"
	ProductManager       = "productManager"
	SpikeActivityManager = "spikeActivityManager"
	SlideShowManager     = "SlideShowManager"
	ArticleManager       = "ArticleManager"
	ProductAuditor       = "ProductAuditor"
	ActivityAuditor      = "ActivityAuditor"
	ArticleAuditor       = "ArticleAuditor"
)

type Authority struct {
	global.Model
	AuthorityName string  `gorm:"default:'user'" json:"authority_name"` // 权限名
	Rank          uint    `json:"rank" gorm:"default:1"`                // 权限等级
	Description   string  `json:"description"`                          // 权限描述
	Admin         []Admin `gorm:"many2many:admin_authority" json:"-" swaggerignore:"true"`
}
