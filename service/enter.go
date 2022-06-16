package service

type ServiceGroup struct {
	AdminService
	CaptchaService
	UserService
	ProductService
	SpikeActivityService
	AuthorityService
	FilterService
	SlideShowService
	ArticleService
	OrderService
	RedisServer
	JWTService
}

type ServicesGroup struct {
	SystemService ServiceGroup
}

var ServiceGroupAPP = new(ServicesGroup)
