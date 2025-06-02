package db

import (
	"context"
	"database/sql"
)

type (
	Config struct {
		Host         string
		Port         int
		User         string
		Password     string
		Name         string
		InternalPath string
		DBType       string
	}

	Adapter interface {
		QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
		QueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
		Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	}
)
