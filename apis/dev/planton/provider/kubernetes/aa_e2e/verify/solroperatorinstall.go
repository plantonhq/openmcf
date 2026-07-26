package verify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// solrOperatorCrds is the CRD set the modules own for the Solr operator:
// the three solr.apache.org CRDs plus the ZookeeperCluster CRD the
// bundled zookeeper-operator dependency serves.
var solrOperatorCrds = []string{
	"solrclouds.solr.apache.org",
	"solrbackups.solr.apache.org",
	"solrprometheusexporters.solr.apache.org",
	"zookeeperclusters.zookeeper.pravega.io",
}

// SolrOperatorInstallVerifier checks an Apache Solr Operator install:
// the operator Deployment Available and the module-owned CRDs
// Established. On destroy only the release is asserted gone — the CRDs
// are KEPT by design (module-owned keep-on-uninstall so SolrCloud
// resources survive an operator uninstall) and are never orphans.
type SolrOperatorInstallVerifier struct {
	Namespace   string
	ReleaseName string
	// ZookeeperOperator asserts the bundled dependency's Deployment when
	// the scenario installs it (the chart default).
	ZookeeperOperator bool
}

// deploymentName derives the operator Deployment name from the chart's
// fullname rule: the release name verbatim when it contains the chart
// name, else `<release>-solr-operator` (verified against the chart's
// _helpers.tpl).
func (v *SolrOperatorInstallVerifier) deploymentName() string {
	if strings.Contains(v.ReleaseName, "solr-operator") {
		return v.ReleaseName
	}
	return v.ReleaseName + "-solr-operator"
}

func (v *SolrOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] solr operator release %q in namespace %q\n", v.ReleaseName, v.Namespace)

	if err := kubectlWait(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace,
		"condition=Available", 5*time.Minute); err != nil {
		return errors.Wrapf(err, "operator deployment %q never became Available", v.deploymentName())
	}
	for _, crd := range solrOperatorCrds {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %q never became Established", crd)
		}
	}
	if v.ZookeeperOperator {
		// The dependency chart names its deployment after its own chart.
		if err := kubectlWait(ctx, kubeconfig, "deployment", v.ReleaseName+"-zookeeper-operator", v.Namespace,
			"condition=Available", 4*time.Minute); err != nil {
			return errors.Wrap(err, "bundled zookeeper-operator deployment never became Available")
		}
	}
	fmt.Printf("  [verify] operator Available and all %d CRDs Established\n", len(solrOperatorCrds))
	return nil
}

func (v *SolrOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// CRDs deliberately NOT asserted absent (kept by design).
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace)
}
