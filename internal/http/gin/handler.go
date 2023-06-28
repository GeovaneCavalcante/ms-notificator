package gin

import (
	"net/http"

	"github.com/GeovaneCavalcantems-notificator/config"

	"github.com/gin-gonic/gin"
)

func Handlers(*config.Environments) *gin.Engine {
	r := gin.Default()

	r.GET("/health", healthHandler)
	v1 := r.Group("/api/v1")

	MakeNotificationHandler(v1)

	return r
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "App is healthy")
}
