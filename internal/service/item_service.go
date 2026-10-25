package service

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/crisywini/pinta-api/internal/domain"
	"github.com/crisywini/pinta-api/internal/port/repository"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9 '\-]+$`)

var validSeasons = map[string]bool{
	"spring": true, "summer": true, "fall": true, "winter": true,
}

var validOccasions = map[string]bool{
	"work": true, "casual": true, "brunch": true, "formal": true,
	"party": true, "sport": true, "date": true,
}

type ItemService struct {
	repository repository.ItemRepository
}

func NewItemService(repository repository.ItemRepository) *ItemService {
	return &ItemService{repository: repository}
}

func (s *ItemService) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	if err := validateItem(item); err != nil {
		return nil, err
	}

	return s.repository.Save(ctx, item)
}

func (s *ItemService) GetById(ctx context.Context, id string) (*domain.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("Missing item id")
	}

	item, err := s.repository.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) GetAll(ctx context.Context) []domain.Item {
	return s.repository.FindAll(ctx)
}

func (s *ItemService) Update(ctx context.Context, id string, updated *domain.Item) error {

	if id == "" {
		return fmt.Errorf("Missing item id")
	}

	if validItem := validateItem(updated); validItem != nil {
		return validItem
	}

	return s.repository.Update(ctx, id, updated)
}

func (s *ItemService) DeleteById(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("Missing item id")
	}
	return s.repository.Delete(ctx, id)
}

func validateItem(item *domain.Item) error {
	var errs []string

	if len(item.Name) < 2 || len(item.Name) > 50 {
		errs = append(errs, "Name must be between 2 and 50 characters")
	} else if !nameRegex.MatchString(item.Name) {
		errs = append(errs, "Name may only contain letters, digits, spaces, hyphens and apostrophes")
	}

	if !item.Category.IsValid() {
		errs = append(errs, "Category must be one of: top, bottom, shoes, hosiery, outerwear, accessory, jewelry, innerwear, fragrance")
	}

	if len(item.PrimaryColor) < 2 || len(item.PrimaryColor) > 30 {
		errs = append(errs, "Primary color must be between 2 and 30 characters")
	}

	if item.Pattern != "" && (len(item.Pattern) < 2 || len(item.Pattern) > 50) {
		errs = append(errs, "Pattern must be between 2 and 50 characters when provided")
	}

	if item.Material != "" && (len(item.Material) < 2 || len(item.Material) > 50) {
		errs = append(errs, "Material must be between 2 and 50 characters when provided")
	}

	for _, s := range item.SeasonTags {
		if !validSeasons[s] {
			errs = append(errs, fmt.Sprintf("invalid season %q: allowed values are spring, summer, fall, winter", s))
		}
	}

	for _, o := range item.Ocassions {
		if !validOccasions[o] {
			errs = append(errs, fmt.Sprintf("invalid occasion %q: allowed values are work, casual, brunch, formal, party, sport, date", o))
		}
	}

	if item.ImageAssetID != "" {
		u, err := url.ParseRequestURI(item.ImageAssetID)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			errs = append(errs, "photo must be a valid http or https URL")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}
	return nil
}
