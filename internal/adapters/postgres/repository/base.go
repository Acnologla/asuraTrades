package repository

import (
	"context"
	"fmt"

	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Entity interface {
	GetID() uuid.UUID
}

type BaseRepository[T Entity] struct {
	db        postgres.Database
	tableName string
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.tableName)

	_, err := r.db.Exec(ctx,
		query,
		id)

	return err
}

func (r *BaseRepository[T]) GetEntitiesByUserID(
	ctx context.Context,
	userID domain.ID,
	query string,
	scanRow func(rows pgx.Rows) (T, error),
) ([]T, error) {
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]T, 0)
	for rows.Next() {
		entity, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func NewBaseRepository[T Entity](db postgres.Database, tableName string) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:        db,
		tableName: tableName,
	}
}
