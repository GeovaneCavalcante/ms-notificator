package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
	"github.com/GeovaneCavalcante/ms-notificator/preference"
)

type Service struct {
	ctx                       context.Context
	messenger                 messenger.Messenger
	notificationRepo          NotificationRepository
	scheduledNotificationRepo ScheduledNotificationRepository
	preferenceService         preference.UseCase
}

func NewService(ctx context.Context, messenger messenger.Messenger, n NotificationRepository, sN ScheduledNotificationRepository, p preference.UseCase) *Service {
	return &Service{
		ctx:                       ctx,
		messenger:                 messenger,
		notificationRepo:          n,
		scheduledNotificationRepo: sN,
		preferenceService:         p,
	}
}

func (s *Service) ManageNotification(ctx context.Context, notification Notification, dateScheduling string) error {

	preference, err := s.preferenceService.GetPreferenceByUser(ctx, notification.UserID)
	if err != nil {
		log.Printf("[Service Notification] Error getting preference for user %s: %v", notification.UserID, err)
		return fmt.Errorf("error getting preference for user %s: %w", notification.UserID, err)
	}

	if !preference.Allow {
		log.Printf("[Service Notification] Notification disabled by user %s preferences", notification.UserID)
		return nil
	}

	if dateScheduling != "" {
		sN, err := s.CreateScheduledNotification(notification, dateScheduling)
		if err != nil {
			return fmt.Errorf("error creating scheduled notification for user %s: %w", notification.UserID, err)
		}
		log.Printf("[Service Notification] Scheduled notification %s created successfully for user %s", sN.ID, notification.UserID)
		return nil
	}

	err = s.SendNoticiation(ctx, notification)

	if err != nil {
		log.Printf("[Service Notification] Error sending notification for user %s: %v", notification.UserID, err)
		return fmt.Errorf("error sending notification for user %s: %w", notification.UserID, err)
	}

	return nil
}

func (s *Service) SendNoticiation(ctx context.Context, notification Notification) error {

	n, err := s.notificationRepo.CreateNotification(&notification)

	if err != nil {
		log.Printf("[Service Notification] Error creating notification for user %s: %v", notification.UserID, err)
		return fmt.Errorf("error creating notification for user %s: %w", notification.UserID, err)
	}

	msg, err := s.messenger.PublishMessage(n.RawMessage)

	if err != nil {
		log.Printf("[Service Notification] Failed to send notification for user %s: %v", notification.UserID, err)
		return fmt.Errorf("failed to send notification for user %s: %w", notification.UserID, err)
	}

	log.Printf("[Service Notification] Notification %s sent successfully for user %s", msg.ID, notification.UserID)

	return nil
}

func (s *Service) CreateScheduledNotification(notification Notification, dateScheduling string) (*ScheduledNotification, error) {
	log.Printf("[Service Notification] Parsing date for scheduled notification: %s", dateScheduling)

	t, err := parseDate(dateScheduling)
	if err != nil {
		log.Printf("[Service Notification] Error parsing date: %v", err)
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	sN := &ScheduledNotification{
		Notification:   notification,
		DateScheduling: t,
	}

	log.Printf("[Service Notification] Creating scheduled notification in repository for user: %s", notification.UserID)

	if _, err = s.scheduledNotificationRepo.CreateScheduledNotification(sN); err != nil {
		log.Printf("[Service Notification] Error creating scheduled notification in repository: %v", err)
		return nil, fmt.Errorf("error creating scheduled notification in repo: %w", err)
	}

	log.Printf("[Service Notification] Scheduled notification successfully created for user: %s", notification.UserID)

	return sN, nil
}

func parseDate(date string) (time.Time, error) {
	const layout = "2006-01-02 15:04:05"
	return time.Parse(layout, date)
}
