package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/GeovaneCavalcante/ms-notificator/notification/mongo/dto"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	scheduledNotificationDto = dto.ScheduledNotificationDTO{
		Notification: notification.Notification{
			ID:         "123",
			RawMessage: "TESTE",
			UserID:     "456",
		},
		DateScheduling: time.Time{},
		Status:         "peding",
		CreatedAt:      time.Time{},
		UpdatedAt:      primitive.Timestamp{},
	}

	scheduledNotificationDto2 = dto.ScheduledNotificationDTO{
		Notification: notification.Notification{
			ID:         "5488",
			RawMessage: "TESTE2",
			UserID:     "5457",
		},
		DateScheduling: time.Time{},
		Status:         "peding",
		CreatedAt:      time.Time{},
		UpdatedAt:      primitive.Timestamp{},
	}
)

func TestCreateScheduledNotification(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.TODO()
		test := NewScheduledNotificationStorage(db, ctx)
		createdNotification := notification.Notification{
			ID:         "123",
			RawMessage: "TESTE",
			UserID:     "456",
		}

		createdScheduledNotification := notification.ScheduledNotification{
			ID:             "123",
			Notification:   createdNotification,
			DateScheduling: time.Time{},
			Status:         "peding",
		}

		result, err := test.CreateScheduledNotification(&createdScheduledNotification)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdScheduledNotification.Status, result.Status)
		assert.Equal(t, createdScheduledNotification.Notification, result.Notification)
	})

	mt.Run("duplicate key error", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		ctx := context.TODO()
		test := NewScheduledNotificationStorage(db, ctx)
		createdNotification := notification.Notification{
			ID:         "123",
			RawMessage: "TESTE",
			UserID:     "456",
		}

		createdScheduledNotification := notification.ScheduledNotification{
			ID:             "123",
			Notification:   createdNotification,
			DateScheduling: time.Time{},
			Status:         "peding",
		}

		result, err := test.CreateScheduledNotification(&createdScheduledNotification)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.True(t, mongo.IsDuplicateKeyError(err))
	})

}

func TestListScheduledNotifications(t *testing.T) {

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("it should return the expected list of Scheduled Notification", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		ctx := context.TODO()
		sut := NewScheduledNotificationStorage(db, ctx)

		objId1 := primitive.NewObjectID()
		objId2 := primitive.NewObjectID()

		scheduledNotificationDto.ID = objId1
		scheduledNotificationDto2.ID = objId2

		expectedScheduledNotifications := []*notification.ScheduledNotification{
			{
				ID:             scheduledNotificationDto.ID.Hex(),
				Notification:   scheduledNotificationDto.Notification,
				DateScheduling: scheduledNotificationDto.DateScheduling,
				Status:         scheduledNotificationDto.Status,
			},
			{
				ID:             scheduledNotificationDto2.ID.Hex(),
				Notification:   scheduledNotificationDto2.Notification,
				DateScheduling: scheduledNotificationDto2.DateScheduling,
				Status:         scheduledNotificationDto2.Status,
			},
		}

		first := mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: scheduledNotificationDto.ID},
			{Key: "notification", Value: scheduledNotificationDto.Notification},
			{Key: "dateScheduling", Value: scheduledNotificationDto.DateScheduling},
			{Key: "status", Value: scheduledNotificationDto.Status},
		})

		second := mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
			{Key: "_id", Value: scheduledNotificationDto2.ID},
			{Key: "notification", Value: scheduledNotificationDto2.Notification},
			{Key: "dateScheduling", Value: scheduledNotificationDto2.DateScheduling},
			{Key: "status", Value: scheduledNotificationDto2.Status},
		})
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		result, err := sut.ListScheduledNotifications("peding")

		assert.Nil(t, err)
		assert.EqualValues(t, expectedScheduledNotifications, result)
	})

	mt.Run("it should return an error when find function fails", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		ctx := context.TODO()
		sut := NewScheduledNotificationStorage(db, ctx)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		result, err := sut.ListScheduledNotifications("")

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	mt.Run("it should return error when cursor returns an error", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		ctx := context.TODO()
		sut := NewScheduledNotificationStorage(db, ctx)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{}))

		result, err := sut.ListScheduledNotifications("")

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestUpdateStatusByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("It should return updated scheduled notification", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		ctx := context.TODO()
		sut := NewScheduledNotificationStorage(db, ctx)
		objId := primitive.NewObjectID()

		scheduledNotificationDto.ID = objId

		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: scheduledNotificationDto.ID},
				{Key: "notification", Value: scheduledNotificationDto.Notification},
				{Key: "dateScheduling", Value: scheduledNotificationDto.DateScheduling},
				{Key: "status", Value: scheduledNotificationDto.Status},
			}},
		})

		err := sut.UpdateStatusByID(ctx, objId.Hex(), notification.Delivered)

		assert.Nil(t, err)

	})

	mt.Run("it should return null scheduled notification and error when update fails because of id", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		ctx := context.TODO()
		sut := NewScheduledNotificationStorage(db, ctx)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		err := sut.UpdateStatusByID(ctx, "", notification.Delivered)

		assert.NotNil(t, err)
		assert.Error(t, err)
	})
	mt.Run("it should return nil scheduled notification and error when update fails", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		ctx := context.TODO()
		sut := NewScheduledNotificationStorage(db, ctx)
		mt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})

		err := sut.UpdateStatusByID(ctx, "5d3a2b1e8c5e5f75108233d4", notification.Delivered)

		assert.NotNil(t, err)
		assert.Error(t, err)
	})
}
