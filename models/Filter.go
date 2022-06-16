package models

import "bank-product-spike-system/global"

type FilterTable struct {
	ID          uint   `json:"id"`          // id
	TableName   string `json:"table_name"`  // 表名
	Description string `json:"description"` // 描述
}

type FilterRoles struct {
	global.Model
	ActivityID    uint        `json:"activity_id"` // 规则所属活动id
	FilterTableID uint        `json:"filter_table_id"`
	FilterTable   FilterTable `json:"filter_table"` // 过滤表
	FiledName     string      `json:"filed_name"`   // 过滤字段名
	ValueRange    string      `json:"value_range"`  // 过滤取值范围
	ErrorTips     string      `json:"error_tips"`   // 错误提示
	Description   string      `json:"description"`  // 描述
	AdminID       uint        `json:"admin_id"`
	Admin         Admin       `json:"admin"` // 规则创建者
}

type FilterNode struct {
	FilterTable string       `json:"filter_table"` // 筛选表
	FilterField string       `json:"filter_field"` // 筛选字段
	ValueScope  string       `json:"value_scope"`  // 取值范围
	ChildNode   []FilterNode `json:"child_node"`   // 子节点
}

type Node struct {
	ActivityID uint   `json:"activity_id"`
	Id         uint   `json:"id"`         // 节点id
	TableName  string `json:"table_name"` // 表名
	FieldName  string `json:"field_name"` // 字段名
	Condition  string `json:"condition"`  // 条件
	Left       int    `json:"left"`
	Top        int    `json:"top"`
}

type Edge struct {
	ActivityID   uint `json:"activity_id"`
	SourceNodeID uint `json:"source_node_id"` // 源节点id
	TargetNodeID uint `json:"target_node_id"` // 目标节点id
}

type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Id <= n[j].Id
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
	return
}
