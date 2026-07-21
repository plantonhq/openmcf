// Package aa_e2e implements the E2E provider harness for Kubernetes.
//
// Two cluster lanes are supported, selected by environment:
//
//   - kind (Kubernetes IN Docker), the default: the harness owns the cluster
//     lifecycle. The cluster PERSISTS across test runs -- Setup reuses a running
//     cluster with the configured name and only creates one when none exists,
//     because cluster create/destroy dominates run time and local runs execute
//     back-to-back. Set PLANTON_E2E_DESTROY_CLUSTER=1 (ephemeral CI runners) to
//     delete it at teardown.
//
//   - external cluster: PLANTON_E2E_KUBECONFIG points at a kubeconfig for a
//     cluster provisioned and owned OUTSIDE the run (batched real EKS/GKE/AKS
//     lanes -- one cluster per wave batch, reused across component lanes). The
//     harness never creates or deletes anything cluster-level in this lane.
//
// Component-level verify semantics (deployed / destroyed) are identical in both
// lanes: verification keys off the kubeconfig path, not the cluster's origin.
package aa_e2e

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/aa_e2e/verify"
	"github.com/plantonhq/planton/e2e/framework/provider"
)

// Environment contract for lane selection.
const (
	// ExternalKubeconfigEnvVar selects the external-cluster lane: a path to the
	// kubeconfig of a cluster the harness does not own.
	ExternalKubeconfigEnvVar = "PLANTON_E2E_KUBECONFIG"

	// DestroyClusterEnvVar opts INTO deleting the kind cluster at teardown ("1").
	// The default keeps it running for the next run (persistent local lane).
	DestroyClusterEnvVar = "PLANTON_E2E_DESTROY_CLUSTER"
)

// Harness manages the test cluster for Kubernetes E2E tests.
type Harness struct {
	clusterName    string
	kubeconfigPath string
	tempDir        string

	// external marks the external-cluster lane: the cluster is provisioned and
	// owned outside the run; Setup/Teardown must not touch its lifecycle.
	external bool
}

// NewHarness creates a Kubernetes test harness with the given kind cluster name.
// The name is ignored when the external-cluster lane is selected via environment.
func NewHarness(clusterName string) *Harness {
	return &Harness{
		clusterName: clusterName,
	}
}

// Setup makes a cluster available and stores its kubeconfig path: adopting the
// external cluster when the lane is selected, reusing a running kind cluster with
// the configured name, or creating one when none exists.
func (h *Harness) Setup(ctx context.Context) error {
	if externalKubeconfig := os.Getenv(ExternalKubeconfigEnvVar); externalKubeconfig != "" {
		if _, err := os.Stat(externalKubeconfig); err != nil {
			return errors.Wrapf(err, "%s points at an unreadable kubeconfig", ExternalKubeconfigEnvVar)
		}
		h.external = true
		h.kubeconfigPath = externalKubeconfig
		fmt.Printf("  [cluster] Using external cluster via %s (%s)\n", ExternalKubeconfigEnvVar, externalKubeconfig)

		os.Setenv("KUBECONFIG", h.kubeconfigPath)
		return nil
	}

	tmpDir, err := os.MkdirTemp("", "planton-e2e-*")
	if err != nil {
		return errors.Wrap(err, "failed to create temp directory for kubeconfig")
	}
	h.tempDir = tmpDir
	h.kubeconfigPath = filepath.Join(tmpDir, "kubeconfig")

	running, err := kindClusterExists(ctx, h.clusterName)
	if err != nil {
		return err
	}

	if running {
		// Reuse: export a fresh kubeconfig for the already-running cluster.
		fmt.Printf("  [kind] Reusing running cluster %q\n", h.clusterName)
		cmd := exec.CommandContext(ctx, "kind", "export", "kubeconfig",
			"--name", h.clusterName, "--kubeconfig", h.kubeconfigPath)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "kind export kubeconfig failed: %s", stderr.String())
		}
	} else {
		args := []string{
			"create", "cluster",
			"--name", h.clusterName,
			"--kubeconfig", h.kubeconfigPath,
			"--wait", "120s",
		}

		cmd := exec.CommandContext(ctx, "kind", args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		cmd.Stdout = os.Stdout

		fmt.Printf("  [kind] Creating cluster %q...\n", h.clusterName)
		start := time.Now()

		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "kind create cluster failed: %s", stderr.String())
		}

		fmt.Printf("  [kind] Cluster %q ready in %s\n", h.clusterName, time.Since(start).Round(time.Second))
	}

	// Set KUBECONFIG globally so Pulumi's kubernetes provider picks it up
	os.Setenv("KUBECONFIG", h.kubeconfigPath)

	return nil
}

// Teardown releases the cluster per lane: external clusters are never touched, and
// the kind cluster persists for the next run unless deletion is opted into.
func (h *Harness) Teardown(ctx context.Context) error {
	defer func() {
		if h.tempDir != "" {
			os.RemoveAll(h.tempDir)
		}
		os.Unsetenv("KUBECONFIG")
	}()

	if h.external {
		fmt.Println("  [cluster] External cluster left untouched (owned outside the run)")
		return nil
	}

	if os.Getenv(DestroyClusterEnvVar) != "1" {
		fmt.Printf("  [kind] Leaving cluster %q running for the next run (set %s=1 to delete)\n",
			h.clusterName, DestroyClusterEnvVar)
		return nil
	}

	fmt.Printf("  [kind] Deleting cluster %q...\n", h.clusterName)

	args := []string{"delete", "cluster", "--name", h.clusterName}
	cmd := exec.CommandContext(ctx, "kind", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "kind delete cluster failed: %s", stderr.String())
	}

	return nil
}

// kindClusterExists reports whether a kind cluster with the given name is running.
func kindClusterExists(ctx context.Context, name string) (bool, error) {
	cmd := exec.CommandContext(ctx, "kind", "get", "clusters")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return false, errors.Wrapf(err, "kind get clusters failed: %s", stderr.String())
	}
	for _, line := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		if strings.TrimSpace(line) == name {
			return true, nil
		}
	}
	return false, nil
}

// KubeconfigPath returns the path to the kubeconfig for the active cluster.
func (h *Harness) KubeconfigPath() string {
	return h.kubeconfigPath
}

// ClusterName returns the kind cluster name.
func (h *Harness) ClusterName() string {
	return h.clusterName
}

// VerifyDeployed delegates to resource-type-specific verification based on the manifest.
func (h *Harness) VerifyDeployed(ctx context.Context, component string, outputs map[string]interface{}) error {
	manifestPath, _ := ctx.Value(provider.ManifestPathKey{}).(string)
	if manifestPath == "" {
		return errors.New("manifest path not found in context -- cannot verify deployment")
	}
	verifier, err := verify.GetVerifierFromManifest(manifestPath)
	if err != nil {
		return err
	}
	return verifier.VerifyExists(ctx, h.kubeconfigPath)
}

// VerifyDestroyed confirms that resources have been removed after destroy.
func (h *Harness) VerifyDestroyed(ctx context.Context, component string) error {
	manifestPath, _ := ctx.Value(provider.ManifestPathKey{}).(string)
	if manifestPath == "" {
		return errors.New("manifest path not found in context -- cannot verify cleanup")
	}
	verifier, err := verify.GetVerifierFromManifest(manifestPath)
	if err != nil {
		return err
	}
	return verifier.VerifyAbsent(ctx, h.kubeconfigPath)
}
