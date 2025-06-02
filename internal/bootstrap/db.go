package bootstrap

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/abdussalamfaqih/wallet-service-dev/internal/appconfig"
	"github.com/abdussalamfaqih/wallet-service-dev/pkg/db"
)

func NewDB(cfg appconfig.Database) *db.Repository {
	session, _ := createSession(&cfg)
	return db.NewRepository(session, nil)
}

func createSession(cfg *appconfig.Database) (*sql.DB, error) {
	session, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name))
	session.SetMaxOpenConns(20)
	session.SetMaxIdleConns(10)
	session.SetConnMaxLifetime(10 * time.Minute)
	return session, err
}
