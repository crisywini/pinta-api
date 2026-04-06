package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/repository"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestItemRepository_SaveAndRetrieve(t *testing.T) {
	ctx := context.Background()

	// Start a real MongoDB container
	mongoContainer, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatalf("failed to start mongodb container: %v", err)
	}
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	uri, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("failed to connect to mongodb: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
	}()

	db := client.Database("pinta_test")
	repo := repository.NewItemRepository(db)

	// Build the item to save
	item := model.NewItemBuilder("White T-Shirt", model.Top).
		WithColor("White").
		WithBrand("Zara").
		WithCondition("new").
		Build()

	// Save
	saved, err := repo.Save(&item)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify the returned item has an ID assigned
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

	// Verify the document was actually persisted in MongoDB
	var found model.Item
	err = db.Collection("items").
		FindOne(ctx, bson.M{"name": "White T-Shirt"}).
		Decode(&found)
	if err != nil {
		t.Fatalf("FindOne() after Save error = %v", err)
	}
	if found.Name != saved.Name {
		t.Errorf("persisted Name = %q, want %q", found.Name, saved.Name)
	}
	if found.Category != saved.Category {
		t.Errorf("persisted Category = %q, want %q", found.Category, saved.Category)
	}
}
