package preference_test

import (
	"context"
	"errors"
	"testing"

	"github.com/GeovaneCavalcante/ms-notificator/preference"
	"github.com/stretchr/testify/assert"
)

func TestGetPreferenceByUser(t *testing.T) {
	t.Run("it should return a reference to the requested user", func(t *testing.T) {
		ctx := context.TODO()

		preferenceData := preference.PreferenceNotification{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     true,
		}

		s := preference.NewService(ctx)

		data, err := s.GetPreferenceByUser(ctx, "123")

		assert.Nil(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, *data, preferenceData)
	})

	t.Run("it should return an empty preference and an error for user preference not found", func(t *testing.T) {
		ctx := context.TODO()

		s := preference.NewService(ctx)

		data, err := s.GetPreferenceByUser(ctx, "8978")
		expectedError := errors.New("no preference found for userID: 8978")
		assert.NotNil(t, err)
		assert.Nil(t, data)
		assert.EqualError(t, err, expectedError.Error())
	})

}
