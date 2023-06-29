package notification

import (
	"context"
	"log"

	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
)

type Service struct {
	ctx       context.Context
	messenger messenger.Messenger
}

func NewService(ctx context.Context, messenger messenger.Messenger) *Service {
	return &Service{
		ctx:       ctx,
		messenger: messenger,
	}
}

func (s *Service) SendNoticiation(ctx context.Context, message map[string]interface{}) (*messenger.MessageResponse, error) {
	msg, err := s.messenger.PublishMessage(message)

	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return nil, err
	}

	return msg, nil
}
