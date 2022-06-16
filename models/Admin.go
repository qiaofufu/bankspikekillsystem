package models

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/utils"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	global.Model
	UUID             string      `gorm:"not null" json:"uuid"`                                // UUID
	CreateByUsername string      `gorm:"default:'admin';not null;" json:"create_by_username"` // 创建者名称
	Username         string      `gorm:"not null" json:"username"`                            // 用户名
	Password         string      `gorm:"not null;size:100" json:"-" swaggerignore:"true"`     // 密码
	LoginTime        time.Time   `gorm:"autoUpdateTime" json:"login_time"`                    // 登陆时间
	LoginIP          string      `gorm:"size:20" json:"login_ip"`                             // 登录ip
	Status           int         `gorm:"default:0;not null" json:"status"`                    // 账号状态
	NickName         string      `gorm:"default:nick_name" json:"nick_name"`                  // 昵称
	Email            string      `gorm:"size:50;" json:"email"`                               // 邮箱
	Phone            string      `gorm:"size:20;" json:"phone"`                               // 电话
	Authorities      []Authority `gorm:"many2many:admin_authority;" json:"authorities"`       // 权限

}

func (receiver *Admin) BeforeCreate(db *gorm.DB) (err error) {
	receiver.UUID, err = utils.GetUUID()
	if err != nil {
		return err
	}
	return nil
}

// SignatureToken 签发Token
func (admin *Admin) SignatureToken() (string, error) {

	return global.GenerateUserToken(map[string]interface{}{
		"username":  admin.Username,
		"id":        admin.Model.ID,
		"uuid":      admin.UUID,
		"authority": admin.Authorities,
	})
}
