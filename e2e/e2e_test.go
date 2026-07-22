//go:build e2e

// Package e2e contains end-to-end tests that deploy real infrastructure
// using Planton IaC modules and verify the results.
//
// These tests require:
//   - kind CLI installed
//   - pulumi CLI installed
//   - kubectl CLI installed
//   - Docker running
//
// Run with: go test -tags=e2e -timeout=30m -v ./e2e/...
package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/google/uuid"
	kubernetese2e "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/aa_e2e"
	"github.com/plantonhq/planton/e2e/framework/runner"
)

var (
	// testHarness is the shared kind cluster harness for all Kubernetes tests.
	testHarness *kubernetese2e.Harness

	// ciliumHarness is the lazily-created harness for scenarios annotated with
	// the cilium-cni cluster profile (a dedicated persistent kind cluster with
	// the default CNI disabled). Nil until the first such scenario runs; runs
	// that touch no profiled scenario never create the second cluster.
	ciliumHarness *kubernetese2e.Harness

	// ciliumHarnessMu guards lazy ciliumHarness construction. Scenarios run
	// serially within a process, but the guard keeps the invariant explicit.
	ciliumHarnessMu sync.Mutex

	// repoRoot is the absolute path to the planton repository root.
	repoRoot string

	// runID is a unique identifier for this test run.
	runID string

	// pulumiBackendURL is the file-based backend for Pulumi stacks.
	pulumiBackendURL string
)

// harnessForScenario routes a scenario to the cluster its manifest asks for
// via the planton.dev/e2e-cluster-profile annotation. The default (no
// annotation) is the shared cluster. In the external-cluster lane profiles
// are ignored — the batched real cluster is assumed suitable for its batch.
// Unknown profile values fail loudly rather than silently running on the
// wrong cluster.
func harnessForScenario(manifestPath string) (*kubernetese2e.Harness, error) {
	profile, err := runner.ManifestAnnotation(manifestPath, kubernetese2e.ClusterProfileAnnotation)
	if err != nil {
		return nil, fmt.Errorf("reading cluster profile from %s: %w", manifestPath, err)
	}
	if profile == "" || testHarness.External() {
		return testHarness, nil
	}
	switch profile {
	case kubernetese2e.ClusterProfileCiliumCni:
		ciliumHarnessMu.Lock()
		defer ciliumHarnessMu.Unlock()
		if ciliumHarness == nil {
			h := kubernetese2e.NewCiliumCniHarness(testHarness.ClusterName())
			if err := h.Setup(context.Background()); err != nil {
				return nil, fmt.Errorf("setting up %s cluster profile: %w", profile, err)
			}
			ciliumHarness = h
		}
		return ciliumHarness, nil
	default:
		return nil, fmt.Errorf("scenario %s names unknown cluster profile %q", manifestPath, profile)
	}
}

func TestMain(m *testing.M) {
	// Resolve repo root (this file lives at e2e/e2e_test.go)
	var err error
	repoRoot, err = filepath.Abs(filepath.Join(".."))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve repo root: %v\n", err)
		os.Exit(1)
	}

	// Generate unique run ID
	runID = uuid.New().String()[:8]

	// Set up local Pulumi backend
	backendDir, err := os.MkdirTemp("", "planton-e2e-pulumi-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp backend dir: %v\n", err)
		os.Exit(1)
	}
	pulumiBackendURL = "file://" + backendDir
	defer os.RemoveAll(backendDir)

	// Log into local backend
	if err := runner.PulumiLogin(pulumiBackendURL); err != nil {
		fmt.Fprintf(os.Stderr, "failed to login to pulumi backend: %v\n", err)
		os.Exit(1)
	}

	// The cluster name is STABLE across runs (not run-ID-suffixed) so the harness's
	// persistent-kind lane can reuse the running cluster from the previous run.
	clusterName := os.Getenv("PLANTON_E2E_KIND_CLUSTER_NAME")
	if clusterName == "" {
		clusterName = "planton-e2e"
	}
	testHarness = kubernetese2e.NewHarness(clusterName)

	ctx := context.Background()
	if err := testHarness.Setup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create kind cluster: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown kind cluster(s) — the profile cluster first (created last).
	if ciliumHarness != nil {
		if err := ciliumHarness.Teardown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to tear down cilium profile cluster: %v\n", err)
		}
	}
	if err := testHarness.Teardown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to delete kind cluster: %v\n", err)
	}

	os.Exit(code)
}
