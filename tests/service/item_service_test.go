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
