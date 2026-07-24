package verify

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// opensearchOperatorCrds is the exact CRD set the operator serves at the
// pinned stable release — the modules own their lifecycle (the chart's
// own CRD templating is disabled), so the install verifier asserts every
// one of them Established.
var opensearchOperatorCrds = []string{
	"opensearchclusters.opensearch.opster.io",
	"opensearchactiongroups.opensearch.opster.io",
	"opensearchcomponenttemplates.opensearch.opster.io",
	"opensearchindextemplates.opensearch.opster.io",
	"opensearchismpolicies.opensearch.opster.io",
	"opensearchroles.opensearch.opster.io",
	"opensearchsnapshotpolicies.opensearch.opster.io",
	"opensearchtenants.opensearch.opster.io",
	"opensearchuserrolebindings.opensearch.opster.io",
	"opensearchusers.opensearch.opster.io",
}

// OpenSearchOperatorInstallVerifier checks an OpenSearch Kubernetes
// Operator install: the controller-manager Deployment Available and the
// module-owned CRDs Established. On destroy only the release is asserted
// gone — the CRDs are KEPT by design (module-owned keep-on-uninstall so
// OpenSearchCluster resources survive an operator uninstall) and are
// never orphans.
type OpenSearchOperatorInstallVerifier struct {
	Namespace   string
	ReleaseName string
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
	for _, crd := range opensearchOperatorCrds {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %q never became Established", crd)
		}
	}
	fmt.Printf("  [verify] operator Available and all %d OpenSearch CRDs Established\n", len(opensearchOperatorCrds))
	return nil
}

func (v *OpenSearchOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// CRDs deliberately NOT asserted absent: the module keeps them on
	// destroy so search clusters and their data outlive the operator.
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace)
}
