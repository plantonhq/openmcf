package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FGABootstrapResult holds the IDs produced by a successful bootstrap.
type FGABootstrapResult struct {
	StoreID              string
	AuthorizationModelID string
}

// EnsureFGABootstrap creates the OpenFGA store and writes the authorization
// model if they do not already exist. The function is idempotent: calling it
// when the store and model already exist is a safe no-op.
//
// The caller is responsible for persisting the returned IDs (typically in a
// ConfigMap) so that downstream services can discover them.
func EnsureFGABootstrap(ctx context.Context, client *http.Client, fgaURL, storeName string, modelJSON []byte) (*FGABootstrapResult, error) {
	storeID, err := ensureStore(ctx, client, fgaURL, storeName)
	if err != nil {
		return nil, fmt.Errorf("ensuring FGA store: %w", err)
	}

	modelID, err := ensureAuthorizationModel(ctx, client, fgaURL, storeID, modelJSON)
	if err != nil {
		return nil, fmt.Errorf("ensuring FGA authorization model: %w", err)
	}

	return &FGABootstrapResult{
		StoreID:              storeID,
		AuthorizationModelID: modelID,
	}, nil
}

// BootstrapConfigMapName returns the ConfigMap name used to persist FGA
// bootstrap state (store ID, model ID) for a given PlantonPlatform CR.
func BootstrapConfigMapName(crName string) string {
	return fmt.Sprintf("%s-fga-bootstrap", crName)
}

// ensureStore lists existing stores and returns the ID of one matching
// storeName, or creates a new store if none exists.
func ensureStore(ctx context.Context, client *http.Client, fgaURL, storeName string) (string, error) {
	id, err := findStore(ctx, client, fgaURL, storeName)
	if err != nil {
		return "", err
	}
	if id != "" {
		return id, nil
	}
	return createStore(ctx, client, fgaURL, storeName)
}

func findStore(ctx context.Context, client *http.Client, fgaURL, storeName string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fgaURL+"/stores", nil)
	if err != nil {
		return "", fmt.Errorf("building list stores request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("listing stores: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("list stores returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Stores []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"stores"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding list stores response: %w", err)
	}

	for _, s := range result.Stores {
		if s.Name == storeName {
			return s.ID, nil
		}
	}
	return "", nil
}

func createStore(ctx context.Context, client *http.Client, fgaURL, storeName string) (string, error) {
	body, err := json.Marshal(map[string]string{"name": storeName})
	if err != nil {
		return "", fmt.Errorf("marshaling create store body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fgaURL+"/stores", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("building create store request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("creating store: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create store returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding create store response: %w", err)
	}

	return result.ID, nil
}

// ensureAuthorizationModel writes the authorization model to the store. It
// always writes a new model version; OpenFGA models are immutable and writing
// a duplicate is harmless (it just creates a new version with the same content).
func ensureAuthorizationModel(ctx context.Context, client *http.Client, fgaURL, storeID string, modelJSON []byte) (string, error) {
	url := fmt.Sprintf("%s/stores/%s/authorization-models", fgaURL, storeID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(modelJSON))
	if err != nil {
		return "", fmt.Errorf("building write model request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("writing authorization model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("write model returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		AuthorizationModelID string `json:"authorization_model_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding write model response: %w", err)
	}

	return result.AuthorizationModelID, nil
}
