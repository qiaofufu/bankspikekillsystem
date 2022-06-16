package response

import "time"

type ArticleBaseInfo struct {
	ArticleID    uint      `json:"article_id"`
	ArticleTitle string    `json:"article_title"`
	FrontCover   Url       `json:"front_cover"`
	CreatedAT    time.Time `json:"created_at"`
}

type GetArticleList struct {
	Total int64             `json:"total"`
	List  []ArticleBaseInfo `json:"list"`
}

type Url struct {
	Host         string `json:"host"`
	RelativePATH string `json:"relative_path"`
}
