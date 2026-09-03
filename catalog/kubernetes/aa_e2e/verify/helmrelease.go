package verify

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// expectCrdsAnnotation names, comma-separated, the CRDs a KubernetesHelmRelease
// scenario expects the module to own for the chart it installs (or, on a
// refusal lane, the CRDs the refusal must name). The chart is arbitrary, so
// the verifier cannot know its CRD set; the scenario declares it, and the
// verifier holds the module to it (a chart with no CRDs simply omits the
// annotation).
const expectCrdsAnnotation = "planton.dev/e2e-expect-crds"

// HelmReleaseVerifier checks a KubernetesHelmRelease: the namespace and a
// running workload (what any chart install must yield) plus the CRD
// lifecycle the module promises for the chart's crds/ directory: the
// declared CRDs Established and stamped with the pinned version; kept or
// deleted on destroy as the manifest says; a refused deploy explained in
// three parts.
type HelmReleaseVerifier struct {
	HelmComponentVerifier
	helmCRDLifecycle
}

// newHelmReleaseVerifier reads the chart identity and the lifecycle fields
// from the scenario. The generic kind's version field is `version` (Helm's
// own name), and the source label the module stamps is the chart name.
func newHelmReleaseVerifier(manifestPath, name, namespace string) *HelmReleaseVerifier {
	spec := manifestSpecMap(manifestPath)
	version, _ := specField(spec, "version").(string)
	chart, _ := specField(spec, "chart").(string)
	var crds []string
	if declared := manifestAnnotation(manifestPath, expectCrdsAnnotation); declared != "" {
		for _, crd := range strings.Split(declared, ",") {
			if crd = strings.TrimSpace(crd); crd != "" {
				crds = append(crds, crd)
			}
		}
	}
	lifecycle := readHelmCRDLifecycle(manifestPath, version, chart, crds)
	// The generic kind pins `version`, not `chartVersion`; the reader's
	// default already carries it, and there is no other spelling to read.
	lifecycle.ChartVersion = version
	return &HelmReleaseVerifier{
		HelmComponentVerifier: HelmComponentVerifier{Namespace: namespace, ComponentName: name},
		helmCRDLifecycle:      lifecycle,
	}
}

// VerifyExists asserts what the module promises for an arbitrary chart: the
// namespace exists, the chart's workload is Running, and the CRD lifecycle
// holds for the CRDs the scenario declares. Whether the chart deploys a
// Service is the chart's business (Flagger deploys none; podinfo deploys
// one), so unlike the Tier-2 Helm kinds the generic verifier does not
// assert one.
func (v *HelmReleaseVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] Helm release %q in namespace %q\n", v.ComponentName, v.Namespace)
	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for helm release %q", v.Namespace, v.ComponentName)
	}
	if err := KubectlPodsRunningInNamespace(ctx, kubeconfig, v.Namespace); err != nil {
		return errors.Wrapf(err, "no running pods in namespace %q for helm release %q", v.Namespace, v.ComponentName)
	}
	return v.verifyEstablishedAndStamped(ctx, kubeconfig)
}

func (v *HelmReleaseVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := v.HelmComponentVerifier.VerifyAbsent(ctx, kubeconfig); err != nil {
		return err
	}
	return v.verifyDestroyPosture(ctx, kubeconfig)
}

// VerifyExpectedDeployFailure pins a refused deploy or upgrade to the
// primitive's three-part refusal (see helmCRDLifecycle).
func (v *HelmReleaseVerifier) VerifyExpectedDeployFailure(ctx context.Context, kubeconfig, expectation string, deployErr error) error {
	return v.verifyRefusal(ctx, kubeconfig, expectation, deployErr)
}
