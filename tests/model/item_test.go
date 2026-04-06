package model_test

import (
	"reflect"
	"testing"

	"github.com/crisywini/pinta-api/internal/model"
)

func TestNewItemBuilder(t *testing.T) {
	item := model.NewItemBuilder("White T-Shirt", model.Top).Build()

	if item.Name != "White T-Shirt" {
		t.Errorf("expected Name %q, got %q", "White T-Shirt", item.Name)
	}
	if item.Category != model.Top {
		t.Errorf("expected Category %q, got %q", model.Top, item.Category)
	}
}

func TestItemBuilder_WithColor(t *testing.T) {
	item := model.NewItemBuilder("Jeans", model.Bottom).
		WithColor("Blue").
		Build()

	if item.Color != "Blue" {
		t.Errorf("expected Color %q, got %q", "Blue", item.Color)
	}
}

func TestItemBuilder_WithBrand(t *testing.T) {
	item := model.NewItemBuilder("Sneakers", model.Shoes).
		WithBrand("Nike").
		Build()

	if item.Brand != "Nike" {
		t.Errorf("expected Brand %q, got %q", "Nike", item.Brand)
	}
}

func TestItemBuilder_WithMaterial(t *testing.T) {
	item := model.NewItemBuilder("Shirt", model.Top).
		WithMaterial("Cotton").
		Build()

	if item.Material != "Cotton" {
		t.Errorf("expected Material %q, got %q", "Cotton", item.Material)
	}
}

func TestItemBuilder_WithSeason(t *testing.T) {
	seasons := []string{"summer", "spring"}
	item := model.NewItemBuilder("Shorts", model.Bottom).
		WithSeason(seasons).
		Build()

	if !reflect.DeepEqual(item.Season, seasons) {
		t.Errorf("expected Season %v, got %v", seasons, item.Season)
	}
}

func TestItemBuilder_WithOccasion(t *testing.T) {
	occasions := []string{"casual", "work"}
	item := model.NewItemBuilder("Blazer", model.Top).
		WithOccasion(occasions).
		Build()

	if !reflect.DeepEqual(item.Occasion, occasions) {
		t.Errorf("expected Occasion %v, got %v", occasions, item.Occasion)
	}
}

func TestItemBuilder_WithPhoto(t *testing.T) {
	item := model.NewItemBuilder("Dress", model.Top).
		WithPhoto("https://example.com/photo.jpg").
		Build()

	if item.Photo != "https://example.com/photo.jpg" {
		t.Errorf("expected Photo %q, got %q", "https://example.com/photo.jpg", item.Photo)
	}
}

func TestItemBuilder_WithCondition(t *testing.T) {
	item := model.NewItemBuilder("Coat", model.Outerwear).
		WithCondition("new").
		Build()

	if item.Condition != "new" {
		t.Errorf("expected Condition %q, got %q", "new", item.Condition)
	}
}

func TestItemBuilder_WithWearCount(t *testing.T) {
	item := model.NewItemBuilder("Hat", model.Accesories).
		WithWearCount(5).
		Build()

	if item.WearCount != 5 {
		t.Errorf("expected WearCount %d, got %d", 5, item.WearCount)
	}
}

func TestItemBuilder_FullChain(t *testing.T) {
	seasons := []string{"winter"}
	occasions := []string{"casual"}

	item := model.NewItemBuilder("Wool Sweater", model.Top).
		WithColor("Grey").
		WithBrand("Zara").
		WithMaterial("Wool").
		WithSeason(seasons).
		WithOccasion(occasions).
		WithPhoto("photo.jpg").
		WithCondition("good").
		WithWearCount(10).
		Build()

	if item.Name != "Wool Sweater" {
		t.Errorf("expected Name %q, got %q", "Wool Sweater", item.Name)
	}
	if item.Category != model.Top {
		t.Errorf("expected Category %q, got %q", model.Top, item.Category)
	}
	if item.Color != "Grey" {
		t.Errorf("expected Color %q, got %q", "Grey", item.Color)
	}
	if item.Brand != "Zara" {
		t.Errorf("expected Brand %q, got %q", "Zara", item.Brand)
	}
	if item.Material != "Wool" {
		t.Errorf("expected Material %q, got %q", "Wool", item.Material)
	}
	if !reflect.DeepEqual(item.Season, seasons) {
		t.Errorf("expected Season %v, got %v", seasons, item.Season)
	}
	if !reflect.DeepEqual(item.Occasion, occasions) {
		t.Errorf("expected Occasion %v, got %v", occasions, item.Occasion)
	}
	if item.Photo != "photo.jpg" {
		t.Errorf("expected Photo %q, got %q", "photo.jpg", item.Photo)
	}
	if item.Condition != "good" {
		t.Errorf("expected Condition %q, got %q", "good", item.Condition)
	}
	if item.WearCount != 10 {
		t.Errorf("expected WearCount %d, got %d", 10, item.WearCount)
	}
}

func TestCategory_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		category model.Category
		want     bool
	}{
		{"Top", model.Top, true},
		{"Bottom", model.Bottom, true},
		{"Shoes", model.Shoes, true},
		{"Hosiery", model.Hosiery, true},
		{"Outerwear", model.Outerwear, true},
		{"Accesories", model.Accesories, true},
		{"Jewelry", model.Jewelry, true},
		{"Innerwear", model.Innerwear, true},
		{"Fragance", model.Fragance, true},
		{"Unknown", model.Category("unknown"), false},
		{"EmptyString", model.Category(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.category.IsValid()
			if got != tt.want {
				t.Errorf("Category(%q).IsValid() = %v, want %v", tt.category, got, tt.want)
			}
		})
	}
}
