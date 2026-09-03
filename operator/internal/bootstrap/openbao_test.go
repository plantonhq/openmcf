package bootstrap

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckOpenBAOHealth_Initialized_Unsealed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/health" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"initialized": true,
			"sealed":      false,
		})
	}))
	defer server.Close()

	status, err := CheckOpenBAOHealth(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Initialized {
		t.Error("expected initialized=true")
	}
	if status.Sealed {
		t.Error("expected sealed=false")
	}
}

func TestCheckOpenBAOHealth_NotInitialized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]any{
			"initialized": false,
			"sealed":      true,
		})
	}))
	defer server.Close()

	status, err := CheckOpenBAOHealth(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Initialized {
		t.Error("expected initialized=false")
	}
	if !status.Sealed {
		t.Error("expected sealed=true")
	}
}

func TestInitializeOpenBAO_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/init" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		var body map[string]int
		json.NewDecoder(r.Body).Decode(&body)
		if body["secret_shares"] != 5 {
			t.Errorf("expected 5 shares, got %d", body["secret_shares"])
		}
		if body["secret_threshold"] != 3 {
			t.Errorf("expected 3 threshold, got %d", body["secret_threshold"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"keys":       []string{"key1", "key2", "key3", "key4", "key5"},
			"root_token": "root-token-abc",
		})
	}))
	defer server.Close()

	result, err := InitializeOpenBAO(context.Background(), server.Client(), server.URL, 5, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.UnsealKeys) != 5 {
		t.Errorf("expected 5 unseal keys, got %d", len(result.UnsealKeys))
	}
	if result.RootToken != "root-token-abc" {
		t.Errorf("expected root-token-abc, got %s", result.RootToken)
	}
}

func TestInitializeOpenBAO_AlreadyInitialized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"errors":["Vault is already initialized"]}`))
	}))
	defer server.Close()

	_, err := InitializeOpenBAO(context.Background(), server.Client(), server.URL, 5, 3)
	if err == nil {
		t.Fatal("expected error for already initialized")
	}
}

func TestUnsealOpenBAO_Success(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/unseal" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		callCount++

		sealed := callCount < 3
		json.NewEncoder(w).Encode(map[string]any{
			"sealed":   sealed,
			"t":        3,
			"n":        5,
			"progress": callCount,
		})
	}))
	defer server.Close()

	err := UnsealOpenBAO(context.Background(), server.Client(), server.URL, []string{"k1", "k2", "k3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 unseal calls, got %d", callCount)
	}
}

func TestUnsealOpenBAO_StillSealed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"sealed": true,
		})
	}))
	defer server.Close()

	err := UnsealOpenBAO(context.Background(), server.Client(), server.URL, []string{"k1", "k2"})
	if err == nil {
		t.Fatal("expected error when still sealed")
	}
}
