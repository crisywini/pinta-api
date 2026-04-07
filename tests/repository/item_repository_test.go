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
