package request

import "bank-product-spike-system/models"

type GetFieldNameList struct {
	TableName string `json:"table_name"` // 数据表名
}

type SetValueRange struct {
	ActivityID    uint   `json:"activity_id"`     // 活动ID
	FilterTableID uint   `json:"filter_table_id"` // 数据表id
	FieldName     string `json:"field_name"`      // 数据字段名
	ValueRange    string `json:"value_range"`     // 取值范围 等于: "= value" 大于："> value" 小于："< value"
	Description   string `json:"description"`     // 规则描述
	ErrorTips     string `json:"error_tips"`      // 错误提示
}

type FilterCheck struct {
	ActivityID uint `json:"activity_id"` // 活动ID
}

type GetFilterRolesByID struct {
	ActivityID uint `json:"activity_id"` // 活动ID
}

type DeleteFilterRoles struct {
	ActivityID uint `json:"activity_id"` // 活动ID
}

type FinishFilterConfiguration struct {
	ActivityID uint `json:"activity_id"` // 活动ID
}

type DeleteFilterRolesByID struct {
	FilterRolesID uint `json:"filter_roles_id"` // 筛选规则ID
}

type FinishFilterConfiguration2 struct {
	ActivityID uint          `json:"activity_id"` // 活动id
	Node       []models.Node `json:"node"`        // 节点
	Edge       []models.Edge `json:"edge"`        // 边
}
