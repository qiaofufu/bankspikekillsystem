package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/unrolled/secure"
)

func Tls() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ":" + viper.GetString("server.port"),
		})
		err := secureMiddleware.Process(context.Writer, context.Request)
		if err != nil {
			return
		}
		context.Next()
	}
}
