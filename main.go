package main

import (
	"bank-product-spike-system/config"
	_ "bank-product-spike-system/docs"
	"bank-product-spike-system/initialize"
	"bank-product-spike-system/utils"
	"bank-product-spike-system/utils/SendMQ"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

const (
	Release = 1
	Debug   = 0
)

// @title Swagger API
// @version 1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	Init()
	StartServer()
}

func Init() {
	config.LoadConfig()
	initialize.InitDatabase()
	initialize.InitRedis()
	initialize.InitLogger()
	utils.InitCredential()
	utils.InitUploadConf()
	SendMQ.MQ.InitRabbitMQ()
}

func StartServer() {
	engine := gin.Default()
	http := gin.Default()
	// swagger 接口文档
	if viper.GetInt("server.mod") == Debug {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	initialize.InitDatabaseTables()
	initialize.InitRouter(engine)
	initialize.InitRouter(http)
	//http.Run(":8001")
	panic(engine.RunTLS(":"+viper.GetString("server.port"), "./ssl/ssl.pem", "./ssl/ssl.key"))
}
