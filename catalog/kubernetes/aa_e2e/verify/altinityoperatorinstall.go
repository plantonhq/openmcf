package verify

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// altinityOperatorCrds is the exact CRD set the Altinity chart ships at
// the pinned release. Unlike module-owned CRD kinds, these are
// CHART-owned: Helm's native crds/ handling installs them on first
// install and the chart's pre-install/pre-upgrade hook job server-side
// applies them on every install and upgrade — the install verifier
// asserts all four Established regardless of which path created them.
var altinityOperatorCrds = []string{
	"clickhouseinstallations.clickhouse.altinity.com",
	"clickhouseinstallationtemplates.clickhouse.altinity.com",
	"clickhouseoperatorconfigurations.clickhouse.altinity.com",
	"clickhousekeeperinstallations.clickhouse-keeper.altinity.com",
}

// AltinityOperatorInstallVerifier checks an Altinity ClickHouse
// operator install: the operator Deployment Available and the four
// CRDs Established. On destroy only the Deployment is asserted gone —
// Helm never deletes crds/-shipped CRDs on uninstall (keep-on-uninstall
// is inherent), so ClickHouse clusters and their data survive an
// operator removal by design.
type AltinityOperatorInstallVerifier struct {
	Namespace   string
	ReleaseName string
}

// deploymentName mirrors the modules' naming contract: the chart
// fullname is pinned to the resource name, and the operator Deployment
// is named exactly the fullname.
func (v *AltinityOperatorInstallVerifier) deploymentName() string {
	return v.ReleaseName
}

func (v *AltinityOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] altinity operator release %q in namespace %q\n", v.ReleaseName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace,
		"condition=Available", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "operator deployment %q never became Available", v.deploymentName())
	}
	for _, crd := range altinityOperatorCrds {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %q never became Established", crd)
		}
	}
	fmt.Printf("  [verify] operator Available and all %d ClickHouse CRDs Established\n", len(altinityOperatorCrds))
	return nil
}

func (v *AltinityOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// CRDs deliberately NOT asserted absent: Helm keeps crds/-shipped
	// CRDs on uninstall so ClickHouse clusters outlive the operator.
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace)
}
