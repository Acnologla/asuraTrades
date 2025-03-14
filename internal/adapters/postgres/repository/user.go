package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db postgres.Database
}

func (r *UserRepository) Get(ctx context.Context, id domain.ID) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.QueryRow(ctx,
		"SELECT id, xp FROM users WHERE id = $1",
		id).Scan(&user.ID, &user.Xp)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository(db postgres.Database) port.UserRepository {
	return &UserRepository{db: db}
}

type TradeTransactionProvider struct {
	db *pgxpool.Pool
}

func (p *TradeTransactionProvider) Transact(ctx context.Context, txFunc func(adapters port.UserTradeTxAdapters, lock func(domain.ID) error) error) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	lock := func(id domain.ID) error {
		_, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", id)
		return err
	}
	adapters := port.UserTradeTxAdapters{
		UserRepository:    NewUserRepository(tx),
		RoosterRepository: NewRoosterRepository(tx),
		ItemRepository:    NewItemRepository(tx),
	}

	err = txFunc(adapters, lock)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func NewTransactionProvider(db *pgxpool.Pool) port.TradeTxProvider {
	return &TradeTransactionProvider{
		db: db,
	}
}
