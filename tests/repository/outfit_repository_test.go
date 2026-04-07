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

func setupOutfitRepo(t *testing.T) (*repository.OutfitRepository, func()) {
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

	repo := repository.NewOutfitRepository(client.Database("pinta_test"))

	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		_ = container.Terminate(ctx)
	}

	return repo, cleanup
}

func TestOutfitRepository_Save(t *testing.T) {
	repo, cleanup := setupOutfitRepo(t)
	defer cleanup()

	shirt := model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").WithBrand("Zara").Build()
	jeans := model.NewItemBuilder("Blue Jeans", model.Bottom).WithColor("Blue").WithBrand("Levi's").Build()

	outfit := model.NewOutfitBuilder().
		WithName("Casual Friday").
		WithMood("relaxed").
		AddSeason("spring").
		AddOccasion("casual").
		WithItems([]model.Item{shirt, jeans}).
		Build()

	saved, err := repo.Save(&outfit)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if saved.ID.IsZero() {
		t.Fatal("expected ID to be set after Save, got zero value")
	}
	if saved.Name != "Casual Friday" {
		t.Errorf("Name = %q, want %q", saved.Name, "Casual Friday")
	}
	if saved.Mood != "relaxed" {
		t.Errorf("Mood = %q, want %q", saved.Mood, "relaxed")
	}
	if len(saved.Items) != 2 {
		t.Errorf("Items count = %d, want 2", len(saved.Items))
	}
}

func TestOutfitRepository_FindAll(t *testing.T) {
	repo, cleanup := setupOutfitRepo(t)
	defer cleanup()

	seeds := []model.Outfit{
		model.NewOutfitBuilder().WithName("Casual Friday").WithMood("relaxed").Build(),
		model.NewOutfitBuilder().WithName("Summer Vibes").AddSeason("summer").Build(),
		model.NewOutfitBuilder().WithName("Winter Formal").WithMood("elegant").AddSeason("winter").Build(),
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
		t.Fatalf("FindAll() returned %d outfits, want %d", len(found), len(seeds))
	}

	nameSet := make(map[string]bool, len(found))
	for _, o := range found {
		nameSet[o.Name] = true
	}
	for _, seed := range seeds {
		if !nameSet[seed.Name] {
			t.Errorf("FindAll() missing outfit %q", seed.Name)
		}
	}
}

func TestOutfitRepository_FindByID(t *testing.T) {
	repo, cleanup := setupOutfitRepo(t)
	defer cleanup()

	shirt := model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build()
	outfit := model.NewOutfitBuilder().
		WithName("Casual Friday").
		WithMood("relaxed").
		AddOccasion("casual").
		WithItems([]model.Item{shirt}).
		Build()

	saved, err := repo.Save(&outfit)
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
		if found.Mood != saved.Mood {
			t.Errorf("Mood = %q, want %q", found.Mood, saved.Mood)
		}
		if len(found.Items) != len(saved.Items) {
			t.Errorf("Items count = %d, want %d", len(found.Items), len(saved.Items))
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

func TestOutfitRepository_Update(t *testing.T) {
	repo, cleanup := setupOutfitRepo(t)
	defer cleanup()

	outfit := model.NewOutfitBuilder().
		WithName("Casual Friday").
		WithMood("relaxed").
		AddSeason("spring").
		Build()

	saved, err := repo.Save(&outfit)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	t.Run("updates fields and persists", func(t *testing.T) {
		sneakers := model.NewItemBuilder("Red Sneakers", model.Shoes).WithColor("Red").Build()
		updated := model.NewOutfitBuilder().
			WithName("Weekend Casual").
			WithMood("energetic").
			AddSeason("summer").
			AddOccasion("outdoor").
			WithItems([]model.Item{sneakers}).
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
		if found.Mood != updated.Mood {
			t.Errorf("Mood = %q, want %q", found.Mood, updated.Mood)
		}
		if len(found.Items) != 1 {
			t.Errorf("Items count = %d, want 1", len(found.Items))
		}
		if found.ID != saved.ID {
			t.Errorf("ID changed after update: got %v, want %v", found.ID, saved.ID)
		}
	})

	t.Run("unknown id returns no error", func(t *testing.T) {
		patch := model.NewOutfitBuilder().WithName("Ghost Outfit").Build()
		if err := repo.Update("000000000000000000000000", &patch); err != nil {
			t.Errorf("Update() unexpected error for unknown ID: %v", err)
		}
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		patch := model.NewOutfitBuilder().WithName("Ghost Outfit").Build()
		if err := repo.Update("not-a-valid-id", &patch); err == nil {
			t.Fatal("Update() expected error for invalid ID, got nil")
		}
	})
}

func TestOutfitRepository_Delete(t *testing.T) {
	repo, cleanup := setupOutfitRepo(t)
	defer cleanup()

	outfit := model.NewOutfitBuilder().
		WithName("Casual Friday").
		WithMood("relaxed").
		Build()

	saved, err := repo.Save(&outfit)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	t.Run("deletes existing outfit", func(t *testing.T) {
		if err := repo.Delete(saved.ID.Hex()); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		_, err := repo.FindByID(saved.ID.Hex())
		if err == nil {
			t.Fatal("FindByID() expected error after Delete, got nil")
		}
	})

	t.Run("unknown id returns no error", func(t *testing.T) {
		if err := repo.Delete("000000000000000000000000"); err != nil {
			t.Errorf("Delete() unexpected error for unknown ID: %v", err)
		}
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		if err := repo.Delete("not-a-valid-id"); err == nil {
			t.Fatal("Delete() expected error for invalid ID, got nil")
		}
	})
}
