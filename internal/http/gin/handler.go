package gin

import (
	"net/http"

	"github.com/GeovaneCavalcante/ms-notificator/config"
	"github.com/GeovaneCavalcante/ms-notificator/notification"

	"github.com/gin-gonic/gin"

	_ "github.com/GeovaneCavalcante/ms-notificator/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Handlers(envs *config.Environments, s notification.UseCase) *gin.Engine {
	r := gin.Default()

	r.GET("/health", healthHandler)
	v1 := r.Group("/api/v1")

	MakeNotificationHandler(v1, s)
	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	return r
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "App is healthy")
}
