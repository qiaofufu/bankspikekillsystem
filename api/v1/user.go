package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"bank-product-spike-system/utils"
	"bank-product-spike-system/utils/captcha"
	"github.com/gin-gonic/gin"
)

type UserAPI struct{}

// Register
// @Summary 用户注册
// @Tags 用户
// @Accept json
// @Product json
// @Param data body request.UserRegister true "传入参数"
// @Success 200 {object} response.DTO{} "注册成功"
// @Router /user/register [post]
func (receiver UserAPI) Register(ctx *gin.Context) {
	var register request.UserRegister
	if err := ctx.ShouldBindJSON(&register); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	if !utils.VerifySMS(register.PhoneNumber, register.VerifyCode) {
		response.FailWithMessage("验证码错误， 请重新输入。", ctx)
		return
	}

	err := UserService.Register(register.PhoneNumber, register.Password, register.RealName, register.IDCard)
	if err != nil {
		response.FailWithMessage("注册失败！"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("注册成功", ctx)
}

// Login
// @Summary 用户登录
// @Tags 用户
// @Accept json
// @Product json
// @Param data body request.UserLogin true "传入参数"
// @Success 200 {object} response.DTO{data=response.UserLoginResponse} "登录成功"
// @Router /user/login [post]
func (receiver UserAPI) Login(ctx *gin.Context) {
	var userLogin request.UserLogin
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	var result bool
	switch userLogin.VerifyType {
	case 0:
		result = captcha.VerifyCaptcha(userLogin.VerifyCodeID, userLogin.VerifyCode)
	case 1:
		result = utils.VerifySMS(userLogin.PhoneNumber, userLogin.VerifyCode)
	}

	if result == false {
		response.FailWithMessage("验证码错误，请重新输入。", ctx)
		return
	}

	ip := ctx.ClientIP()
	user := models.User{
		PhoneNumber: userLogin.PhoneNumber,
		Password:    userLogin.Password,
		LoginIP:     ip,
	}
	var err error
	switch userLogin.VerifyType {
	case 0:
		_, err = UserService.LoginByPassword(&user)
	case 1:
		_, err = UserService.LoginByPhone(&user)
	}
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	token, err := user.SignatureToken()
	if err != nil {
		response.FailWithMessage("token签发失败", ctx)
		return
	}
	res := response.UserLoginResponse{
		User:  user,
		Token: token,
	}
	response.Success(res, "登录成功", ctx)
}

// ChangeSelfPassword
// @Summary 更改自身密码
// @Tags 用户
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ChangeSelfPassword true "请求参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "成功修改密码"
// @Router /user/changeSelfPassword [put]
func (receiver UserAPI) ChangeSelfPassword(ctx *gin.Context) {
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
	user := models.User{
		UUID:     uuid.(string),
		Password: requestDTO.Password,
	}
	_, err = UserService.ChangePassword(&user, requestDTO.NewPassword)
	if err != nil {
		response.FailWithMessage("更改密码失败，"+err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("更改密码成功", ctx)
}

// ChangePhone
// @Summary 换绑手机
// @Tags 用户
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ChangePhone true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更换成功"
// @Router /user/changePhone [put]
func (receiver UserAPI) ChangePhone(ctx *gin.Context) {
	var requestDTO request.ChangePhone
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	if !utils.VerifySMS(requestDTO.NewPhone, requestDTO.VerifyCode) {
		response.FailWithMessage("验证码错误， 请重新获取。", ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	u, err := UserService.GetUserByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	user := models.User{
		Model:       global.Model{ID: u.ID},
		PhoneNumber: requestDTO.NewPhone,
	}
	err = UserService.UpdateUserInfo(user)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("更换成功，请重新登陆", ctx)
}

// GetInformation
// @Summary 获取个人信息
// @Tags 用户
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=models.User,msg=string} "获取成功"
// @Router /user/getInformation [get]
func (receiver UserAPI) GetInformation(ctx *gin.Context) {
	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	user, err := UserService.GetUserByToken(token)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	response.Success(user, "获取成功", ctx)
}

// UpdateUserSelfInfo
// @Summary 更新用户自身信息
// @Tags 用户
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.UpdateUserSelfInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更新成功"
// @Router /user/updateSelfInfo [put]
func (receiver UserAPI) UpdateUserSelfInfo(ctx *gin.Context) {
	var requestDTO request.UpdateUserSelfInfo
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	tokenStr, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	token, err := JWTService.ParseToken(tokenStr)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	userID, _ := token.Get("user_id")

	user := models.User{
		Model:      global.Model{ID: uint(userID.(float64))},
		NickName:   requestDTO.NickName,
		Occupation: requestDTO.Occupation,
	}
	err = UserService.UpdateUserInfo(user)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("更新成功", ctx)
}
