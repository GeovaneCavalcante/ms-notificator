package notification

import (
	"context"

	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger"
)

type UseCase interface {
	SendNoticiation(ctx context.Context, message map[string]interface{}) (*messenger.MessageResponse, error)
}
