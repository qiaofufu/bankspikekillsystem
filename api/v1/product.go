package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"github.com/gin-gonic/gin"
)

type ProductAPI struct{}

// ReleaseProduct
// @Summary 发布产品
// @Tags 产品
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ReleaseProduct true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.Product,msg=string} "发布成功"
// @Router /product/release [post]
func (receiver ProductAPI) ReleaseProduct(ctx *gin.Context) {
	var requestDTO request.ReleaseProduct
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.ProductManager, token) {
		response.AuthFail("没有操作权限，非法操作, 行为已记录.", ctx)
	}

	product := models.Product{
		ProductName:         requestDTO.ProductName,
		ProductDescription:  requestDTO.ProductDescription,
		ProductNumber:       requestDTO.ProductNumber,
		SoldNumber:          0,
		MinHoldTime:         requestDTO.MinHoldTime,
		ProductTags:         getTags(requestDTO.ProductTags),
		ProductInterestRate: requestDTO.ProductInterestRate,
		InterestRateType:    requestDTO.InterestRateType,
		ProductPrice:        requestDTO.ProductPrice,
		Audit:               models.Audit{AuditStatus: models.WaitAudit},
	}

	err := ProductService.ReleaseProduct(&product)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.Success(product, "成功生成产品", ctx)
}

// UpdateProduct
// @Summary 更新产品
// @Tags 产品
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.UpdateProduct true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.Product} "更新成功"
// @Router /product/update [put]
func (receiver ProductAPI) UpdateProduct(ctx *gin.Context) {
	var requestDTO request.UpdateProduct
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, _ := getToken(ctx)
	if !AuthorityService.VerifyAuthority(models.ProductManager, token) {
		response.AuthFail("没有操作权限，非法访问，行为已记录。", ctx)
		return
	}

	product := models.Product{
		Model: global.Model{
			ID: requestDTO.ProductID,
		},
		ProductName:         requestDTO.ProductName,
		ProductDescription:  requestDTO.ProductDescription,
		ProductNumber:       requestDTO.ProductNumber,
		MinHoldTime:         requestDTO.MinHoldTime,
		ProductTags:         getTags(requestDTO.ProductTags),
		ProductInterestRate: requestDTO.ProductInterestRate,
		InterestRateType:    requestDTO.InterestRateType,
		ProductPrice:        requestDTO.ProductPrice,
		Audit:               models.Audit{AuditStatus: models.WaitAudit, AuditMessage: "等待审核中"},
	}

	if err := ProductService.UpdateProduct(&product); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(product, "更新成功", ctx)
}

// DeleteProduct
// @Summary 删除产品
// @Tags 产品
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.DeleteProduct true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "删除成功"
// @Router /product/delete [delete]
func (receiver ProductAPI) DeleteProduct(ctx *gin.Context) {
	var requestDTO request.DeleteProduct
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}
	product := models.Product{
		Model: global.Model{ID: requestDTO.ProductID},
	}
	err := ProductService.DeleteProduct(&product)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(product, "删除成功", ctx)
}

// GetAllProduct
// @Summary 获取所有产品
// @Tags 产品
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetAllProduct true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.ProductList,msg=string} "获取成功"
// @Router /product/all [post]
func (receiver ProductAPI) GetAllProduct(ctx *gin.Context) {
	var requestDTO request.GetAllProduct
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ProductManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	products, cnt, err := ProductService.GetAllProduct(requestDTO.ProductName, requestDTO.InterestRateType, requestDTO.AuditStatus, requestDTO.Page, requestDTO.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	responseDto := response.ProductList{
		Total: cnt,
		List:  products,
	}
	response.Success(responseDto, "获取成功", ctx)
}

// GetFeaturedProduct
// @Summary 获取精选产品
// @Tags 前台-产品
// @Accept json
// @Param data body request.PageInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.GetFeaturedProduct,msg=string} "获取成功"
// @Router /product/featuredProduct [post]
func (receiver ProductAPI) GetFeaturedProduct(ctx *gin.Context) {
	var requestDTO request.PageInfo
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}
	list, err := ProductService.GetVisibleProductListOrderByWeight(requestDTO.Page, requestDTO.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	var responseDto response.GetFeaturedProduct
	for _, v := range list {
		elem := response.ProductBaseInfo{
			ProductID:           v.ID,
			ProductName:         v.ProductName,
			ProductInterestRate: v.ProductInterestRate,
		}
		responseDto.List = append(responseDto.List, elem)
	}
	response.Success(responseDto, "获取成功", ctx)
}

// GetVisibleProduct
// @Summary 根据ID获取可见产品
// @Tags 前台-产品
// @Accept json
// @Param data body request.GetVisibleProduct true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.Product,msg=string} "获取成功"
// @Router /product/getProduct [post]
func (receiver ProductAPI) GetVisibleProduct(ctx *gin.Context) {
	var requestDTO request.GetVisibleProduct
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	product, err := ProductService.GetVisibleProduct(requestDTO.ProductID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.Success(product, "获取成功", ctx)
}

// GetProductTagsList
// @Summary 获取产品标签列表
// @Tags 产品
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.PageInfo true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=[]models.Tag} "获取成功"
// @Router /product/tag/list [post]
func (receiver ProductAPI) GetProductTagsList(ctx *gin.Context) {
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
	if !AuthorityService.VerifyAuthority(models.ProductManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}
	list, cnt, err := ProductService.GetTagsList(requestDTO.Page, requestDTO.PageSize)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	responseDto := response.TagsList{
		Total: cnt,
		List:  list,
	}
	response.Success(responseDto, "获取列表成功", ctx)
}
