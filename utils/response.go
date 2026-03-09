package utils

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Result(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: nil,
		Msg:  "Operation successful",
	})
}

func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Msg:  "Operation successful",
	})
}

func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: nil,
		Msg:  message,
	})
}

func Error(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Data: nil,
		Msg:  "Operation failed",
	})
}

func ErrorWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Data: nil,
		Msg:  message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: 401,
		Data: nil,
		Msg:  message,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code: 403,
		Data: nil,
		Msg:  message,
	})
}

func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Code: 429,
		Data: nil,
		Msg:  message,
	})
}

// ParseID parses ID from string
func ParseID(id string) uint {
	result, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0
	}
	return uint(result)
}
