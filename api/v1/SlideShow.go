package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"

	"bank-product-spike-system/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type SlideShowAPI struct{}

// GetSlideShow
// @Summary 获取轮播图
// @Tags 前台-轮播图
// @Product json
// @Success 200 {object} response.DTO{data=response.SlideShowBaseList,msg=string} "获取成功"
// @Router /SlideShow/list [get]
func (receiver SlideShowAPI) GetSlideShow(ctx *gin.Context) {
	list, err := SlideShowService.GetSlideShowBaseList()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	baseList := response.SlideShowBaseList{}
	for _, v := range list {
		elem := response.SlideShowBaseInfo{
			SlideShowID:  v.ID,
			Host:         v.Host,
			RelativePATH: v.RelativePATH,
			Title:        v.Title,
			Description:  v.Description,
		}
		baseList.List = append(baseList.List, elem)
	}

	response.Success(baseList, "获取成功", ctx)
}

// AddSlideShow
// @Summary 添加轮播图
// @Tags 轮播图
// @Security ApiKeyAuth
// @Param file formData file true "图片"
// @Param title formData string true "轮播图标题"
// @Param description formData string true "轮播图描述"
// @Param weight formData int true "轮播图权重"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "添加成功"
// @Router /SlideShow/add [post]
func (receiver SlideShowAPI) AddSlideShow(ctx *gin.Context) {

	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	weight, _ := strconv.Atoi(ctx.PostForm("weight"))

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SlideShowManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	admin, err := AdminService.GetAdminByToken(token)
	if err != nil {
		response.AuthFail("获取管理员信息失败", ctx)
		return
	}
	res, err := utils.UploadFile(ctx, "file")

	slideShow := models.SlideShow{
		Host:         res.Host,
		RelativePATH: res.RelativePath,
		Title:        title,
		Description:  description,
		Weight:       weight,
		Admin:        *admin,
	}
	err = SlideShowService.AddSlideShow(slideShow)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("添加成功", ctx)
}

// UpdateSlideShow
// @Summary 更新轮播图
// @Tags 轮播图
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.UpdateSlideShow true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更新成功"
// @Router /SlideShow/update [put]
func (receiver SlideShowAPI) UpdateSlideShow(ctx *gin.Context) {
	var requestDTO request.UpdateSlideShow
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SlideShowManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	slideShow := models.SlideShow{
		Model: global.Model{
			ID: requestDTO.SlideShowID,
		},
		Title:       requestDTO.Title,
		Description: requestDTO.Description,
		Weight:      requestDTO.Weight,
	}

	err = SlideShowService.UpdateSlideShow(slideShow)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.SuccessWithMessage("更新成功", ctx)
}

// DeleteSlideShow
// @Summary 删除轮播图
// @Tags 轮播图
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.DeleteSlideShow true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "删除成功"
// @Router /SlideShow/delete [delete]
func (receiver SlideShowAPI) DeleteSlideShow(ctx *gin.Context) {
	var requestDTO request.DeleteSlideShow
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SlideShowManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	err = SlideShowService.DeleteSlideShow(requestDTO.SlideShowID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("删除成功", ctx)
}

// GetSlideShowListByAdmin
// @Summary 获取轮播图列表
// @Tags 轮播图
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=response.SlideShowList,msg=string} "获取成功"
// @Router /SlideShow/listAdmin [get]
func (receiver SlideShowAPI) GetSlideShowListByAdmin(ctx *gin.Context) {

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.SlideShowManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	list, err := SlideShowService.GetSlideShowList()
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	var responseDto response.SlideShowList
	for _, v := range list {
		elem := response.SlideShowInfo{
			Model: global.Model{
				ID:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
			},
			Host:         v.Host,
			RelativePATH: v.RelativePATH,
			Title:        v.Title,
			Description:  v.Description,
			Weight:       v.Weight,
			AdminID:      v.AdminID,
			AdminName:    v.Admin.NickName,
		}
		responseDto.List = append(responseDto.List, elem)
	}
	response.Success(responseDto, "获取成功", ctx)
}
