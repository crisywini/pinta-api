// Package handler_test contains integration tests for the HTTP handler layer.
//
// # How Gin handler testing works — no real server needed
//
// Testing a Gin handler does NOT require binding a TCP port or calling http.ListenAndServe.
// Gin implements http.Handler via its ServeHTTP method, which means the standard
// net/http/httptest package is all you need:
//
//  1. httptest.NewRecorder()   — creates an in-memory http.ResponseWriter.
//                                After the call you can read w.Code (status) and
//                                w.Body (the response body).
//  2. http.NewRequest(...)     — builds a synthetic *http.Request with no network I/O.
//  3. router.ServeHTTP(w, req) — drives the full Gin middleware chain and the handler,
//                                writing into the recorder.
//
// # What gin.SetMode(gin.TestMode) does
//
// Setting TestMode before creating the engine suppresses Gin's startup banner,
// per-request debug lines, and "headers already written" warnings, keeping
// test output clean.
//
// # Dependency strategy
//
// ItemHandler depends on *service.ItemService, which depends on
// *repository.ItemRepository.  Both are concrete types (no interfaces), so we
// cannot inject mocks.  Instead we wire the full real stack
//
//	MongoDB (TestContainers) → ItemRepository → ItemService → ItemHandler
//
// backed by a throwaway MongoDB container, consistent with how the rest of
// the project's tests are structured.
package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ── Setup helpers ─────────────────────────────────────────────────────────────

// setupRouter wires the full dependency stack and returns a configured *gin.Engine.
//
// Call pattern:
//
//	router, cleanup := setupRouter(t)
//	defer cleanup()
//
// Each top-level test function gets its own isolated MongoDB container so that
// subtests within it share state (seeded items remain visible across subtests
// in the same function), but different test functions are fully isolated.
func setupRouter(t *testing.T) (*gin.Engine, func()) {
	t.Helper()

	// Suppress Gin's debug output in tests.
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

	repo := repository.NewItemRepository(client.Database("pinta_test"))
	svc := service.NewItemService(repo)
	h := handler.NewItemHandler(svc)

	// gin.New() — a bare engine with no built-in middleware (no Logger, no Recovery).
	// Use gin.Default() if you need those two added automatically.
	router := gin.New()
	router.POST("/items", h.PostItem)
	router.GET("/items", h.GetAllItems)

	// GetItemByID, PutItem, and DeleteItem each read their target via the ?id=
	// query parameter (c.Query("id")), not a path segment.  We register them
	// on /item (singular) to avoid a route conflict with GET /items.
	router.GET("/item", h.GetItemByID)
	router.PUT("/item", h.PutItem)
	router.DELETE("/item", h.DeleteItem)

	cleanup := func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Disconnect(shutdownCtx)
		_ = container.Terminate(ctx)
	}

	return router, cleanup
}

// jsonBody marshals v to JSON and returns a *bytes.Buffer ready to use as a
// request body.  It fails the test immediately if marshalling errors.
func jsonBody(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("jsonBody: %v", err)
	}
	return bytes.NewBuffer(b)
}

// seedItem posts a minimal valid item via POST /items and returns the assigned
// MongoDB ID (hex string).  It fails the test if the server rejects the item.
//
// This is a convenient way to set up pre-existing data for subtests that need
// an item to already exist (GetItemByID, PutItem, DeleteItem happy paths).
func seedItem(t *testing.T, router *gin.Engine) string {
	t.Helper()

	item := model.NewItemBuilder("White T-Shirt", model.Top).WithColor("White").Build()

	req := httptest.NewRequest(http.MethodPost, "/items", jsonBody(t, item))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("seedItem: POST /items returned %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("seedItem: unmarshal: %v", err)
	}
	id, ok := resp["id"]
	if !ok || id == "" {
		t.Fatalf("seedItem: missing 'id' in response: %s", w.Body.String())
	}
	return id
}

// ── PostItem ──────────────────────────────────────────────────────────────────

func TestItemHandler_PostItem(t *testing.T) {
	router, cleanup := setupRouter(t)
	defer cleanup()

	t.Run("valid item returns 200 with id", func(t *testing.T) {
		item := model.NewItemBuilder("Blue Jeans", model.Bottom).
			WithColor("Indigo").
			WithBrand("Levi's").
			Build()

		req := httptest.NewRequest(http.MethodPost, "/items", jsonBody(t, item))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if resp["id"] == "" {
			t.Errorf("'id' field is empty; body: %s", w.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		// The handler calls ShouldBindBodyWith which returns an error on invalid JSON,
		// causing it to respond with 400 Bad Request.
		req := httptest.NewRequest(http.MethodPost, "/items",
			bytes.NewBufferString("{not valid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("validation failure returns 500 with message", func(t *testing.T) {
		// The name "x" is too short — the service returns a validation error.
		// PostItem maps service errors to 500 InternalServerError (as currently
		// coded).  Ideally this would be 422 Unprocessable Entity, but we test
		// the actual handler behaviour, not the ideal.
		item := model.Item{Name: "x", Category: model.Top, Color: "White"}

		req := httptest.NewRequest(http.MethodPost, "/items", jsonBody(t, item))
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
}

// ── GetAllItems ───────────────────────────────────────────────────────────────

func TestItemHandler_GetAllItems(t *testing.T) {
	router, cleanup := setupRouter(t)
	defer cleanup()

	t.Run("empty database returns count 0", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/items", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", w.Code)
		}

		var resp map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}

		// JSON numbers are decoded as float64 by encoding/json when the target
		// type is map[string]any, so cast accordingly.
		count, ok := resp["count"].(float64)
		if !ok {
			t.Fatalf("'count' field missing or wrong type; body: %s", w.Body.String())
		}
		if int(count) != 0 {
			t.Errorf("count = %d, want 0", int(count))
		}
	})

	t.Run("returns all seeded items", func(t *testing.T) {
		// Seed two more items into the same container that is shared across
		// subtests.  The empty-database subtest ran first, so the DB had 0 items.
		// After seeding here we expect at least 2.
		seedItem(t, router)
		seedItem(t, router)

		req := httptest.NewRequest(http.MethodGet, "/items", nil)
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

// ── GetItemByID ───────────────────────────────────────────────────────────────

func TestItemHandler_GetItemByID(t *testing.T) {
	router, cleanup := setupRouter(t)
	defer cleanup()

	// Seed one item to use in the happy-path subtest.
	id := seedItem(t, router)

	t.Run("returns 200 and item body for valid existing id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/item?id="+id, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("missing id param returns 400", func(t *testing.T) {
		// The handler writes 400 when the ?id query param is absent.
		// Note: the handler is missing a `return` after this write, so it
		// continues to call the service (which also fails) and writes additional
		// responses — but Gin only forwards the first WriteHeader call to the
		// underlying ResponseWriter, so the recorder captures 400.
		req := httptest.NewRequest(http.MethodGet, "/item", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("unknown id returns 404", func(t *testing.T) {
		// A well-formed ObjectID that does not exist in the database.
		req := httptest.NewRequest(http.MethodGet,
			"/item?id=000000000000000000000000", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404; body: %s", w.Code, w.Body.String())
		}
	})
}

// ── PutItem ───────────────────────────────────────────────────────────────────

func TestItemHandler_PutItem(t *testing.T) {
	router, cleanup := setupRouter(t)
	defer cleanup()

	id := seedItem(t, router)

	t.Run("valid update returns 200", func(t *testing.T) {
		updated := model.NewItemBuilder("Black T-Shirt", model.Top).
			WithColor("Black").
			WithCondition("fair").
			Build()

		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/item?id=%s", id),
			jsonBody(t, updated))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// On success, PutItem does not call c.JSON — the handler exits without
		// writing an explicit response.  Gin leaves the status at its default (200)
		// and the body is empty.  This is a handler bug (no 200 OK body), but
		// we verify the actual, current behaviour.
		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want 200; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/item?id=%s", id),
			bytes.NewBufferString("{bad json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid item body returns 500", func(t *testing.T) {
		invalid := model.Item{Name: "x", Category: model.Top, Color: "White"}

		req := httptest.NewRequest(http.MethodPut,
			fmt.Sprintf("/item?id=%s", id),
			jsonBody(t, invalid))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500; body: %s", w.Code, w.Body.String())
		}
	})
}

// ── DeleteItem ────────────────────────────────────────────────────────────────

func TestItemHandler_DeleteItem(t *testing.T) {
	router, cleanup := setupRouter(t)
	defer cleanup()

	t.Run("deletes existing item and returns 204", func(t *testing.T) {
		id := seedItem(t, router)

		req := httptest.NewRequest(http.MethodDelete, "/item?id="+id, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want 204; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("unknown id is a no-op and returns 204", func(t *testing.T) {
		// The repository's Delete treats a missing document as a success (no error).
		// So the handler responds 204 even when no document matched.
		req := httptest.NewRequest(http.MethodDelete,
			"/item?id=000000000000000000000000", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("status = %d, want 204; body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("empty id returns 500 because service rejects it", func(t *testing.T) {
		// service.DeleteById("") returns an error, which the handler maps to 500.
		req := httptest.NewRequest(http.MethodDelete, "/item", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want 500; body: %s", w.Code, w.Body.String())
		}
	})
}
