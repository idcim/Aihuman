package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Envelope{
		Code:    400,
		Message: message,
		Data:    gin.H{},
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Envelope{
		Code:    404,
		Message: message,
		Data:    gin.H{},
	})
}
