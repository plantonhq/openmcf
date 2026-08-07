package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// kubeRayOperatorCrds are the three CRDs the chart ships from its crds/
// directory: Helm installs them once and KEEPS them on uninstall, so
// destroy asserts their survival, not their deletion.
var kubeRayOperatorCrds = []string{
	"rayclusters.ray.io",
	"rayjobs.ray.io",
	"rayservices.ray.io",
}

// KubeRayOperatorVerifier checks a KubeRay operator install to the point
// a Ray declaration kind could be applied against it: the operator
// Deployment rolled out, all three ray.io CRDs established — and THE
// DESIGN INVARIANT proven: installing the operator alone deploys NO Ray
// cluster (the declaration kind owns every cluster). On the fenced
// posture it additionally proves the watch fence: a RayCluster applied
// OUTSIDE the watched namespaces is never reconciled — no head pod
// appears and the CR's status stays empty.
type KubeRayOperatorVerifier struct {
	// Namespace is the release namespace; the module pins the chart
	// fullname to the resource name, so the Deployment is
	// deployment/<Name> here.
	Namespace string
	Name      string
	// WatchNamespaces is the watch fence from the spec. Non-empty means
	// the operator reconciles Ray CRs ONLY in these namespaces; empty
	// means cluster-wide watch.
	WatchNamespaces []string
}

func (v *KubeRayOperatorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] kuberay-operator %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the operator deployment never rolled out")
	}
	if err := waitForCrdsEstablished(ctx, kubeconfig, kubeRayOperatorCrds); err != nil {
		return err
	}

	// THE DESIGN INVARIANT: the operator install deploys no Ray cluster
	// — a RayCluster CR appearing here would mean the chart grew an
	// auto-provision path the two-kind grain does not expect.
	clusters, err := listCustomResourcesAllNamespaces(ctx, kubeconfig, "rayclusters.ray.io", "")
	if err != nil {
		return errors.Wrap(err, "listing RayCluster CRs")
	}
	if clusters != "" {
		return errors.Errorf("a RayCluster CR exists after installing the operator alone (found: %s) — the declaration kind owns every cluster", clusters)
	}
	fmt.Printf("  [verify] INVARIANT: operator rolled out, 3 CRDs established, and NO Ray cluster deployed — the declaration kind owns clusters\n")

	if len(v.WatchNamespaces) > 0 {
		return v.proveWatchFence(ctx, kubeconfig)
	}
	return nil
}

// proveWatchFence is THE FENCE PROOF: a minimal RayCluster in a
// verifier-owned namespace OUTSIDE the watch set must be IGNORED by the
// operator — no head pod is created and the CR's status stays entirely
// empty. A head pod (or any status) appearing would mean the fence leaks
// and the single-tenant watch posture is a lie. The CR body is
// upstream's ray-cluster.sample.yaml (kuberay ray-operator
// config/samples) shrunk to head-only with small resources: the point is
// to witness NON-reconciliation, so nothing ever needs to actually run.
func (v *KubeRayOperatorVerifier) proveWatchFence(ctx context.Context, kubeconfig string) error {
	fenceNamespace := v.Name + "-fence"
	proofName := v.Name + "-fence-proof"

	// Sweep the CR first (confirmed) and only then the namespace — a
	// namespace deletion racing a live CR can wedge on finalizers.
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		_ = kubectlDeleteResource(cleanupCtx, kubeconfig, "rayclusters.ray.io", proofName, fenceNamespace)
		_ = KubectlResourceAbsent(cleanupCtx, kubeconfig, "rayclusters.ray.io", proofName, fenceNamespace)
		_ = kubectlDeleteResource(cleanupCtx, kubeconfig, "namespace", fenceNamespace, "")
	}()

	if err := ensureNamespace(ctx, kubeconfig, fenceNamespace); err != nil {
		return err
	}

	manifest := fmt.Sprintf(`apiVersion: ray.io/v1
kind: RayCluster
metadata:
  name: %s
  namespace: %s
spec:
  rayVersion: '2.52.0'
  headGroupSpec:
    rayStartParams: {}
    template:
      spec:
        containers:
          - name: ray-head
            image: rayproject/ray:2.52.0
            resources:
              requests:
                cpu: 500m
                memory: 1Gi
              limits:
                cpu: 500m
                memory: 1Gi
`, proofName, fenceNamespace)
	if out, err := applyManifestString(ctx, kubeconfig, manifest); err != nil {
		return errors.Errorf("applying the fence-probe RayCluster: %v: %s", err, firstLines(out, 3))
	}

	// 60 seconds is several reconcile intervals — long enough that a
	// leaking watch would have created the head pod and stamped status.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(60 * time.Second):
	}

	// No head pod: the operator materializes a RayCluster by creating
	// pods labeled ray.io/node-type=head — their absence is the fence.
	pods, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "pods", "-n", fenceNamespace, "-l", "ray.io/node-type=head", "-o", "name").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "listing head pods in the fence namespace: %s", firstLines(string(pods), 3))
	}
	if strings.TrimSpace(string(pods)) != "" {
		return errors.Errorf("the operator created a head pod OUTSIDE the watch namespaces (%s) — the watch fence leaks", strings.TrimSpace(string(pods)))
	}
	status, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "rayclusters.ray.io", proofName, "-n", fenceNamespace,
		"-o", "jsonpath={.status}").Output()
	if err != nil {
		return errors.Wrap(err, "reading the fence-probe status")
	}
	if strings.TrimSpace(string(status)) != "" {
		return errors.Errorf("the operator reconciled a RayCluster OUTSIDE the watch namespaces (status: %s) — the watch fence leaks", firstLines(string(status), 2))
	}
	fmt.Printf("  [verify] THE FENCE PROOF: RayCluster outside the watched namespaces stayed unreconciled (no head pod, empty status) after 60s — the watch fence holds\n")
	return nil
}

// VerifyAbsent asserts the destroy posture: the operator Deployment gone,
// all three CRDs SURVIVING (the crds/-directory keep — a designed
// outcome, asserted, not tolerated), and no RayCluster CRs anywhere.
func (v *KubeRayOperatorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	for _, crd := range kubeRayOperatorCrds {
		if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %q was DELETED on destroy — the crds/-directory keep posture regressed", crd)
		}
	}
	// The CRDs survive, so the API is still queryable: any RayCluster
	// remaining would be verifier residue or an unexpected tenant.
	clusters, err := listCustomResourcesAllNamespaces(ctx, kubeconfig, "rayclusters.ray.io", "")
	if err != nil {
		return errors.Wrap(err, "listing leftover RayCluster CRs")
	}
	if clusters != "" {
		return errors.Errorf("RayCluster CRs survived the destroy: %s", clusters)
	}
	fmt.Printf("  [verify] DESTROY: operator deployment gone, all 3 CRDs RETAINED by design, no RayCluster CRs anywhere\n")
	return nil
}
