package system

type RouterGroup struct {
	AdminRouter
	BaseRouter
	UserRouter
	ProductRouter
	SpikeActivityRouter
	FilterRouter
	SlideShowRouter
	ArticleRouter
	AuditRouter
	OrderRouter
}
