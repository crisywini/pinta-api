package repository

import (
	"github.com/crisywini/pinta-api/internal/domain"
)

type ItemRepository interface {
	Save(item *domain.Item) (*domain.Item, error)
	FindByID(id string) (*domain.Item, error)
	ListByUser(userID string) ([]domain.Item, error)
	FindAll() []domain.Item
	Update(id string, updated *domain.Item) error
	Delete(id string) error
}
