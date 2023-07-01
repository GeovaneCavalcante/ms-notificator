package mongo

import (
	"context"
	"testing"

	"github.com/GeovaneCavalcante/ms-notificator/notification"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateNotification(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		ctx := context.TODO()
		test := NewNotificationStorage(db, ctx)
		createdNotification := notification.Notification{
			ID:         "123",
			RawMessage: "TESTE",
			UserID:     "456",
		}

		result, err := test.CreateNotification(&createdNotification)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdNotification.RawMessage, result.RawMessage)
		assert.Equal(t, createdNotification.UserID, result.UserID)
	})

	mt.Run("duplicate key error", func(mt *mtest.T) {
		db := mt.Client.Database("test-db")
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		ctx := context.TODO()
		test := NewNotificationStorage(db, ctx)
		createdNotification := notification.Notification{
			ID:         "123",
			RawMessage: "TESTE",
			UserID:     "456",
		}

		result, err := test.CreateNotification(&createdNotification)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.True(t, mongo.IsDuplicateKeyError(err))
	})

}
