package response

import "bank-product-spike-system/models"

type RolesInfo struct {
	ID            uint   `json:"id"` // id
	FilterTableID uint   `json:"filter_table_id"`
	FiledName     string `json:"filed_name"`  // 过滤字段名
	ValueRange    string `json:"value_range"` // 过滤取值范围
	ErrorTips     string `json:"error_tips"`  // 错误提示
	Description   string `json:"description"` // 描述
	NextID        uint   `json:"next_id"`
}
type RolesList struct {
	List []RolesInfo
}

type FilterTree struct {
	Node []models.Node
	Edge []models.Edge
}
