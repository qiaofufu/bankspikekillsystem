package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"bank-product-spike-system/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type SpikeActivityAPI struct{}

// ReleaseSpikeActivity
// @Summary 发布秒杀活动
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ReleaseSpikeActivity true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.SpikeActivity,msg=string} "发布成功"
// @Router /spikeActivity/release [post]
func (receiver SpikeActivityAPI) ReleaseSpikeActivity(ctx *gin.Context) {

	var requestDTO request.ReleaseSpikeActivity
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	admin, err := AdminService.GetAdminByToken(token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	product, err := ProductService.GetProductByID(requestDTO.ProductID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	startTime, err := utils.GetTime(requestDTO.StartTime)
	if err != nil {
		response.BindJSONError(err, ctx)
		return
	}
	endingTime, err := utils.GetTime(requestDTO.EndingTime)
	if err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	activity := models.SpikeActivity{
		ActivityName:          requestDTO.ActivityName,
		StartTime:             startTime,
		EndingTime:            endingTime,
		ActivityStatus:        models.ActivityFilterRoleConfiguration,
		Admin:                 *admin,
		Product:               product,
		DescriptionText:       requestDTO.Description,
		RolesIntroductionText: requestDTO.RolesIntroduction,
	}

	err = SpikeActivityService.ReleaseSpikeActivity(&activity)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(activity, "发布成功", ctx)
}

// CacheSpike
// @Summary 缓存秒杀相关信息
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.CacheSpike true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "缓存成功"
// @Router /spikeActivity/cache [post]
func (receiver SpikeActivityAPI) CacheSpike(ctx *gin.Context) {
	var requestDTO request.CacheSpike
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	// 缓存存在1天
	err = SpikeActivityService.Cache(requestDTO.ActivityID, time.Hour*24)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("缓存成功", ctx)
}

// GetCacheSpike
// @Summary 获取秒杀相关缓存信息
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.CacheSpike true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.ActivityBaseInfoCache,msg=string} "获取成功"
// @Router /spikeActivity/getCacheSpike [post]
func (receiver SpikeActivityAPI) GetCacheSpike(ctx *gin.Context) {
	var requestDTO request.CacheSpike
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	dto, err := SpikeActivityService.GetSpikeCache(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(dto, "获取成功", ctx)
}

// UpdateSpikeActivity
// @Summary 更新秒杀活动
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.UpdateActivity true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更新成功"
// @Router /spikeActivity/update [put]
func (receiver SpikeActivityAPI) UpdateSpikeActivity(ctx *gin.Context) {
	var requestDTO request.UpdateActivity
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	var startTime time.Time
	var endingTime time.Time
	if requestDTO.StartTime != "" {
		startTime, err = utils.GetTime(requestDTO.StartTime)
		if err != nil {
			response.BindJSONError(err, ctx)
			return
		}
		fmt.Println(startTime)
	}
	if requestDTO.EndingTime != "" {
		endingTime, err = utils.GetTime(requestDTO.EndingTime)
		if err != nil {
			response.BindJSONError(err, ctx)
			return
		}
	}

	activity := models.SpikeActivity{
		Model:                 global.Model{ID: requestDTO.ActivityID},
		ActivityName:          requestDTO.ActivityName,
		StartTime:             startTime,
		EndingTime:            endingTime,
		DescriptionText:       requestDTO.Description,
		RolesIntroductionText: requestDTO.RolesIntroduction,
		Product: models.Product{
			Model: global.Model{ID: requestDTO.ProductID},
		},
		ActivityStatus: models.ActivityWaitAudit,
		Audit:          models.Audit{AuditStatus: models.WaitAudit, AuditMessage: "等待审核中"},
	}

	err = SpikeActivityService.UpdateSpikeActivity(&activity)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("更新成功", ctx)
}

// DeleteActivity
// @Summary 删除秒杀活动
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.DeleteActivity true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "删除成功"
// @Router /spikeActivity/delete [delete]
func (receiver SpikeActivityAPI) DeleteActivity(ctx *gin.Context) {
	var requestDTO request.DeleteActivity
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	activity := models.SpikeActivity{
		Model: global.Model{ID: requestDTO.ActivityID},
	}
	err = SpikeActivityService.DeleteSpikeActivity(&activity)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("删除成功", ctx)
}

// GetActivityList
// @Summary 获取活动基本信息列表
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetActivityList true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.SpikeActivityBaseInfoList,msg=string} "获取成功"
// @Router /spikeActivity/list [post]
func (receiver SpikeActivityAPI) GetActivityList(ctx *gin.Context) {

	var requestDTO request.GetActivityList
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	list, cnt, err := SpikeActivityService.GetAllActivityBaseInfo(requestDTO.StartTime, requestDTO.EndingTime, requestDTO.ActivityStatus, requestDTO.AuditStatus, requestDTO.Page, requestDTO.PageSize)
	responseDto := response.SpikeActivityBaseInfoList{
		Total: cnt,
	}

	for _, v := range list {
		elem := response.SpikeActivityBaseInfo{
			ActivityID:            v.ID,
			ActivityName:          v.ActivityName,
			ActivityStatus:        v.ActivityStatus,
			AuditStatus:           v.AuditStatus,
			AuditMessage:          v.AuditMessage,
			StartTime:             v.StartTime,
			EndingTime:            v.EndingTime,
			NumberOfParticipants:  v.NumberOfParticipants,
			NumberOfFavorites:     v.NumberOfFavorites,
			AdminID:               v.AdminID,
			AdminName:             v.Admin.NickName,
			DescriptionText:       v.DescriptionText,
			RolesIntroductionText: v.RolesIntroductionText,
		}

		responseDto.List = append(responseDto.List, elem)
	}

	response.Success(responseDto, "获取成功", ctx)
}

// AddActivityGuidanceInformation
// @Summary 添加活动引导相关信息
// @Tags 活动
// @Security ApiKeyAuth
// @Param activity_id formData int true "活动id"
// @Param fileNumber formData int true "文件个数"
// @Param file formData file true "文件"
// @Product json
// @Success 200 {object} response.DTO{msg=string} ""
// @Router /spikeActivity/addGuidanceInformation [post]
func (receiver SpikeActivityAPI) AddActivityGuidanceInformation(ctx *gin.Context) {

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	err = SpikeActivityService.AddActivityGuidanceInformation(ctx, "file")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("添加成功", ctx)
}

// UpdateActivityGuidanceInformation
// @Summary 更新活动引导相关信息
// @Tags 活动
// @Security ApiKeyAuth
// @Param activity_id formData int true "活动id"
// @Param fileNumber formData int true "文件个数"
// @Param file formData file true "文件"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更新成功"
// @Router /spikeActivity/updateGuidanceInformation [post]
func (receiver SpikeActivityAPI) UpdateActivityGuidanceInformation(ctx *gin.Context) {
	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	err = SpikeActivityService.UpdateActivityGuidanceInformation(ctx, "file")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("更新成功", ctx)
}

// GetActivityUserRecord
// @Summary 获取活动用户记录
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ActivityID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=[]models.ActivityUser,msg=string} "获取成功"
// @Router /spikeActivity/getActivityUserRecord [post]
func (receiver SpikeActivityAPI) GetActivityUserRecord(ctx *gin.Context) {
	var requestDTO request.ActivityID
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	list, err := SpikeActivityService.GetActivityUserRecord(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(list, "获取成功", ctx)
}

// GetActivityInfoByID
// @Summary 获取活动全部信息通过活动id
// @Tags 活动
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ActivityID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.SpikeActivity,msg=string} "获取成功"
// @Router /spikeActivity/getActivityInfoByID [post]
func (receiver SpikeActivityAPI) GetActivityInfoByID(ctx *gin.Context) {
	var requestDTO request.ActivityID
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	dto, err := SpikeActivityService.GetActivityInfoByID(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(dto, "获取成功", ctx)
}
