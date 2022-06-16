package request

type UserLogin struct {
	PhoneNumber  string `json:"phone_number" binding:"required,alphanum"` // 手机号
	Password     string `json:"password"`                                 // 密码
	VerifyCode   string `json:"verify_code" binding:"required"`           // 验证码
	VerifyCodeID string `json:"verify_code_id"`                           // 验证码ID
	VerifyType   int    `json:"verify_type" binding:""`                   // 验证码类型 0: 图片验证码 1：手机验证码
}

type UserRegister struct {
	PhoneNumber string `json:"phone_number" binding:"required"` // 手机号 带地区号 例如+8618547304726
	VerifyCode  string `json:"verify_code" binding:"required"`  // 手机验证码
	Password    string `json:"password" binding:"required"`     // 密码
	RealName    string `json:"real_name" binding:"required"`    // 真实名字
	IDCard      string `json:"id_card" binding:"required"`      // 身份证
}

type GetUserList struct {
	PageInfo
	Age    int `json:"age"`    // 年龄 为空不进行筛选
	Gender int `json:"gender"` // 性别 1女 2男 0不进行筛选
	Status int `json:"status"` // 状态 1-激活 -1封禁
}

type UpdateUserInfo struct {
	UserUUID   string `json:"user_uuid"`  // UUID
	NickName   string `json:"nick_name"`  // 昵称
	Occupation string `json:"occupation"` // 职业
	Status     int    `json:"status"`     // 状态 1-激活 -1封禁
}

type ChangePhone struct {
	NewPhone   string `json:"new_phone" binding:"required,alphanum"` // 新手机号 例如+8618547304726
	VerifyCode string `json:"verify_code" binding:"required"`        // 手机验证码
}

type UpdateUserSelfInfo struct {
	NickName   string `json:"nick_name"`  // 昵称
	Occupation string `json:"occupation"` // 职业
}
