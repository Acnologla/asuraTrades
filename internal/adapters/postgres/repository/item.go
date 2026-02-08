package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ItemRepository struct {
	*BaseRepository[*domain.Item]
}

func (r *ItemRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	item := &domain.Item{}

	err := r.db.QueryRow(ctx,
		"SELECT id, userid, type, itemID, quantity FROM item WHERE id = $1",
		id).Scan(&item.ID, &item.UserID, &item.Type, &item.ItemID, &item.Quantity)

	return item, err
}

func (r *ItemRepository) GetUserItems(ctx context.Context, userID domain.ID) ([]*domain.Item, error) {
	return r.GetEntitiesByUserID(ctx, userID,
		"SELECT id, userid, quantity, itemid, type FROM item WHERE userid = $1",
		func(rows pgx.Rows) (*domain.Item, error) {
			item := &domain.Item{}
			err := rows.Scan(&item.ID, &item.UserID, &item.Quantity, &item.ItemID, &item.Type)
			return item, err
		})
}

func (r *ItemRepository) Add(ctx context.Context, item *domain.Item, quantity int) error {
	cmdTag, err := r.db.Exec(ctx,
		`UPDATE item
		SET quantity = quantity + $4 
		WHERE id = (
			SELECT id 
			FROM item 
			WHERE userid = $1 AND itemid = $2 AND type = $3
			LIMIT 1
		)`,
		item.UserID, item.ItemID, item.Type, quantity)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		_, err = r.db.Exec(ctx,
			`INSERT INTO item (userid, quantity, itemid, type)
			VALUES ($1, $2, $3, $4)`,
			item.UserID, quantity, item.ItemID, item.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ItemRepository) Remove(ctx context.Context, id uuid.UUID, quantity int) error {
	cmdTag, err := r.db.Exec(ctx,
		`UPDATE item 
		SET quantity = quantity - $2 
		WHERE id = $1 AND quantity > $2`,
		id, quantity)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {

		_, err = r.db.Exec(ctx,
			`DELETE FROM item where id = $1`,
			id)

		if err != nil {
			return err
		}
	}
	return nil
}

func NewItemRepository(db postgres.Database) port.ItemRepository {
	return &ItemRepository{
		NewBaseRepository[*domain.Item](db, "item"),
	}
}
