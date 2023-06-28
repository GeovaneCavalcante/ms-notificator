package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListNotification() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.JSON(http.StatusOK, "hello world")
	}
}

func MakeNotificationHandler(r *gin.RouterGroup) {
	r.Handle("GET", "/notifications", ListNotification())
}
