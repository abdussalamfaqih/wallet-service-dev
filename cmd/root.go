package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdussalamfaqih/wallet-service-dev/cmd/http"
	"github.com/abdussalamfaqih/wallet-service-dev/cmd/migrations"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/spf13/cobra"
)

var configFile string

func Run() {
	rootCmd := &cobra.Command{}

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
		os.Exit(1)
	}()

	cmd := []*cobra.Command{
		{
			Use:   "run-http",
			Short: "Run HTTP Server",
			Run: func(cmd *cobra.Command, args []string) {
				cfg := appconfig.LoadConfig(configFile)
				http.Start(ctx, cfg)
			},
		},
		{
			Use:   "run-migration",
			Short: "Run HTTP Server",
			Run: func(cmd *cobra.Command, args []string) {
				cfg := appconfig.LoadConfig(configFile)
				migrations.RunDBMigration(ctx, cfg)
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	rootCmd.PersistentFlags().StringVar(&configFile, "config_file", "config.json", "Path to the config file (e.g. config.json)")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
