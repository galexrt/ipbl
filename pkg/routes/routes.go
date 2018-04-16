package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response structure
type Response struct {
	Code    int16
	Error   error
	Message string
	Result  interface{}
}

// Register registers routes
func Register(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, Response{
			Code:    http.StatusOK,
			Message: "Welcome to galexrt/ipbl server! The result key contains the available IP lists.",
			Result:  nil,
		})
	})
	registerIP(e)
	registerList(e)
}
