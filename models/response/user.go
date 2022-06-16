package response

import (
	"bank-product-spike-system/models"
	"time"
)

type UserLoginResponse struct {
	//UUID              string    `json:"UUID" gorm:"not null"`
	//PhoneNumber       string    `gorm:"not null;size:20" json:"phone_number"`
	//LoginTime         time.Time `gorm:"autoUpdateTime" json:"login_time"`
	//LoginIP           string    `gorm:"size:20" json:"login_ip"`
	//Status            int       `gorm:"default:0;not null" json:"status"`
	//NickName          string    `json:"nick_name"`
	//Gender            int       `gorm:"default:2" json:"gender"`
	//Age               int       `gorm:"" json:"age"`
	//ProfilePictureUrl string    `gorm:"default:" json:"profile_picture_url"`
	//Point             int       `gorm:"default:0" json:"point"`
	////Address           Address   `gorm:"embedded" json:"address"`
	//Occupation string `gorm:"default:未知" json:"occupation"`
	User  models.User
	Token string `json:"token"` // Token
}

type UserBaseInfo struct {
	UserID            uint      `json:"user_id"`             // 用户id
	UserUUID          string    `json:"user_uuid"`           // uuid
	PhoneNumber       string    `json:"phone_number"`        // 手机号
	LoginTime         time.Time `json:"login_time"`          // 登陆时间
	LoginIP           string    `json:"login_ip"`            // 登录ip
	Status            int       `json:"status"`              // 账号状态 1:激活 -1:冻结
	NickName          string    `json:"nick_name"`           // 昵称
	Gender            int       `json:"gender"`              // 性别 1:女 2:男
	Age               int       `json:"age"`                 // 年龄
	ProfilePictureUrl string    `json:"profile_picture_url"` // 头像url
	Occupation        string    `json:"occupation"`          // 职业
}

type GetUserList struct {
	Total int64 `json:"total"` // 总数
	List  []UserBaseInfo
}
