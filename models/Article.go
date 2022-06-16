package models

import "bank-product-spike-system/global"

type Url struct {
	Host         string `json:"host"`
	RelativePATH string `json:"relative_path"`
}

type Article struct {
	global.Model
	FrontCover Url    `gorm:"embedded" json:"front_cover"` // 封面
	Title      string `gorm:"varchar(20)" json:"title"`    // 标题
	Content    string `gorm:"" json:"content"`             // 内容
	Author     string `json:"author"`                      // 作者
	LikeNumber int64  `json:"like_number"`                 // 喜欢数
	Audit
}
