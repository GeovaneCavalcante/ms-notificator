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

func (s *Service) SendNoticiation(ctx context.Context, message string) (*messenger.MessageResponse, error) {
	msg, err := s.messenger.PublishMessage(message)

	if err != nil {
		log.Printf("[Service Notification] Failed to send notification: %v", err)
		return nil, err
	}

	log.Printf("[Service Notification] Notification sent successfully: %s", msg.ID)

	return msg, nil
}
