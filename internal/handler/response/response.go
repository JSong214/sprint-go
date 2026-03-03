package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResSuccess 成功响应结构
type ResSuccess struct {
	Code    int         `json:"code"`           // 业务状态码，成功时为 200
	Message string      `json:"message"`        // 响应消息，默认 "success"
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// ResError 错误响应结构
type ResError struct {
	Code    int         `json:"code"`              // 业务错误码
	Message string      `json:"message"`           // 错误消息
	Error   string      `json:"error,omitempty"`   // 详细错误信息（可选，仅开发环境）
	Details interface{} `json:"details,omitempty"` // 错误详情（可选）
}

// Success 返回成功响应（HTTP 200）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ResSuccess{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 返回带自定义消息的成功响应（HTTP 200）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, ResSuccess{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// BadRequest 返回参数验证失败响应（HTTP 400）
func BadRequest(c *gin.Context, bizCode int, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, ResError{
		Code:    bizCode,
		Message: message,
		Details: details,
	})
}

// Unauthorized 返回认证失败响应（HTTP 401）
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ResError{
		Code:    40001,
		Message: message,
	})
}

// InternalServerError 返回服务器内部错误响应（HTTP 500）
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ResError{
		Code:    50000,
		Message: message,
	})
}

// InternalServerErrorWithDetails 返回带详细错误的服务器内部错误响应（仅开发环境）
func InternalServerErrorWithDetails(c *gin.Context, message string, err error) {
	c.JSON(http.StatusInternalServerError, ResError{
		Code:    50000,
		Message: message,
		Error:   err.Error(),
	})
}
