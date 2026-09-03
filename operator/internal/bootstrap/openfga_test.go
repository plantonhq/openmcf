package bootstrap_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plantonhq/planton/operator/internal/bootstrap"
)

// The operator provides the store and nothing inside it: no request ever
// reaches an authorization-models path, because the control plane owns its
// model and writes it itself at boot.
func TestEnsureFGABootstrap_CreatesStoreOnly(t *testing.T) {
	storeCreated := false

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/stores":
			json.NewEncoder(w).Encode(map[string]any{"stores": []any{}})

		case r.Method == http.MethodPost && r.URL.Path == "/stores":
			storeCreated = true
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": "store-123", "name": "planton"})

		default:
			t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	result, err := bootstrap.EnsureFGABootstrap(context.Background(), srv.Client(), srv.URL, "planton")
	if err != nil {
		t.Fatalf("EnsureFGABootstrap failed: %v", err)
	}

	if !storeCreated {
		t.Error("Expected store to be created")
	}
	if result.StoreID != "store-123" {
		t.Errorf("StoreID = %q, want %q", result.StoreID, "store-123")
	}
}

func TestEnsureFGABootstrap_FindsExistingStore(t *testing.T) {
	storeCreateCalls := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/stores":
			json.NewEncoder(w).Encode(map[string]any{
				"stores": []any{
					map[string]any{"id": "existing-store", "name": "planton"},
				},
			})

		case r.Method == http.MethodPost && r.URL.Path == "/stores":
			storeCreateCalls++
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"id": "new-store", "name": "planton"})

		default:
			t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	result, err := bootstrap.EnsureFGABootstrap(context.Background(), srv.Client(), srv.URL, "planton")
	if err != nil {
		t.Fatalf("EnsureFGABootstrap failed: %v", err)
	}

	if storeCreateCalls > 0 {
		t.Errorf("Store should not have been created (already exists), got %d create calls", storeCreateCalls)
	}
	if result.StoreID != "existing-store" {
		t.Errorf("StoreID = %q, want %q", result.StoreID, "existing-store")
	}
}

func TestEnsureFGABootstrap_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer srv.Close()

	_, err := bootstrap.EnsureFGABootstrap(context.Background(), srv.Client(), srv.URL, "planton")
	if err == nil {
		t.Error("Expected error when server returns 500, got nil")
	}
}

func TestBootstrapConfigMapName(t *testing.T) {
	if got := bootstrap.BootstrapConfigMapName("myp"); got != "myp-fga-bootstrap" {
		t.Errorf("BootstrapConfigMapName = %q, want %q", got, "myp-fga-bootstrap")
	}
}
