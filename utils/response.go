package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS      = 0
	ERROR        = 1
	UNAUTHORIZED = 401
	FORBIDDEN    = 403
	NOT_FOUND    = 404
)

func Result(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Success(c *gin.Context) {
	Result(c, SUCCESS, nil, "操作成功")
}

func SuccessWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "操作成功")
}

func SuccessWithMessage(c *gin.Context, message string) {
	Result(c, SUCCESS, nil, message)
}

func Error(c *gin.Context) {
	Result(c, ERROR, nil, "操作失败")
}

func ErrorWithMessage(c *gin.Context, message string) {
	Result(c, ERROR, nil, message)
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: UNAUTHORIZED,
		Data: nil,
		Msg:  message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code: FORBIDDEN,
		Data: nil,
		Msg:  message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code: NOT_FOUND,
		Data: nil,
		Msg:  message,
	})
}
