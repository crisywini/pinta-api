package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/mapper"
	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/model"
	"github.com/crisywini/pinta-api/internal/domain"
	"github.com/crisywini/pinta-api/internal/port/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ repository.OutfitRepository = (*OutfitPostgresRepository)(nil)

type OutfitPostgresRepository struct {
	pool           *pgxpool.Pool
	itemRepository *ItemPostgresRepository
}

func NewPostgresOutfitRepository(pool *pgxpool.Pool, itemRepository *ItemPostgresRepository) *OutfitPostgresRepository {
	return &OutfitPostgresRepository{pool: pool, itemRepository: itemRepository}
}

const outfitColumns = `name, occasions, seasons, mood,
	fragrance, notes, favorite, wear_count, photo, last_worn`

const outfitColumnsWithId = `id, name, occasions, seasons, mood,
	fragrance, notes, favorite, wear_count, photo, last_worn`

const outfit_itemsColumns = `outfit_id, item_id`

func (r *OutfitPostgresRepository) Save(ctx context.Context, outfit *domain.Outfit) (*domain.Outfit, error) {
	row := mapper.ToRowOutfit(outfit)
	err := r.pool.QueryRow(ctx, `
		INSERT INTO outfits (`+outfitColumns+`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING id`,
		row.Name, row.Occasions, row.Seasons, row.Mood, row.Fragrance,
		row.Notes, row.Favorite, row.WearCount, row.Photo, row.LastWorn,
	).Scan(&row.ID)

	if err != nil {
		return nil, err
	}

	for _, item := range outfit.Items {
		_, err := r.pool.Exec(ctx, `
		INSERT INTO outfit_items (`+outfit_itemsColumns+`)
		VALUES ($1,$2)`,
			row.ID, item.ID,
		)
		if err != nil {
			return nil, err
		}
	}
	saved := mapper.ToDomainOutfit(row)
	saved.Items = outfit.Items
	return &saved, err
}

func (r *OutfitPostgresRepository) FindByID(ctx context.Context, id string) (*domain.Outfit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+outfitColumnsWithId+` FROM outfits WHERE id = $1`, id)

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
	outfit.Items, err = r.itemRepository.FindAllByOutfitID(ctx, id)

	if err != nil {
		return nil, err
	}
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

	outfits := mapper.ToDomainListOutfits(dbRows)
	return outfits, nil
}

func (r *OutfitPostgresRepository) Update(ctx context.Context, id string, outfit *domain.Outfit) error {
	row := mapper.ToRowOutfit(outfit)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		UPDATE outfits SET
			name=$2, occasions=$3, seasons=$4,
			mood=$5, fragrance=$6, notes=$7, favorite=$8,
			wear_count=$9, photo=$10, last_worn=$11
		WHERE id=$1`,
		id, row.Name, row.Occasions, row.Seasons, row.Mood, row.Fragrance,
		row.Notes, row.Favorite, row.WearCount, row.Photo, row.LastWorn,
	)
	if err != nil {
		return fmt.Errorf("update outfit: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	if _, err := tx.Exec(ctx,
		`DELETE FROM outfit_items WHERE outfit_id = $1`, id); err != nil {
		return fmt.Errorf("clear outfit_items: %w", err)
	}
	for _, item := range outfit.Items {
		if _, err := tx.Exec(ctx,
			`INSERT INTO outfit_items (outfit_id, item_id) VALUES ($1, $2)`,
			id, item.ID,
		); err != nil {
			return fmt.Errorf("insert outfit_item %s: %w", item.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (r *OutfitPostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM outfit WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
