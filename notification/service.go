package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
)

type Service struct {
	ctx                       context.Context
	messenger                 messenger.Messenger
	notificationRepo          NotificationRepository
	scheduledNotificationRepo ScheduledNotificationRepository
}

func NewService(ctx context.Context, messenger messenger.Messenger, n NotificationRepository, sN ScheduledNotificationRepository) *Service {
	return &Service{
		ctx:                       ctx,
		messenger:                 messenger,
		notificationRepo:          n,
		scheduledNotificationRepo: sN,
	}
}

func (s *Service) SendNoticiation(ctx context.Context, notification Notification, dateScheduling string) error {

	if dateScheduling != "" {
		sN, err := s.createScheduledNotification(notification, dateScheduling)
		if err != nil {
			return fmt.Errorf("error creating scheduled notification: %w", err)
		}
		log.Printf("[Service] Scheduled notification created successfully: %s", sN.ID)

		return nil
	}

	n, err := s.notificationRepo.CreateNotification(&notification)

	if err != nil {
		log.Printf("[Service] Create notification storage error: %v", err)
		return err
	}

	msg, err := s.messenger.PublishMessage(n.RawMessage)

	if err != nil {
		log.Printf("[Service Notification] Failed to send notification: %v", err)
		return err
	}

	log.Printf("[Service Notification] Notification sent successfully: %s", msg.ID)

	return nil
}

func (s *Service) createScheduledNotification(notification Notification, dateScheduling string) (*ScheduledNotification, error) {
	t, err := parseDate(dateScheduling)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	sN := &ScheduledNotification{
		Notification:   notification,
		DateScheduling: t,
	}

	if _, err = s.scheduledNotificationRepo.CreateScheduledNotification(sN); err != nil {
		return nil, fmt.Errorf("error creating scheduled notification in repo: %w", err)
	}

	return sN, nil
}

func parseDate(date string) (time.Time, error) {
	const layout = "2006-01-02 15:04:05"
	return time.Parse(layout, date)
}
