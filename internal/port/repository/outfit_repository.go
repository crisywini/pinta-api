package repository

import (
	"github.com/crisywini/pinta-api/internal/domain"
)

type OutfitRepository interface {
	Save(item *domain.Outfit) (*domain.Outfit, error)
	FindByID(id string) (*domain.Outfit, error)
	ListByUser(userID string) ([]domain.Outfit, error)
	FindAll() []domain.Outfit
	Update(id string, updated *domain.Outfit) error
	Delete(id string) error
}
