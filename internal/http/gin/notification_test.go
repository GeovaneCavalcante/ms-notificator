package gin

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GeovaneCavalcante/ms-notificator/notification/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSendNotification(t *testing.T) {
	t.Run("it should return 200 for sending an instant notification", func(t *testing.T) {
		notificationData := NotificationData{

			RawMessage:     "teste",
			DateScheduling: "",
			UserID:         "124",
		}

		notificationBytes, _ := json.Marshal(notificationData)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		service := mock.NewMockUseCase(ctrl)

		service.EXPECT().
			ManageNotification(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		r := gin.Default()

		v1 := r.Group("/api/v1")

		MakeNotificationHandler(v1, service)

		req, err := http.NewRequest("POST", "/api/v1/notifications", bytes.NewBuffer(notificationBytes))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("it should return bad request when invalid body is passed", func(t *testing.T) {
		r := gin.Default()

		v1 := r.Group("/api/v1")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		service := mock.NewMockUseCase(ctrl)

		MakeNotificationHandler(v1, service)

		invalidBody := []byte(`111`)

		req, err := http.NewRequest("POST", "/api/v1/notifications", bytes.NewBuffer(invalidBody))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("it should return 500 when ManageNotification returns error", func(t *testing.T) {
		notificationData := NotificationData{

			RawMessage:     "teste",
			DateScheduling: "",
			UserID:         "124",
		}

		notificationBytes, _ := json.Marshal(notificationData)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		service := mock.NewMockUseCase(ctrl)

		service.EXPECT().
			ManageNotification(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("error to sending notification"))

		r := gin.Default()

		v1 := r.Group("/api/v1")

		MakeNotificationHandler(v1, service)

		req, err := http.NewRequest("POST", "/api/v1/notifications", bytes.NewBuffer(notificationBytes))

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
