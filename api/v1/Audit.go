package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"

	"github.com/gin-gonic/gin"
)

type AuditAPI struct {
}

// AuditProduct
// @Summary 审核产品
// @Tags 审核
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.AuditProduct true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "审核成功"
// @Router /audit/product [post]
func (receiver AuditAPI) AuditProduct(ctx *gin.Context) {
	var requestDTO request.AuditProduct
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ProductAuditor, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	product := models.Product{
		Model: global.Model{
			ID: requestDTO.ProductID,
		},
		Audit: models.Audit{
			AuditStatus:  requestDTO.Result,
			AuditMessage: requestDTO.Message,
		},
	}

	err = ProductService.UpdateProduct(&product)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("审核成功", ctx)
}

// AuditActivity
// @Summary 审核活动
// @Tags 审核
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.AuditActivity true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "审核成功"
// @Router /audit/activity [post]
func (receiver AuditAPI) AuditActivity(ctx *gin.Context) {
	var requestDTO request.AuditActivity
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ActivityAuditor, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	var status int
	switch requestDTO.Result {
	case models.AuditFailed:
		status = models.ActivityWaitAudit
	case models.AuditPassed:
		status = models.ActivityWaitStart
	}
	activity := models.SpikeActivity{
		Model:          global.Model{ID: requestDTO.ActivityID},
		ActivityStatus: status,
		Audit: models.Audit{
			AuditStatus:  requestDTO.Result,
			AuditMessage: requestDTO.Message,
		},
	}

	err = SpikeActivityService.UpdateSpikeActivity(&activity)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("审核成功", ctx)
}

// AuditArticle
// @Summary 审核文章
// @Tags 审核
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.AuditArticle true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "审核成功"
// @Router /audit/article [post]
func (receiver AuditAPI) AuditArticle(ctx *gin.Context) {
	var requestDTO request.AuditArticle
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ArticleAuditor, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	article := models.Article{
		Model: global.Model{
			ID: requestDTO.ArticleID,
		},
		Audit: models.Audit{
			AuditStatus:  requestDTO.Result,
			AuditMessage: requestDTO.Message,
		},
	}

	err = ArticleService.UpdateArticle(article)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("审核成功", ctx)
}
