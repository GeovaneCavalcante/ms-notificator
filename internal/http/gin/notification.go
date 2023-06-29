package gin

import (
	"net/http"

	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/gin-gonic/gin"
)

func ListNotification(s notification.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload := map[string]interface{}{}
		s.SendNoticiation(c, payload)
		c.JSON(http.StatusOK, "hello world")
	}
}

func MakeNotificationHandler(r *gin.RouterGroup, s notification.UseCase) {
	r.Handle("GET", "/notifications", ListNotification(s))
}
