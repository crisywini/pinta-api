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

func seedItem(t *testing.T, repo *repository.ItemRepository, item model.Item) model.Item {
	t.Helper()
	saved, err := repo.Save(&item)
	if err != nil {
		t.Fatalf("seedItem: failed to save item %q: %v", item.Name, err)
	}
	return *saved
}

func minimalOutfit(name string, top, bottom, shoes model.Item) model.Outfit {
	return model.NewOutfitBuilder().
		WithName(name).
		WithItems([]model.Item{top, bottom, shoes}).
		Build()
}

func TestOutfitService_Create(t *testing.T) {
	svc, itemRepo, cleanup := setupOutfitService(t)
	defer cleanup()

	top := seedItem(t, itemRepo,
		model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build())
	bottom := seedItem(t, itemRepo,
		model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build())
	shoes := seedItem(t, itemRepo,
		model.NewItemBuilder("White Sneakers", model.Shoes).WithColor("White").Build())

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

	t.Run("items/zero id rejected", func(t *testing.T) {
		zeroTop := model.Item{Name: "Unsaved Top", Category: model.Top, Color: "White"}
		outfit := minimalOutfit("Zero ID", zeroTop, bottom, shoes)
		assertOutfitValidationError(t, svc, outfit, "must be saved to the closet before being added")
	})

	t.Run("items/duplicate item rejected", func(t *testing.T) {
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
			WithItems([]model.Item{top, bottom, shoes}).
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
			AddOccasion("work").
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
			WithItems([]model.Item{top, bottom, shoes}).
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

func assertWarning(t *testing.T, warnings []string, wantFragment string) {
	t.Helper()
	for _, w := range warnings {
		if strings.Contains(w, wantFragment) {
			return
		}
	}
	t.Errorf("expected a warning containing %q\ngot warnings: %v", wantFragment, warnings)
}

func assertUpdateOutfitValidationError(t *testing.T, svc *service.OutfitService, id string, outfit model.Outfit, wantFragment string) {
	t.Helper()
	_, err := svc.Update(id, &outfit)
	if err == nil {
		t.Fatalf("expected validation error containing %q, got nil", wantFragment)
	}
	if !strings.Contains(err.Error(), wantFragment) {
		t.Errorf("error = %q\nwant to contain %q", err.Error(), wantFragment)
	}
}

func seedOutfit(t *testing.T, svc *service.OutfitService, outfit model.Outfit) *model.Outfit {
	t.Helper()
	saved, _, err := svc.Create(&outfit)
	if err != nil {
		t.Fatalf("seedOutfit: failed to save outfit %q: %v", outfit.Name, err)
	}
	return saved
}

// ── GetById ───────────────────────────────────────────────────────────────────

func TestOutfitService_GetById(t *testing.T) {
	svc, itemRepo, cleanup := setupOutfitService(t)
	defer cleanup()

	top := seedItem(t, itemRepo,
		model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build())
	bottom := seedItem(t, itemRepo,
		model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build())
	shoes := seedItem(t, itemRepo,
		model.NewItemBuilder("White Sneakers", model.Shoes).WithColor("White").Build())

	saved := seedOutfit(t, svc, minimalOutfit("City Walk", top, bottom, shoes))

	t.Run("returns correct outfit for valid id", func(t *testing.T) {
		found, err := svc.GetById(saved.ID.Hex())
		if err != nil {
			t.Fatalf("GetById() error = %v", err)
		}
		if found.ID != saved.ID {
			t.Errorf("ID = %v, want %v", found.ID, saved.ID)
		}
		if found.Name != saved.Name {
			t.Errorf("Name = %q, want %q", found.Name, saved.Name)
		}
		if len(found.Items) != len(saved.Items) {
			t.Errorf("Items count = %d, want %d", len(found.Items), len(saved.Items))
		}
	})

	t.Run("empty id returns error", func(t *testing.T) {
		_, err := svc.GetById("")
		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if !strings.Contains(err.Error(), "Missing outfit id") {
			t.Errorf("error = %q, want to contain 'Missing outfit id'", err.Error())
		}
	})

	t.Run("unknown id returns error", func(t *testing.T) {
		_, err := svc.GetById("000000000000000000000000")
		if err == nil {
			t.Fatal("expected error for unknown id, got nil")
		}
	})

	t.Run("invalid id format returns error", func(t *testing.T) {
		_, err := svc.GetById("not-a-valid-id")
		if err == nil {
			t.Fatal("expected error for invalid id format, got nil")
		}
	})
}

// ── GetAll ────────────────────────────────────────────────────────────────────

func TestOutfitService_GetAll(t *testing.T) {
	svc, itemRepo, cleanup := setupOutfitService(t)
	defer cleanup()

	t.Run("returns empty slice when no outfits exist", func(t *testing.T) {
		outfits, err := svc.GetAll()
		if err != nil {
			t.Fatalf("GetAll() unexpected error: %v", err)
		}
		if len(outfits) != 0 {
			t.Errorf("GetAll() returned %d outfits, want 0", len(outfits))
		}
	})

	t.Run("returns all saved outfits", func(t *testing.T) {
		top := seedItem(t, itemRepo,
			model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build())
		bottom := seedItem(t, itemRepo,
			model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build())
		shoes := seedItem(t, itemRepo,
			model.NewItemBuilder("White Sneakers", model.Shoes).WithColor("White").Build())

		seedOutfit(t, svc, minimalOutfit("Outfit One", top, bottom, shoes))
		seedOutfit(t, svc, minimalOutfit("Outfit Two", top, bottom, shoes))
		seedOutfit(t, svc, minimalOutfit("Outfit Three", top, bottom, shoes))

		outfits, err := svc.GetAll()
		if err != nil {
			t.Fatalf("GetAll() unexpected error: %v", err)
		}
		if len(outfits) != 3 {
			t.Errorf("GetAll() returned %d outfits, want 3", len(outfits))
		}
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestOutfitService_Update(t *testing.T) {
	svc, itemRepo, cleanup := setupOutfitService(t)
	defer cleanup()

	top := seedItem(t, itemRepo,
		model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build())
	bottom := seedItem(t, itemRepo,
		model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build())
	shoes := seedItem(t, itemRepo,
		model.NewItemBuilder("White Sneakers", model.Shoes).WithColor("White").Build())

	saved := seedOutfit(t, svc, minimalOutfit("Original Name", top, bottom, shoes))
	// Use a stable valid ID for validation-only subtests (the repo is never reached).
	anyValidID := saved.ID.Hex()

	t.Run("updates fields and persists changes", func(t *testing.T) {
		updated := model.NewOutfitBuilder().
			WithName("Updated Name").
			WithItems([]model.Item{top, bottom, shoes}).
			WithMood("confident").
			AddSeason("spring").
			AddOccasion("casual").
			Build()

		warnings, err := svc.Update(anyValidID, &updated)
		if err != nil {
			t.Fatalf("Update() unexpected error: %v", err)
		}
		if len(warnings) != 0 {
			t.Errorf("expected no warnings, got: %v", warnings)
		}

		found, err := svc.GetById(anyValidID)
		if err != nil {
			t.Fatalf("GetById() after Update error = %v", err)
		}
		if found.Name != "Updated Name" {
			t.Errorf("Name = %q, want %q", found.Name, "Updated Name")
		}
		if found.Mood != "confident" {
			t.Errorf("Mood = %q, want %q", found.Mood, "confident")
		}
		if found.ID != saved.ID {
			t.Errorf("ID changed after update: got %v, want %v", found.ID, saved.ID)
		}
	})

	t.Run("empty id returns error", func(t *testing.T) {
		u := minimalOutfit("Valid Name", top, bottom, shoes)
		_, err := svc.Update("", &u)
		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if !strings.Contains(err.Error(), "Missing outfit id") {
			t.Errorf("error = %q, want to contain 'Missing outfit id'", err.Error())
		}
	})

	t.Run("unknown id returns no error", func(t *testing.T) {
		u := minimalOutfit("Valid Name", top, bottom, shoes)
		if _, err := svc.Update("000000000000000000000000", &u); err != nil {
			t.Errorf("Update() unexpected error for unknown id: %v", err)
		}
	})

	t.Run("invalid id format returns error", func(t *testing.T) {
		u := minimalOutfit("Valid Name", top, bottom, shoes)
		if _, err := svc.Update("not-a-valid-id", &u); err == nil {
			t.Fatal("expected error for invalid id format, got nil")
		}
	})

	// Validation rules — all fail before the repository is called.

	t.Run("name/too short", func(t *testing.T) {
		u := minimalOutfit("x", top, bottom, shoes)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "name must be between 2 and 50 characters")
	})

	t.Run("name/too long", func(t *testing.T) {
		u := minimalOutfit(strings.Repeat("a", 51), top, bottom, shoes)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "name must be between 2 and 50 characters")
	})

	t.Run("name/invalid characters", func(t *testing.T) {
		u := minimalOutfit("Look#1!", top, bottom, shoes)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "name may only contain letters")
	})

	t.Run("items/too few", func(t *testing.T) {
		u := model.NewOutfitBuilder().
			WithName("Too Few").
			WithItems([]model.Item{top, bottom}).
			Build()
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "outfit must have at least 3 items")
	})

	t.Run("items/missing required category", func(t *testing.T) {
		u := model.NewOutfitBuilder().
			WithName("No Shoes").
			WithItems([]model.Item{
				top,
				bottom,
				{Name: "Socks", Category: model.Hosiery, Color: "White"},
			}).
			Build()
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "outfit must have exactly 1 pair of shoes, got 0")
	})

	t.Run("occasion/invalid value", func(t *testing.T) {
		u := minimalOutfit("Bad Occasion", top, bottom, shoes)
		u.Occasion = []string{"casual", "gala"}
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, `invalid occasion "gala"`)
	})

	t.Run("season/invalid value", func(t *testing.T) {
		u := minimalOutfit("Bad Season", top, bottom, shoes)
		u.Season = []string{"summer", "monsoon"}
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, `invalid season "monsoon"`)
	})

	t.Run("mood/too long", func(t *testing.T) {
		u := minimalOutfit("Moody", top, bottom, shoes)
		u.Mood = strings.Repeat("a", 51)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "mood must be at most 50 characters")
	})

	t.Run("fragrance/too long", func(t *testing.T) {
		u := minimalOutfit("Fragrant", top, bottom, shoes)
		u.Fragrance = strings.Repeat("a", 101)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "fragrance must be at most 100 characters")
	})

	t.Run("notes/too long", func(t *testing.T) {
		u := minimalOutfit("Noted", top, bottom, shoes)
		u.Notes = strings.Repeat("a", 301)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "notes must be at most 300 characters")
	})

	t.Run("items/zero id rejected by closet check", func(t *testing.T) {
		zeroTop := model.Item{Name: "Unsaved Top", Category: model.Top, Color: "White"}
		u := minimalOutfit("Zero ID", zeroTop, bottom, shoes)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "must be saved to the closet before being added")
	})

	t.Run("items/duplicate item rejected by closet check", func(t *testing.T) {
		acc := seedItem(t, itemRepo,
			model.NewItemBuilder("Brown Belt", model.Accesories).WithColor("Brown").Build())
		u := model.NewOutfitBuilder().
			WithName("Duplicate Acc").
			WithItems([]model.Item{top, bottom, shoes, acc, acc}).
			Build()
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "appears more than once in the outfit")
	})

	t.Run("items/not in closet rejected by closet check", func(t *testing.T) {
		ghostTop := model.Item{Name: "Ghost Top", Category: model.Top, Color: "White"}
		ghostTop.ID = primitive.NewObjectID()
		u := minimalOutfit("Ghost Outfit", ghostTop, bottom, shoes)
		assertUpdateOutfitValidationError(t, svc, anyValidID, u, "was not found in the closet")
	})

	t.Run("warnings/season mismatch still updates", func(t *testing.T) {
		winterTop := seedItem(t, itemRepo,
			model.NewItemBuilder("Wool Turtleneck", model.Top).
				WithColor("Grey").
				WithSeason([]string{"winter"}).
				Build())

		u := model.NewOutfitBuilder().
			WithName("Summer Wool Update").
			WithItems([]model.Item{winterTop, bottom, shoes}).
			AddSeason("summer").
			Build()

		warnings, err := svc.Update(anyValidID, &u)
		if err != nil {
			t.Fatalf("Update() unexpected error: %v", err)
		}
		assertWarning(t, warnings, "season mismatch")
	})

	t.Run("multiple validation errors are collected", func(t *testing.T) {
		u := model.Outfit{
			Name:      "x",
			Items:     []model.Item{},
			Mood:      strings.Repeat("a", 51),
			Fragrance: strings.Repeat("a", 101),
		}
		_, err := svc.Update(anyValidID, &u)
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

// ── DeleteById ────────────────────────────────────────────────────────────────

func TestOutfitService_DeleteById(t *testing.T) {
	svc, itemRepo, cleanup := setupOutfitService(t)
	defer cleanup()

	top := seedItem(t, itemRepo,
		model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build())
	bottom := seedItem(t, itemRepo,
		model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build())
	shoes := seedItem(t, itemRepo,
		model.NewItemBuilder("White Sneakers", model.Shoes).WithColor("White").Build())

	t.Run("deletes existing outfit", func(t *testing.T) {
		saved := seedOutfit(t, svc, minimalOutfit("Temporary Outfit", top, bottom, shoes))

		if err := svc.DeleteById(saved.ID.Hex()); err != nil {
			t.Fatalf("DeleteById() error = %v", err)
		}

		_, err := svc.GetById(saved.ID.Hex())
		if err == nil {
			t.Fatal("expected error after deletion, got nil")
		}
	})

	t.Run("empty id returns error", func(t *testing.T) {
		err := svc.DeleteById("")
		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if !strings.Contains(err.Error(), "Missing outfit id") {
			t.Errorf("error = %q, want to contain 'Missing outfit id'", err.Error())
		}
	})

	t.Run("unknown id returns no error", func(t *testing.T) {
		if err := svc.DeleteById("000000000000000000000000"); err != nil {
			t.Errorf("DeleteById() unexpected error for unknown id: %v", err)
		}
	})

	t.Run("invalid id format returns error", func(t *testing.T) {
		if err := svc.DeleteById("not-a-valid-id"); err == nil {
			t.Fatal("expected error for invalid id format, got nil")
		}
	})
}
