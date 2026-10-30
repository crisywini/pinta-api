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

var _ repository.ItemRepository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresItemRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const itemColumns = `id, user_id, name, category, sub_category,
	primary_color, secondary_color, pattern, material, fit, silhouette,
	style_tags, occasions, season_tags, image_asset_id, available`

func (r *PostgresRepository) Save(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	row := mapper.ToRowItem(item)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO items (`+itemColumns+`)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		row.ID, row.UserID, row.Name, row.Category, row.SubCategory,
		row.PrimaryColor, row.SecondaryColor, row.Pattern, row.Material,
		row.Fit, row.Silhouette, row.StyleTags, row.Occasions,
		row.SeasonTags, row.ImageAssetID, row.AvailabilityStatus,
	)
	saved := mapper.ToDomainItem(row)
	return &saved, err
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*domain.Item, error) {
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
	item := mapper.ToDomainItem(row)
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
		items[i] = mapper.ToDomainItem(row)
	}
	return items, nil
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]domain.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT * FROM items WHERE id = $1`)
	if err != nil {
		return nil, err
	}

	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ItemRow])
	if err != nil {
		return nil, err
	}

	items := mapper.ToDomainListItem(dbRows)
	return items, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, item *domain.Item) error {
	row := mapper.ToRowItem(item)
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
