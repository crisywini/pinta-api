package model_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/crisywini/pinta-api/internal/model"
)

func TestNewOutfitBuilder_InitializesEmptySlices(t *testing.T) {
	outfit := model.NewOutfitBuilder().Build()

	if outfit.Occasion == nil {
		t.Error("expected Occasion to be initialized, got nil")
	}
	if len(outfit.Occasion) != 0 {
		t.Errorf("expected Occasion to be empty, got %v", outfit.Occasion)
	}
	if outfit.Season == nil {
		t.Error("expected Season to be initialized, got nil")
	}
	if len(outfit.Season) != 0 {
		t.Errorf("expected Season to be empty, got %v", outfit.Season)
	}
	if outfit.Items == nil {
		t.Error("expected Items to be initialized, got nil")
	}
	if len(outfit.Items) != 0 {
		t.Errorf("expected Items to be empty, got %v", outfit.Items)
	}
}

func TestOutfitBuilder_WithName(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		WithName("Summer Look").
		Build()

	if outfit.Name != "Summer Look" {
		t.Errorf("expected Name %q, got %q", "Summer Look", outfit.Name)
	}
}

func TestOutfitBuilder_WithOccasion(t *testing.T) {
	occasions := []string{"casual", "party"}
	outfit := model.NewOutfitBuilder().
		WithOccasion(occasions).
		Build()

	if !reflect.DeepEqual(outfit.Occasion, occasions) {
		t.Errorf("expected Occasion %v, got %v", occasions, outfit.Occasion)
	}
}

func TestOutfitBuilder_WithOccasion_IsCopy(t *testing.T) {
	occasions := []string{"casual"}
	outfit := model.NewOutfitBuilder().
		WithOccasion(occasions).
		Build()

	occasions[0] = "modified"

	if outfit.Occasion[0] == "modified" {
		t.Error("WithOccasion should store a copy, not a reference to the original slice")
	}
}

func TestOutfitBuilder_AddOccasion(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		AddOccasion("casual").
		AddOccasion("party").
		Build()

	want := []string{"casual", "party"}
	if !reflect.DeepEqual(outfit.Occasion, want) {
		t.Errorf("expected Occasion %v, got %v", want, outfit.Occasion)
	}
}

func TestOutfitBuilder_WithSeason(t *testing.T) {
	seasons := []string{"summer", "spring"}
	outfit := model.NewOutfitBuilder().
		WithSeason(seasons).
		Build()

	if !reflect.DeepEqual(outfit.Season, seasons) {
		t.Errorf("expected Season %v, got %v", seasons, outfit.Season)
	}
}

func TestOutfitBuilder_WithSeason_IsCopy(t *testing.T) {
	seasons := []string{"summer"}
	outfit := model.NewOutfitBuilder().
		WithSeason(seasons).
		Build()

	seasons[0] = "modified"

	if outfit.Season[0] == "modified" {
		t.Error("WithSeason should store a copy, not a reference to the original slice")
	}
}

func TestOutfitBuilder_AddSeason(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		AddSeason("summer").
		AddSeason("spring").
		Build()

	want := []string{"summer", "spring"}
	if !reflect.DeepEqual(outfit.Season, want) {
		t.Errorf("expected Season %v, got %v", want, outfit.Season)
	}
}

func TestOutfitBuilder_WithMood(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		WithMood("confident").
		Build()

	if outfit.Mood != "confident" {
		t.Errorf("expected Mood %q, got %q", "confident", outfit.Mood)
	}
}

func TestOutfitBuilder_WithFragrance(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		WithFragrance("Chanel No. 5").
		Build()

	if outfit.Fragrance != "Chanel No. 5" {
		t.Errorf("expected Fragrance %q, got %q", "Chanel No. 5", outfit.Fragrance)
	}
}

func TestOutfitBuilder_WithNotes(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		WithNotes("Great for summer evenings").
		Build()

	if outfit.Notes != "Great for summer evenings" {
		t.Errorf("expected Notes %q, got %q", "Great for summer evenings", outfit.Notes)
	}
}

func TestOutfitBuilder_WithFavorite(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		WithFavorite(true).
		Build()

	if !outfit.Favorite {
		t.Error("expected Favorite to be true")
	}
}

func TestOutfitBuilder_WithWearCount(t *testing.T) {
	outfit := model.NewOutfitBuilder().
		WithWearCount(3).
		Build()

	if outfit.WearCount != 3 {
		t.Errorf("expected WearCount %d, got %d", 3, outfit.WearCount)
	}
}

func TestOutfitBuilder_WithLastWorn(t *testing.T) {
	lastWorn := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)
	outfit := model.NewOutfitBuilder().
		WithLastWorn(lastWorn).
		Build()

	if !outfit.LastWorn.Equal(lastWorn) {
		t.Errorf("expected LastWorn %v, got %v", lastWorn, outfit.LastWorn)
	}
}

func TestOutfitBuilder_WithItems(t *testing.T) {
	shirt := model.NewItemBuilder("White Shirt", model.Top).Build()
	jeans := model.NewItemBuilder("Blue Jeans", model.Bottom).Build()
	items := []model.Item{shirt, jeans}

	outfit := model.NewOutfitBuilder().
		WithItems(items).
		Build()

	if !reflect.DeepEqual(outfit.Items, items) {
		t.Errorf("expected Items %v, got %v", items, outfit.Items)
	}
}

func TestOutfitBuilder_WithItems_IsCopy(t *testing.T) {
	shirt := model.NewItemBuilder("White Shirt", model.Top).Build()
	items := []model.Item{shirt}

	outfit := model.NewOutfitBuilder().
		WithItems(items).
		Build()

	items[0] = model.NewItemBuilder("Different Shirt", model.Bottom).Build()

	if outfit.Items[0].Name == "Different Shirt" {
		t.Error("WithItems should store a copy, not a reference to the original slice")
	}
}

func TestOutfitBuilder_AddItem(t *testing.T) {
	shirt := model.NewItemBuilder("White Shirt", model.Top).Build()
	jeans := model.NewItemBuilder("Blue Jeans", model.Bottom).Build()

	outfit := model.NewOutfitBuilder().
		AddItem(shirt).
		AddItem(jeans).
		Build()

	if len(outfit.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(outfit.Items))
	}
	if outfit.Items[0].Name != "White Shirt" {
		t.Errorf("expected first item Name %q, got %q", "White Shirt", outfit.Items[0].Name)
	}
	if outfit.Items[1].Name != "Blue Jeans" {
		t.Errorf("expected second item Name %q, got %q", "Blue Jeans", outfit.Items[1].Name)
	}
}

func TestOutfitBuilder_FullChain(t *testing.T) {
	lastWorn := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	shirt := model.NewItemBuilder("White Shirt", model.Top).Build()
	jeans := model.NewItemBuilder("Blue Jeans", model.Bottom).Build()

	outfit := model.NewOutfitBuilder().
		WithName("Casual Friday").
		WithOccasion([]string{"casual"}).
		WithSeason([]string{"summer"}).
		WithMood("relaxed").
		WithFragrance("Sauvage").
		WithNotes("go-to summer outfit").
		WithFavorite(true).
		WithWearCount(7).
		WithLastWorn(lastWorn).
		WithItems([]model.Item{shirt, jeans}).
		Build()

	if outfit.Name != "Casual Friday" {
		t.Errorf("expected Name %q, got %q", "Casual Friday", outfit.Name)
	}
	if !reflect.DeepEqual(outfit.Occasion, []string{"casual"}) {
		t.Errorf("unexpected Occasion: %v", outfit.Occasion)
	}
	if !reflect.DeepEqual(outfit.Season, []string{"summer"}) {
		t.Errorf("unexpected Season: %v", outfit.Season)
	}
	if outfit.Mood != "relaxed" {
		t.Errorf("expected Mood %q, got %q", "relaxed", outfit.Mood)
	}
	if outfit.Fragrance != "Sauvage" {
		t.Errorf("expected Fragrance %q, got %q", "Sauvage", outfit.Fragrance)
	}
	if outfit.Notes != "go-to summer outfit" {
		t.Errorf("expected Notes %q, got %q", "go-to summer outfit", outfit.Notes)
	}
	if !outfit.Favorite {
		t.Error("expected Favorite to be true")
	}
	if outfit.WearCount != 7 {
		t.Errorf("expected WearCount %d, got %d", 7, outfit.WearCount)
	}
	if !outfit.LastWorn.Equal(lastWorn) {
		t.Errorf("expected LastWorn %v, got %v", lastWorn, outfit.LastWorn)
	}
	if len(outfit.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(outfit.Items))
	}
}
