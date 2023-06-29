package gin

import (
	"net/http"

	"github.com/GeovaneCavalcante/ms-notificator/config"
	"github.com/GeovaneCavalcante/ms-notificator/notification"

	"github.com/gin-gonic/gin"
)

func Handlers(envs *config.Environments, s notification.UseCase) *gin.Engine {
	r := gin.Default()

	r.GET("/health", healthHandler)
	v1 := r.Group("/api/v1")

	MakeNotificationHandler(v1, s)

	return r
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "App is healthy")
}
