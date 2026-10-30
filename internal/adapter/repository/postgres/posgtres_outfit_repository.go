package postgres

import (
	"github.com/crisywini/pinta-api/internal/port/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.OutfitRepository = (*OutfitPostgresRepository)(nil)

type OutfitPostgresRepository struct {
	pool *pgxpool.Pool
}
