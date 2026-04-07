package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/repository"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9 '\-]+$`)

var validSeasons = map[string]bool{
	"spring": true, "summer": true, "fall": true, "winter": true,
}

var validOccasions = map[string]bool{
	"work": true, "casual": true, "brunch": true, "formal": true,
	"party": true, "sport": true, "date": true,
}

var validConditions = map[string]bool{
	"new": true, "good": true, "fair": true, "retired": true,
}

type ItemService struct {
	repository *repository.ItemRepository
}

func NewItemService(repository *repository.ItemRepository) *ItemService {
	return &ItemService{repository: repository}
}

func (s *ItemService) Create(item *model.Item) (*model.Item, error) {
	if err := validateItem(item); err != nil {
		return nil, err
	}
	applyItemDefaults(item)
	return s.repository.Save(item)
}

func (s *ItemService) GetById(id string) (*model.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("Missing item id")
	}

	item, err := s.repository.FindByID(id)

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) GetAll() ([]model.Item, error) {
	return s.repository.FindAll()
}

func (s *ItemService) Update(id string, updated *model.Item) error {

	if id == "" {
		return fmt.Errorf("Missing item id")
	}

	if validItem := validateItem(updated); validItem != nil {
		return validItem
	}

	return s.repository.Update(id, updated)
}

func (s *ItemService) DeleteById(id string) error {
	if id == "" {
		return fmt.Errorf("Missing item id")
	}
	return s.repository.Delete(id)
}

func applyItemDefaults(item *model.Item) {
	if item.Condition == "" {
		item.Condition = "good"
	}
}

func validateItem(item *model.Item) error {
	var errs []string

	if len(item.Name) < 2 || len(item.Name) > 50 {
		errs = append(errs, "name must be between 2 and 50 characters")
	} else if !nameRegex.MatchString(item.Name) {
		errs = append(errs, "name may only contain letters, digits, spaces, hyphens and apostrophes")
	}

	if !item.Category.IsValid() {
		errs = append(errs, "category must be one of: top, bottom, shoes, hosiery, outerwear, accessory, jewelry, innerwear, fragrance")
	}

	if len(item.Color) < 2 || len(item.Color) > 30 {
		errs = append(errs, "color must be between 2 and 30 characters")
	}

	if item.Brand != "" && (len(item.Brand) < 2 || len(item.Brand) > 50) {
		errs = append(errs, "brand must be between 2 and 50 characters when provided")
	}

	if item.Material != "" && (len(item.Material) < 2 || len(item.Material) > 50) {
		errs = append(errs, "material must be between 2 and 50 characters when provided")
	}

	for _, s := range item.Season {
		if !validSeasons[s] {
			errs = append(errs, fmt.Sprintf("invalid season %q: allowed values are spring, summer, fall, winter", s))
		}
	}

	for _, o := range item.Occasion {
		if !validOccasions[o] {
			errs = append(errs, fmt.Sprintf("invalid occasion %q: allowed values are work, casual, brunch, formal, party, sport, date", o))
		}
	}

	if item.Photo != "" {
		u, err := url.ParseRequestURI(item.Photo)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
			errs = append(errs, "photo must be a valid http or https URL")
		}
	}

	if item.Condition != "" && !validConditions[item.Condition] {
		errs = append(errs, "condition must be one of: new, good, fair, retired")
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}
	return nil
}
