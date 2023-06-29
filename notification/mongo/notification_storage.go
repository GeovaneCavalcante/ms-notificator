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

const NotificationCollection = "notifications"

type NotificationStorage struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewNotificationStorage(client *mongo.Database, ctx context.Context) *NotificationStorage {
	return &NotificationStorage{
		collection: client.Collection(NotificationCollection),
		ctx:        ctx,
	}
}

func fromNotification(n notification.Notification) dto.NotificationDTO {
	nDTO := dto.NotificationDTO{
		RawMessage: n.RawMessage,
		UserID:     n.UserID,
	}
	objId, err := primitive.ObjectIDFromHex(n.ID)
	if err == nil {
		nDTO.ID = objId
	}
	return nDTO
}

func (nS *NotificationStorage) CreateNotification(notification *notification.Notification) (*notification.Notification, error) {

	log.Printf("[Repository Notification] Create notification repository starting")

	nDTO := fromNotification(*notification)
	nDTO.CreatedAt = time.Now()
	result, err := nS.collection.InsertOne(nS.ctx, nDTO)
	if err != nil {
		log.Printf("[Repository Notification] Create notification error: %v", err)
		return nil, err
	}
	notification.ID = result.InsertedID.(primitive.ObjectID).Hex()

	log.Printf("[Repository Notification] Create notification succeeded")

	return notification, nil
}
