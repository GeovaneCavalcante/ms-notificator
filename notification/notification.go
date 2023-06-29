package notification

import (
	"context"

	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
)

type UseCase interface {
	SendNoticiation(ctx context.Context, message string) (*messenger.MessageResponse, error)
}
