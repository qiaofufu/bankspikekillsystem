package models

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/utils"
	"gorm.io/gorm"
	"time"
)

type RealInformation struct {
	RealName string `gorm:""` // 真实名字
	IDCard   string `gorm:""` // 身份证
}

type ExtendInformation struct {
}

type User struct {
	global.Model
	UUID              string          `json:"UUID" gorm:"not null"`
	PhoneNumber       string          `gorm:"not null;size:20" json:"phone_number"`            // 手机号
	Password          string          `gorm:"not null;size:100" json:"-" swaggerignore:"true"` // 密码
	LoginTime         time.Time       `gorm:"autoUpdateTime" json:"login_time"`                // 登陆时间
	LoginIP           string          `gorm:"size:20" json:"login_ip"`                         // 登录ip
	Status            int             `gorm:"default:1;not null" json:"status"`                // 账号状态 1:激活 -1:冻结
	NickName          string          `json:"nick_name"`                                       //昵称
	Gender            int             `gorm:"default:2" json:"gender"`                         // 性别 1:女 2:男
	Age               int             `gorm:"" json:"age"`                                     // 年龄
	ProfilePictureUrl string          `gorm:"default:" json:"profile_picture_url"`             // 头像url
	Point             int             `gorm:"default:0" json:"point"`                          // 用户积分
	RealInformation   RealInformation `gorm:"embedded"`                                        // 真实信息
	Address           Address         `gorm:"embedded" json:"address"`                         // 地址
	Occupation        string          `gorm:"default:未知" json:"occupation"`                    // 职业
	Activity          []SpikeActivity `gorm:"many2many:activity_user"`
}

func (receiver *User) BeforeCreate(db *gorm.DB) (err error) {
	receiver.UUID, err = utils.GetUUID()
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SignatureToken() (string, error) {
	return global.GenerateUserToken(map[string]interface{}{
		"uuid":    user.UUID,
		"user_id": user.ID,
	})
}
