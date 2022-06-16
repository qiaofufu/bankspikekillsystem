package v1

import (
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"bank-product-spike-system/utils"
	"bank-product-spike-system/utils/captcha"
	"github.com/gin-gonic/gin"
)

type AdminAPI struct{}

// Login
// @Summary 后台登录
// @Tags 后台
// @Accept json
// @Param data body request.Login true "请求参数"
// @Product json
// @Success 200 {object} response.DTO{code=int,data=response.LoginResponse,msg=string} "成功登录"
// @Router /admin/login [post]
func (receiver AdminAPI) Login(ctx *gin.Context) {
	var loginDTO request.Login
	if err := ctx.ShouldBindJSON(&loginDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}
	if !captcha.VerifyCaptcha(loginDTO.VerifyCodeID, loginDTO.VerifyCode) {
		response.AuthFail("验证码错误，请重新验证", ctx)
		return
	}

	ip := ctx.ClientIP()
	admin := models.Admin{
		Username: loginDTO.Username,
		Password: loginDTO.Password,
		LoginIP:  ip,
	}

	returnAdmin, err := AdminService.Login(&admin)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	token, err := returnAdmin.SignatureToken()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	responseLogin := response.LoginResponse{
		User:  *returnAdmin,
		Token: token,
	}

	response.Success(responseLogin, "登录成功", ctx)
}

// GetInformation
// @Summary 获取个人信息
// @Tags 后台
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=models.Admin,msg=string} "获取成功"
// @Router /admin/getInformation [get]
func (receiver AdminAPI) GetInformation(ctx *gin.Context) {
	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	admin, err := AdminService.GetAdminByToken(token)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	response.Success(admin, "获取成功", ctx)
}

// ChangePassword
// @Summary 更改密码
// @Tags 管理员管理
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ChangePassword true "请求参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "成功修改密码"
// @Router /admin/changePassword [put]
func (receiver AdminAPI) ChangePassword(ctx *gin.Context) {
	var changePasswordDTO request.ChangePassword
	if err := ctx.ShouldBindJSON(&changePasswordDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.AdminManager, token) {
		response.AuthFail("非法操作，没有操作权限，行为已记录.", ctx)
		return
	}

	admin := models.Admin{
		Password: changePasswordDTO.Password,
		UUID:     changePasswordDTO.UUID,
	}
	_, err := AdminService.ChangePassword(&admin, changePasswordDTO.NewPassword)
	if err != nil {
		response.FailWithMessage("更改密码失败"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("更改密码成功", ctx)
}

// ChangeSelfPassword
// @Summary 更改自身密码
// @Tags 后台
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ChangeSelfPassword true "请求参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "成功修改密码"
// @Router /admin/changeSelfPassword [put]
func (receiver AdminAPI) ChangeSelfPassword(ctx *gin.Context) {
	var requestDTO request.ChangeSelfPassword
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	uuid, err := utils.GetUUIDByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	admin := models.Admin{
		UUID:     uuid.(string),
		Password: requestDTO.Password,
	}
	_, err = AdminService.ChangePassword(&admin, requestDTO.NewPassword)
	if err != nil {
		response.FailWithMessage("更改密码失败，"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("更改密码成功", ctx)
}

// GetAdminInfoList
// @Summary 获取管理员分页信息
// @Tags 管理员管理
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.PageInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.PageInfoResponse{total=int,data=[]models.Admin}} "成功获取数据"
// @Router /admin/getAdminInfoList [post]
func (receiver AdminAPI) GetAdminInfoList(ctx *gin.Context) {
	var pageInfo request.PageInfo
	if err := ctx.ShouldBindJSON(&pageInfo); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.AdminManager, token) {
		response.AuthFail("非法操作，没有操作权限，行为已记录.", ctx)
		return
	}

	adminList, total, err := AdminService.GetAdminInfoList(pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	pageInfoResponse := response.PageInfoResponse{
		Total: total,
		Data:  adminList,
	}
	response.Success(pageInfoResponse, "获取分页信息成功", ctx)
}

// GenerateAdminAccount
// @Summary 生成管理员账号
// @Tags 管理员管理
// @Security ApiKeyAuth
// @Accept json
// @Product json
// @Param data body request.GenerateAdminAccount true "传入参数"
// @Success 200 {object} response.DTO{data=response.GenerateAdminAccountResponse,msg=string,code=int} "成功生成账号"
// @Router /admin/generateAdminAccount [post]
func (receiver AdminAPI) GenerateAdminAccount(ctx *gin.Context) {
	var generateAdminAccountDTO request.GenerateAdminAccount
	if err := ctx.ShouldBindJSON(&generateAdminAccountDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.AdminManager, token) {
		response.AuthFail("非法操作，没有操作权限，行为已记录.", ctx)
		return
	}

	admin, err := AdminService.GetAdminByToken(token)
	if err != nil {
		response.AuthFail("用户信息获取失败， 请重新登录", ctx)
		return
	}

	additionalAdmin := &models.Admin{
		Username:         generateAdminAccountDTO.Username,
		Password:         generateAdminAccountDTO.Password,
		Status:           generateAdminAccountDTO.Status,
		Email:            generateAdminAccountDTO.Email,
		Phone:            generateAdminAccountDTO.Phone,
		CreateByUsername: admin.Username,
	}
	returnAdmin, err := AdminService.GenerateAdminAccount(additionalAdmin)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	generateAdminAccountResponse := response.GenerateAdminAccountResponse{
		User: *returnAdmin,
	}
	response.Success(generateAdminAccountResponse, "生成成功", ctx)
}

// SetAdminAuthorities
// @Summary 设置管理员权限
// @Tags 管理员管理
// @Security ApiKeyAuth
// @Accept json
// @Product json
// @Param data body request.SetAdminAuthorities true "传入参数"
// @Success 200 {object} response.DTO{data=models.Admin} "成功设置权限"
// @Router /admin/setAdminAuthorities [put]
func (receiver AdminAPI) SetAdminAuthorities(ctx *gin.Context) {
	var requestDTO request.SetAdminAuthorities
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.AdminManager, token) {
		response.FailWithMessage("没有操作权限，非法访问", ctx)
		return
	}

	admin := models.Admin{
		UUID:        requestDTO.UUID,
		Authorities: getAuthority(requestDTO.AuthorityID),
	}
	err := AdminService.SetAdminAuthorities(&admin) //

	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(admin, "设置成功", ctx)
}

// SetAdminInfo
// @Summary 设置管理员信息
// @Tags 管理员管理
// @Security ApiKeyAuth
// @Accept json
// @Product json
// @Param data body request.SetAdminInfo true "传入参数"
// @Success 200 {object} response.DTO "成功设置信息"
// @Router /admin/setAdminInfo [put]
func (receiver AdminAPI) SetAdminInfo(ctx *gin.Context) {
	var info request.SetAdminInfo
	if err := ctx.ShouldBindJSON(&info); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.AdminManager, token) {
		response.AuthFail("没有操作权限", ctx)
		return
	}

	admin := models.Admin{
		UUID:     info.UUID,
		Status:   info.Status,
		NickName: info.NickName,
		Email:    info.Email,
		Phone:    info.Phone,
	}
	err := AdminService.SetAdminInfo(&admin)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("设置成功", ctx)
}

// SetSelfInfo
// @Summary 设置个人信息
// @Tags 后台
// @Security ApiKeyAuth
// @Accept json
// @Product json
// @Param data body request.SetSelfInfo true "传入参数"
// @Success 200 {object} response.DTO "成功设置"
// @Router /admin/setSelfInfo [put]
func (receiver AdminAPI) SetSelfInfo(ctx *gin.Context) {
	var info request.SetSelfInfo
	if err := ctx.ShouldBindJSON(&info); err != nil {
		response.BindJSONError(err, ctx)
	}
	token, _ := getToken(ctx)
	uuid, err := utils.GetUUIDByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	admin := models.Admin{
		UUID:     uuid.(string),
		NickName: info.NickName,
		Email:    info.Email,
		Phone:    info.Phone,
	}
	err = AdminService.SetAdminInfo(&admin)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("设置成功", ctx)
	return
}

// GetAllAuthority
// @Summary 获取所有权限列表
// @Tags 管理员管理
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.PageInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.AuthorityList,msg=string} "获取成功"
// @Router /admin/authority/all [post]
func (receiver AdminAPI) GetAllAuthority(ctx *gin.Context) {
	var requestDTO request.PageInfo
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.AdminManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	authorityArr, cnt, err := AuthorityService.GetAllAuthority(requestDTO.Page, requestDTO.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	responseDto := response.AuthorityList{
		Total: cnt,
		List:  authorityArr,
	}
	response.Success(responseDto, "获取成功", ctx)
}
