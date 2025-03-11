package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoosterRepository struct {
	db *pgxpool.Pool
}

func (r *RoosterRepository) Get(ctx context.Context, id domain.ID) (*domain.Rooster, error) {
	rooster := &domain.Rooster{}

	err := r.db.QueryRow(ctx,
		"SELECT id, userid, type FROM rooster WHERE id = $1",
		id).Scan(&rooster.ID, &rooster.UserID)

	if err != nil {
		return nil, err
	}

	return rooster, nil
}

func (r *RoosterRepository) GetUserRoosters(ctx context.Context, userID domain.ID) ([]*domain.Rooster, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, userid, type FROM rooster WHERE userid = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roosters := make([]*domain.Rooster, 0)

	for rows.Next() {
		rooster := &domain.Rooster{}
		if err := rows.Scan(&rooster.ID, &rooster.UserID, &rooster.Type); err != nil {
			return nil, err
		}
		roosters = append(roosters, rooster)
	}

	return roosters, nil
}

func (r *RoosterRepository) Delete(ctx context.Context, id domain.ID) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM rooster WHERE id = $1",
		id)
	if err != nil {
		return err
	}

	return nil
}

func (r *RoosterRepository) Create(ctx context.Context, rooster *domain.Rooster) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO rooster (userid, type, origin) VALUES ($1, $2, $3)",
		rooster.UserID, rooster.Type, rooster.Origin)

	if err != nil {
		return err
	}

	return nil
}

func NewRoosterRepository(db *pgxpool.Pool) port.RoosterRepository {
	return &RoosterRepository{db: db}
}
