package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (r *UserRepository) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	row, err := r.db.Query(ctx, "SELECT id, xp FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[*domain.User])

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) LockUpdate(ctx context.Context, id domain.ID) (func(context.Context) error, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil
	}
	defer tx.Rollback(ctx)

	tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", id)
	return tx.Commit, nil
}

func NewUserRepository(db *pgxpool.Pool) port.UserRepository {
	return &UserRepository{db: db}
}
