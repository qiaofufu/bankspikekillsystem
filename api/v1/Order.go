package v1

import (
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"github.com/gin-gonic/gin"
)

type OrderAPI struct {
}

// GetSelfOrder
// @Summary 获取自身订单
// @Tags 订单
// @Security ApiKeyAuth
// @Product json
// @Success 200 {object} response.DTO{data=response.OrderList,msg=string} "获取成功"
// @Router /order/getSelfOrder [get]
func (receiver OrderAPI) GetSelfOrder(ctx *gin.Context) {

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

	list, err := OrderService.GetSelfOrder(user.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	var dto response.OrderList
	for _, v := range list {
		product, err := ProductService.GetProductByID(v.ProductID)
		if err != nil {
			//response.FailWithMessage(err.Error(), ctx)
			continue
		}
		elem := response.OrderInfo{
			Order:               v,
			ProductName:         product.ProductName,
			ProductDescription:  product.ProductDescription,
			ProductType:         product.ProductType,
			ProductInterestRate: product.ProductInterestRate,
			InterestRateType:    product.InterestRateType,
			MinHoldTime:         product.MinHoldTime,
			ProductPrice:        product.ProductPrice,
		}
		dto.List = append(dto.List, elem)
	}

	response.Success(dto, "获取成功", ctx)
}

// GetOrderListByActivityID
// @Summary 获取订单列表通过活动id
// @Tags 订单
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetOrderListByActivityID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.GetOrderListByActivityID,msg=string} "获取成功"
// @Router /order/getOrderListByActivityID [post]
func (receiver OrderAPI) GetOrderListByActivityID(ctx *gin.Context) {
	var requestDTO request.GetOrderListByActivityID
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
	cnt, list, err := OrderService.GetOrderList(requestDTO.Page, requestDTO.PageSize, requestDTO.ActivityID, requestDTO.OrderStatus)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	dto := response.GetOrderListByActivityID{
		Total: cnt,
	}
	for _, v := range list {
		elem := response.OrderBaseInfo{
			Order: v,
		}
		user, err := UserService.GetUserInfo(v.UserID)
		if err != nil {
			elem.RealName = "查询失败"
		} else {
			elem.RealName = user.RealInformation.RealName
		}
		dto.List = append(dto.List, elem)
	}

	response.Success(dto, "获取成功", ctx)
}
