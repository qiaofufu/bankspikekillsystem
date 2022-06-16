package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"time"
)

type FilterService struct{}

const (
	exist   = true
	unexist = false

	Refuse = "refuse"
	Pass   = "pass"
)

// GetTableNameList 获取数据表列表
func (f FilterService) GetTableNameList() (tables []models.FilterTable, err error) {
	err = global.DB.Find(&tables).Error
	return
}

// IsExist 检验是否存在
func (f FilterService) IsExist(name string) bool {
	var cnt int64
	global.DB.Model(&models.FilterTable{}).Where("table_name = ?", name).Count(&cnt)
	fmt.Println(cnt)
	return cnt > 0
}

func (f FilterService) GetFilterTableByID(id uint) (filterTable models.FilterTable, err error) {
	err1 := global.DB.Where("id = ?", id).First(&filterTable).Error
	if err1 != nil {
		err = errors.New("获取FilterTable失败")
		return
	}
	return
}

// GetFieldListByTableName 获取字段名列表
func (f FilterService) GetFieldListByTableName(tableName string) (fieldList []string, err error) {
	if f.IsExist(tableName) == unexist {
		err = errors.New("数据表名错误")
		return
	}
	sqlStr := "select * from " + tableName
	rows, err1 := global.DB.Raw(sqlStr).Rows()
	if err1 != nil {
		err = errors.New("获取数据表field错误")
		return
	}

	fieldList, err = rows.Columns()
	return
}

// SetValueRange 设置取值范围
func (f FilterService) SetValueRange(roles *models.FilterRoles) (err error) {
	err = global.DB.Create(&roles).Error
	return
}

// GetFilterRolesByActivity 通过活动ID获取筛选规则
func (f FilterService) GetFilterRolesByActivity(activityID uint) (filterList []models.FilterRoles, err error) {
	err = global.DB.Where("activity_id = ?", activityID).Find(&filterList).Error
	return
}

// GetFilterRolesList 获取筛选规则列表
func (f FilterService) GetFilterRolesList(activityID uint) (list []models.FilterRoles, err error) {
	result := global.DB
	if activityID != 0 {
		result.Where("activity_id = ?", activityID)
	}
	err = result.Find(&list).Error
	return
}

// DeleteFilterRoles 删除筛选规则
func (f FilterService) DeleteFilterRoles(activityID uint) (err error) {
	return global.DB.Where("activity_id = ?", activityID).Delete(&models.FilterRoles{}).Error
}

func (f FilterService) DeleteFilterRolesByID(id uint) (err error) {
	return global.DB.Where("id = ?", id).Delete(&models.FilterRoles{}).Error
}

// FinishFilterConfiguration 完成筛选配置
func (f FilterService) FinishFilterConfiguration(activityID uint) (err error) {
	err = global.DB.Model(&models.SpikeActivity{}).Where("activity_status >= ? and id = ?", models.ActivityFilterRoleConfiguration, activityID).Update("activity_status", models.ActivityResourceConfiguration).Error
	if err != nil {
		return
	}
	return global.DB.Model(&models.SpikeActivity{}).Where("activity_status >= ? and id = ?", models.ActivityFilterRoleConfiguration, activityID).Update("audit_status", models.WaitAudit).Error

}

// Check 检查是否符合规则
func (f FilterService) Check(token string, activityID uint) (err error) {
	user, err1 := UserService{}.GetUserByToken(token)
	if err1 != nil {
		err = err1
		return
	}
	// 判断是否筛选过
	var record models.ActivityUser
	if err1 := global.DB.Where("user_id = ? and spike_activity_id = ?", user.ID, activityID).First(&record).Error; err1 == gorm.ErrRecordNotFound {

		list, err1 := f.GetFilterRolesByActivity(activityID)
		if err1 != nil {
			err = err1
			return
		}

		for _, v := range list {
			var cnt int64
			filterTable, err1 := f.GetFilterTableByID(v.FilterTableID)
			if err1 != nil {
				err = err1
				return
			}
			err1 = global.DB.Where(fmt.Sprintf("uuid = '%s' and %s %s", user.UUID, v.FiledName, v.ValueRange)).Table(filterTable.TableName).Count(&cnt).Error
			if err1 != nil || cnt == 0 {
				global.DB.Create(&models.ActivityUser{SpikeActivityID: activityID, UserID: user.ID, Status: models.AccessDenied, ApplyTime: time.Now()})
				err = errors.New("初筛未通过")
				return
			}
		}

		global.DB.Create(&models.ActivityUser{SpikeActivityID: activityID, UserID: user.ID, Status: models.AccessAllowed, ApplyTime: time.Now()})
		return
	} else if err1 == nil {
		if record.Status == models.AccessAllowed {
			return nil
		} else {
			err = errors.New("初筛未通过")
			return
		}
	} else {
		err = errors.New("内部错误准入判断失败")
		return
	}
}

// BuildTree 构造树
func (f FilterService) BuildTree(node models.Nodes, edge []models.Edge) (tree models.FilterNode, err error) {
	sort.Sort(node)
	var fnode []models.FilterNode
	for _, v := range node {
		elem := models.FilterNode{
			FilterTable: v.TableName,
			FilterField: v.FieldName,
			ValueScope:  v.Condition,
			ChildNode:   nil,
		}
		fnode = append(fnode, elem)
	}

	for _, v := range edge {
		fnode[v.SourceNodeID-1].ChildNode = append(fnode[v.SourceNodeID-1].ChildNode, fnode[v.TargetNodeID-1])
	}
	return fnode[0], nil
}

// FinishFilterConfiguration2 完成筛选规则配置
func (f FilterService) FinishFilterConfiguration2(node []models.Node, edge []models.Edge, activityID uint) (err error) {
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("activity_id = ?", activityID).Delete(&node)
		tx.Where("activity_id = ?", activityID).Delete(&edge)
		if tx.Create(node).Error != nil {
			return errors.New("创建节点失败，内部错误")
		}
		if tx.Create(edge).Error != nil {
			return errors.New("创建边失败，内部错误")
		}
		if f.FinishFilterConfiguration(activityID) != nil {
			return errors.New("更新活动状态失败，内部错误")
		}
		return nil
	})
	return
}

func (f FilterService) GetFilterTree(ActivityID uint) (tree models.FilterNode, err error) {
	nodes, err := f.GetFilterNode(ActivityID)
	if err != nil {
		return
	}
	edge, err := f.GetFilterEdge(ActivityID)
	if err != nil {
		return
	}
	tree, err = f.BuildTree(nodes, edge)
	return
}

func (f FilterService) GetFilterNode(ActivityID uint) (node []models.Node, err error) {
	err1 := global.DB.Where("activity_id = ?", ActivityID).Find(&node).Error
	if err1 != nil {
		err = errors.New("获取筛选规则节点错误")
		return
	}
	return
}

func (f FilterService) GetFilterEdge(ActivityID uint) (edge []models.Edge, err error) {
	err1 := global.DB.Where("activity_id = ?", ActivityID).Find(&edge).Error
	if err1 != nil {
		err = errors.New("获取筛选规则边集错误")
		return
	}
	return
}

// Check2 检查是否符合规则
func (f FilterService) Check2(token string, activityID uint) (err error) {
	user, err1 := UserService{}.GetUserByToken(token)
	if err1 != nil {
		err = err1
		return
	}
	// 判断是否筛选过
	var record models.ActivityUser
	if err1 := global.DB.Where("user_id = ? and spike_activity_id = ?", user.ID, activityID).First(&record).Error; err1 == gorm.ErrRecordNotFound {
		tree, err := f.GetFilterTree(activityID)
		if err != nil {
			return err
		}
		if f.ValidationRules(tree, &user) == false {
			global.DB.Create(&models.ActivityUser{SpikeActivityID: activityID, UserID: user.ID, Status: models.AccessDenied, ApplyTime: time.Now()})

			return errors.New("初筛未通过")
		}
		global.DB.Create(&models.ActivityUser{SpikeActivityID: activityID, UserID: user.ID, Status: models.AccessAllowed, ApplyTime: time.Now()})
		return nil
	} else if err1 == nil {
		if record.Status == models.AccessAllowed {
			return nil
		} else {
			err = errors.New("初筛未通过")
			return
		}
	} else {
		err = errors.New("内部错误准入判断失败")
		return
	}
}

func (f FilterService) ValidationRules(tree models.FilterNode, user *models.User) bool {
	if len(tree.ChildNode) == 0 {
		return true
	}
	var cnt int64
	err1 := global.DB.Where(fmt.Sprintf("uuid = '%s' and %s %s", user.UUID, tree.FilterField, tree.FilterField)).Table(tree.FilterTable).Count(&cnt).Error
	if err1 != nil || cnt == 0 {
		return false
	}
	for _, v := range tree.ChildNode {
		if f.ValidationRules(v, user) == true {
			return true
		}
	}
	return false
}
