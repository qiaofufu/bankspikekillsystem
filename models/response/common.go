package response

type PageInfoResponse struct {
	Total int64       `json:"total"` // 总数量
	Data  interface{} `json:"data"`  // 获取的数据
}

type CaptchaResponse struct {
	CaptchaType   string `json:"captcha_type"`    // 验证类型
	CaptchaBase64 string `json:"captcha_content"` // 验证码Base64
	CaptchaID     string `json:"captcha_id"`      // 验证码ID
}
