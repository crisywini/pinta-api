package postgres

import (
	"context"
	"errors"

	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/mapper"
	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/model"
	"github.com/crisywini/pinta-api/internal/domain"
	"github.com/crisywini/pinta-api/internal/port/repository"
	"github.com/jackc/pgx/v5"
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

const outfits_itemsColumns = `outfit_id, item_id`

func (r *OutfitPostgresRepository) Save(ctx context.Context, outfit *domain.Outfit) (*domain.Outfit, error) {
	row := mapper.ToRowOutfit(outfit)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO outfits (`+outfitColumns+`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		row.ID, row.Name, row.Occasion, row.Season, row.Mood, row.Fragrance,
		row.Notes, row.Favorite, row.WearCount, row.Photo, row.LastWorn,
	)
	// Now we need to save the outfit_items

	for _, item := range outfit.Items {
		_, err := r.pool.Exec(ctx, `
		INSERT INTO outfits_items (`+outfits_itemsColumns+`)
		VALUES ($1,$2)`,
			outfit.ID, item.ID,
		)
		if err != nil {

		}
	}
	saved := mapper.ToDomainOutfit(row)
	return &saved, err
}

func (r *OutfitPostgresRepository) FindByID(ctx context.Context, id string) (*domain.Outfit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+itemColumns+` FROM outfits WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.OutfitRow])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	outfit := mapper.ToDomainOutfit(row)
	return &outfit, nil
}

func (r *OutfitPostgresRepository) FindAll(ctx context.Context) ([]domain.Outfit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT * FROM outfits`)
	if err != nil {
		return nil, err
	}

	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.OutfitRow])
	if err != nil {
		return nil, err
	}

	items := mapper.ToDomainListOutfits(dbRows)
	return items, nil
}

func (r *OutfitPostgresRepository) Update(ctx context.Context, id string, outfit *domain.Outfit) error {
	row := mapper.ToRowOutfit(outfit)
	tag, err := r.pool.Exec(ctx, `
		UPDATE items SET
			name=$2, ocassions=$3, season=$4, season=$5,
			mood=$6, fragrance=$7, notes=$8, favorite=$9,
			wear_count=$10, photo=$11, last_worn=$12
		WHERE id=$1`,
		row.ID, row.Name, row.Occasion, row.Season, row.Mood, row.Fragrance, row.Notes,
		row.Favorite, row.WearCount, row.Photo, row.LastWorn,
	)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
