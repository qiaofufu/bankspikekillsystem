package request

type DeleteSlideShow struct {
	SlideShowID uint `json:"slide_show_id"`
}

type UpdateSlideShow struct {
	SlideShowID uint   `json:"slide_show_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}
