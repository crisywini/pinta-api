package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/crisywini/pinta-api/internal/domain"
	"github.com/crisywini/pinta-api/internal/port/repository"
)

type OutfitService struct {
	repository     repository.OutfitRepository
	itemRepository repository.ItemRepository
}

func NewOutfitService(outfitRepository repository.OutfitRepository, itemRepository repository.ItemRepository) *OutfitService {
	return &OutfitService{
		repository:     outfitRepository,
		itemRepository: itemRepository,
	}
}

func (s *OutfitService) GetById(ctx context.Context, id string) (*domain.Outfit, error) {
	if id == "" {
		return nil, fmt.Errorf("Missing outfit id")
	}

	outfit, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return outfit, nil
}

func (s *OutfitService) GetAll(ctx context.Context) []domain.Outfit {
	outfits, _ := s.repository.FindAll(ctx)
	return outfits
}

func (s *OutfitService) Update(ctx context.Context, id string, updated *domain.Outfit) ([]string, error) {
	if id == "" {
		return nil, fmt.Errorf("Missing outfit id")
	}

	if err := validateOutfit(updated); err != nil {
		return nil, err
	}

	if err := s.verifyItemsInCloset(ctx, updated.Items); err != nil {
		return nil, err
	}

	warnings := crossFieldWarnings(updated)

	return warnings, s.repository.Update(ctx, id, updated)
}

func (s *OutfitService) DeleteById(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("Missing outfit id")
	}

	return s.repository.Delete(ctx, id)
}

func (s *OutfitService) Create(ctx context.Context, outfit *domain.Outfit) (*domain.Outfit, []string, error) {
	if err := validateOutfit(outfit); err != nil {
		return nil, nil, err
	}

	if err := s.verifyItemsInCloset(ctx, outfit.Items); err != nil {
		return nil, nil, err
	}

	warnings := crossFieldWarnings(outfit)

	saved, err := s.repository.Save(ctx, outfit)
	if err != nil {
		return nil, warnings, err
	}
	return saved, warnings, nil
}

func (s *OutfitService) verifyItemsInCloset(ctx context.Context, items []domain.Item) error {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if len(item.ID) == 0 {
			return fmt.Errorf("item %q must be saved to the closet before being added to an outfit", item.Name)
		}
		id := item.ID
		if seen[id] {
			return fmt.Errorf("item %q appears more than once in the outfit", item.Name)
		}
		seen[id] = true
		if _, err := s.itemRepository.FindByID(ctx, id); err != nil {
			return fmt.Errorf("item %q (id: %s) was not found in the closet", item.Name, id)
		}
	}
	return nil
}

func validateOutfit(outfit *domain.Outfit) error {
	var errs []string

	if len(outfit.Name) < 2 || len(outfit.Name) > 50 {
		errs = append(errs, "name must be between 2 and 50 characters")
	} else if !nameRegex.MatchString(outfit.Name) {
		errs = append(errs, "name may only contain letters, digits, spaces, hyphens and apostrophes")
	}

	errs = append(errs, validateOutfitItems(outfit.Items)...)

	for _, o := range outfit.Occasion {
		if !validOccasions[o] {
			errs = append(errs, fmt.Sprintf("invalid occasion %q: allowed values are work, casual, brunch, formal, party, sport, date", o))
		}
	}

	for _, s := range outfit.Season {
		if !validSeasons[s] {
			errs = append(errs, fmt.Sprintf("invalid season %q: allowed values are spring, summer, fall, winter", s))
		}
	}

	if len(outfit.Mood) > 50 {
		errs = append(errs, "mood must be at most 50 characters")
	}

	if len(outfit.Fragrance) > 100 {
		errs = append(errs, "fragrance must be at most 100 characters")
	}

	if len(outfit.Notes) > 300 {
		errs = append(errs, "notes must be at most 300 characters")
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, "; "))
	}
	return nil
}

func validateOutfitItems(items []domain.Item) []string {
	var errs []string

	total := len(items)

	if total < 3 {
		errs = append(errs, "outfit must have at least 3 items (top, bottom and shoes are required)")
		return errs
	}
	if total > 10 {
		errs = append(errs, "outfit must have at most 10 items")
	}

	counts := make(map[domain.Category]int, total)
	for _, item := range items {
		counts[item.Category]++
	}

	if counts[domain.Top] != 1 {
		errs = append(errs, fmt.Sprintf("outfit must have exactly 1 top, got %d", counts[domain.Top]))
	}
	if counts[domain.Bottom] != 1 {
		errs = append(errs, fmt.Sprintf("outfit must have exactly 1 bottom, got %d", counts[domain.Bottom]))
	}
	if counts[domain.Shoes] != 1 {
		errs = append(errs, fmt.Sprintf("outfit must have exactly 1 pair of shoes, got %d", counts[domain.Shoes]))
	}

	if counts[domain.Outerwear] > 1 {
		errs = append(errs, fmt.Sprintf("outfit may have at most 1 outerwear piece, got %d", counts[domain.Outerwear]))
	}
	if counts[domain.Hosiery] > 1 {
		errs = append(errs, fmt.Sprintf("outfit may have at most 1 hosiery piece, got %d", counts[domain.Hosiery]))
	}
	if counts[domain.Innerwear] > 1 {
		errs = append(errs, fmt.Sprintf("outfit may have at most 1 innerwear piece, got %d", counts[domain.Innerwear]))
	}
	if counts[domain.Accesories] > 3 {
		errs = append(errs, fmt.Sprintf("outfit may have at most 3 accessories, got %d", counts[domain.Accesories]))
	}
	if counts[domain.Jewelry] > 3 {
		errs = append(errs, fmt.Sprintf("outfit may have at most 3 jewelry pieces, got %d", counts[domain.Jewelry]))
	}

	return errs
}

func crossFieldWarnings(outfit *domain.Outfit) []string {
	var warnings []string

	mandatory := mandatoryItems(outfit.Items)

	if len(outfit.Season) > 0 {
		outfitSeasons := toSet(outfit.Season)
		for _, item := range mandatory {
			if len(item.SeasonTags) > 0 && !setsOverlap(outfitSeasons, toSet(item.SeasonTags)) {
				warnings = append(warnings, fmt.Sprintf(
					"season mismatch: outfit is tagged %v but %q (%s) is only tagged %v",
					outfit.Season, item.Name, item.Category, item.SeasonTags,
				))
			}
		}
	}

	if len(outfit.Occasion) > 0 {
		itemOccasions := make(map[string]bool)
		for _, item := range mandatory {
			for _, o := range item.Ocassions {
				itemOccasions[o] = true
			}
		}
		if len(itemOccasions) > 0 && !setsOverlap(toSet(outfit.Occasion), itemOccasions) {
			warnings = append(warnings, fmt.Sprintf(
				"occasion mismatch: outfit is tagged %v but none of the mandatory items (top, bottom, shoes) support these occasions",
				outfit.Occasion,
			))
		}
	}

	return warnings
}

func mandatoryItems(items []domain.Item) []domain.Item {
	var out []domain.Item
	for _, item := range items {
		if item.Category == domain.Top || item.Category == domain.Bottom || item.Category == domain.Shoes {
			out = append(out, item)
		}
	}
	return out
}

func toSet(values []string) map[string]bool {
	s := make(map[string]bool, len(values))
	for _, v := range values {
		s[v] = true
	}
	return s
}

func setsOverlap(a, b map[string]bool) bool {
	for k := range a {
		if b[k] {
			return true
		}
	}
	return false
}
