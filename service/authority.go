package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"errors"
	"fmt"
)

type AuthorityService struct{}

// VerifyAuthority
// 权限验证, 通过返回 true, 未通过返回 false
func (receiver AuthorityService) VerifyAuthority(authority string, tokenStr string) bool {
	admin, err := AdminService{}.GetAdminByToken(tokenStr)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, v := range admin.Authorities {
		if v.AuthorityName == authority || v.AuthorityName == "root" {
			return true
		}
	}
	return false
}

// GetAuthorityByID 根据权限ID获取权限
func (receiver AuthorityService) GetAuthorityByID(id uint) (authority models.Authority, err error) {
	result := global.DB.Where("id = ?", id).First(&authority)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取权限信息失败，%v", result.Error))
		return
	}
	return
}

// GetAllAuthority 获取所有的权限
func (receiver AuthorityService) GetAllAuthority(page int, pageSize int) (authorities []models.Authority, total int64, err error) {
	global.DB.Model(&models.Authority{}).Count(&total)
	result := global.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&authorities)
	if result.Error != nil {
		err = errors.New(fmt.Sprintf("获取权限列表失败，%v", result.Error))
		return
	}
	return
}
