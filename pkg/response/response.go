package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   *Error      `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{Success: false, Error: &Error{Code: code, Message: message}})
}

func InternalFail(c *gin.Context) {
	Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "there is an internal error!")
}
