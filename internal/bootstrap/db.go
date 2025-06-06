package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/abdussalamfaqih/wallet-service-dev/pkg/db"
)

func NewDB(cfg appconfig.Database) *db.Repository {
	session, err := CreateSession(&cfg)
	if err != nil {
		panic(err)
	}
	return db.NewRepository(session)
}

func CreateSession(cfg *appconfig.Database) (*sql.DB, error) {
	session, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	session.SetMaxOpenConns(20)
	session.SetMaxIdleConns(10)
	session.SetConnMaxLifetime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := session.PingContext(ctx); err != nil {
		session.Close()
		return nil, fmt.Errorf("database unreachable: %w", err)
	}
	return session, err
}
