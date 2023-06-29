package main

import (
	"context"
	"fmt"

	"github.com/GeovaneCavalcante/ms-notificator/cmd/api"
	"github.com/GeovaneCavalcante/ms-notificator/config"
	"github.com/GeovaneCavalcante/ms-notificator/internal/http/gin"
	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger/sns"
	"github.com/GeovaneCavalcante/ms-notificator/internal/mongo"
	"github.com/GeovaneCavalcante/ms-notificator/notification"
	notificationRepo "github.com/GeovaneCavalcante/ms-notificator/notification/mongo"

	"github.com/spf13/cobra"
)

func apiCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Inicialização da api",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			envs := config.LoadEnvVars()
			snsArn := envs.AwsSnsArn
			snsBroker := sns.New(snsArn)
			dbConn, _ := mongo.Open(envs.MongoAddress)

			db := dbConn.Database(envs.DbName)

			notificationRepository := notificationRepo.NewNotificationStorage(db, ctx)
			scheduledNotificationRepository := notificationRepo.NewScheduledNotificationStorage(db, ctx)

			notificationService := notification.NewService(ctx, snsBroker, notificationRepository, scheduledNotificationRepository)

			h := gin.Handlers(envs, notificationService)
			err := api.Start(envs.APIPort, h)
			if err != nil {
				fmt.Println("error running api", err)
			}
		},
	}
}

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(apiCommand())
	rootCmd.Execute()
}
