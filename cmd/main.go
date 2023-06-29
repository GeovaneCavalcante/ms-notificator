package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/GeovaneCavalcante/ms-notificator/cmd/api"
	"github.com/GeovaneCavalcante/ms-notificator/config"
	"github.com/GeovaneCavalcante/ms-notificator/internal/http/gin"
	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger/sns"
	"github.com/GeovaneCavalcante/ms-notificator/internal/mongo"
	"github.com/GeovaneCavalcante/ms-notificator/notification"
	notificationRepo "github.com/GeovaneCavalcante/ms-notificator/notification/mongo"
	"github.com/GeovaneCavalcante/ms-notificator/preference"
	"github.com/go-co-op/gocron"

	"github.com/spf13/cobra"
)

func setUpServices(ctx context.Context, envs config.Environments) (*notification.Service, *notificationRepo.ScheduledNotificationStorage, error) {
	snsArn := envs.AwsSnsArn
	snsBroker := sns.New(snsArn)
	dbConn, err := mongo.Open(envs.MongoAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open database connection: %v", err)
	}

	db := dbConn.Database(envs.DbName)
	notificationRepository := notificationRepo.NewNotificationStorage(db, ctx)
	scheduledNotificationRepository := notificationRepo.NewScheduledNotificationStorage(db, ctx)
	preferenceService := preference.NewService(ctx)

	notificationService := notification.NewService(ctx, snsBroker, notificationRepository, scheduledNotificationRepository, preferenceService)

	return notificationService, scheduledNotificationRepository, nil
}

func apiCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Inicialização da api",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			envs := config.LoadEnvVars()
			notificationService, _, err := setUpServices(ctx, *envs)
			if err != nil {
				log.Fatalf("Failed to set up services: %v", err)
			}

			h := gin.Handlers(envs, notificationService)
			if err := api.Start(envs.APIPort, h); err != nil {
				log.Fatalf("Failed to start API: %v", err)
			}
		},
	}
}

func workerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Inicialização do worker",
		Run: func(cmd *cobra.Command, args []string) {
			s := gocron.NewScheduler(time.UTC)
			ctx := context.Background()
			envs := config.LoadEnvVars()
			notificationService, scheduledNotificationRepository, err := setUpServices(ctx, *envs)
			if err != nil {
				log.Fatalf("Failed to set up services: %v", err)
			}

			_, err = s.Every(5).Seconds().Do(func() {
				notifications, err := scheduledNotificationRepository.ListScheduledNotifications(string(notification.Pending))
				if err != nil {
					log.Printf("Failed to list scheduled notifications: %v", err)
					return
				}

				var wg sync.WaitGroup
				wg.Add(len(notifications))

				for _, sN := range notifications {
					go func(n *notification.ScheduledNotification) {
						defer wg.Done()
						err := scheduledNotificationRepository.UpdateStatusByID(ctx, n.ID, notification.Delivered)
						if err != nil {
							log.Printf("Failed to update scheduled notification status to Delivered: %v", err)
							return
						}
						err = notificationService.SendNoticiation(ctx, n.Notification)
						if err != nil {
							log.Printf("Failed to send scheduled notification: %v", err)
							err = scheduledNotificationRepository.UpdateStatusByID(ctx, n.ID, notification.Error)
							if err != nil {
								log.Printf("Failed to update scheduled notification status to Error: %v", err)
							}
						}
					}(sN)
				}

				wg.Wait()
			})
			if err != nil {
				log.Fatalf("Failed to set up job: %v", err)
			}

			s.StartBlocking()
		},
	}
}

func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(apiCommand())
	rootCmd.AddCommand(workerCommand())
	rootCmd.Execute()
}
