package verify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// helmFullname reproduces Helm's standard fullname template: when the
// release name already contains the chart name the fullname collapses to
// the release name; otherwise it is "<release>-<chart>". Both Percona
// operator charts use the standard helper, so the operator Deployment's
// name is computable from the release name the modules pin to
// metadata.name.
func helmFullname(release, chart string) string {
	if strings.Contains(release, chart) {
		return release
	}
	return release + "-" + chart
}

// PsmdbOperatorInstallVerifier checks a Percona Operator for MongoDB
// installation: the operator Deployment Available and the
// PerconaServerMongoDB CRDs Established — the preconditions every
// KubernetesMongodb apply depends on.
type PsmdbOperatorInstallVerifier struct {
	Namespace string
	// ReleaseName is the Helm release (the module pins it to
	// metadata.name); the Deployment name derives from it via the chart's
	// standard fullname helper.
	ReleaseName string
}

func (v *PsmdbOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	deployment := helmFullname(v.ReleaseName, "psmdb-operator")
	fmt.Printf("  [verify] percona mongodb operator %q in namespace %q\n", deployment, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for the psmdb operator", v.Namespace)
	}
	if err := kubectlWait(ctx, kubeconfig, "deployment", deployment, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "psmdb operator deployment not available")
	}

	// The CRDs KubernetesMongodb renders against: the cluster itself plus
	// the backup/restore kinds the E2E backup proof drives.
	for _, crd := range []string{
		"perconaservermongodbs.psmdb.percona.com",
		"perconaservermongodbbackups.psmdb.percona.com",
		"perconaservermongodbrestores.psmdb.percona.com",
	} {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", 2*time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %s not established", crd)
		}
	}
	return nil
}

func (v *PsmdbOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The CRDs intentionally SURVIVE uninstall (the chart ships them in
	// its Helm-native crds/ directory, which Helm installs once and never
	// deletes) — only the Deployment's absence is asserted.
	deployment := helmFullname(v.ReleaseName, "psmdb-operator")
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", deployment, v.Namespace)
}

// PxcOperatorInstallVerifier checks a Percona Operator for MySQL (PXC)
// installation: the operator Deployment Available and the
// PerconaXtraDBCluster CRDs Established — the preconditions every
// KubernetesMysql apply depends on. Widened-watch installs additionally
// assert the module-owned validation webhook: present (and pointing at
// THIS install's namespace) while deployed, and GONE after destroy — a
// stranded Fail-closed webhook would brick every future
// PerconaXtraDBCluster admission in the cluster, which is exactly the
// failure the module's ownership of the object exists to prevent.
type PxcOperatorInstallVerifier struct {
	Namespace   string
	ReleaseName string
	// WatchWidened mirrors the manifest's watch block: cluster-wide or a
	// namespace fence — the arms in which the webhook is module-owned.
	WatchWidened bool
}

// pxcWebhookName is the operator's fixed cluster-scoped webhook object
// name (upstream registers exactly one per cluster).
const pxcWebhookName = "percona-xtradbcluster-webhook"

// pxcWatchWidened reads the manifest's watch block: cluster-wide or a
// non-empty namespace fence both widen the operator's scope (the arms in
// which the modules own the validation webhook).
func pxcWatchWidened(spec map[string]interface{}) bool {
	watch, ok := spec["watch"].(map[string]interface{})
	if !ok {
		return false
	}
	if clusterWide, ok := watch["cluster_wide"].(bool); ok && clusterWide {
		return true
	}
	if namespaces, ok := watch["namespaces"].([]interface{}); ok && len(namespaces) > 0 {
		return true
	}
	return false
}

func (v *PxcOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	deployment := helmFullname(v.ReleaseName, "pxc-operator")
	fmt.Printf("  [verify] percona mysql (pxc) operator %q in namespace %q\n", deployment, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for the pxc operator", v.Namespace)
	}
	if err := kubectlWait(ctx, kubeconfig, "deployment", deployment, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "pxc operator deployment not available")
	}

	for _, crd := range []string{
		"perconaxtradbclusters.pxc.percona.com",
		"perconaxtradbclusterbackups.pxc.percona.com",
		"perconaxtradbclusterrestores.pxc.percona.com",
	} {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", 2*time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %s not established", crd)
		}
	}

	if v.WatchWidened {
		// The module-owned webhook must exist and route to THIS install's
		// namespace (the upstream update path never corrects a stale
		// service pointer, so the pointer is the assertion that matters).
		out, err := kubectlGetJSONPath(ctx, kubeconfig,
			"validatingwebhookconfiguration", pxcWebhookName, "",
			"{.webhooks[0].clientConfig.service.namespace}")
		if err != nil {
			return errors.Wrapf(err, "module-owned webhook %s not found", pxcWebhookName)
		}
		if strings.TrimSpace(out) != v.Namespace {
			return errors.Errorf("webhook %s routes to namespace %q, want %q",
				pxcWebhookName, strings.TrimSpace(out), v.Namespace)
		}
	}
	return nil
}

func (v *PxcOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// CRDs survive by the Helm-native crds/ posture — Deployment only.
	deployment := helmFullname(v.ReleaseName, "pxc-operator")
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", deployment, v.Namespace); err != nil {
		return err
	}

	if v.WatchWidened {
		// The webhook must be GONE with the resource — a survivor is the
		// cluster-bricking orphan the module's ownership prevents.
		if err := KubectlResourceAbsent(ctx, kubeconfig,
			"validatingwebhookconfiguration", pxcWebhookName, ""); err != nil {
			return errors.Wrap(err, "module-owned validation webhook survived destroy")
		}
	}
	return nil
}
