package mongo

import (
	"context"
	"log"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/GeovaneCavalcante/ms-notificator/notification/mongo/dto"

	"go.mongodb.org/mongo-driver/bson"
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
		Status:         notification.Pending,
	}
	objId, err := primitive.ObjectIDFromHex(n.ID)
	if err == nil {
		nDTO.ID = objId
	}
	return nDTO
}

func toScheduledNotification(nsDTO dto.ScheduledNotificationDTO) notification.ScheduledNotification {
	return notification.ScheduledNotification{
		ID:             nsDTO.ID.Hex(),
		DateScheduling: nsDTO.DateScheduling,
		Notification:   nsDTO.Notification,
		Status:         nsDTO.Status,
	}
}

func (nS *ScheduledNotificationStorage) CreateScheduledNotification(scheduleNotification *notification.ScheduledNotification) (*notification.ScheduledNotification, error) {

	log.Printf("[Repository ScheduledNotification] Create scheduled notification repository starting")

	nDTO := fromScheduledNotification(*scheduleNotification)
	nDTO.CreatedAt = time.Now()
	result, err := nS.collection.InsertOne(nS.ctx, nDTO)
	if err != nil {
		log.Printf("[Repository ScheduledNotification] Create scheduled notification error: %v", err)
		return nil, err
	}
	scheduleNotification.ID = result.InsertedID.(primitive.ObjectID).Hex()

	log.Printf("[Repository ScheduledNotification] Create scheduled notification succeeded")

	return scheduleNotification, nil
}

func (nS *ScheduledNotificationStorage) ListScheduledNotifications(status string) ([]*notification.ScheduledNotification, error) {

	filter := bson.M{
		"status": status,
		"dateScheduling": bson.M{
			"$lte": time.Now(),
		},
	}

	cursor, err := nS.collection.Find(nS.ctx, filter)
	if err != nil {
		log.Printf("[Repository ScheduledNotification] Failed to list scheduled notifications: %v", err)
		return nil, err
	}

	var results []*dto.ScheduledNotificationDTO
	if err := cursor.All(nS.ctx, &results); err != nil {
		log.Printf("[Repository ScheduledNotification] Failed to decode scheduled notifications: %v", err)
		return nil, err
	}

	var scheduledNotifications []*notification.ScheduledNotification
	for _, sNDTO := range results {
		p := toScheduledNotification(*sNDTO)
		scheduledNotifications = append(scheduledNotifications, &p)
	}

	return scheduledNotifications, nil
}

func (nS *ScheduledNotificationStorage) UpdateStatusByID(ctx context.Context, ID string, status notification.Status) error {
	objID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Printf("[Repository ScheduledNotification] Failed to convert ID to ObjectID: %v", err)
		return err
	}

	filter := bson.M{"_id": objID}

	update := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "status",
					Value: status,
				},
			},
		},
	}

	_, err = nS.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("[Repository ScheduledNotification] Failed to update status: %v", err)
		return err
	}

	log.Printf("[Repository ScheduledNotification] Successfully changed status")

	return nil
}
