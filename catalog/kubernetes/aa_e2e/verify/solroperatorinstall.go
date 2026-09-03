package verify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// solrOperatorCrds are the CRDs the modules derive from the pinned chart:
// the three solr.apache.org CRDs from the chart's crds/ directory, plus the
// ZookeeperCluster CRD the bundled zookeeper-operator subchart templates
// when it is installed (the chart default). The set follows the chart's own
// behaviour: with the subchart off, its CRD is not rendered either.
var solrOperatorCrds = []string{
	"solrclouds.solr.apache.org",
	"solrbackups.solr.apache.org",
	"solrprometheusexporters.solr.apache.org",
}

const solrZookeeperClusterCrd = "zookeeperclusters.zookeeper.pravega.io"

// solrDefaultChartVersion mirrors the spec's default pin. A scenario that
// leaves chartVersion unset installs this version, and the CRDs must carry
// it in their stamp.
const solrDefaultChartVersion = "0.9.1"

// SolrOperatorInstallVerifier checks an Apache Solr Operator install: the
// operator Deployment Available, the bundled zookeeper-operator Available
// when the scenario installs it, and the module-owned CRDs Established and
// stamped with the manifest's chart version. Destroy asserts the CRD posture
// the manifest declares: kept (the default) or deleted
// (crds.keepOnUninstall: false). The upgrade lane re-runs VerifyExists
// against the upgraded manifest, so the stamp check is what proves "a chart
// bump re-applied the CRDs".
type SolrOperatorInstallVerifier struct {
	Namespace   string
	ReleaseName string
	// ZookeeperOperator asserts the bundled dependency's Deployment when
	// the scenario installs it (the chart default).
	ZookeeperOperator bool
	// The CRD-lifecycle assertions every kind on the primitive shares.
	helmCRDLifecycle
}

// newSolrOperatorInstallVerifier reads the lifecycle-relevant spec fields
// so the same verifier serves the plain, upgrade, cleanup, reinstall and
// refusal lanes.
func newSolrOperatorInstallVerifier(manifestPath, releaseName, namespace string) *SolrOperatorInstallVerifier {
	zookeeper := solrZookeeperOperatorInstalled(manifestSpecMap(manifestPath))
	crds := append([]string{}, solrOperatorCrds...)
	if zookeeper {
		crds = append(crds, solrZookeeperClusterCrd)
	}
	return &SolrOperatorInstallVerifier{
		Namespace:         namespace,
		ReleaseName:       releaseName,
		ZookeeperOperator: zookeeper,
		helmCRDLifecycle:  readHelmCRDLifecycle(manifestPath, solrDefaultChartVersion, "solr-operator", crds),
	}
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
	if v.ZookeeperOperator {
		// The dependency chart names its deployment after its own chart.
		if err := kubectlWait(ctx, kubeconfig, "deployment", v.ReleaseName+"-zookeeper-operator", v.Namespace,
			"condition=Available", 4*time.Minute); err != nil {
			return errors.Wrap(err, "bundled zookeeper-operator deployment never became Available")
		}
	}
	return v.verifyEstablishedAndStamped(ctx, kubeconfig)
}

func (v *SolrOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.deploymentName(), v.Namespace); err != nil {
		return err
	}
	return v.verifyDestroyPosture(ctx, kubeconfig)
}

// VerifyExpectedDeployFailure pins a refused deploy or upgrade to the
// primitive's three-part refusal (see helmCRDLifecycle).
func (v *SolrOperatorInstallVerifier) VerifyExpectedDeployFailure(ctx context.Context, kubeconfig, expectation string, deployErr error) error {
	return v.verifyRefusal(ctx, kubeconfig, expectation, deployErr)
}
