package request

type PublishArticle struct {
	Title   string `json:"title" binding:"required"`   // 标题
	Author  string `json:"author" binding:"required"`  // 作者
	Content string `json:"content" binding:"required"` // 内容
}

type UpdateArticle struct {
	ArticleID uint   `json:"article_id" binding:"required"` // 文章id
	Title     string `json:"title"`                         // 标题
	Author    string `json:"author"`                        // 作者
	Content   string `json:"content"`                       // 内容
}

type DeleteArticle struct {
	ArticleID uint `json:"article_id"` // 文章id
}

type GetArticleList struct {
	PageInfo
}

type GetArticleListAdmin struct {
	PageInfo
	AuditStatus int `json:"status"` // 文章状态 1-等待审核 2-审核成功 -1-审核失败
}

type GetArticleContent struct {
	ArticleID uint `json:"article_id"` // 文章ID
}

type ArticleID struct {
	ArticleID uint `json:"article_id"` // 文章id
}
