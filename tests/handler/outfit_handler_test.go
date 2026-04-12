package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/crisywini/pinta-api/internal/handler"
	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/repository"
	"github.com/crisywini/pinta-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// setupOutfitRouter wires the full dependency stack for outfit handler tests.
// Both /items and /outfits routes are registered on the same router because:
//   - OutfitService.verifyItemsInCloset queries ItemRepository directly, so items
//     must exist in the same DB that the outfit service uses.
//   - seedItemByCategory seeds items through POST /items on this router.
func setupOutfitRouter(t *testing.T) (*gin.Engine, func()) {
	t.Helper()
	gin.SetMode(gin.TestMode)

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
	itemSvc := service.NewItemService(itemRepo)
	itemH := handler.NewItemHandler(itemSvc)

	outfitRepo := repository.NewOutfitRepository(db)
	outfitSvc := service.NewOutfitService(outfitRepo, itemRepo)
	outfitH := handler.NewOutfitHandler(outfitSvc)

	router := gin.New()
	router.POST("/items", itemH.PostItem)
	router.GET("/items/:id", itemH.GetItemByID)
	router.PUT("/items/:id", itemH.PutItem)
	router.DELETE("/items/:id", itemH.DeleteItem)
	router.GET("/items", itemH.GetAllItems)

	router.POST("/outfits", outfitH.PostOutfit)
	router.GET("/outfits/:id", outfitH.GetOutfitById)
	router.GET("/outfits", outfitH.GetAllOutfits)
	router.PUT("/outfits/:id", outfitH.PutOutfit)
	router.DELETE("/outfits/:id", outfitH.DeleteOutfitById)

	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		_ = container.Terminate(ctx)
	}

	return router, cleanup
}

// mustObjectID parses a hex string into a primitive.ObjectID and fails the test
// immediately if the string is not valid.
func mustObjectID(t *testing.T, hex string) primitive.ObjectID {
	t.Helper()
	oid, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		t.Fatalf("mustObjectID: invalid hex %q: %v", hex, err)
	}
	return oid
}

// seedItemByCategory posts an item with the given category via POST /items and
// returns the assigned MongoDB ID (hex string).
func seedItemByCategory(t *testing.T, router *gin.Engine, name string, category model.Category) string {
	t.Helper()

	item := model.NewItemBuilder(name, category).WithColor("White").Build()

	req := httptest.NewRequest(http.MethodPost, "/items", jsonBody(t, item))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("seedItemByCategory: POST /items returned %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("seedItemByCategory: unmarshal: %v", err)
	}

	id, ok := resp["id"]
	if !ok || id == "" {
		t.Fatalf("seedItemByCategory: missing 'id' in response: %s", w.Body.String())
	}
	return id
}

// seedOutfit seeds the three mandatory items (top, bottom, shoes) into the DB
// and then posts a minimal valid outfit referencing them.  It returns the
// outfit ID (hex string).
func seedOutfit(t *testing.T, router *gin.Engine) string {
	t.Helper()

	topID := seedItemByCategory(t, router, "White T-Shirt", model.Top)
	bottomID := seedItemByCategory(t, router, "Blue Jeans", model.Bottom)
	shoesID := seedItemByCategory(t, router, "White Sneakers", model.Shoes)

	outfit := model.NewOutfitBuilder().
		WithName("Casual Look").
		AddItem(model.Item{ID: mustObjectID(t, topID), Name: "White T-Shirt", Category: model.Top}).
		AddItem(model.Item{ID: mustObjectID(t, bottomID), Name: "Blue Jeans", Category: model.Bottom}).
		AddItem(model.Item{ID: mustObjectID(t, shoesID), Name: "White Sneakers", Category: model.Shoes}).
		Build()

	req := httptest.NewRequest(http.MethodPost, "/outfits", jsonBody(t, outfit))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("seedOutfit: POST /outfits returned %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("seedOutfit: unmarshal: %v", err)
	}

	id, ok := resp["outfit"].(string)
	if !ok || id == "" {
		t.Fatalf("seedOutfit: missing 'outfit' in response: %s", w.Body.String())
	}
	return id
}

// ── PostOutfit ────────────────────────────────────────────────────────────────

func TestOutfitHandler_PostOutfit(t *testing.T) {
	router, cleanup := setupOutfitRouter(t)
	defer cleanup()

	t.Run("valid outfit returns 201 with id and warnings", func(t *testing.T) {
		topID := seedItemByCategory(t, router, "Black Tee", model.Top)
		bottomID := seedItemByCategory(t, router, "Chinos", model.Bottom)
		shoesID := seedItemByCategory(t, router, "Loafers", model.Shoes)

		outfit := model.NewOutfitBuilder().
			WithName("Smart Casual").
			AddItem(model.Item{ID: mustObjectID(t, topID), Name: "Black Tee", Category: model.Top}).
			AddItem(model.Item{ID: mustObjectID(t, bottomID), Name: "Chinos", Category: model.Bottom}).
			AddItem(model.Item{ID: mustObjectID(t, shoesID), Name: "Loafers", Category: model.Shoes}).
			Build()

		req := httptest.NewRequest(http.MethodPost, "/outfits", jsonBody(t, outfit))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want 201; body: %s", w.Code, w.Body.String())
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if id, _ := resp["outfit"].(string); id == "" {
			t.Errorf("'outfit' field is empty; body: %s", w.Body.String())
		}
		if _, ok := resp["warnings"]; !ok {
			t.Errorf("'warnings' field missing; body: %s", w.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/outfits",
			bytes.NewBufferString("{not valid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid outfit name returns 500", func(t *testing.T) {
		// Name "x" is 1 character — service validates min length of 2 and maps
		// the error to 500 InternalServerError.
		outfit := model.Outfit{Name: "x"}

		req := httptest.NewRequest(http.MethodPost, "/outfits", jsonBody(t, outfit))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500; body: %s", w.Code, w.Body.String())
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal error response: %v", err)
		}
		if resp["message"] == "" {
			t.Errorf("expected 'message' in error body; got: %s", w.Body.String())
		}
	})

	t.Run("items not saved to closet returns 500", func(t *testing.T) {
		// Items have valid (non-zero) ObjectIDs that do not exist in the DB.
		// verifyItemsInCloset returns an error → 500.
		outfit := model.NewOutfitBuilder().
			WithName("Ghost Outfit").
			AddItem(model.Item{ID: primitive.NewObjectID(), Name: "Ghost Top", Category: model.Top}).
			AddItem(model.Item{ID: primitive.NewObjectID(), Name: "Ghost Bottom", Category: model.Bottom}).
			AddItem(model.Item{ID: primitive.NewObjectID(), Name: "Ghost Shoes", Category: model.Shoes}).
			Build()

		req := httptest.NewRequest(http.MethodPost, "/outfits", jsonBody(t, outfit))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500; body: %s", w.Code, w.Body.String())
		}
	})
}

// ── GetOutfitById ─────────────────────────────────────────────────────────────

func TestOutfitHandler_GetOutfitById(t *testing.T) {
	router, cleanup := setupOutfitRouter(t)
	defer cleanup()

	id := seedOutfit(t, router)

	t.Run("existing outfit returns 200 with body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/outfits/"+id, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
		}

		var outfit model.Outfit
		if err := json.Unmarshal(w.Body.Bytes(), &outfit); err != nil {
			t.Fatalf("unmarshal outfit: %v", err)
		}
		if outfit.Name == "" {
			t.Errorf("expected non-empty outfit name; body: %s", w.Body.String())
		}
	})

	t.Run("unknown valid id returns 400", func(t *testing.T) {
		// Repository returns mongo.ErrNoDocuments; handler maps all service
		// errors to 400 BadRequest (unlike ItemHandler which uses 404).
		req := httptest.NewRequest(http.MethodGet, "/outfits/000000000000000000000000", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid id format returns 400", func(t *testing.T) {
		// Repository cannot parse the hex string → error propagated to handler → 400.
		req := httptest.NewRequest(http.MethodGet, "/outfits/not-a-valid-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
		}
	})
}

// ── GetAllOutfits ─────────────────────────────────────────────────────────────

func TestOutfitHandler_GetAllOutfits(t *testing.T) {
	router, cleanup := setupOutfitRouter(t)
	defer cleanup()

	t.Run("empty database returns count 0", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/outfits", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		// JSON numbers decode as float64 when the target type is map[string]any.
		count, ok := resp["count"].(float64)
		if !ok {
			t.Fatalf("'count' field missing or wrong type; body: %s", w.Body.String())
		}
		if int(count) != 0 {
			t.Errorf("count = %d, want 0", int(count))
		}
	})

	t.Run("returns all seeded outfits", func(t *testing.T) {
		seedOutfit(t, router)
		seedOutfit(t, router)

		req := httptest.NewRequest(http.MethodGet, "/outfits", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		count, ok := resp["count"].(float64)
		if !ok {
			t.Fatalf("'count' field missing; body: %s", w.Body.String())
		}
		if int(count) < 2 {
			t.Errorf("count = %d, want >= 2", int(count))
		}
	})
}

// ── DeleteOutfitById ──────────────────────────────────────────────────────────

func TestOutfitHandler_DeleteOutfitById(t *testing.T) {
	router, cleanup := setupOutfitRouter(t)
	defer cleanup()

	t.Run("deletes existing outfit and returns 204", func(t *testing.T) {
		id := seedOutfit(t, router)

		req := httptest.NewRequest(http.MethodDelete, "/outfits/"+id, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want 204; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("unknown valid id is a no-op and returns 204", func(t *testing.T) {
		// DeleteOne on a missing document returns no error, so the handler
		// responds with 204 just as it would for a real deletion.
		req := httptest.NewRequest(http.MethodDelete, "/outfits/000000000000000000000000", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want 204; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid id format returns 500", func(t *testing.T) {
		// Repository fails to parse the hex string → error → 500.
		req := httptest.NewRequest(http.MethodDelete, "/outfits/not-a-valid-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500; body: %s", w.Code, w.Body.String())
		}
	})
}

// ── PutOutfit ─────────────────────────────────────────────────────────────────

func TestOutfitHandler_PutOutfit(t *testing.T) {
	router, cleanup := setupOutfitRouter(t)
	defer cleanup()

	id := seedOutfit(t, router)

	t.Run("valid update returns 201 with id and warnings", func(t *testing.T) {
		topID := seedItemByCategory(t, router, "Red Shirt", model.Top)
		bottomID := seedItemByCategory(t, router, "Black Jeans", model.Bottom)
		shoesID := seedItemByCategory(t, router, "Boots", model.Shoes)

		updated := model.NewOutfitBuilder().
			WithName("Updated Look").
			AddItem(model.Item{ID: mustObjectID(t, topID), Name: "Red Shirt", Category: model.Top}).
			AddItem(model.Item{ID: mustObjectID(t, bottomID), Name: "Black Jeans", Category: model.Bottom}).
			AddItem(model.Item{ID: mustObjectID(t, shoesID), Name: "Boots", Category: model.Shoes}).
			Build()

		req := httptest.NewRequest(http.MethodPut, "/outfits/"+id, jsonBody(t, updated))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want 201; body: %s", w.Code, w.Body.String())
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if outfitID, _ := resp["outfit"].(string); outfitID != id {
			t.Errorf("'outfit' = %q, want %q; body: %s", outfitID, id, w.Body.String())
		}
		if _, ok := resp["warnings"]; !ok {
			t.Errorf("'warnings' field missing; body: %s", w.Body.String())
		}
	})

	t.Run("malformed JSON is silently ignored and returns 200", func(t *testing.T) {
		// PutOutfit has no else-branch after ShouldBindBodyWithJSON fails, so
		// the handler exits without writing any response.  Gin defaults to 200.
		req := httptest.NewRequest(http.MethodPut, "/outfits/"+id,
			bytes.NewBufferString("{bad json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want 200 (silent no-op on bind failure); body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid outfit body returns 500", func(t *testing.T) {
		// Name "x" fails service validation → serviceError != nil → 500.
		invalid := model.Outfit{Name: "x"}

		req := httptest.NewRequest(http.MethodPut, "/outfits/"+id, jsonBody(t, invalid))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500; body: %s", w.Code, w.Body.String())
		}
	})
}
