package repository

import (
	"context"

	"github.com/acnologla/asuraTrades/internal/adapters/postgres"
	"github.com/acnologla/asuraTrades/internal/core/domain"
	"github.com/acnologla/asuraTrades/internal/core/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PetRepository struct {
	*BaseRepository[*domain.Pet]
}

func (r *PetRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Pet, error) {
	pet := &domain.Pet{}

	err := r.db.QueryRow(ctx,
		"SELECT id, userid, type, level FROM pet WHERE id = $1",
		id).Scan(&pet.ID, &pet.UserID, &pet.Type, &pet.Level)

	return pet, err
}

func (r *PetRepository) GetUserPets(ctx context.Context, userID domain.ID) ([]*domain.Pet, error) {
	return r.GetEntitiesByUserID(ctx, userID,
		"SELECT id, userid, type, level FROM pet WHERE userid = $1 and (equip = false or equip is null)",
		func(rows pgx.Rows) (*domain.Pet, error) {
			pet := &domain.Pet{}
			err := rows.Scan(&pet.ID, &pet.UserID, &pet.Type, &pet.Level)
			return pet, err
		})
}

func (r *PetRepository) Create(ctx context.Context, pet *domain.Pet) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO pet (userid, type, level) VALUES ($1, $2, $3)",
		pet.UserID, pet.Type, pet.Level)

	return err
}

func NewPetRepository(db postgres.Database) port.PetRepository {
	return &PetRepository{
		NewBaseRepository[*domain.Pet](db, "pet"),
	}
}
