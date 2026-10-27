package repository

import "context"

type Repositories struct {
	Items   ItemRepository
	Outfits OutfitRepository
	Close   func(ctx context.Context) error
}
