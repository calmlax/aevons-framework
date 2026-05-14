package utils

import (
	"net/http"

	"github.com/calmlax/aevons-framework/response"

	"github.com/gin-gonic/gin"
)

// BadRequest 返回 400 响应
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, response.Response{
		Code:    400,
		Message: msg,
	})
}

// Unauthorized 返回 401 响应
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, response.Response{
		Code:    401,
		Message: msg,
	})
}
