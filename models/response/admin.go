package response

import "bank-product-spike-system/models"

type LoginResponse struct {
	User  models.Admin `json:"user"`  // 用户信息
	Token string       `json:"token"` // JWT
}

type GenerateAdminAccountResponse struct {
	User models.Admin `json:"user"` // 用户信息

}

type AuthorityList struct {
	Total int64              // 总数
	List  []models.Authority // 权限列表
}
