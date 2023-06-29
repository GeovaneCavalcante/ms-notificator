package notification

import (
	"context"
	"time"
)

type Status string

const (
	Delivered Status = "delivered"
	Pending   Status = "pending"
	Error     Status = "error"
)

type Notification struct {
	ID         string `json:"id"`
	RawMessage string `json:"rawMessage"`
	UserID     string `json:"userId"`
}
type ScheduledNotification struct {
	ID             string       `json:"id"`
	Notification   Notification `json:"notification"`
	DateScheduling time.Time    `json:"dateScheduling"`
	Status         Status       `json:"status"`
}

type NotificationRepository interface {
	CreateNotification(n *Notification) (*Notification, error)
}

type ScheduledNotificationRepository interface {
	CreateScheduledNotification(n *ScheduledNotification) (*ScheduledNotification, error)
	ListScheduledNotifications(status string) ([]*ScheduledNotification, error)
	UpdateStatusByID(ctx context.Context, ID string, status Status) error
}

type UseCase interface {
	SendNoticiation(ctx context.Context, notification Notification) error
	ManageNotification(ctx context.Context, notification Notification, dateScheduling string) error
	CreateScheduledNotification(notification Notification, dateScheduling string) (*ScheduledNotification, error)
}
