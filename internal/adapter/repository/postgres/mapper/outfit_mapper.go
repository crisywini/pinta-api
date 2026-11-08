package mapper

import (
	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/model"
	"github.com/crisywini/pinta-api/internal/domain"
)

func ToRowOutfit(o *domain.Outfit) model.OutfitRow {
	return model.OutfitRow{
		Name:      o.Name,
		Occasions: o.Occasions,
		Seasons:   o.Seasons,
		Mood:      o.Mood,
		Fragrance: o.Fragrance,
		Notes:     o.Notes,
		Favorite:  o.Favorite,
		WearCount: o.WearCount,
		Photo:     o.Photo,
		LastWorn:  o.LastWorn,
	}
}

func ToDomainOutfit(o model.OutfitRow) domain.Outfit {
	return domain.Outfit{
		ID:        o.ID,
		Name:      o.Name,
		Occasions: o.Occasions,
		Seasons:   o.Seasons,
		Mood:      o.Mood,
		Fragrance: o.Fragrance,
		Notes:     o.Notes,
		Favorite:  o.Favorite,
		WearCount: o.WearCount,
		Photo:     o.Photo,
		LastWorn:  o.LastWorn,
	}
}

func ToDomainListOutfits(rows []model.OutfitRow) []domain.Outfit {
	result := make([]domain.Outfit, 0, len(rows))
	for _, row := range rows {

		result = append(result, ToDomainOutfit(row))
	}
	return result
}
