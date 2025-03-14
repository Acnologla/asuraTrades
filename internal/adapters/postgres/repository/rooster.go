package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
)

type RoosterRepository struct {
	db postgres.Database
}

func (r *RoosterRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Rooster, error) {
	rooster := &domain.Rooster{}

	err := r.db.QueryRow(ctx,
		"SELECT id, userid, type, equip FROM rooster WHERE id = $1",
		id).Scan(&rooster.ID, &rooster.UserID, &rooster.Type, &rooster.Equip)

	if err != nil {
		return nil, err
	}

	return rooster, nil
}

func (r *RoosterRepository) GetUserRoosters(ctx context.Context, userID domain.ID) ([]*domain.Rooster, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, userid, type, equip FROM rooster WHERE userid = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roosters := make([]*domain.Rooster, 0)

	for rows.Next() {
		rooster := &domain.Rooster{}
		if err := rows.Scan(&rooster.ID, &rooster.UserID, &rooster.Type, &rooster.Equip); err != nil {
			return nil, err
		}
		roosters = append(roosters, rooster)
	}

	return roosters, nil
}

func (r *RoosterRepository) Delete(ctx context.Context, id uuid.UUID) error {
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

func NewRoosterRepository(db postgres.Database) port.RoosterRepository {
	return &RoosterRepository{db: db}
}
