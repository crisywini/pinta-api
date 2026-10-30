package repository

import (
	"context"

	"github.com/crisywini/pinta-api/internal/domain"
)

type OutfitRepository interface {
	Save(ctx context.Context, item *domain.Outfit) (*domain.Outfit, error)
	FindByID(ctx context.Context, id string) (*domain.Outfit, error)
	FindAll(ctx context.Context) ([]domain.Outfit, error)
	Update(ctx context.Context, id string, updated *domain.Outfit) error
	Delete(ctx context.Context, id string) error
}
