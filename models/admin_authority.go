package models

type AdminAuthority struct {
	AdminID     uint `json:"admin_id"`     // 管理员id
	AuthorityID uint `json:"authority_id"` // 权限id
}

func (receiver AdminAuthority) TableName() string {
	return "admin_authority"
}
