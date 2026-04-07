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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupItemService(t *testing.T) (*service.ItemService, func()) {
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

	repo := repository.NewItemRepository(client.Database("pinta_test"))
	svc := service.NewItemService(repo)

	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		_ = container.Terminate(ctx)
	}

	return svc, cleanup
}

func baseItem() model.Item {
	return model.NewItemBuilder("White T-Shirt", model.Top).
		WithColor("White").
		Build()
}

func TestItemService_Create(t *testing.T) {
	svc, cleanup := setupItemService(t)
	defer cleanup()

	t.Run("saves valid item and assigns id", func(t *testing.T) {
		item := model.NewItemBuilder("White T-Shirt", model.Top).
			WithColor("White").
			WithBrand("Zara").
			WithCondition("new").
			Build()

		saved, err := svc.Create(&item)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if saved.ID.IsZero() {
			t.Fatal("expected ID to be set after Create, got zero value")
		}
	})

	t.Run("applies default condition when not provided", func(t *testing.T) {
		item := baseItem()

		saved, err := svc.Create(&item)
		if err != nil {
			t.Fatalf("Create() unexpected error: %v", err)
		}
		if saved.Condition != "good" {
			t.Errorf("Condition = %q, want %q", saved.Condition, "good")
		}
	})

	t.Run("accepts name with hyphens and apostrophes", func(t *testing.T) {
		item := model.NewItemBuilder("O'Brien's T-Shirt", model.Top).
			WithColor("Navy").
			Build()

		if _, err := svc.Create(&item); err != nil {
			t.Errorf("Create() unexpected error for valid name: %v", err)
		}
	})

	t.Run("name/too short", func(t *testing.T) {
		item := baseItem()
		item.Name = "x"
		assertValidationError(t, svc, item, "name must be between 2 and 50 characters")
	})

	t.Run("name/too long", func(t *testing.T) {
		item := baseItem()
		item.Name = strings.Repeat("a", 51)
		assertValidationError(t, svc, item, "name must be between 2 and 50 characters")
	})

	t.Run("name/invalid characters", func(t *testing.T) {
		item := baseItem()
		item.Name = "Cool@Shirt!"
		assertValidationError(t, svc, item, "name may only contain letters")
	})

	t.Run("category/invalid", func(t *testing.T) {
		item := baseItem()
		item.Category = model.Category("hat")
		assertValidationError(t, svc, item, "category must be one of")
	})

	t.Run("color/too short", func(t *testing.T) {
		item := baseItem()
		item.Color = "R"
		assertValidationError(t, svc, item, "color must be between 2 and 30 characters")
	})

	t.Run("color/too long", func(t *testing.T) {
		item := baseItem()
		item.Color = strings.Repeat("a", 31)
		assertValidationError(t, svc, item, "color must be between 2 and 30 characters")
	})

	t.Run("color/empty", func(t *testing.T) {
		item := baseItem()
		item.Color = ""
		assertValidationError(t, svc, item, "color must be between 2 and 30 characters")
	})

	t.Run("brand/too short when provided", func(t *testing.T) {
		item := baseItem()
		item.Brand = "Z"
		assertValidationError(t, svc, item, "brand must be between 2 and 50 characters when provided")
	})

	t.Run("brand/too long when provided", func(t *testing.T) {
		item := baseItem()
		item.Brand = strings.Repeat("a", 51)
		assertValidationError(t, svc, item, "brand must be between 2 and 50 characters when provided")
	})

	t.Run("material/too short when provided", func(t *testing.T) {
		item := baseItem()
		item.Material = "C"
		assertValidationError(t, svc, item, "material must be between 2 and 50 characters when provided")
	})

	t.Run("material/too long when provided", func(t *testing.T) {
		item := baseItem()
		item.Material = strings.Repeat("a", 51)
		assertValidationError(t, svc, item, "material must be between 2 and 50 characters when provided")
	})

	t.Run("season/invalid value", func(t *testing.T) {
		item := baseItem()
		item.Season = []string{"spring", "monsoon"}
		assertValidationError(t, svc, item, `invalid season "monsoon"`)
	})

	t.Run("season/all valid values accepted", func(t *testing.T) {
		item := baseItem()
		item.Season = []string{"spring", "summer", "fall", "winter"}
		if _, err := svc.Create(&item); err != nil {
			t.Errorf("Create() unexpected error for valid seasons: %v", err)
		}
	})

	t.Run("occasion/invalid value", func(t *testing.T) {
		item := baseItem()
		item.Occasion = []string{"casual", "gala"}
		assertValidationError(t, svc, item, `invalid occasion "gala"`)
	})

	t.Run("occasion/all valid values accepted", func(t *testing.T) {
		item := baseItem()
		item.Occasion = []string{"work", "casual", "brunch", "formal", "party", "sport", "date"}
		if _, err := svc.Create(&item); err != nil {
			t.Errorf("Create() unexpected error for valid occasions: %v", err)
		}
	})

	t.Run("photo/not a URL", func(t *testing.T) {
		item := baseItem()
		item.Photo = "not-a-url"
		assertValidationError(t, svc, item, "photo must be a valid http or https URL")
	})

	t.Run("photo/non http scheme rejected", func(t *testing.T) {
		item := baseItem()
		item.Photo = "ftp://cdn.example.com/photo.jpg"
		assertValidationError(t, svc, item, "photo must be a valid http or https URL")
	})

	t.Run("photo/valid https URL accepted", func(t *testing.T) {
		item := baseItem()
		item.Photo = "https://cdn.example.com/photo.jpg"
		if _, err := svc.Create(&item); err != nil {
			t.Errorf("Create() unexpected error for valid photo URL: %v", err)
		}
	})

	t.Run("condition/invalid value", func(t *testing.T) {
		item := baseItem()
		item.Condition = "worn-out"
		assertValidationError(t, svc, item, "condition must be one of: new, good, fair, retired")
	})

	t.Run("multiple validation errors are collected", func(t *testing.T) {
		item := model.Item{
			Name:      "x",
			Category:  model.Category("hat"),
			Color:     "R",
			Condition: "worn-out",
		}

		_, err := svc.Create(&item)
		if err == nil {
			t.Fatal("expected validation errors, got nil")
		}
		for _, fragment := range []string{
			"name must be between 2 and 50 characters",
			"category must be one of",
			"color must be between 2 and 30 characters",
			"condition must be one of",
		} {
			if !strings.Contains(err.Error(), fragment) {
				t.Errorf("error missing %q\nfull error: %s", fragment, err.Error())
			}
		}
	})
}

func assertValidationError(t *testing.T, svc *service.ItemService, item model.Item, wantFragment string) {
	t.Helper()
	_, err := svc.Create(&item)
	if err == nil {
		t.Fatalf("expected validation error containing %q, got nil", wantFragment)
	}
	if !strings.Contains(err.Error(), wantFragment) {
		t.Errorf("error = %q\nwant to contain %q", err.Error(), wantFragment)
	}
}

func assertUpdateValidationError(t *testing.T, svc *service.ItemService, id string, item model.Item, wantFragment string) {
	t.Helper()
	err := svc.Update(id, &item)
	if err == nil {
		t.Fatalf("expected validation error containing %q, got nil", wantFragment)
	}
	if !strings.Contains(err.Error(), wantFragment) {
		t.Errorf("error = %q\nwant to contain %q", err.Error(), wantFragment)
	}
}

// ── GetById ──────────────────────────────────────────────────────────────────

func TestItemService_GetById(t *testing.T) {
	svc, cleanup := setupItemService(t)
	defer cleanup()

	item := model.NewItemBuilder("Blue Jeans", model.Bottom).
		WithColor("Blue").
		WithBrand("Levi's").
		Build()
	saved, err := svc.Create(&item)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("returns correct item for valid id", func(t *testing.T) {
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
		if found.Color != saved.Color {
			t.Errorf("Color = %q, want %q", found.Color, saved.Color)
		}
	})

	t.Run("empty id returns error", func(t *testing.T) {
		_, err := svc.GetById("")
		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if !strings.Contains(err.Error(), "Missing item id") {
			t.Errorf("error = %q, want to contain 'Missing item id'", err.Error())
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

// ── GetAll ───────────────────────────────────────────────────────────────────

func TestItemService_GetAll(t *testing.T) {
	svc, cleanup := setupItemService(t)
	defer cleanup()

	t.Run("returns empty slice when no items exist", func(t *testing.T) {
		items, err := svc.GetAll()
		if err != nil {
			t.Fatalf("GetAll() unexpected error: %v", err)
		}
		if len(items) != 0 {
			t.Errorf("GetAll() returned %d items, want 0", len(items))
		}
	})

	t.Run("returns all saved items", func(t *testing.T) {
		seeds := []model.Item{
			model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build(),
			model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").Build(),
			model.NewItemBuilder("Red Sneakers", model.Shoes).WithColor("Red").Build(),
		}
		for i := range seeds {
			if _, err := svc.Create(&seeds[i]); err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		items, err := svc.GetAll()
		if err != nil {
			t.Fatalf("GetAll() unexpected error: %v", err)
		}
		if len(items) != len(seeds) {
			t.Errorf("GetAll() returned %d items, want %d", len(items), len(seeds))
		}
	})
}

// ── Update ───────────────────────────────────────────────────────────────────

func TestItemService_Update(t *testing.T) {
	svc, cleanup := setupItemService(t)
	defer cleanup()

	item := model.NewItemBuilder("White T-Shirt", model.Top).
		WithColor("White").
		WithBrand("Zara").
		Build()
	saved, err := svc.Create(&item)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Use a stable valid ObjectID for validation-only subtests (repo never reached).
	anyValidID := saved.ID.Hex()

	t.Run("updates fields and persists changes", func(t *testing.T) {
		updated := model.NewItemBuilder("Black T-Shirt", model.Top).
			WithColor("Black").
			WithBrand("H&M").
			WithCondition("fair").
			Build()

		if err := svc.Update(saved.ID.Hex(), &updated); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		found, err := svc.GetById(saved.ID.Hex())
		if err != nil {
			t.Fatalf("GetById() after Update error = %v", err)
		}
		if found.Name != "Black T-Shirt" {
			t.Errorf("Name = %q, want %q", found.Name, "Black T-Shirt")
		}
		if found.Color != "Black" {
			t.Errorf("Color = %q, want %q", found.Color, "Black")
		}
		if found.Brand != "H&M" {
			t.Errorf("Brand = %q, want %q", found.Brand, "H&M")
		}
		if found.Condition != "fair" {
			t.Errorf("Condition = %q, want %q", found.Condition, "fair")
		}
		if found.ID != saved.ID {
			t.Errorf("ID changed after update: got %v, want %v", found.ID, saved.ID)
		}
	})

	t.Run("empty id returns error", func(t *testing.T) {
		u := baseItem()
		err := svc.Update("", &u)
		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if !strings.Contains(err.Error(), "Missing item id") {
			t.Errorf("error = %q, want to contain 'Missing item id'", err.Error())
		}
	})

	t.Run("unknown id returns no error", func(t *testing.T) {
		u := baseItem()
		if err := svc.Update("000000000000000000000000", &u); err != nil {
			t.Errorf("Update() unexpected error for unknown id: %v", err)
		}
	})

	t.Run("invalid id format returns error", func(t *testing.T) {
		u := baseItem()
		if err := svc.Update("not-a-valid-id", &u); err == nil {
			t.Fatal("expected error for invalid id format, got nil")
		}
	})

	// Validation rules — all fail before the repository is called.

	t.Run("name/too short", func(t *testing.T) {
		u := baseItem()
		u.Name = "x"
		assertUpdateValidationError(t, svc, anyValidID, u, "name must be between 2 and 50 characters")
	})

	t.Run("name/too long", func(t *testing.T) {
		u := baseItem()
		u.Name = strings.Repeat("a", 51)
		assertUpdateValidationError(t, svc, anyValidID, u, "name must be between 2 and 50 characters")
	})

	t.Run("name/invalid characters", func(t *testing.T) {
		u := baseItem()
		u.Name = "Cool@Shirt!"
		assertUpdateValidationError(t, svc, anyValidID, u, "name may only contain letters")
	})

	t.Run("category/invalid", func(t *testing.T) {
		u := baseItem()
		u.Category = model.Category("hat")
		assertUpdateValidationError(t, svc, anyValidID, u, "category must be one of")
	})

	t.Run("color/too short", func(t *testing.T) {
		u := baseItem()
		u.Color = "R"
		assertUpdateValidationError(t, svc, anyValidID, u, "color must be between 2 and 30 characters")
	})

	t.Run("color/too long", func(t *testing.T) {
		u := baseItem()
		u.Color = strings.Repeat("a", 31)
		assertUpdateValidationError(t, svc, anyValidID, u, "color must be between 2 and 30 characters")
	})

	t.Run("color/empty", func(t *testing.T) {
		u := baseItem()
		u.Color = ""
		assertUpdateValidationError(t, svc, anyValidID, u, "color must be between 2 and 30 characters")
	})

	t.Run("brand/too short when provided", func(t *testing.T) {
		u := baseItem()
		u.Brand = "Z"
		assertUpdateValidationError(t, svc, anyValidID, u, "brand must be between 2 and 50 characters when provided")
	})

	t.Run("brand/too long when provided", func(t *testing.T) {
		u := baseItem()
		u.Brand = strings.Repeat("a", 51)
		assertUpdateValidationError(t, svc, anyValidID, u, "brand must be between 2 and 50 characters when provided")
	})

	t.Run("material/too short when provided", func(t *testing.T) {
		u := baseItem()
		u.Material = "C"
		assertUpdateValidationError(t, svc, anyValidID, u, "material must be between 2 and 50 characters when provided")
	})

	t.Run("material/too long when provided", func(t *testing.T) {
		u := baseItem()
		u.Material = strings.Repeat("a", 51)
		assertUpdateValidationError(t, svc, anyValidID, u, "material must be between 2 and 50 characters when provided")
	})

	t.Run("season/invalid value", func(t *testing.T) {
		u := baseItem()
		u.Season = []string{"spring", "monsoon"}
		assertUpdateValidationError(t, svc, anyValidID, u, `invalid season "monsoon"`)
	})

	t.Run("occasion/invalid value", func(t *testing.T) {
		u := baseItem()
		u.Occasion = []string{"casual", "gala"}
		assertUpdateValidationError(t, svc, anyValidID, u, `invalid occasion "gala"`)
	})

	t.Run("photo/invalid URL", func(t *testing.T) {
		u := baseItem()
		u.Photo = "not-a-url"
		assertUpdateValidationError(t, svc, anyValidID, u, "photo must be a valid http or https URL")
	})

	t.Run("photo/non http scheme rejected", func(t *testing.T) {
		u := baseItem()
		u.Photo = "ftp://cdn.example.com/photo.jpg"
		assertUpdateValidationError(t, svc, anyValidID, u, "photo must be a valid http or https URL")
	})

	t.Run("condition/invalid value", func(t *testing.T) {
		u := baseItem()
		u.Condition = "worn-out"
		assertUpdateValidationError(t, svc, anyValidID, u, "condition must be one of: new, good, fair, retired")
	})

	t.Run("multiple validation errors are collected", func(t *testing.T) {
		u := model.Item{
			Name:      "x",
			Category:  model.Category("hat"),
			Color:     "R",
			Condition: "worn-out",
		}
		err := svc.Update(anyValidID, &u)
		if err == nil {
			t.Fatal("expected validation errors, got nil")
		}
		for _, fragment := range []string{
			"name must be between 2 and 50 characters",
			"category must be one of",
			"color must be between 2 and 30 characters",
			"condition must be one of",
		} {
			if !strings.Contains(err.Error(), fragment) {
				t.Errorf("error missing %q\nfull error: %s", fragment, err.Error())
			}
		}
	})
}

// ── DeleteById ───────────────────────────────────────────────────────────────

func TestItemService_DeleteById(t *testing.T) {
	svc, cleanup := setupItemService(t)
	defer cleanup()

	t.Run("deletes existing item", func(t *testing.T) {
		item := model.NewItemBuilder("Red Sneakers", model.Shoes).WithColor("Red").Build()
		saved, err := svc.Create(&item)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if err := svc.DeleteById(saved.ID.Hex()); err != nil {
			t.Fatalf("DeleteById() error = %v", err)
		}

		_, err = svc.GetById(saved.ID.Hex())
		if err == nil {
			t.Fatal("expected error after deletion, got nil")
		}
	})

	t.Run("empty id returns error", func(t *testing.T) {
		err := svc.DeleteById("")
		if err == nil {
			t.Fatal("expected error for empty id, got nil")
		}
		if !strings.Contains(err.Error(), "Missing item id") {
			t.Errorf("error = %q, want to contain 'Missing item id'", err.Error())
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
