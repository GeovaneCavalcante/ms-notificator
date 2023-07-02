package gin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/gin-gonic/gin"
)

type NotificationData struct {
	RawMessage     string `json:"rawMessage"`
	DateScheduling string `json:"dateScheduling"`
	UserID         string `json:"userId"`
}

func SendNotification(s notification.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {

		var notificationData NotificationData

		if err := c.BindJSON(&notificationData); err != nil {

			log.Printf("[Handler Notification] Error during notification request deserialization: %v", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("Notification request desserialization error: %v", err))
			return
		}
		fmt.Println("aaa")
		log.Printf("[Handler Notification] Event Notification payload: %v", notificationData)

		n := notification.Notification{
			RawMessage: notificationData.RawMessage,
			UserID:     notificationData.UserID,
		}

		dateScheduling := notificationData.DateScheduling

		err := s.ManageNotification(c, n, dateScheduling)
		if err != nil {
			log.Printf("[Handler] Error during notification service: %v", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("Notification service error: %v", err))
			return
		}

		log.Printf("[Handler Notification] Successfully sent notification")

		c.Status(http.StatusOK)
	}
}

func MakeNotificationHandler(r *gin.RouterGroup, s notification.UseCase) {
	r.Handle("POST", "/notifications", SendNotification(s))
}
