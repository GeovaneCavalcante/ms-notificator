package notification_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
	mockMessenger "github.com/GeovaneCavalcante/ms-notificator/internal/messenger/mock"
	"github.com/GeovaneCavalcante/ms-notificator/notification"
	"github.com/GeovaneCavalcante/ms-notificator/notification/mock"
	"github.com/GeovaneCavalcante/ms-notificator/preference"
	mockPreference "github.com/GeovaneCavalcante/ms-notificator/preference/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestManageNotification(t *testing.T) {
	t.Run("it should return an empty error for an instant notification case", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		preferences := preference.PreferenceNotification{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     true,
		}

		msg := messenger.MessageResponse{
			ID: "123",
		}

		ctrl := gomock.NewController(t)

		nRepo := mock.NewMockNotificationRepository(ctrl)
		messager := mockMessenger.NewMockMessenger(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		messager.EXPECT().PublishMessage(gomock.Any()).Return(&msg, nil).Times(1)

		preferenceService.EXPECT().
			GetPreferenceByUser(ctx, gomock.Any()).
			Return(&preferences, nil).
			Times(1)

		nRepo.EXPECT().
			CreateNotification(&expectedNotification).
			Return(&expectedNotification, nil).
			Times(1)

		s := notification.NewService(ctx, messager, nRepo, nil, preferenceService)

		err := s.ManageNotification(ctx, expectedNotification, "")

		assert.Nil(t, err)
	})

	t.Run("it should return an error for a instant notification case", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		preferences := preference.PreferenceNotification{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     true,
		}

		ctrl := gomock.NewController(t)

		nRepo := mock.NewMockNotificationRepository(ctrl)
		messager := mockMessenger.NewMockMessenger(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		preferenceService.EXPECT().
			GetPreferenceByUser(ctx, gomock.Any()).
			Return(&preferences, nil).
			Times(1)

		nRepo.EXPECT().
			CreateNotification(gomock.Any()).
			Return(nil, errors.New("error sending notification")).
			Times(1)

		expectedError := errors.New("error sending notification for user 123: error creating notification for user 123: error sending notification")

		s := notification.NewService(ctx, messager, nRepo, nil, preferenceService)

		err := s.ManageNotification(ctx, expectedNotification, "")

		assert.NotNil(t, err)
		assert.EqualError(t, expectedError, err.Error())
	})

	t.Run("it should return an empty error for an scheduled notification case", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		expectedScheduleNotification := notification.ScheduledNotification{
			ID:             "123",
			Notification:   expectedNotification,
			DateScheduling: time.Time{},
			Status:         "peding",
		}

		preferences := preference.PreferenceNotification{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     true,
		}

		ctrl := gomock.NewController(t)

		sNrepo := mock.NewMockScheduledNotificationRepository(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		preferenceService.EXPECT().
			GetPreferenceByUser(ctx, gomock.Any()).
			Return(&preferences, nil).
			Times(1)

		sNrepo.EXPECT().
			CreateScheduledNotification(gomock.Any()).
			Return(&expectedScheduleNotification, nil).
			Times(1)

		s := notification.NewService(ctx, nil, nil, sNrepo, preferenceService)

		err := s.ManageNotification(ctx, expectedNotification, "2023-09-12 22:22:22")

		assert.Nil(t, err)
	})

	t.Run("it should return an error for a scheduled notification case", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		preferences := preference.PreferenceNotification{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     true,
		}

		ctrl := gomock.NewController(t)

		sNrepo := mock.NewMockScheduledNotificationRepository(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		preferenceService.EXPECT().
			GetPreferenceByUser(ctx, gomock.Any()).
			Return(&preferences, nil).
			Times(1)

		expectedError := errors.New("error creating scheduled notification for user 123: error creating scheduled notification in repo: error creating schedule")
		sNrepo.EXPECT().
			CreateScheduledNotification(gomock.Any()).
			Return(nil, errors.New("error creating schedule")).
			Times(1)

		s := notification.NewService(ctx, nil, nil, sNrepo, preferenceService)

		err := s.ManageNotification(ctx, expectedNotification, "2023-09-12 22:22:22")

		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("it should return an error for a notification case by querying the user's preferred service", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		ctrl := gomock.NewController(t)

		nRepo := mock.NewMockNotificationRepository(ctrl)
		messager := mockMessenger.NewMockMessenger(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		preferenceService.EXPECT().
			GetPreferenceByUser(ctx, gomock.Any()).
			Return(nil, errors.New("error sending notification")).
			Times(1)

		expectedError := errors.New("error getting preference for user 123: error sending notification")
		s := notification.NewService(ctx, messager, nRepo, nil, preferenceService)

		err := s.ManageNotification(ctx, expectedNotification, "")

		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedError.Error())
	})

	t.Run("it should tests when the user doesn't have the notification permission on their profile enabled", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		preferences := preference.PreferenceNotification{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     false,
		}

		ctrl := gomock.NewController(t)

		sNrepo := mock.NewMockScheduledNotificationRepository(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		preferenceService.EXPECT().
			GetPreferenceByUser(ctx, gomock.Any()).
			Return(&preferences, nil).
			Times(1)

		s := notification.NewService(ctx, nil, nil, sNrepo, preferenceService)

		err := s.ManageNotification(ctx, expectedNotification, "2023-09-12 22:22:22")

		assert.Nil(t, err)
	})

}

func TestSendNoticiation(t *testing.T) {
	t.Run("it should return an error for sending the notification", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		ctrl := gomock.NewController(t)

		nRepo := mock.NewMockNotificationRepository(ctrl)
		messager := mockMessenger.NewMockMessenger(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		messager.EXPECT().PublishMessage(gomock.Any()).Return(nil, errors.New("error sending notification")).Times(1)

		nRepo.EXPECT().
			CreateNotification(&expectedNotification).
			Return(&expectedNotification, nil).
			Times(1)

		expectedError := errors.New("failed to send notification for user 123: error sending notification")

		s := notification.NewService(ctx, messager, nRepo, nil, preferenceService)

		err := s.SendNoticiation(ctx, expectedNotification)

		assert.NotNil(t, err)
		assert.EqualError(t, expectedError, err.Error())
	})

	t.Run("it should return a sending notification", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		msg := messenger.MessageResponse{
			ID: "123",
		}

		ctrl := gomock.NewController(t)

		nRepo := mock.NewMockNotificationRepository(ctrl)
		messager := mockMessenger.NewMockMessenger(ctrl)
		sNrepo := mock.NewMockScheduledNotificationRepository(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)
		messager.EXPECT().PublishMessage(gomock.Any()).Return(&msg, nil).Times(1)

		nRepo.EXPECT().
			CreateNotification(&expectedNotification).
			Return(&expectedNotification, nil).
			Times(1)

		s := notification.NewService(ctx, messager, nRepo, sNrepo, preferenceService)

		err := s.SendNoticiation(ctx, expectedNotification)

		assert.Nil(t, err)
	})
}

func TestCreateScheduledNotification(t *testing.T) {
	t.Run("it should return an error for formatter date", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		ctrl := gomock.NewController(t)

		nRepo := mock.NewMockNotificationRepository(ctrl)
		messager := mockMessenger.NewMockMessenger(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		expectedError := errors.New("error parsing date: parsing time \"a\" as \"2006-01-02 15:04:05\": cannot parse \"a\" as \"2006\"")

		s := notification.NewService(ctx, messager, nRepo, nil, preferenceService)

		_, err := s.CreateScheduledNotification(expectedNotification, "a")

		assert.NotNil(t, err)
		assert.EqualError(t, expectedError, err.Error())
	})

	t.Run("it should return a scheduled notification", func(t *testing.T) {
		ctx := context.TODO()

		expectedNotification := notification.Notification{
			ID:         "123",
			RawMessage: "SomeProduct",
			UserID:     "123",
		}

		expectedScheduleNotification := notification.ScheduledNotification{
			ID:             "123",
			Notification:   expectedNotification,
			DateScheduling: time.Time{},
			Status:         "peding",
		}

		ctrl := gomock.NewController(t)

		sNrepo := mock.NewMockScheduledNotificationRepository(ctrl)
		preferenceService := mockPreference.NewMockUseCase(ctrl)

		sNrepo.EXPECT().
			CreateScheduledNotification(gomock.Any()).
			Return(&expectedScheduleNotification, nil).
			Times(1)

		s := notification.NewService(ctx, nil, nil, sNrepo, preferenceService)

		data, err := s.CreateScheduledNotification(expectedNotification, "2023-09-12 22:22:22")

		assert.Nil(t, err)
		assert.Equal(t, data.Notification, expectedNotification)
	})
}
