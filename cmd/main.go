package main

import (
	"context"
	"fmt"

	"github.com/GeovaneCavalcante/ms-notificator/cmd/api"
	"github.com/GeovaneCavalcante/ms-notificator/config"
	"github.com/GeovaneCavalcante/ms-notificator/internal/http/gin"
	"github.com/GeovaneCavalcante/ms-notificator/internal/messenger/sns"
	"github.com/GeovaneCavalcante/ms-notificator/notification"

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
			notificationService := notification.NewService(ctx, snsBroker)

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
