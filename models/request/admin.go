package request

type Login struct {
	Username     string `json:"username" binding:"required,alphanum"`      // 用户名
	Password     string `json:"password" binding:"required,min=8,max=100"` // 密码
	VerifyCode   string `json:"verify_code" binding:"required"`            // 验证码
	VerifyCodeID string `json:"verify_code_id" binding:"required"`         // 验证码ID
}

type ChangePassword struct {
	UUID        string `json:"uuid" binding:"required,uuid4"`                                 // 用户UUID
	Password    string `json:"password" binding:"required,nefield=NewPassword,min=8,max=100"` // 旧密码
	NewPassword string `json:"newPassword" binding:"required,min=8,max=100"`                  // 新密码
}

type ChangeSelfPassword struct {
	Password    string `json:"password" binding:"required,nefield=NewPassword,min=8,max=100"` // 旧密码
	NewPassword string `json:"new_password" binding:"required,min=8,max=100"`                 // 新密码
}

type GenerateAdminAccount struct {
	Username         string `json:"username" binding:"required,alphanum"`   // 用户名
	Password         string `json:"password" binding:"required"`            // 密码
	IsRandomPassword bool   `json:"is_random_password"`                     // 是否随机密码
	Status           int    `json:"status"`                                 // 状态 0激活 1封禁
	Email            string `json:"email" binding:"required,email"`         // 邮箱
	Phone            string `json:"phone" binding:"required,max=11,min=11"` // 电话
}

type SetAdminAuthorities struct {
	UUID        string `json:"uuid" binding:"required,uuid4"`   // 用户UUID
	AuthorityID []uint `json:"authority_id" binding:"required"` // 新设置的权限ID
}

type SetAdminInfo struct {
	UUID     string `json:"UUID" binding:"uuid4"`            // UUID
	NickName string `json:"nick_name" binding:"omitempty"`   // 昵称
	Email    string `json:"email" binding:"omitempty,email"` // 邮箱
	Phone    string `json:"phone" binding:"omitempty"`       // 头像
	Status   int    `json:"status" binding:"omitempty"`      // 用户状态 0激活 1封禁
}

type SetSelfInfo struct {
	NickName string `json:"nick_name" binding:"omitempty"`   // 昵称
	Email    string `json:"email" binding:"omitempty,email"` // 邮箱
	Phone    string `json:"phone" binding:"omitempty"`       // 头像
}
