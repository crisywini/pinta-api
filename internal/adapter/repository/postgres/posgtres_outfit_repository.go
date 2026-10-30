package postgres

import (
	"context"

	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/mapper"
	"github.com/crisywini/pinta-api/internal/domain"
	"github.com/crisywini/pinta-api/internal/port/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.OutfitRepository = (*OutfitPostgresRepository)(nil)

type OutfitPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOutfitRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const outfitColumns = `id, name, ocassions, season, mood,
	fragrance, notes, favorite, wear_count, photo, last_worn`

const outfit_itemsColumns = `outfit_id, item_id`

func (r *OutfitPostgresRepository) Save(ctx context.Context, outfit *domain.Outfit) (*domain.Outfit, error) {
	row := mapper.ToRowOutfit(outfit)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO outfit (`+outfitColumns+`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		row.ID, row.Name, row.Occasion, row.Season, row.Mood, row.Fragrance,
		row.Notes, row.Favorite, row.WearCount, row.Photo, row.LastWorn,
	)
	// Now we need to save the outfit_items

	for _, item := range outfit.Items {
		_, err := r.pool.Exec(ctx, `
		INSERT INTO outfit_items (`+outfit_itemsColumns+`)
		VALUES ($1,$2)`,
			outfit.ID, item.ID,
		)
		if err != nil {

		}
	}
	saved := mapper.ToDomainOutfit(row)
	return &saved, err
}
