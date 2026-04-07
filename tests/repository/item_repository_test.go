package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/repository"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongoRepo(t *testing.T) (*repository.ItemRepository, func()) {
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

	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		_ = container.Terminate(ctx)
	}

	return repo, cleanup
}

func TestItemRepository_SaveAndRetrieve(t *testing.T) {
	repo, cleanup := setupMongoRepo(t)
	defer cleanup()

	item := model.NewItemBuilder("White T-Shirt", model.Top).
		WithColor("White").
		WithBrand("Zara").
		WithCondition("new").
		Build()

	saved, err := repo.Save(&item)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if saved.ID.IsZero() {
		t.Fatal("expected ID to be set after Save, got zero value")
	}
	if saved.Name != "White T-Shirt" {
		t.Errorf("Name = %q, want %q", saved.Name, "White T-Shirt")
	}
	if saved.Category != model.Top {
		t.Errorf("Category = %q, want %q", saved.Category, model.Top)
	}
	if saved.Color != "White" {
		t.Errorf("Color = %q, want %q", saved.Color, "White")
	}
	if saved.Brand != "Zara" {
		t.Errorf("Brand = %q, want %q", saved.Brand, "Zara")
	}
	if saved.Condition != "new" {
		t.Errorf("Condition = %q, want %q", saved.Condition, "new")
	}
}

func TestItemRepository_FindAll(t *testing.T) {
	repo, cleanup := setupMongoRepo(t)
	defer cleanup()

	seeds := []model.Item{
		model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").WithBrand("Zara").Build(),
		model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").WithBrand("Levi's").Build(),
		model.NewItemBuilder("Red Sneakers", model.Shoes).WithColor("Red").WithBrand("Nike").Build(),
	}

	for i := range seeds {
		if _, err := repo.Save(&seeds[i]); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	found, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if len(found) != len(seeds) {
		t.Fatalf("FindAll() returned %d items, want %d", len(found), len(seeds))
	}

	nameSet := make(map[string]bool, len(found))
	for _, it := range found {
		nameSet[it.Name] = true
	}
	for _, seed := range seeds {
		if !nameSet[seed.Name] {
			t.Errorf("FindAll() missing item %q", seed.Name)
		}
	}
}

func TestItemRepository_FindByID(t *testing.T) {
	repo, cleanup := setupMongoRepo(t)
	defer cleanup()

	item := model.NewItemBuilder("Black Hoodie", model.Outerwear).
		WithColor("Black").
		WithBrand("Nike").
		WithCondition("used").
		Build()

	saved, err := repo.Save(&item)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	t.Run("found", func(t *testing.T) {
		found, err := repo.FindByID(saved.ID.Hex())
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found.ID != saved.ID {
			t.Errorf("ID = %v, want %v", found.ID, saved.ID)
		}
		if found.Name != saved.Name {
			t.Errorf("Name = %q, want %q", found.Name, saved.Name)
		}
		if found.Category != saved.Category {
			t.Errorf("Category = %q, want %q", found.Category, saved.Category)
		}
		if found.Color != saved.Color {
			t.Errorf("Color = %q, want %q", found.Color, saved.Color)
		}
		if found.Brand != saved.Brand {
			t.Errorf("Brand = %q, want %q", found.Brand, saved.Brand)
		}
		if found.Condition != saved.Condition {
			t.Errorf("Condition = %q, want %q", found.Condition, saved.Condition)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.FindByID("000000000000000000000000")
		if err == nil {
			t.Fatal("FindByID() expected error for unknown ID, got nil")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		_, err := repo.FindByID("not-a-valid-id")
		if err == nil {
			t.Fatal("FindByID() expected error for invalid ID, got nil")
		}
	})
}

func TestItemRepository_Update(t *testing.T) {
	repo, cleanup := setupMongoRepo(t)
	defer cleanup()

	item := model.NewItemBuilder("White T-Shirt", model.Top).
		WithColor("White").
		WithBrand("Zara").
		WithCondition("new").
		Build()

	saved, err := repo.Save(&item)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	t.Run("updates fields and persists", func(t *testing.T) {
		updated := model.NewItemBuilder("Black T-Shirt", model.Top).
			WithColor("Black").
			WithBrand("H&M").
			WithCondition("used").
			Build()

		if err := repo.Update(saved.ID.Hex(), &updated); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		found, err := repo.FindByID(saved.ID.Hex())
		if err != nil {
			t.Fatalf("FindByID() after Update error = %v", err)
		}
		if found.Name != updated.Name {
			t.Errorf("Name = %q, want %q", found.Name, updated.Name)
		}
		if found.Color != updated.Color {
			t.Errorf("Color = %q, want %q", found.Color, updated.Color)
		}
		if found.Brand != updated.Brand {
			t.Errorf("Brand = %q, want %q", found.Brand, updated.Brand)
		}
		if found.Condition != updated.Condition {
			t.Errorf("Condition = %q, want %q", found.Condition, updated.Condition)
		}
		if found.ID != saved.ID {
			t.Errorf("ID changed after update: got %v, want %v", found.ID, saved.ID)
		}
	})

	t.Run("unknown id returns no error", func(t *testing.T) {
		patch := model.NewItemBuilder("Ghost Item", model.Top).Build()
		if err := repo.Update("000000000000000000000000", &patch); err != nil {
			t.Errorf("Update() unexpected error for unknown ID: %v", err)
		}
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		patch := model.NewItemBuilder("Ghost Item", model.Top).Build()
		if err := repo.Update("not-a-valid-id", &patch); err == nil {
			t.Fatal("Update() expected error for invalid ID, got nil")
		}
	})
}
