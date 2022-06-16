package middleware

import (
	"bank-product-spike-system/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Log() gin.HandlerFunc {
	return func(context *gin.Context) {
		startTime := time.Now()
		context.Next()
		endingTime := time.Now()
		latencyTime := endingTime.Sub(startTime)
		reqMethod := context.Request.Method
		reqUrl := context.Request.URL
		statusCode := context.Writer.Status()
		clientIP := context.ClientIP()
		global.Logger.RequestInfo(fmt.Sprintf("| %3d | %13v | %15s | %s | %s |\n", statusCode, latencyTime, clientIP, reqMethod, reqUrl))
	}
}
