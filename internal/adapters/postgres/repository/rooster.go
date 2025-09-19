package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RoosterRepository struct {
	*BaseRepository[*domain.Rooster]
}

func (r *RoosterRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Rooster, error) {
	rooster := &domain.Rooster{}

	err := r.db.QueryRow(ctx,
		"SELECT id, userid, type, COALESCE(equip, false), special FROM rooster WHERE id = $1",
		id).Scan(&rooster.ID, &rooster.UserID, &rooster.Type, &rooster.Equip, &rooster.Special)

	return rooster, err
}

func (r *RoosterRepository) GetUserRoosters(ctx context.Context, userID domain.ID) ([]*domain.Rooster, error) {
	return r.GetEntitiesByUserID(ctx, userID,
		"SELECT id, userid, type, COALESCE(equip,false), special FROM rooster WHERE userid = $1 and equip = false",
		func(rows pgx.Rows) (*domain.Rooster, error) {
			rooster := &domain.Rooster{}
			err := rows.Scan(&rooster.ID, &rooster.UserID, &rooster.Type, &rooster.Equip, &rooster.Special)
			return rooster, err
		})
}

func (r *RoosterRepository) GetUserRoosterQuantity(ctx context.Context, userID domain.ID) (int, error) {
	quantity := 1
	err := r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM rooster WHERE userid = $1",
		userID).Scan(&quantity)

	if err != nil {
		return 0, err
	}

	return quantity, err
}

func (r *RoosterRepository) Create(ctx context.Context, rooster *domain.Rooster) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO rooster (userid, type, origin, equip, special) VALUES ($1, $2, $3, false, $4)",
		rooster.UserID, rooster.Type, rooster.Origin, rooster.Special)

	return err
}

func NewRoosterRepository(db postgres.Database) port.RoosterRepository {
	return &RoosterRepository{
		NewBaseRepository[*domain.Rooster](db, "rooster"),
	}
}
