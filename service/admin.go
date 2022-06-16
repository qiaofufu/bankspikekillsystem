package service

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"

	"bank-product-spike-system/utils"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AdminService struct{}

// Login
// 登录服务
func (receiver AdminService) Login(admin *models.Admin) (*models.Admin, error) {
	admin.Password = utils.MD5(admin.Password)
	ip := admin.LoginIP
	result := global.DB.Where("username = ? and password = ?", admin.Username, admin.Password).Preload("Authorities").First(&admin)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return nil, errors.New("账号或密码错误")
	}
	admin.LoginIP = ip
	global.DB.Save(admin)
	return admin, nil

}

// ChangePassword
// 更改密码服务
func (receiver AdminService) ChangePassword(admin *models.Admin, newPassword string) (*models.Admin, error) {
	password := utils.MD5(admin.Password)
	result := global.DB.Where("uuid = ? and password = ?", admin.UUID, password).Preload("Authorities").First(&admin)
	if result.Error != nil {
		return nil, errors.New("账号或密码错误")
	}
	admin.Password = utils.MD5(newPassword)
	result = global.DB.Save(admin)
	if result.Error != nil {
		return nil, errors.New("内部错误更新密码失败")
	}
	return admin, nil
}

// GenerateAdminAccount
// 生成管理员账号服务
func (receiver AdminService) GenerateAdminAccount(admin *models.Admin) (*models.Admin, error) {
	admin.Password = utils.MD5(admin.Password)
	if err := global.DB.Where("username = ?", admin.Username).First(&admin).Error; err != gorm.ErrRecordNotFound {
		return nil, errors.New("用户名已存在，请重新设置")
	}
	result := global.DB.Create(admin)
	if result.Error != nil {
		return nil, errors.New("生成管理员账号失败 " + result.Error.Error())
	}
	return admin, nil
}

// GetAdminInfoList
// 获取管理员信息列表服务
func (receiver AdminService) GetAdminInfoList(page int, pageSize int) (list []models.Admin, total int64, err error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var adminList []models.Admin
	result := global.DB.Model(&models.Admin{}).Count(&total)
	if result.Error != nil {
		return nil, 0, errors.New("内部错误：获取数据数量失败")
	}
	result = result.Offset(offset).Limit(limit).Preload("Authorities").Find(&adminList)
	if result.Error != nil {
		return nil, 0, errors.New("内部错误：获取分页数据失败")
	}
	return adminList, total, nil
}

// GetAdminByToken
// 通过Token获取Admin
func (receiver AdminService) GetAdminByToken(tokenStr string) (*models.Admin, error) {
	token, err := utils.GetTokenPayload(tokenStr)
	if err != nil {
		return nil, err
	}

	username, ok := token.Get("username")
	if !ok {
		err = errors.New("token信息不完整，请重新登录")
		return nil, err
	}
	id, ok := token.Get("id")
	if !ok {
		err = errors.New("token信息不完整，请重新登录")
		return nil, err
	}
	uuid, ok := token.Get("uuid")
	if !ok {
		err = errors.New("token信息不完整，请重新登录")
		return nil, err
	}

	var admin models.Admin
	result := global.DB.Where("username = ? and id = ? and uuid = ?", username, id, uuid).Preload("Authorities").First(&admin)
	if result.Error != nil {
		return nil, err
	}
	return &admin, nil
}

// SetAdminAuthorities
// 设置admin权限
func (receiver AdminService) SetAdminAuthorities(admin *models.Admin) (err error) {

	result := global.DB.Where("uuid = ? ", admin.UUID).First(&admin)
	if result.Error != nil {
		err = errors.New("用户信息不存在")
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("admin_id = ?", admin.ID).Delete(&models.AdminAuthority{})
		if result.Error != nil {
			return result.Error
		}

		var auth []models.AdminAuthority

		for _, v := range admin.Authorities {
			if err := global.DB.Model(&models.Authority{}).Where("id = ?", v.ID).First(&models.Authority{}).Error; err != nil {
				return errors.New("非法操作，添加未存在权限")
			}
			auth = append(auth, models.AdminAuthority{
				AdminID:     admin.ID,
				AuthorityID: v.ID,
			})
		}

		result = tx.Create(&auth)
		if result.Error != nil {
			return errors.New("添加权限失败")
		}

		result = tx.Where("id = ?", admin.ID).Preload("Authorities").First(&admin)
		if result.Error != nil {
			return errors.New("获取信息失败")
		}
		return nil
	})

	return err
}

// SetAdminInfo
// 设置admin信息
func (receiver AdminService) SetAdminInfo(admin *models.Admin) error {
	result := global.DB.Where("uuid = ?", admin.UUID).Updates(admin)
	if result.Error != nil {
		return errors.New("更新信息失败")
	}
	return nil
}

// relatedUpdate 关联更新， 只更新非空字段
func (receiver AdminService) relatedUpdate(admin *models.Admin) (err error) {
	result := global.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&admin)
	return result.Error
}
