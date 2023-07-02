package preference

import "context"

type PreferenceNotification struct {
	ID        string   `json:"id"`
	UserID    string   `json:"userId"`
	RateLimit int32    `json:"rateLimit"`
	Channels  []string `json:"channels"`
	Allow     bool     `json:"allow"`
}

type UseCase interface {
	GetPreferenceByUser(ctx context.Context, userID string) (*PreferenceNotification, error)
}
