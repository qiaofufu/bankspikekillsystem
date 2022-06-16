package models

import "bank-product-spike-system/global"

type SlideShow struct {
	global.Model
	Host         string `json:"host"`          // 源
	RelativePATH string `json:"relative_path"` // 相对路径
	Title        string `json:"title"`         // 轮播图标题
	Description  string `json:"description"`   // 轮播图描述
	Weight       int    `json:"weight"`        // 权重
	AdminID      uint   // 创建者id
	Admin        Admin
}
