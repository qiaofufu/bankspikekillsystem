package v1

import (
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"errors"
	"github.com/gin-gonic/gin"
)

// GetActivity
// @Summary 获取当天活动列表
// @Tags 活动-前台
// @Product json
// @Security ApiKeyAuth
// @Success 200 {object} response.DTO{data=response.GetActivityList,msg=string} "获取成功"
// @Router /spikeActivity/getActivityList [get]
func (receiver SpikeActivityAPI) GetActivity(ctx *gin.Context) {

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	List1, err := SpikeActivityService.GetRunningActivity()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	var list response.GetActivityList
	for _, v := range List1 {
		if err1 := FilterService.Check(token, v.ID); err1 != nil {
			continue
		}
		elem := response.ActivityBaseInfo{
			ActivityName:         v.ActivityName,
			ActivityID:           v.ID,
			StartTime:            v.StartTime,
			NumberOfParticipants: v.NumberOfParticipants,
			NumberOfFavorites:    v.NumberOfFavorites,
			ProductTotal:         v.Product.ProductNumber,
			ProductNowNumber:     v.Product.ProductNumber - v.Product.SoldNumber,
			ProductRate:          v.Product.ProductInterestRate,
		}
		list.RunningActivity = append(list.RunningActivity, elem)
	}

	activityList, err := SpikeActivityService.GetUpcomingActivity()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	for _, v := range activityList {
		if err1 := FilterService.Check(token, v.ID); err1 != nil {
			continue
		}
		elem := response.ActivityBaseInfo{
			ActivityName:         v.ActivityName,
			ActivityID:           v.ID,
			StartTime:            v.StartTime,
			NumberOfParticipants: v.NumberOfParticipants,
			NumberOfFavorites:    v.NumberOfFavorites,
			ProductTotal:         v.Product.ProductNumber,
			ProductNowNumber:     v.Product.ProductNumber - v.Product.SoldNumber,
			ProductRate:          v.Product.ProductInterestRate,
		}
		list.UpcomingActivity = append(list.UpcomingActivity, elem)
	}
	response.Success(list, "获取成功", ctx)
}

// GetActivityGuidanceInformation
// @Summary 获取活动相关引导信息
// @Tags 活动-前台
// @Accept json
// @Security ApiKeyAuth
// @Param data body request.GetActivityGuidanceInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.GetActivityGuidanceInformation,msg=string} "获取成功"
// @Router /spikeActivity/getActivityGuidanceInformation [post]
func (receiver SpikeActivityAPI) GetActivityGuidanceInformation(ctx *gin.Context) {
	var requestDTO request.GetActivityGuidanceInfo
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	baseInfo, err := SpikeActivityService.GetActivityInfoByID(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	guidanceList, err := SpikeActivityService.GetActivityGuidanceInformation(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	info := response.GetActivityGuidanceInformation{
		ActivityID:           baseInfo.ID,
		ActivityName:         baseInfo.ActivityName,
		DescriptionText:      baseInfo.DescriptionText,
		ActivityStatus:       baseInfo.ActivityStatus,
		StartTime:            baseInfo.StartTime,
		EndingTime:           baseInfo.EndingTime,
		NumberOfParticipants: baseInfo.NumberOfParticipants,
		NumberOfFavorites:    baseInfo.NumberOfFavorites,
		AdminID:              baseInfo.AdminID,
		AdminName:            baseInfo.Admin.NickName,
		Product:              baseInfo.Product,
		GuidanceInformation:  guidanceList,
	}

	response.Success(info, "获取成功", ctx)
}

// Attention
// @Summary 关注活动
// @Tags 活动-前台
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.Attention true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "关注成功"
// @Router /spikeActivity/attention [post]
func (receiver SpikeActivityAPI) Attention(ctx *gin.Context) {
	var requestDTO request.Attention
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	user, err := UserService.GetUserByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	att := models.ActivityAttention{
		ActivityID: requestDTO.ActivityID,
		UserID:     user.ID,
	}
	err = SpikeActivityService.Attention(att)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("关注成功", ctx)
}

// UnAttention
// @Summary 取消关注活动
// @Tags 活动-前台
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.Attention true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "取消成功"
// @Router /spikeActivity/attention [delete]
func (receiver SpikeActivityAPI) UnAttention(ctx *gin.Context) {
	var requestDTO request.Attention
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
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	userID, _ := token.Get("user_id")

	err = SpikeActivityService.UnAttention(requestDTO.ActivityID, uint(userID.(float64)))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("取消关注成功", ctx)
}

// GetAttentionList
// @Summary 获取关注列表
// @Tags 活动-前台
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=response.AttentionList,msg=string} "获取成功"
// @Router /spikeActivity/attention [get]
func (receiver SpikeActivityAPI) GetAttentionList(ctx *gin.Context) {

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	user, err := UserService.GetUserByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	list, err := SpikeActivityService.GetAttentionListByUserID(user.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	var dto response.AttentionList
	for _, v := range list {
		act, err := SpikeActivityService.GetActivityBaseInfoByID(v.ActivityID)
		if err != nil {
			response.FailWithMessage(err.Error(), ctx)
			return
		}
		elem := response.ActivityBaseInfo{
			ActivityName:         act.ActivityName,
			ActivityID:           v.ActivityID,
			StartTime:            act.StartTime,
			NumberOfParticipants: act.NumberOfParticipants,
			NumberOfFavorites:    act.NumberOfFavorites,
			ProductTotal:         act.Product.ProductNumber,
			ProductNowNumber:     act.Product.ProductNumber - act.Product.SoldNumber,
		}
		dto.List = append(dto.List, elem)
	}

	response.Success(dto, "获取成功", ctx)
}

// IsAttention
// @Summary 是否关注
// @Tags 活动-前台
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ActivityID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.IsAttention,msg=string} "判定成功"
// @Router /spikeActivity/isAttention [post]
func (receiver SpikeActivityAPI) IsAttention(ctx *gin.Context) {
	var requestDTO request.ActivityID
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
		response.FailWithMessage("Failed to Parse Token", ctx)
	}

	userID, ok := token.Get("user_id")
	if !ok {
		err = errors.New("token parse failed")
		if err != nil {
			response.FailWithMessage("Failed to get userID", ctx)
		}
	}

	isAttention := SpikeActivityService.IsAttention(requestDTO.ActivityID, uint(userID.(float64)))

	dto := response.IsAttention{
		IsAttention: isAttention,
	}
	response.Success(dto, "判定成功", ctx)
}

// ParticipateActivity
// @Summary 参与秒杀活动
// @Tags 秒杀
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ParticipateActivity true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "参与成功"
// @Router /spikeActivity/participate/:verifyCode [post]
func (receiver SpikeActivityAPI) ParticipateActivity(ctx *gin.Context) {
	// 数据绑定
	var requestDTO request.ParticipateActivity
	verifyCode := ctx.Param("verifyCode")
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
	userID, ok := token.Get("user_id")
	if !ok {
		response.AuthFail("token 信息 携带错误", ctx)
		return
	}
	// 资格验证
	if err := SpikeActivityService.CheckEligibility(requestDTO.ActivityID, requestDTO.ProductID, uint(userID.(float64)), verifyCode, 1); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	// 参与秒杀
	if err := SpikeActivityService.ParticipateSpike(requestDTO.ActivityID, requestDTO.ProductID, uint(userID.(float64))); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("参与成功", ctx)
}

func (receiver SpikeActivityAPI) ParticipateActivityTest(ctx *gin.Context) {
	var requestDTO request.ParticipateActivityTest
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}
	// 数据绑定
	//var requestDTO request.ParticipateActivity
	//verifyCode := ctx.Param("verifyCode")
	//if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
	//	response.BindJSONError(err, ctx)
	//	return
	//}
	//tokenStr, err := getToken(ctx)
	//if err != nil {
	//	response.AuthFail(err.Error(), ctx)
	//	return
	//}
	//
	//token, err := JWTService.ParseToken(tokenStr)
	//if err != nil {
	//	response.AuthFail(err.Error(), ctx)
	//	return
	//}
	//userID, ok := token.Get("user_id")
	//if !ok {
	//	response.AuthFail("token 信息 携带错误", ctx)
	//	return
	//}
	// 资格验证
	if err := SpikeActivityService.CheckEligibility(requestDTO.ActivityID, requestDTO.ProductID, 0, "", 0); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	// 参与秒杀
	if err := SpikeActivityService.ParticipateSpike(requestDTO.ActivityID, requestDTO.ProductID, requestDTO.UserID); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("参与成功", ctx)
}

// GetSpikeResult
// @Summary 获取秒杀结果
// @Tags 秒杀
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.QuerySpikeResult true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.SpikeResult,msg=string} "获取成功"
// @Router /spikeActivity/getSpikeResult [post]
func (receiver SpikeActivityAPI) GetSpikeResult(ctx *gin.Context) {
	var requestDTO request.QuerySpikeResult
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
		response.FailWithMessage("Failed to Parse Token", ctx)
		return
	}

	userID, ok := token.Get("user_id")
	if !ok {
		err = errors.New("token parse failed")
		if err != nil {
			response.FailWithMessage("Failed to get userID", ctx)
			return
		}
	}

	dto, err := SpikeActivityService.QuerySpikeResult(requestDTO.ActivityID, requestDTO.ProductID, uint(userID.(float64)))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(dto, "获取成功", ctx)
}

// GetSpikeVerifyCode
// @Summary 获取秒杀验证码
// @Tags 秒杀
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetSpikeVerifyCode true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.GetSpikeVerifyCode,msg=string} "获取成功"
// @Router /spikeActivity/getSpikeVerifyCode [post]
func (receiver SpikeActivityAPI) GetSpikeVerifyCode(ctx *gin.Context) {
	var requestDTO request.GetSpikeVerifyCode
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
		response.FailWithMessage("Failed to Parse Token", ctx)
		return
	}

	userID, ok := token.Get("user_id")
	if !ok {
		err = errors.New("token parse failed")
		if err != nil {
			response.FailWithMessage("Failed to get userID", ctx)
			return
		}
	}

	verifyCode := SpikeActivityService.GetSpikeVerifyCode(requestDTO.ActivityID, uint(userID.(float64)))
	dto := response.GetSpikeVerifyCode{VerifyCode: verifyCode}
	response.Success(dto, "获取成功", ctx)
}
