package posgtres

import (
	"context"
	"errors"

	"github.com/crisywini/pinta-api/internal/adapter/repository/posgtres/mapper"
	"github.com/crisywini/pinta-api/internal/adapter/repository/posgtres/model"
	"github.com/crisywini/pinta-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresItemRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const itemColumns = `id, user_id, name, category, sub_category,
	primary_color, secondary_color, pattern, material, fit, silhouette,
	style_tags, occasions, season_tags, image_asset_id, available`

func (r *PostgresRepository) Create(ctx context.Context, item *domain.Item) error {
	row := mapper.ToRow(item)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO items (`+itemColumns+`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		row.ID, row.UserID, row.Name, row.Category, row.SubCategory,
		row.PrimaryColor, row.SecondaryColor, row.Pattern, row.Material,
		row.Fit, row.Silhouette, row.StyleTags, row.Occasions,
		row.SeasonTags, row.ImageAssetID, row.AvailabilityStatus,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+itemColumns+` FROM items WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ItemRow])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound // traducción del error
		}
		return nil, err
	}
	item := mapper.ToDomain(row)
	return &item, nil
}

func (r *PostgresRepository) ListByUser(ctx context.Context, userID string) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+itemColumns+` FROM items WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ItemRow])
	if err != nil {
		return nil, err
	}
	items := make([]domain.Item, len(dbRows))
	for i, row := range dbRows {
		items[i] = mapper.ToDomain(row)
	}
	return items, nil
}

func (r *PostgresRepository) Update(ctx context.Context, item *domain.Item) error {
	row := mapper.ToRow(item)
	tag, err := r.pool.Exec(ctx, `
		UPDATE items SET
			name=$2, category=$3, sub_category=$4, primary_color=$5,
			secondary_color=$6, pattern=$7, material=$8, fit=$9,
			silhouette=$10, style_tags=$11, occasions=$12, season_tags=$13,
			image_asset_id=$14, available=$15
		WHERE id=$1`,
		row.ID, row.Name, row.Category, row.SubCategory, row.PrimaryColor,
		row.SecondaryColor, row.Pattern, row.Material, row.Fit,
		row.Silhouette, row.StyleTags, row.Occasions, row.SeasonTags,
		row.ImageAssetID, row.AvailabilityStatus,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM items WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
