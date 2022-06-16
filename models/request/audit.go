package request

type AuditProduct struct {
	ProductID uint   `json:"productID"` // 产品ID
	Result    int    `json:"result"`    // 审核结果
	Message   string `json:"message"`   // 审核消息
}

type AuditActivity struct {
	ActivityID uint   `json:"activityID"` // 活动id
	Result     int    `json:"result"`     // 审核结果
	Message    string `json:"message"`    // 审核消息
}

type AuditArticle struct {
	ArticleID uint   `json:"articleID"` // 文章id
	Result    int    `json:"result"`    // 审核结果
	Message   string `json:"message"`   // 审核消息
}
