package main

import (
	"fmt"

	"github.com/GeovaneCavalcantems-notificator/cmd/api"
	"github.com/GeovaneCavalcantems-notificator/config"
	"github.com/GeovaneCavalcantems-notificator/internal/http/gin"

	"github.com/spf13/cobra"
)

func apiCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Inicialização da api",
		Run: func(cmd *cobra.Command, args []string) {
			envs := config.LoadEnvVars()

			h := gin.Handlers(envs)
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
