package verify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// QdrantVerifier checks a Qdrant vector database to the point clients can
// rely on it: the StatefulSet's replicas ready, the main Service present,
// and a LIVE vector round-trip through the REST API — create a run-unique
// collection, upsert points, run a similarity search and assert the
// nearest neighbour (a vector database that cannot answer a similarity
// query is not a vector database). When an API key is declared, every
// request carries it — which also proves the auth wiring end to end.
//
// The behavioral-persistence scenario (recognized by name) additionally
// DELETES pod 0 after the upsert and re-runs the search once the pod
// returns — vectors surviving pod loss through the PVC is the proof.
type QdrantVerifier struct {
	Namespace string
	Name      string
	Replicas  int
	// ApiKeySecretName is the chart-owned `<name>-apikey` Secret (key
	// api-key). Empty = the cluster is unauthenticated; requests carry no
	// key.
	ApiKeySecretName string
	Persistence      bool
}

func (v *QdrantVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] qdrant cluster %q in namespace %q\n", v.Name, v.Namespace)

	if err := v.waitForReady(ctx, kubeconfig, 10*time.Minute); err != nil {
		return err
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "qdrant service not found")
	}
	return v.proveVectorRoundTrip(ctx, kubeconfig)
}

func (v *QdrantVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.Name, v.Namespace)
}

func (v *QdrantVerifier) waitForReady(ctx context.Context, kubeconfig string, budget time.Duration) error {
	want := fmt.Sprintf("%d", v.Replicas)
	deadline := time.Now().Add(budget)
	var lastReady string
	for time.Now().Before(deadline) {
		ready, _ := kubectlGetJSONPath(ctx, kubeconfig, "statefulset", v.Name, v.Namespace, "{.status.readyReplicas}")
		lastReady = ready
		if ready == want {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("qdrant statefulset never reached %s ready replicas (last %q)", want, lastReady)
}

// apiKey reads the read-write key from the chart-owned `<name>-apikey`
// Secret, or "" when the cluster is unauthenticated.
func (v *QdrantVerifier) apiKey(ctx context.Context, kubeconfig string) (string, error) {
	if v.ApiKeySecretName == "" {
		return "", nil
	}
	keyB64, err := kubectlGetJSONPath(ctx, kubeconfig, "secret", v.ApiKeySecretName, v.Namespace, "{.data.api-key}")
	if err != nil {
		return "", errors.Wrapf(err, "reading secret %q api-key", v.ApiKeySecretName)
	}
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(key)), nil
}

// proveVectorRoundTrip drives the REST API over a port-forward: create a
// run-unique collection, upsert two labelled vectors, search with a query
// vector near the first and assert it comes back as the nearest
// neighbour. The persistence variant kills pod 0 between upsert and
// search.
func (v *QdrantVerifier) proveVectorRoundTrip(ctx context.Context, kubeconfig string) error {
	const localPort = "16333"

	pfCtx, cancel := context.WithCancel(ctx)
	pf := exec.CommandContext(pfCtx, "kubectl", "--kubeconfig", kubeconfig,
		"port-forward", "svc/"+v.Name, localPort+":6333", "-n", v.Namespace)
	var pfOut strings.Builder
	pf.Stdout = &pfOut
	pf.Stderr = &pfOut
	if err := pf.Start(); err != nil {
		cancel()
		return errors.Wrap(err, "starting port-forward to the qdrant service")
	}
	// ONE deferred func, cancel FIRST — Wait blocks forever on a
	// port-forward that is never told to exit.
	defer func() {
		cancel()
		_ = pf.Wait()
	}()

	apiKey, err := v.apiKey(ctx, kubeconfig)
	if err != nil {
		return err
	}
	base := "http://127.0.0.1:" + localPort

	collection := fmt.Sprintf("e2e_proof_%d", time.Now().Unix())

	// Create the collection (4-dim cosine space keeps the proof cheap and
	// deterministic). Retried across the port-forward/boot warm-up.
	createBody := `{"vectors": {"size": 4, "distance": "Cosine"}}`
	if out, err := v.request(ctx, http.MethodPut, base+"/collections/"+collection, apiKey, createBody, 4*time.Minute); err != nil {
		return errors.Wrapf(err, "creating the proof collection: %s", out)
	}
	fmt.Printf("  [verify] VECTOR: collection %q created\n", collection)

	// Two well-separated unit-ish vectors; the query vector sits next to
	// point 1, so cosine similarity must rank it first.
	upsertBody := `{"points": [
		{"id": 1, "vector": [0.95, 0.05, 0.05, 0.05], "payload": {"label": "alpha"}},
		{"id": 2, "vector": [0.05, 0.95, 0.05, 0.05], "payload": {"label": "beta"}}
	]}`
	if out, err := v.request(ctx, http.MethodPut, base+"/collections/"+collection+"/points?wait=true", apiKey, upsertBody, 2*time.Minute); err != nil {
		return errors.Wrapf(err, "upserting the proof vectors: %s", out)
	}
	fmt.Printf("  [verify] VECTOR: 2 points upserted\n")

	if v.Persistence {
		pod := v.Name + "-0"
		fmt.Printf("  [verify] PERSISTENCE: deleting pod %q\n", pod)
		if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"delete", "pod", pod, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
			return errors.Wrapf(err, "deleting the qdrant pod: %s", string(out))
		}
		if err := v.waitForReady(ctx, kubeconfig, 10*time.Minute); err != nil {
			return errors.Wrap(err, "the pod never returned after deletion")
		}
	}

	searchBody := `{"vector": [1.0, 0.0, 0.0, 0.0], "limit": 1}`
	out, err := v.request(ctx, http.MethodPost, base+"/collections/"+collection+"/points/search", apiKey, searchBody, 4*time.Minute)
	if err != nil {
		return errors.Wrapf(err, "the similarity search never succeeded: %s", out)
	}
	var search struct {
		Result []struct {
			Id json.Number `json:"id"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(out), &search); err != nil {
		return errors.Wrapf(err, "parsing the search response: %s", firstLines(out, 3))
	}
	if len(search.Result) == 0 || search.Result[0].Id.String() != "1" {
		return errors.Errorf("the nearest neighbour was not point 1: %s", firstLines(out, 3))
	}
	if v.Persistence {
		fmt.Printf("  [verify] PERSISTENCE: similarity search answered correctly AFTER pod loss — vectors survived on the PVC\n")
	} else {
		fmt.Printf("  [verify] VECTOR: similarity search returned the expected nearest neighbour\n")
	}
	return nil
}

// request performs one JSON request with retries across the warm-up
// window, returning the response body. Non-2xx responses are errors.
func (v *QdrantVerifier) request(ctx context.Context, method, url, apiKey, body string, budget time.Duration) (string, error) {
	deadline := time.Now().Add(budget)
	var lastOut string
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader([]byte(body)))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		if apiKey != "" {
			req.Header.Set("api-key", apiKey)
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(resp.Body)
			resp.Body.Close()
			lastOut = buf.String()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return lastOut, nil
			}
			lastErr = errors.Errorf("HTTP %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(10 * time.Second)
	}
	return lastOut, lastErr
}
