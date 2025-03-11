package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	db *pgxpool.Pool
}

func (r *ItemRepository) Get(ctx context.Context, id domain.ID) (*domain.Item, error) {
	item := &domain.Item{}

	err := r.db.QueryRow(ctx,
		"SELECT id, userid, type FROM items WHERE id = $1",
		id).Scan(&item.ID, &item.UserID)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepository) GetUserItems(ctx context.Context, userID domain.ID) ([]*domain.Item, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, userid, quantity, itemid, type FROM items WHERE userid = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[*domain.Item])
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) Add(ctx context.Context, item *domain.Item) error {
	cmdTag, err := r.db.Exec(ctx,
		`UPDATE items 
		SET quantity = quantity + 1 
		WHERE userid = $1 AND itemid = $2 AND type = $3`,
		item.UserID, item.ItemID, item.Type)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {

		_, err = r.db.Exec(ctx,
			`INSERT INTO items (userid, quantity, itemid, type)
			VALUES ($1, $2, $3, $4)`,
			item.UserID, 1, item.ItemID, item.Type)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ItemRepository) Remove(ctx context.Context, id domain.ID) error {
	cmdTag, err := r.db.Exec(ctx,
		`UPDATE items 
		SET quantity = quantity - 1 
		WHERE id = $1 AND quantity > 1`,
		id)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {

		_, err = r.db.Exec(ctx,
			`DELETE FROM items where id = $1`,
			id)

		if err != nil {
			return err
		}
	}
	return nil
}

func NewItemRepository(db *pgxpool.Pool) port.ItemRepository {
	return &ItemRepository{db: db}
}
