package gin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GeovaneCavalcante/ms-notificator/internal/http/presenter"
	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/gin-gonic/gin"
)

type NotificationData struct {
	RawMessage string `json:"rawMessage"`
}

func ListNotification(s notification.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {

		var notificationData NotificationData

		if err := c.BindJSON(&notificationData); err != nil {
			log.Printf("[Handler Notification] Error during notification request deserialization: %v", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Notification request desserialization error: %v", err))
			return
		}

		notification, err := s.SendNoticiation(c, notificationData.RawMessage)
		if err != nil {
			log.Printf("[Handler] Error during notification service: %v", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("Notification service error: %v", err))
			return
		}

		notificationResponse := &presenter.NotificationPresenter{}

		notificationResponse = notificationResponse.Parse(notification)

		log.Printf("[Handler Notification] Successfully sent notification")

		c.JSON(http.StatusOK, notificationResponse)
	}
}

func MakeNotificationHandler(r *gin.RouterGroup, s notification.UseCase) {
	r.Handle("POST", "/notifications", ListNotification(s))
}
