package middleware

import (
	"bank-product-spike-system/global"
	"bank-product-spike-system/models/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			response.AuthFail("jwt为空，验证失败", context)
			context.Abort()
			return
		}
		if global.VerifyToken(token, global.SignatureAlgorithm, global.Secret) == false {
			response.AuthFail("jwt验证失败, 请重新登录", context)
			context.Abort()
			return
		}
		context.Next()
	}
}
