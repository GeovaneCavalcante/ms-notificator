package mongo

import (
	"context"
	"log"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/GeovaneCavalcante/ms-notificator/notification/mongo/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ScheduledNotificationCollection = "scheduledNotifications"

type ScheduledNotificationStorage struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewScheduledNotificationStorage(client *mongo.Database, ctx context.Context) *ScheduledNotificationStorage {
	return &ScheduledNotificationStorage{
		collection: client.Collection(ScheduledNotificationCollection),
		ctx:        ctx,
	}
}

func fromScheduledNotification(n notification.ScheduledNotification) dto.ScheduledNotificationDTO {

	nDTO := dto.ScheduledNotificationDTO{
		Notification:   n.Notification,
		DateScheduling: n.DateScheduling,
	}
	objId, err := primitive.ObjectIDFromHex(n.ID)
	if err == nil {
		nDTO.ID = objId
	}
	return nDTO
}

func (nS *ScheduledNotificationStorage) CreateScheduledNotification(scheduleNotification *notification.ScheduledNotification) (*notification.ScheduledNotification, error) {

	log.Printf("[Repository] Create scheduled notification repository starting")

	nDTO := fromScheduledNotification(*scheduleNotification)
	nDTO.CreatedAt = time.Now()
	result, err := nS.collection.InsertOne(nS.ctx, nDTO)
	if err != nil {
		log.Printf("[Repository] Create scheduled notification error: %v", err)
		return nil, err
	}
	scheduleNotification.ID = result.InsertedID.(primitive.ObjectID).Hex()

	log.Printf("[Repository] Create scheduled notification succeeded")

	return scheduleNotification, nil
}

func (nS *ScheduledNotificationStorage) ListScheduledNotifications(status string) ([]*notification.ScheduledNotification, error) {
	return nil, nil
}
