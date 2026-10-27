package postgres

import (
	"context"

	"github.com/crisywini/pinta-api/config"
	"github.com/crisywini/pinta-api/internal/port/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, cfg config.Config) (*repository.Repositories, error) {
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &repository.Repositories{
		Items:   NewPostgresItemRepository(pool),
		Outfits: NewOutfitRepository(pool),
		Close:   func(context.Context) error { pool.Close(); return nil },
	}, nil
}
