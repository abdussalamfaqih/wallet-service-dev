package migrations

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/abdussalamfaqih/wallet-service-dev/internal/bootstrap"
	"github.com/pressly/goose"
)

func RunDBMigration(ctx context.Context, cfg appconfig.Config) error {

	session, err := bootstrap.CreateSession(&cfg.Database)
	if err != nil {
		return err
	}
	defer session.Close()

	if err := goose.Up(session, "db/migrations"); err != nil {
		slog.Error(fmt.Sprintf("goose up failed: %v", err))
	}

	slog.Info("migrations applied successfully")

	return nil
}
