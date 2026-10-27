package mapper

import (
	"github.com/crisywini/pinta-api/internal/adapter/repository/postgres/model"
	"github.com/crisywini/pinta-api/internal/domain"
)

func ToRow(i *domain.Item) model.ItemRow {
	return model.ItemRow{
		ID:                 i.ID,
		UserID:             i.UserID,
		Name:               i.Name,
		Category:           string(i.Category),
		SubCategory:        i.SubCategory,
		PrimaryColor:       i.PrimaryColor,
		SecondaryColor:     i.SecondaryColor,
		Pattern:            i.Pattern,
		Material:           i.Material,
		Fit:                i.Fit,
		Silhouette:         i.Silhouette,
		StyleTags:          i.StyleTags,
		Occasions:          i.Ocassions,
		SeasonTags:         i.SeasonTags,
		ImageAssetID:       i.ImageAssetID,
		AvailabilityStatus: i.AvailabilityStatus,
	}
}

func ToDomain(r model.ItemRow) domain.Item {
	return domain.Item{
		ID:                 r.ID,
		UserID:             r.UserID,
		Name:               r.Name,
		Category:           domain.Category(r.Category),
		SubCategory:        r.SubCategory,
		PrimaryColor:       r.PrimaryColor,
		SecondaryColor:     r.SecondaryColor,
		Pattern:            r.Pattern,
		Material:           r.Material,
		Fit:                r.Fit,
		Silhouette:         r.Silhouette,
		StyleTags:          r.StyleTags,
		Ocassions:          r.Occasions,
		SeasonTags:         r.SeasonTags,
		ImageAssetID:       r.ImageAssetID,
		AvailabilityStatus: r.AvailabilityStatus,
	}
}

func ToDomainList(rows []model.ItemRow) []domain.Item {
	result := make([]domain.Item, 0, len(rows))
	for _, row := range rows {

		result = append(result, ToDomain(row))
	}
	return result
}
