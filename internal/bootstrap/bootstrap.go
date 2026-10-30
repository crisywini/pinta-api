package bootstrap

import (
	"context"
	"fmt"

	"github.com/crisywini/pinta-api/config"
	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres"
	"github.com/crisywini/pinta-api/internal/port/repository"
)

type factory func(ctx context.Context, cfg config.Config) (*repository.Repositories, error)

var factories = map[string]factory{
	"mongo":    postgres.New,
	"postgres": postgres.New,
}

func NewRepositories(ctx context.Context, cfg config.Config) (*repository.Repositories, error) {
	f, ok := factories[cfg.DBDriver]
	if !ok {
		return nil, fmt.Errorf("unsupported DB %q", cfg.DBDriver)
	}
	return f(ctx, cfg)
}
