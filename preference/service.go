package preference

import (
	"context"
	"fmt"
	"log"
)

type Service struct {
	ctx         context.Context
	preferences []*PreferenceNotification
}

func NewService(ctx context.Context) *Service {

	preferences := []*PreferenceNotification{
		{
			ID:        "cdefe532-0b9e-4895-8063-fa7b827f06a3",
			UserID:    "123",
			RateLimit: 10,
			Channels:  []string{},
			Allow:     true,
		},
		{
			ID:        "d5ac76de-f600-44c2-972f-70a65b846e96",
			UserID:    "456",
			RateLimit: 50,
			Channels:  []string{},
			Allow:     false,
		},
	}

	return &Service{
		ctx:         ctx,
		preferences: preferences,
	}
}

func (s *Service) GetPreferenceByUser(ctx context.Context, userID string) (*PreferenceNotification, error) {
	log.Printf("[Service Preference] Getting preference for user: %s", userID)

	for _, preference := range s.preferences {
		if preference.UserID == userID {
			log.Printf("[Service Preference] Preference found for user: %s", userID)
			return preference, nil
		}
	}

	err := fmt.Errorf("no preference found for userID: %s", userID)
	log.Printf("[Service Preference] Error: %v", err)
	return nil, err
}
