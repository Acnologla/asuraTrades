package postgres

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(config config.PostgresConfig) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	return pool
}
