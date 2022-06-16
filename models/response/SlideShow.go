package response

import "bank-product-spike-system/global"

type SlideShowBaseInfo struct {
	SlideShowID  uint   `json:"slide_show_id"` // 轮播图ID
	Host         string `json:"host"`          // 源
	RelativePATH string `json:"relative_path"` // 相对路径
	Title        string `json:"title"`         // 标题
	Description  string `json:"description"`   // 描述
}

type SlideShowInfo struct {
	global.Model
	Host         string `json:"host"`          // 源
	RelativePATH string `json:"relative_path"` // 相对路径
	Title        string `json:"title"`         // 标题
	Description  string `json:"description"`   // 描述
	Weight       int    `json:"weight"`        // 权重
	AdminID      uint   `json:"adminID"`       // 发布者id
	AdminName    string `json:"adminName"`     // 发布者名称
}

type SlideShowBaseList struct {
	List []SlideShowBaseInfo `json:"list"`
}

type SlideShowList struct {
	List []SlideShowInfo `json:"list"`
}
