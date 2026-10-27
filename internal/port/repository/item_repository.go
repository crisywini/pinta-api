package repository

import (
	"context"

	"github.com/crisywini/pinta-api/internal/domain"
)

type ItemRepository interface {
	Save(ctx context.Context, item *domain.Item) (*domain.Item, error)
	FindByID(ctx context.Context, id string) (*domain.Item, error)
	ListByUser(ctx context.Context, userID string) ([]domain.Item, error)
	FindAll(ctx context.Context) ([]domain.Item, error)
	Update(ctx context.Context, id string, updated *domain.Item) error
	Delete(ctx context.Context, id string) error
}
