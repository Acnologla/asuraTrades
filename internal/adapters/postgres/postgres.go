package postgres

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func NewConnection(ctx context.Context, config config.PostgresConfig) Database {
	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	return pool
}
