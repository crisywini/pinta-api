package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/repository"
	"github.com/crisywini/pinta-api/internal/service"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// setupOutfitService spins up a real MongoDB container and returns the outfit service,
// the raw item repository (for pre-seeding closet items), and a cleanup func.
func setupOutfitService(t *testing.T) (*service.OutfitService, *repository.ItemRepository, func()) {
	t.Helper()
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatalf("failed to start mongodb container: %v", err)
	}

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to get connection string: %v", err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatalf("failed to connect to mongodb: %v", err)
	}

	db := client.Database("pinta_test")
	itemRepo := repository.NewItemRepository(db)
	outfitRepo := repository.NewOutfitRepository(db)
	svc := service.NewOutfitService(outfitRepo, itemRepo)

	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		_ = container.Terminate(ctx)
	}

	return svc, itemRepo, cleanup
}

// seedItem saves an item directly to the item repo and returns the persisted copy (with ID set).
func seedItem(t *testing.T, repo *repository.ItemRepository, item model.Item) model.Item {
	t.Helper()
	saved, err := repo.Save(&item)
	if err != nil {
		t.Fatalf("seedItem: failed to save item %q: %v", item.Name, err)
	}
	return *saved
}

// minimalOutfit builds the smallest valid outfit (top + bottom + shoes).
func minimalOutfit(name string, top, bottom, shoes model.Item) model.Outfit {
	return model.NewOutfitBuilder().
		WithName(name).
		WithItems([]model.Item{top, bottom, shoes}).
		Build()
}

func TestOutfitService_Create(t *testing.T) {
	svc, itemRepo, cleanup := setupOutfitService(t)
	defer cleanup()

	// Seed the three mandatory items once — reused across all subtests.
	top := seedItem(t, itemRepo,
		model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build())
	bottom := seedItem(t, itemRepo,
		model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build())
	shoes := seedItem(t, itemRepo,
		model.NewItemBuilder("White Sneakers", model.Shoes).WithColor("White").Build())

	// ── happy path ────────────────────────────────────────────────────────────

	t.Run("saves minimal valid outfit and assigns id", func(t *testing.T) {
		outfit := minimalOutfit("Casual Friday", top, bottom, shoes)
		saved, _, err := svc.Create(&outfit)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if saved.ID.IsZero() {
			t.Fatal("expected ID to be set after Create, got zero value")
		}
	})

	t.Run("saves outfit with all optional items and fields", func(t *testing.T) {
		outerwear := seedItem(t, itemRepo,
			model.NewItemBuilder("Black Jacket", model.Outerwear).WithColor("Black").Build())
		hosiery := seedItem(t, itemRepo,
			model.NewItemBuilder("White Socks", model.Hosiery).WithColor("White").Build())
		acc := seedItem(t, itemRepo,
			model.NewItemBuilder("Brown Belt", model.Accesories).WithColor("Brown").Build())

		outfit := model.NewOutfitBuilder().
			WithName("Full Look").
			WithItems([]model.Item{top, bottom, shoes, outerwear, hosiery, acc}).
			WithMood("sharp").
			AddOccasion("work").
			AddSeason("fall").
			WithNotes("My go-to power outfit.").
			Build()

		saved, warnings, err := svc.Create(&outfit)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if saved.ID.IsZero() {
			t.Fatal("expected ID to be set, got zero value")
		}
		if len(warnings) != 0 {
			t.Errorf("expected no warnings, got: %v", warnings)
		}
	})

	t.Run("accepts name with hyphens and apostrophes", func(t *testing.T) {
		outfit := minimalOutfit("O'Brien's Day Look", top, bottom, shoes)
		if _, _, err := svc.Create(&outfit); err != nil {
			t.Errorf("Create() unexpected error for valid name: %v", err)
		}
	})

	// ── name validation ───────────────────────────────────────────────────────

	t.Run("name/too short", func(t *testing.T) {
		outfit := minimalOutfit("x", top, bottom, shoes)
		assertOutfitValidationError(t, svc, outfit, "name must be between 2 and 50 characters")
	})

	t.Run("name/too long", func(t *testing.T) {
		outfit := minimalOutfit(strings.Repeat("a", 51), top, bottom, shoes)
		assertOutfitValidationError(t, svc, outfit, "name must be between 2 and 50 characters")
	})

	t.Run("name/invalid characters", func(t *testing.T) {
		outfit := minimalOutfit("Look#1!", top, bottom, shoes)
		assertOutfitValidationError(t, svc, outfit, "name may only contain letters")
	})

	// ── item list: total count ────────────────────────────────────────────────

	t.Run("items/empty list", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().WithName("Empty").Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have at least 3 items")
	})

	t.Run("items/fewer than 3", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Too Few").
			WithItems([]model.Item{
				{Name: "White T-Shirt", Category: model.Top, Color: "White"},
				{Name: "Blue Jeans", Category: model.Bottom, Color: "Blue"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have at least 3 items")
	})

	t.Run("items/more than 10", func(t *testing.T) {
		// Build a structurally valid but oversized list (no real IDs needed — count check fires first).
		items := []model.Item{
			{Name: "Top", Category: model.Top, Color: "White"},
			{Name: "Bottom", Category: model.Bottom, Color: "Blue"},
			{Name: "Shoes", Category: model.Shoes, Color: "White"},
			{Name: "Outerwear", Category: model.Outerwear, Color: "Black"},
			{Name: "Hosiery", Category: model.Hosiery, Color: "White"},
			{Name: "Innerwear", Category: model.Innerwear, Color: "White"},
			{Name: "Acc 1", Category: model.Accesories, Color: "Gold"},
			{Name: "Acc 2", Category: model.Accesories, Color: "Silver"},
			{Name: "Acc 3", Category: model.Accesories, Color: "Bronze"},
			{Name: "Jewelry 1", Category: model.Jewelry, Color: "Gold"},
			{Name: "Jewelry 2", Category: model.Jewelry, Color: "Silver"},
		}
		outfit := model.NewOutfitBuilder().WithName("Too Many").WithItems(items).Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have at most 10 items")
	})

	// ── item list: mandatory categories ──────────────────────────────────────

	t.Run("items/missing top", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("No Top").
			WithItems([]model.Item{
				{Name: "Blue Jeans", Category: model.Bottom, Color: "Blue"},
				{Name: "Sneakers", Category: model.Shoes, Color: "White"},
				{Name: "Socks", Category: model.Hosiery, Color: "White"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have exactly 1 top, got 0")
	})

	t.Run("items/missing bottom", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("No Bottom").
			WithItems([]model.Item{
				{Name: "White T-Shirt", Category: model.Top, Color: "White"},
				{Name: "Sneakers", Category: model.Shoes, Color: "White"},
				{Name: "Socks", Category: model.Hosiery, Color: "White"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have exactly 1 bottom, got 0")
	})

	t.Run("items/missing shoes", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("No Shoes").
			WithItems([]model.Item{
				{Name: "White T-Shirt", Category: model.Top, Color: "White"},
				{Name: "Blue Jeans", Category: model.Bottom, Color: "Blue"},
				{Name: "Socks", Category: model.Hosiery, Color: "White"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have exactly 1 pair of shoes, got 0")
	})

	t.Run("items/two tops", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Two Tops").
			WithItems([]model.Item{
				{Name: "White T-Shirt", Category: model.Top, Color: "White"},
				{Name: "Blue T-Shirt", Category: model.Top, Color: "Blue"},
				{Name: "Blue Jeans", Category: model.Bottom, Color: "Blue"},
				{Name: "Sneakers", Category: model.Shoes, Color: "White"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have exactly 1 top, got 2")
	})

	t.Run("items/two bottoms", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Two Bottoms").
			WithItems([]model.Item{
				{Name: "White T-Shirt", Category: model.Top, Color: "White"},
				{Name: "Blue Jeans", Category: model.Bottom, Color: "Blue"},
				{Name: "Black Trousers", Category: model.Bottom, Color: "Black"},
				{Name: "Sneakers", Category: model.Shoes, Color: "White"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have exactly 1 bottom, got 2")
	})

	t.Run("items/two pairs of shoes", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Two Shoes").
			WithItems([]model.Item{
				{Name: "White T-Shirt", Category: model.Top, Color: "White"},
				{Name: "Blue Jeans", Category: model.Bottom, Color: "Blue"},
				{Name: "Sneakers", Category: model.Shoes, Color: "White"},
				{Name: "Boots", Category: model.Shoes, Color: "Black"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit must have exactly 1 pair of shoes, got 2")
	})

	// ── item list: optional category limits ───────────────────────────────────

	t.Run("items/two outerwear pieces", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Double Coat").
			WithItems([]model.Item{
				{Name: "Top", Category: model.Top, Color: "White"},
				{Name: "Bottom", Category: model.Bottom, Color: "Blue"},
				{Name: "Shoes", Category: model.Shoes, Color: "White"},
				{Name: "Jacket", Category: model.Outerwear, Color: "Black"},
				{Name: "Coat", Category: model.Outerwear, Color: "Beige"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit may have at most 1 outerwear piece")
	})

	t.Run("items/two hosiery pieces", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Double Socks").
			WithItems([]model.Item{
				{Name: "Top", Category: model.Top, Color: "White"},
				{Name: "Bottom", Category: model.Bottom, Color: "Blue"},
				{Name: "Shoes", Category: model.Shoes, Color: "White"},
				{Name: "Socks 1", Category: model.Hosiery, Color: "White"},
				{Name: "Socks 2", Category: model.Hosiery, Color: "Black"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit may have at most 1 hosiery piece")
	})

	t.Run("items/two innerwear pieces", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Double Innerwear").
			WithItems([]model.Item{
				{Name: "Top", Category: model.Top, Color: "White"},
				{Name: "Bottom", Category: model.Bottom, Color: "Blue"},
				{Name: "Shoes", Category: model.Shoes, Color: "White"},
				{Name: "Tank Top", Category: model.Innerwear, Color: "White"},
				{Name: "Bralette", Category: model.Innerwear, Color: "Beige"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit may have at most 1 innerwear piece")
	})

	t.Run("items/four accessories", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Too Many Accessories").
			WithItems([]model.Item{
				{Name: "Top", Category: model.Top, Color: "White"},
				{Name: "Bottom", Category: model.Bottom, Color: "Blue"},
				{Name: "Shoes", Category: model.Shoes, Color: "White"},
				{Name: "Belt", Category: model.Accesories, Color: "Brown"},
				{Name: "Bag", Category: model.Accesories, Color: "Black"},
				{Name: "Scarf", Category: model.Accesories, Color: "Red"},
				{Name: "Hat", Category: model.Accesories, Color: "Beige"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit may have at most 3 accessories")
	})

	t.Run("items/four jewelry pieces", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Too Much Jewelry").
			WithItems([]model.Item{
				{Name: "Top", Category: model.Top, Color: "White"},
				{Name: "Bottom", Category: model.Bottom, Color: "Blue"},
				{Name: "Shoes", Category: model.Shoes, Color: "White"},
				{Name: "Ring", Category: model.Jewelry, Color: "Gold"},
				{Name: "Necklace", Category: model.Jewelry, Color: "Gold"},
				{Name: "Bracelet", Category: model.Jewelry, Color: "Silver"},
				{Name: "Earrings", Category: model.Jewelry, Color: "Gold"},
			}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "outfit may have at most 3 jewelry pieces")
	})

	// ── closet verification ───────────────────────────────────────────────────

	t.Run("items/zero id rejected", func(t *testing.T) {
		zeroTop := model.Item{Name: "Unsaved Top", Category: model.Top, Color: "White"}
		outfit := minimalOutfit("Zero ID", zeroTop, bottom, shoes)
		assertOutfitValidationError(t, svc, outfit, "must be saved to the closet before being added")
	})

	t.Run("items/duplicate item rejected", func(t *testing.T) {
		// Use two occurrences of the same accessory (category allows 0-3 so counts pass).
		acc := seedItem(t, itemRepo,
			model.NewItemBuilder("Brown Belt", model.Accesories).WithColor("Brown").Build())

		outfit := model.NewOutfitBuilder().
			WithName("Duplicate Acc").
			WithItems([]model.Item{top, bottom, shoes, acc, acc}).
			Build()
		assertOutfitValidationError(t, svc, outfit, "appears more than once in the outfit")
	})

	t.Run("items/not in closet rejected", func(t *testing.T) {
		ghostTop := model.Item{
			Name:     "Ghost Top",
			Category: model.Top,
			Color:    "White",
		}
		ghostTop.ID = primitive.NewObjectID() // has ID but was never saved

		outfit := minimalOutfit("Ghost Outfit", ghostTop, bottom, shoes)
		assertOutfitValidationError(t, svc, outfit, "was not found in the closet")
	})

	// ── optional outfit fields ────────────────────────────────────────────────

	t.Run("occasion/invalid value", func(t *testing.T) {
		outfit := minimalOutfit("Bad Occasion", top, bottom, shoes)
		outfit.Occasion = []string{"casual", "gala"}
		assertOutfitValidationError(t, svc, outfit, `invalid occasion "gala"`)
	})

	t.Run("season/invalid value", func(t *testing.T) {
		outfit := minimalOutfit("Bad Season", top, bottom, shoes)
		outfit.Season = []string{"summer", "monsoon"}
		assertOutfitValidationError(t, svc, outfit, `invalid season "monsoon"`)
	})

	t.Run("mood/too long", func(t *testing.T) {
		outfit := minimalOutfit("Long Mood", top, bottom, shoes)
		outfit.Mood = strings.Repeat("a", 51)
		assertOutfitValidationError(t, svc, outfit, "mood must be at most 50 characters")
	})

	t.Run("fragrance/too long", func(t *testing.T) {
		outfit := minimalOutfit("Long Fragrance", top, bottom, shoes)
		outfit.Fragrance = strings.Repeat("a", 101)
		assertOutfitValidationError(t, svc, outfit, "fragrance must be at most 100 characters")
	})

	t.Run("notes/too long", func(t *testing.T) {
		outfit := minimalOutfit("Long Notes", top, bottom, shoes)
		outfit.Notes = strings.Repeat("a", 301)
		assertOutfitValidationError(t, svc, outfit, "notes must be at most 300 characters")
	})

	// ── cross-field warnings (warn, don't block) ──────────────────────────────

	t.Run("warnings/season mismatch triggers warning but still saves", func(t *testing.T) {
		winterTop := seedItem(t, itemRepo,
			model.NewItemBuilder("Wool Turtleneck", model.Top).
				WithColor("Grey").
				WithSeason([]string{"winter"}).
				Build())

		outfit := model.NewOutfitBuilder().
			WithName("Summer Wool").
			WithItems([]model.Item{winterTop, bottom, shoes}).
			AddSeason("summer").
			Build()

		saved, warnings, err := svc.Create(&outfit)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if saved.ID.IsZero() {
			t.Fatal("expected outfit to be saved despite warning")
		}
		assertWarning(t, warnings, "season mismatch")
	})

	t.Run("warnings/no season warning when mandatory items have no seasons", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Seasonless").
			WithItems([]model.Item{top, bottom, shoes}). // top/bottom/shoes have no seasons set
			AddSeason("summer").
			Build()

		_, warnings, err := svc.Create(&outfit)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if len(warnings) != 0 {
			t.Errorf("expected no warnings when items have no seasons, got: %v", warnings)
		}
	})

	t.Run("warnings/occasion mismatch triggers warning but still saves", func(t *testing.T) {
		casualTop := seedItem(t, itemRepo,
			model.NewItemBuilder("Linen Shirt", model.Top).
				WithColor("Beige").
				WithOccasion([]string{"casual", "brunch"}).
				Build())
		casualBottom := seedItem(t, itemRepo,
			model.NewItemBuilder("Linen Trousers", model.Bottom).
				WithColor("White").
				WithOccasion([]string{"casual"}).
				Build())
		casualShoes := seedItem(t, itemRepo,
			model.NewItemBuilder("Sandals", model.Shoes).
				WithColor("Beige").
				WithOccasion([]string{"casual"}).
				Build())

		outfit := model.NewOutfitBuilder().
			WithName("Casual Outfit at Work").
			WithItems([]model.Item{casualTop, casualBottom, casualShoes}).
			AddOccasion("work"). // outfit tagged work, but all items tagged casual/brunch
			Build()

		saved, warnings, err := svc.Create(&outfit)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if saved.ID.IsZero() {
			t.Fatal("expected outfit to be saved despite warning")
		}
		assertWarning(t, warnings, "occasion mismatch")
	})

	t.Run("warnings/no occasion warning when mandatory items have no occasions", func(t *testing.T) {
		outfit := model.NewOutfitBuilder().
			WithName("Free Spirit").
			WithItems([]model.Item{top, bottom, shoes}). // top/bottom/shoes have no occasions set
			AddOccasion("work").
			Build()

		_, warnings, err := svc.Create(&outfit)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if len(warnings) != 0 {
			t.Errorf("expected no warnings when items have no occasions, got: %v", warnings)
		}
	})

	// ── multiple errors collected at once ─────────────────────────────────────

	t.Run("multiple validation errors are collected", func(t *testing.T) {
		outfit := model.Outfit{
			Name:      "x",
			Items:     []model.Item{},
			Mood:      strings.Repeat("a", 51),
			Fragrance: strings.Repeat("a", 101),
		}

		_, _, err := svc.Create(&outfit)
		if err == nil {
			t.Fatal("expected validation errors, got nil")
		}
		for _, fragment := range []string{
			"name must be between 2 and 50 characters",
			"outfit must have at least 3 items",
			"mood must be at most 50 characters",
			"fragrance must be at most 100 characters",
		} {
			if !strings.Contains(err.Error(), fragment) {
				t.Errorf("error missing %q\nfull error: %s", fragment, err.Error())
			}
		}
	})
}

// assertOutfitValidationError fails the test if Create does not return an error containing wantFragment.
func assertOutfitValidationError(t *testing.T, svc *service.OutfitService, outfit model.Outfit, wantFragment string) {
	t.Helper()
	_, _, err := svc.Create(&outfit)
	if err == nil {
		t.Fatalf("expected validation error containing %q, got nil", wantFragment)
	}
	if !strings.Contains(err.Error(), wantFragment) {
		t.Errorf("error = %q\nwant to contain %q", err.Error(), wantFragment)
	}
}

// assertWarning fails the test if none of the warnings contain wantFragment.
func assertWarning(t *testing.T, warnings []string, wantFragment string) {
	t.Helper()
	for _, w := range warnings {
		if strings.Contains(w, wantFragment) {
			return
		}
	}
	t.Errorf("expected a warning containing %q\ngot warnings: %v", wantFragment, warnings)
}
