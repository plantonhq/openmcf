package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenBAOInitResult holds the secrets produced by a successful initialization.
type OpenBAOInitResult struct {
	UnsealKeys []string
	RootToken  string
}

// OpenBAOHealthStatus represents the state of an OpenBAO instance.
type OpenBAOHealthStatus struct {
	Initialized bool
	Sealed      bool
}

// CheckOpenBAOHealth queries /v1/sys/health to determine initialization and
// seal status. The OpenBAO health endpoint returns different HTTP status codes:
//   - 200: initialized, unsealed, active
//   - 429: unsealed, standby
//   - 472: data recovery mode
//   - 501: not initialized
//   - 503: sealed
//
// All status codes return a JSON body with "initialized" and "sealed" fields.
func CheckOpenBAOHealth(ctx context.Context, client *http.Client, apiAddr string) (*OpenBAOHealthStatus, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiAddr+"/v1/sys/health", nil)
	if err != nil {
		return nil, fmt.Errorf("building health request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("checking health: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Initialized bool `json:"initialized"`
		Sealed      bool `json:"sealed"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding health response: %w", err)
	}

	return &OpenBAOHealthStatus{
		Initialized: result.Initialized,
		Sealed:      result.Sealed,
	}, nil
}

// InitializeOpenBAO calls /v1/sys/init to initialize a fresh OpenBAO instance
// with the given number of secret shares and threshold. Returns the unseal keys
// and root token. This is a one-time operation; calling it on an already
// initialized instance returns an error.
func InitializeOpenBAO(ctx context.Context, client *http.Client, apiAddr string, secretShares, secretThreshold int) (*OpenBAOInitResult, error) {
	body, err := json.Marshal(map[string]int{
		"secret_shares":    secretShares,
		"secret_threshold": secretThreshold,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling init request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiAddr+"/v1/sys/init", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("building init request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("initializing OpenBAO: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("init returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Keys      []string `json:"keys"`
		RootToken string   `json:"root_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding init response: %w", err)
	}

	return &OpenBAOInitResult{
		UnsealKeys: result.Keys,
		RootToken:  result.RootToken,
	}, nil
}

// EnsureOpenBAOMounts idempotently enables the two secrets engines the
// platform expects: KV v2 at "secret/" and Transit at "transit/" (the mounts
// PlatformVaultClient hardcodes). A fresh server-mode OpenBAO starts with NO
// secrets engines -- dev mode auto-mounts a KV engine, the Helm chart runs
// server mode, and Transit is never auto-mounted in any mode -- so without
// this every platform vault call fails with "no handler for route".
// Enabling engines is infrastructure provisioning and therefore the
// operator's job, never the control plane's.
func EnsureOpenBAOMounts(ctx context.Context, client *http.Client, apiAddr, rootToken string) error {
	mounts, err := listMounts(ctx, client, apiAddr, rootToken)
	if err != nil {
		return fmt.Errorf("listing OpenBAO mounts: %w", err)
	}

	if _, ok := mounts["secret/"]; !ok {
		if err := enableMount(ctx, client, apiAddr, rootToken, "secret", map[string]any{
			"type":    "kv",
			"options": map[string]string{"version": "2"},
		}); err != nil {
			return fmt.Errorf("enabling KV v2 engine at secret/: %w", err)
		}
	}

	if _, ok := mounts["transit/"]; !ok {
		if err := enableMount(ctx, client, apiAddr, rootToken, "transit", map[string]any{
			"type": "transit",
		}); err != nil {
			return fmt.Errorf("enabling Transit engine at transit/: %w", err)
		}
	}

	return nil
}

func listMounts(ctx context.Context, client *http.Client, apiAddr, rootToken string) (map[string]json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiAddr+"/v1/sys/mounts", nil)
	if err != nil {
		return nil, fmt.Errorf("building mounts request: %w", err)
	}
	req.Header.Set("X-Vault-Token", rootToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("listing mounts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sys/mounts returned %d: %s", resp.StatusCode, string(respBody))
	}

	// Mount names arrive as top-level keys ("secret/", "transit/", "sys/", ...);
	// newer API versions nest them under "data" as well -- reading the top level
	// works for both because the legacy shape is preserved.
	var result map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding mounts response: %w", err)
	}
	return result, nil
}

func enableMount(ctx context.Context, client *http.Client, apiAddr, rootToken, path string, config map[string]any) error {
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshaling mount config: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiAddr+"/v1/sys/mounts/"+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building mount request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", rootToken)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("enabling mount %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sys/mounts/%s returned %d: %s", path, resp.StatusCode, string(respBody))
	}

	return nil
}

// UnsealOpenBAO calls /v1/sys/unseal with each key until the threshold is met
// and the instance is unsealed. Returns nil on success.
func UnsealOpenBAO(ctx context.Context, client *http.Client, apiAddr string, keys []string) error {
	for i, key := range keys {
		body, err := json.Marshal(map[string]string{"key": key})
		if err != nil {
			return fmt.Errorf("marshaling unseal request %d: %w", i, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiAddr+"/v1/sys/unseal", bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("building unseal request %d: %w", i, err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("unsealing (key %d): %w", i, err)
		}

		var result struct {
			Sealed bool `json:"sealed"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return fmt.Errorf("decoding unseal response %d: %w", i, err)
		}
		resp.Body.Close()

		if !result.Sealed {
			return nil
		}
	}

	return fmt.Errorf("still sealed after applying all %d keys", len(keys))
}
