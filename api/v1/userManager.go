package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"github.com/gin-gonic/gin"
)

// UpdateUserInfo
// @Summary 更新用户信息
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.UpdateUserInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更新成功"
// @Router /user/update [put]
func (receiver UserAPI) UpdateUserInfo(ctx *gin.Context) {
	var requestDTO request.UpdateUserInfo
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.UserManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	id, err := UserService.GetUserIDByUUID(requestDTO.UserUUID)
	user := models.User{
		Model:      global.Model{ID: id},
		NickName:   requestDTO.NickName,
		Occupation: requestDTO.Occupation,
		Status:     requestDTO.Status,
	}

	err = UserService.UpdateUserInfo(user)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("更新成功", ctx)
}

// GetUserList
// @Summary 获取用户列表
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetUserList true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.GetUserList,msg=string} "获取成功"
// @Router /user/list [post]
func (receiver UserAPI) GetUserList(ctx *gin.Context) {
	var requestDTO request.GetUserList
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.UserManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	res, total, err := UserService.GetUserList(requestDTO.Page, requestDTO.PageSize, requestDTO.Age, requestDTO.Gender, requestDTO.Status)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	list := response.GetUserList{
		Total: total,
	}

	for _, v := range res {
		elem := response.UserBaseInfo{
			UserID:            v.ID,
			UserUUID:          v.UUID,
			PhoneNumber:       v.PhoneNumber,
			LoginTime:         v.LoginTime,
			LoginIP:           v.LoginIP,
			Status:            v.Status,
			NickName:          v.NickName,
			Gender:            v.Gender,
			Age:               v.Age,
			ProfilePictureUrl: v.ProfilePictureUrl,
			Occupation:        v.Occupation,
		}
		list.List = append(list.List, elem)
	}

	response.Success(list, "获取成功", ctx)
}
