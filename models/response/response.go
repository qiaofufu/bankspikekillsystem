package response

import (
	"bank-product-spike-system/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS = 0
	ERROR   = -1
	Fail    = -2
)

type DTO struct {
	Code int         `json:"code"` // 响应代码
	Data interface{} `json:"data"` // 相应数据
	Msg  string      `json:"msg"`  // 响应消息
}

func Response(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, DTO{Code: code, Data: data, Msg: msg})
}

// BindJSONError 绑定JSON失败
func BindJSONError(err error, ctx *gin.Context) {
	Response(ERROR, map[string]interface{}{}, "数据绑定失败 "+err.Error(), ctx)
}

// AuthFail 认证失败
func AuthFail(msg string, ctx *gin.Context) {
	global.Logger.SpecialWarning(fmt.Sprintf("[权限认证失败] ip:%s router:%s", ctx.ClientIP(), ctx.Request.RequestURI))
	Response(Fail, map[string]interface{}{}, msg, ctx)
}

// SuccessWithMessage 成功并发送消息
func SuccessWithMessage(msg string, ctx *gin.Context) {
	Response(SUCCESS, map[string]interface{}{}, msg, ctx)
}

// FailWithMessage 失败并发送消息
func FailWithMessage(msg string, ctx *gin.Context) {
	Response(Fail, map[string]interface{}{}, msg, ctx)
}

// Success 成功响应
func Success(data interface{}, msg string, ctx *gin.Context) {
	Response(SUCCESS, data, msg, ctx)
}
