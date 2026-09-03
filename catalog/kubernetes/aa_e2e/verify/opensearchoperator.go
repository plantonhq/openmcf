package verify

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
)

// opensearchOperatorCrds are the CRDs the chart templates at every version
// the catalog installs; the modules derive them from the pinned chart and
// own their lifecycle.
var opensearchOperatorCrds = []string{
	"opensearchclusters.opensearch.opster.io",
	"opensearchactiongroups.opensearch.opster.io",
	"opensearchcomponenttemplates.opensearch.opster.io",
	"opensearchindextemplates.opensearch.opster.io",
	"opensearchismpolicies.opensearch.opster.io",
	"opensearchroles.opensearch.opster.io",
	"opensearchtenants.opensearch.opster.io",
	"opensearchuserrolebindings.opensearch.opster.io",
	"opensearchusers.opensearch.opster.io",
}

// opensearchSnapshotPolicyCrd joined the chart at 2.8.0; a lane installing an
// older chart must not expect it, and a lane upgrading across 2.8.0 proves
// that a chart bump ADDS a CRD, not only re-applies the existing ones.
const opensearchSnapshotPolicyCrd = "opensearchsnapshotpolicies.opensearch.opster.io"

var opensearchSnapshotPolicySince = semver.MustParse("2.8.0")

// opensearchDefaultChartVersion mirrors the spec's default pin.
const opensearchDefaultChartVersion = "2.8.0"

// opensearchOperatorCrdsAt is the CRD set the chart carries at a version.
func opensearchOperatorCrdsAt(version string) []string {
	crds := append([]string{}, opensearchOperatorCrds...)
	if v, err := semver.NewVersion(version); err == nil && !v.LessThan(opensearchSnapshotPolicySince) {
		crds = append(crds, opensearchSnapshotPolicyCrd)
	}
	return crds
}

// OpenSearchOperatorInstallVerifier checks an OpenSearch Kubernetes Operator
// install: the controller-manager Deployment Available and the module-owned
// CRDs Established and stamped with the manifest's chart version. Destroy
// asserts the CRD posture the manifest declares: kept (the default) or
// deleted (crds.keepOnUninstall: false). The upgrade lane re-runs
// VerifyExists against the upgraded manifest, so the stamp check is what
// proves "a chart bump re-applied the CRDs".
type OpenSearchOperatorInstallVerifier struct {
	Namespace   string
	ReleaseName string
	// The CRD-lifecycle assertions every kind on the primitive shares.
	helmCRDLifecycle
}

// newOpenSearchOperatorInstallVerifier reads the lifecycle-relevant spec
// fields so the same verifier serves the plain, upgrade, cleanup, reinstall
// and refusal lanes.
func newOpenSearchOperatorInstallVerifier(manifestPath, releaseName, namespace string) *OpenSearchOperatorInstallVerifier {
	lifecycle := readHelmCRDLifecycle(manifestPath, opensearchDefaultChartVersion, "opensearch-operator", nil)
	lifecycle.CRDs = opensearchOperatorCrdsAt(lifecycle.ChartVersion)
	return &OpenSearchOperatorInstallVerifier{
		Namespace:        namespace,
		ReleaseName:      releaseName,
		helmCRDLifecycle: lifecycle,
	}
}

// deploymentName mirrors the modules' naming contract: they pin the
// chart's fullnameOverride to the resource name (without the pin the
// chart's default fullname pushes its metrics Service name past the
// 63-character limit — caught live), so the controller-manager
// Deployment is always `<name>-controller-manager`.
func (v *OpenSearchOperatorInstallVerifier) deploymentName() string {
	return v.ReleaseName + "-controller-manager"
}

func (v *OpenSearchOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] opensearch operator release %q in namespace %q\n", v.ReleaseName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace,
		"condition=Available", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "operator deployment %q never became Available", v.deploymentName())
	}
	return v.verifyEstablishedAndStamped(ctx, kubeconfig)
}

func (v *OpenSearchOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace); err != nil {
		return err
	}
	return v.verifyDestroyPosture(ctx, kubeconfig)
}

// VerifyExpectedDeployFailure pins a refused deploy or upgrade to the
// primitive's three-part refusal (see helmCRDLifecycle).
func (v *OpenSearchOperatorInstallVerifier) VerifyExpectedDeployFailure(ctx context.Context, kubeconfig, expectation string, deployErr error) error {
	return v.verifyRefusal(ctx, kubeconfig, expectation, deployErr)
}
