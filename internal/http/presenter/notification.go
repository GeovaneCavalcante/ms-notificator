package presenter

import "github.com/GeovaneCavalcante/ms-notificator/internal/messenger"

type NotificationPresenter struct {
	ID string `json:"id"`
}

func (nR *NotificationPresenter) Parse(messageResponse *messenger.MessageResponse) *NotificationPresenter {
	return &NotificationPresenter{
		ID: messageResponse.ID,
	}
}
