package v1

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models"
	"bank-product-spike-system/models/request"
	"bank-product-spike-system/models/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ArticleAPI struct{}

// GetArticleList
// @Summary 获取文章列表
// @Tags 前台-文章
// @Accept json
// @Param data body request.GetArticleList true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=response.GetArticleList,msg=string} "获取成功"
// @Router /article/list [post]
func (receiver ArticleAPI) GetArticleList(ctx *gin.Context) {
	var requestDTO request.GetArticleList
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}
	list, total, err := ArticleService.GetArticleBaseList(requestDTO.Page, requestDTO.PageSize, models.AuditPassed)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	var responseDto response.GetArticleList
	responseDto.Total = total
	for _, v := range list {
		elem := response.ArticleBaseInfo{
			ArticleID:    v.ID,
			ArticleTitle: v.Title,
			FrontCover: response.Url{
				Host:         v.FrontCover.Host,
				RelativePATH: v.FrontCover.RelativePATH,
			},
			CreatedAT: v.CreatedAt,
		}

		responseDto.List = append(responseDto.List, elem)
	}
	response.Success(responseDto, "获取成功", ctx)
}

// GetArticleContent
// @Summary 获取文章内容根据ID
// @Tags 前台-文章
// @Accept json
// @Param data body request.GetArticleContent true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.Article,msg=string} "获取成功"
// @Router /article/getArticleContent [post]
func (receiver ArticleAPI) GetArticleContent(ctx *gin.Context) {
	var requestDTO request.GetArticleContent
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	content, err := ArticleService.GetArticleContent(requestDTO.ArticleID, models.AuditPassed)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(content, "获取成功", ctx)
}

// GetArticleListAdmin
// @Summary 获取文章列表
// @Tags 文章
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.GetArticleListAdmin true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=[]models.Article,msg=string} "获取成功"
// @Router /article/getList [post]
func (receiver ArticleAPI) GetArticleListAdmin(ctx *gin.Context) {
	var requestDTO request.GetArticleListAdmin
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}

	if !AuthorityService.VerifyAuthority(models.ArticleManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	list, total, err := ArticleService.GetArticleBaseList(requestDTO.Page, requestDTO.PageSize, requestDTO.AuditStatus)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	var responseDto response.GetArticleList
	responseDto.Total = total
	for _, v := range list {
		elem := response.ArticleBaseInfo{
			ArticleID:    v.ID,
			ArticleTitle: v.Title,
			FrontCover: response.Url{
				Host:         v.FrontCover.Host,
				RelativePATH: v.FrontCover.RelativePATH,
			},
			CreatedAT: v.CreatedAt,
		}

		responseDto.List = append(responseDto.List, elem)
	}
	response.Success(responseDto, "获取成功", ctx)
}

// GetArticleContentAdmin
// @Summary 获取文章内容根据ID
// @Tags 文章
// @Accept json
// @Param data body request.GetArticleContent true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{data=models.Article,msg=string} "获取成功"
// @Router /article/getArticleContentAdmin [post]
func (receiver ArticleAPI) GetArticleContentAdmin(ctx *gin.Context) {
	var requestDTO request.GetArticleContent
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	content, err := ArticleService.GetArticleContent(requestDTO.ArticleID, 0)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Success(content, "获取成功", ctx)
}

// PublishArticle
// @Summary 发表文章
// @Tags 文章
// @Security ApiKeyAuth
// @Param title formData string true "标题"
// @Param author formData string true "作者"
// @Param content formData string true "内容"
// @Param front_cover formData file true "封面"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "发表成功"
// @Router /article/publish [post]
func (receiver ArticleAPI) PublishArticle(ctx *gin.Context) {
	title, err := getFormData(ctx, "title")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	author, err := getFormData(ctx, "author")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	content, err := getFormData(ctx, "content")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ArticleManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	article := models.Article{
		Title:   title,
		Content: content,
		Author:  author,
		Audit: models.Audit{
			AuditStatus: models.WaitAudit,
		},
	}

	err = ArticleService.PublishArticle(ctx, article)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("发表成功, 等待审核", ctx)
}

// UpdateArticle
// @Summary 更新文章
// @Tags 文章
// @Security ApiKeyAuth
// @Param article_id formData int true "更新文章id"
// @Param title formData string false "标题"
// @Param author formData string false "作者"
// @Param content formData string false "内容"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "更新成功"
// @Router /article/update [put]
func (receiver ArticleAPI) UpdateArticle(ctx *gin.Context) {
	sArticleID, err := getFormData(ctx, "article_id")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	articleID, _ := strconv.Atoi(sArticleID)
	title, err := getFormData(ctx, "title")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	author, err := getFormData(ctx, "author")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	content, err := getFormData(ctx, "content")
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ArticleManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	article := models.Article{
		Model:   global.Model{ID: uint(articleID)},
		Title:   title,
		Content: content,
		Author:  author,
		Audit:   models.Audit{AuditStatus: models.WaitAudit, AuditMessage: "等待审核中"},
	}

	err = ArticleService.UpdateArticle(article)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("更新成功", ctx)
}

// DeleteArticle
// @Summary 删除文章
// @Tags 文章
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.DeleteArticle true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "删除成功"
// @Router /article/delete [delete]
func (receiver ArticleAPI) DeleteArticle(ctx *gin.Context) {
	var requestDTO request.DeleteArticle
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	token, err := getToken(ctx)
	if err != nil {
		response.AuthFail(err.Error(), ctx)
		return
	}
	if !AuthorityService.VerifyAuthority(models.ArticleManager, token) {
		response.AuthFail("权限认证失败，非法访问，行为已记录", ctx)
		return
	}

	err = ArticleService.DeleteArticle(requestDTO.ArticleID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("删除成功", ctx)
}

// LikeArticle
// @Summary 喜欢文章
// @Tags 前台-文章
// @Security ApiKeyAuth
// @Accept json
// @Param data body request.ArticleID true "传入参数"
// @Product json
// @Success 200 {object} response.DTO{msg=string} "操作成功"
// @Router /article/like [post]
func (receiver ArticleAPI) LikeArticle(ctx *gin.Context) {
	var requestDTO request.ArticleID
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		response.BindJSONError(err, ctx)
		return
	}

	err := ArticleService.LikeArticle(requestDTO.ArticleID)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.SuccessWithMessage("操作成功", ctx)
}
