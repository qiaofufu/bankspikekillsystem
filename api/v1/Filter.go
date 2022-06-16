package v1

import (
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

type FilterAPI struct{}

// GetTableNameList
// @Summary 获取数据表列表
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=[]models.FilterTable,msg=string} "获取成功"
// @Router /filter/tableList [get]
func (receiver FilterAPI) GetTableNameList(ctx *gin.Context) {

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	table, err := FilterService.GetTableNameList()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(table, "获取成功", ctx)
}

// GetFieldNameList
// @Summary 获取字段列表
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetFieldNameList true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=[]string,msg=string} "获取成功"
// @Router /filter/fieldList [post]
func (receiver FilterAPI) GetFieldNameList(ctx *gin.Context) {
	var requestDTO request.GetFieldNameList
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

	list, err := FilterService.GetFieldListByTableName(requestDTO.TableName)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.Success(list, "获取成功", ctx)
}

// SetValueRange
// @Summary 设置筛选取值范围
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.SetValueRange true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "设置成功"
// @Router /filter/setValueRange [post]
func (receiver FilterAPI) SetValueRange(ctx *gin.Context) {
	var requestDTO request.SetValueRange
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

	value := models.FilterRoles{
		ActivityID:    requestDTO.ActivityID,
		FilterTableID: requestDTO.FilterTableID,
		FiledName:     requestDTO.FieldName,
		ValueRange:    requestDTO.ValueRange,
		Description:   requestDTO.Description,
		ErrorTips:     requestDTO.ErrorTips,
		Admin:         *admin,
	}
	err = FilterService.SetValueRange(&value)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("设置成功", ctx)
}

// GetFilterRolesByID
// @Summary 根据id获取筛选规则
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetFilterRolesByID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=[]models.FilterRoles,msg=string} "获取成功"
// @Router /filter/getFilterRolesByID [post]
func (receiver FilterAPI) GetFilterRolesByID(ctx *gin.Context) {
	var requestDTO request.GetFilterRolesByID
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

	list, err := FilterService.GetFilterRolesByActivity(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.Success(list, "获取成功", ctx)
}

// DeleteFilterRoles
// @Summary 删除筛选规则
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.DeleteFilterRoles true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "删除成功"
// @Router /filter/delete [delete]
func (receiver FilterAPI) DeleteFilterRoles(ctx *gin.Context) {
	var requestDTO request.DeleteFilterRoles
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

	err = FilterService.DeleteFilterRoles(requestDTO.ActivityID)

	response.SuccessWithMessage("删除成功", ctx)
}

// DeleteFilterRolesByID
// @Summary 删除一条筛选规则
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.DeleteFilterRolesByID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "删除成功"
// @Router /filter/deleteByID [delete]
func (receiver FilterAPI) DeleteFilterRolesByID(ctx *gin.Context) {
	var requestDTO request.DeleteFilterRolesByID
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

	err = FilterService.DeleteFilterRolesByID(requestDTO.FilterRolesID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("删除成功", ctx)
}

// GetRolesList
// @Summary 获取已有规则列表
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=response.RolesList,msg=string} "获取成功"
// @Router /filter/list [get]
func (receiver FilterAPI) GetRolesList(ctx *gin.Context) {
	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SpikeActivityManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	list, err := FilterService.GetFilterRolesList(0)
	var dto response.RolesList
	for _, v := range list {
		elem := response.RolesInfo{
			ID:            v.ID,
			FilterTableID: v.FilterTableID,
			FiledName:     v.FiledName,
			ValueRange:    v.ValueRange,
			ErrorTips:     v.ErrorTips,
			Description:   v.Description,
		}
		dto.List = append(dto.List, elem)
	}

	response.Success(dto, "获取成功", ctx)
}

// FinishFilterConfiguration
// @Summary 完成筛选配置
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.FinishFilterConfiguration true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "完成配置"
// @Router /filter/finishFilterConfiguration [post]
func (receiver FilterAPI) FinishFilterConfiguration(ctx *gin.Context) {
	var requestDTO request.FinishFilterConfiguration
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

	err = FilterService.FinishFilterConfiguration(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("完成配置", ctx)
}

// Check
// @Summary 准入校验
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.FilterCheck true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "初筛通过"
// @Router /filter/check [post]
func (receiver FilterAPI) Check(ctx *gin.Context) {
	var requestDTO request.FilterCheck
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	err = FilterService.Check2(token, requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("初筛通过", ctx)
}

// Check2
// @Summary 准入校验
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.FilterCheck true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "初筛通过"
// @Router /filter/check2 [post]
func (receiver FilterAPI) Check2(ctx *gin.Context) {
	var requestDTO request.FilterCheck
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	err = FilterService.Check2(token, requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("初筛通过", ctx)
}

// FinishFilterConfiguration2
// @Summary 完成配置
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.FinishFilterConfiguration2 true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "完成配置"
// @Router /filter/finishFilterConfiguration2 [post]
func (receiver FilterAPI) FinishFilterConfiguration2(ctx *gin.Context) {
	var requestDTO request.FinishFilterConfiguration2
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
	fmt.Println(requestDTO)
	err = FilterService.FinishFilterConfiguration2(requestDTO.Node, requestDTO.Edge, requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("完成配置", ctx)
}

// GetFilterTree
// @Summary 获取筛选规则树
// @Tags 筛选系统
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ActivityID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.FilterTree,msg=string} "获取成功"
// @Router /filter/getFilterTree [post]
func (receiver FilterAPI) GetFilterTree(ctx *gin.Context) {
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
	var tree response.FilterTree
	tree.Node, err = FilterService.GetFilterNode(requestDTO.ActivityID)
	tree.Edge, err = FilterService.GetFilterEdge(requestDTO.ActivityID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(tree, "获取成功", ctx)
}
