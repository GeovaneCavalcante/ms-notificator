package dto

import (
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationDTO struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	RawMessage string              `bson:"rawMessage"`
	UserID     string              `bson:"userId"`
	CreatedAt  time.Time           `bson:"createdAt"`
	UpdatedAt  primitive.Timestamp `bson:"updatedAt"`
}

type ScheduledNotificationDTO struct {
	ID             primitive.ObjectID        `bson:"_id,omitempty"`
	Notification   notification.Notification `bson:"notification"`
	DateScheduling time.Time                 `bson:"dateScheduling"`
	CreatedAt      time.Time                 `bson:"createdAt"`
	UpdatedAt      primitive.Timestamp       `bson:"updatedAt"`
}
