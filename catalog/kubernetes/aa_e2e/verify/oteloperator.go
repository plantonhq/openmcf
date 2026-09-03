package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// otelCrds are the CRDs the module derives from the pinned chart at its
// default values and applies OUTSIDE the release (kept on destroy unless the
// spec says otherwise -- an operator uninstall never deletes the fleet's
// declarations). The feature-gated fifth CRD (clusterobservabilities) is
// absent at default values because its gate is off.
var otelCrds = []string{
	"opentelemetrycollectors.opentelemetry.io",
	"instrumentations.opentelemetry.io",
	"opampbridges.opentelemetry.io",
	"targetallocators.opentelemetry.io",
}

// otelDefaultChartVersion mirrors the spec's default pin. A scenario that
// leaves chartVersion unset installs this version, and the CRDs must carry
// it in their stamp.
const otelDefaultChartVersion = "0.120.0"

// Stamp keys the module writes on every derived CRD; the source of truth is
// pkg/kubernetes/helmcrds, spelled out here so the verifier reads exactly
// what a user would read with kubectl.
const (
	crdSourceVersionAnnotation = "planton.ai/crd-source-version"
	crdSourceLabel             = "planton.ai/crd-source"
)

// The failure classes the OTel lanes pin. Each names one of the module's
// three-part refusals; the verifier asserts the engine's error carries all
// three parts, so a mechanism-only message can never pass.
const (
	otelFailureChartVersionNotPublished = "chart-version-not-published"
	otelFailureCrdSchemaDowngrade       = "crd-schema-downgrade"
)

// OtelOperatorVerifier checks an OpenTelemetry Operator install to the
// point a KubernetesOtelCollector declaration could be applied against
// it: the manager Deployment rolled out (its pod mounts the
// cert-manager-issued webhook Secret, so a rollout IS the
// cert-issuance proof), all four opentelemetry.io CRDs established and
// stamped with the manifest's chart version -- and the TWO webhook proofs
// that only a live cluster can give:
//
//   - THE ADMISSION GATE: an invalid collector CR (a pipeline
//     referencing an undeclared receiver) is REJECTED by the
//     fail-closed validating webhook.
//   - THE CONVERSION PROOF: a v1beta1 collector CR written with
//     managementState Unmanaged (admitted, never reconciled into a
//     workload) reads back through the v1alpha1 served version -- the
//     conversion webhook call only succeeds when cert-manager's CA
//     injector has patched the MODULE-owned (kept) CRD's conversion
//     caBundle, which is the exact trust seam this component's design
//     hangs on.
//
// Destroy asserts the CRD posture the manifest declares: kept (the
// default) or deleted (crds.keepOnUninstall: false). The upgrade lane
// re-runs VerifyExists against the upgraded manifest, so the stamp check
// is what proves "a chart bump re-applied the CRDs".
type OtelOperatorVerifier struct {
	// Namespace is the install namespace from the spec.
	Namespace string
	// ReleaseName is metadata.name -- the module pins the chart's
	// fullname to it, so the manager Deployment shares the name.
	ReleaseName string
	// ChartVersion is the version the manifest pins (or the default), the
	// value every derived CRD must carry in its source-version stamp.
	ChartVersion string
	// KeepOnUninstall is the manifest's crds.keepOnUninstall (default
	// true): whether destroy must leave the CRDs or take them.
	KeepOnUninstall bool
	// InstallCrds is the manifest's crds.install (default true).
	InstallCrds bool
	// DeployRefused marks a scenario whose deploy is DESIGNED to be refused
	// (the expect-deploy-failure lane): nothing was created, so the destroy
	// assertion has no keep or delete to check, only that no workload exists.
	DeployRefused bool
}

// newOtelOperatorVerifier reads the lifecycle-relevant spec fields so the
// same verifier serves the plain, upgrade, cleanup, and reinstall lanes.
func newOtelOperatorVerifier(manifestPath, releaseName, namespace string) *OtelOperatorVerifier {
	v := &OtelOperatorVerifier{
		Namespace:       namespace,
		ReleaseName:     releaseName,
		ChartVersion:    otelDefaultChartVersion,
		KeepOnUninstall: true,
		InstallCrds:     true,
	}
	spec := manifestSpecMap(manifestPath)
	if version, _ := spec["chartVersion"].(string); version != "" {
		v.ChartVersion = version
	}
	if crds, ok := spec["crds"].(map[string]interface{}); ok {
		if keep, ok := crds["keepOnUninstall"].(bool); ok {
			v.KeepOnUninstall = keep
		}
		if install, ok := crds["install"].(bool); ok {
			v.InstallCrds = install
		}
	}
	v.DeployRefused = manifestAnnotation(manifestPath, "planton.dev/e2e-expect-deploy-failure") != ""
	return v
}

func (v *OtelOperatorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] opentelemetry-operator %q in namespace %q\n", v.ReleaseName, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.ReleaseName, v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the operator deployment never rolled out")
	}
	for _, crd := range otelCrds {
		if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"wait", "--for=condition=Established", "crd/"+crd, "--timeout=120s").CombinedOutput(); err != nil {
			return errors.Wrapf(err, "CRD %q never became Established: %s", crd, firstLines(string(out), 3))
		}
		if v.InstallCrds {
			if err := v.assertCrdStamp(ctx, kubeconfig, crd); err != nil {
				return err
			}
		}
	}
	fmt.Printf("  [verify] operator rolled out; all 4 opentelemetry.io CRDs Established and stamped with chart %s\n", v.ChartVersion)

	if err := v.admissionGate(ctx, kubeconfig); err != nil {
		return err
	}
	return v.conversionProof(ctx, kubeconfig)
}

// admissionGate proves the fail-closed validating webhook is LIVE: a
// collector CR whose pipeline references an undeclared receiver must be
// REJECTED at apply time (the webhook validates pipeline wiring).
func (v *OtelOperatorVerifier) admissionGate(ctx context.Context, kubeconfig string) error {
	invalid := fmt.Sprintf(`apiVersion: opentelemetry.io/v1beta1
kind: OpenTelemetryCollector
metadata:
  name: e2e-invalid-probe
  namespace: %s
spec:
  managementState: unmanaged
  config:
    exporters:
      debug: {}
    service:
      pipelines:
        traces:
          receivers: [otlp]
          exporters: [debug]
`, v.Namespace)

	out, err := applyManifestString(ctx, kubeconfig, invalid)
	if err == nil {
		// It went through — remove it and fail loudly.
		_, _ = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"delete", "opentelemetrycollectors.opentelemetry.io", "e2e-invalid-probe",
			"-n", v.Namespace, "--ignore-not-found").CombinedOutput()
		return errors.New("ADMISSION GATE FAILED: a collector CR with a pipeline referencing an undeclared receiver was ADMITTED — the validating webhook is not enforcing")
	}
	fmt.Printf("  [verify] ADMISSION GATE: invalid collector CR REJECTED by the webhook (%s)\n", firstLines(out, 1))
	return nil
}

// conversionProof exercises the CRD's version-conversion webhook through
// the CA-injected trust chain: write v1beta1 (structured config), read
// back v1alpha1 (string config). Unmanaged, so the operator never
// creates a workload for the probe.
func (v *OtelOperatorVerifier) conversionProof(ctx context.Context, kubeconfig string) error {
	probe := fmt.Sprintf(`apiVersion: opentelemetry.io/v1beta1
kind: OpenTelemetryCollector
metadata:
  name: e2e-conversion-probe
  namespace: %s
spec:
  managementState: unmanaged
  config:
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
    exporters:
      debug: {}
    service:
      pipelines:
        traces:
          receivers: [otlp]
          exporters: [debug]
`, v.Namespace)

	if out, err := applyManifestString(ctx, kubeconfig, probe); err != nil {
		return errors.Wrapf(err, "the valid conversion-probe CR was not admitted: %s", firstLines(out, 3))
	}
	defer func() {
		_, _ = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"delete", "opentelemetrycollectors.opentelemetry.io", "e2e-conversion-probe",
			"-n", v.Namespace, "--ignore-not-found").CombinedOutput()
	}()

	// GET through the v1alpha1 served version — this call round-trips
	// through the conversion webhook (storage is v1beta1), which only
	// works when the kept CRD's conversion caBundle is current.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "opentelemetrycollectors.v1alpha1.opentelemetry.io", "e2e-conversion-probe",
		"-n", v.Namespace, "-o", "jsonpath={.spec.config}").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "CONVERSION PROOF FAILED — the v1alpha1 read did not convert (is the CA injector patching the kept CRD?): %s", firstLines(string(out), 3))
	}
	// v1alpha1 carries config as a STRING — the converted document must
	// contain the receiver we declared.
	if !strings.Contains(string(out), "otlp") {
		return errors.Errorf("CONVERSION PROOF FAILED — the v1alpha1 config did not carry the declared receiver: %s", firstLines(string(out), 3))
	}
	fmt.Printf("  [verify] CONVERSION PROOF: v1beta1-written CR read back through v1alpha1 (the CA-injected conversion webhook is live)\n")
	return nil
}

func (v *OtelOperatorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.ReleaseName, v.Namespace); err != nil {
		return err
	}
	if !v.InstallCrds {
		fmt.Printf("  [verify] DESTROY: operator workloads gone (bring-your-own CRDs: nothing to assert)\n")
		return nil
	}
	if v.DeployRefused {
		fmt.Printf("  [verify] DESTROY: nothing to tear down -- the deploy was refused before anything was created\n")
		return nil
	}
	if v.KeepOnUninstall {
		// The designed keep: the module-owned CRDs survive the destroy so an
		// operator uninstall never cascade-deletes the collector fleet.
		for _, crd := range otelCrds {
			if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
				return errors.Wrapf(err, "CRD %q was DELETED on destroy -- the module-owned keep posture broke", crd)
			}
		}
		fmt.Printf("  [verify] DESTROY: operator workloads gone; all 4 CRDs RETAINED by design\n")
		return nil
	}
	// crds.keepOnUninstall: false -- the destroy must take the CRDs with it.
	for _, crd := range otelCrds {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %q survived a destroy that declared keepOnUninstall: false", crd)
		}
	}
	fmt.Printf("  [verify] DESTROY: operator workloads gone; all 4 CRDs DELETED as the manifest asked\n")
	return nil
}

// assertCrdStamp reads the CRD's source-version annotation and source label:
// the module's own record of which chart version derived it. After an
// upgrade the stamp must have moved to the new version; after a reinstall it
// must be present on the re-adopted CRD.
func (v *OtelOperatorVerifier) assertCrdStamp(ctx context.Context, kubeconfig, crd string) error {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "crd", crd, "-o",
		fmt.Sprintf("jsonpath={.metadata.annotations.%s}|{.metadata.labels.%s}",
			strings.ReplaceAll(crdSourceVersionAnnotation, ".", "\\."),
			strings.ReplaceAll(crdSourceLabel, ".", "\\."))).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading the source stamp of CRD %q: %s", crd, firstLines(string(out), 2))
	}
	parts := strings.SplitN(strings.TrimSpace(string(out)), "|", 2)
	if len(parts) != 2 || parts[0] != v.ChartVersion {
		return errors.Errorf("CRD %q carries source version %q, the manifest pins %q -- the derived CRDs did not follow chart_version", crd, parts[0], v.ChartVersion)
	}
	if parts[1] != "opentelemetry-operator" {
		return errors.Errorf("CRD %q carries source label %q, expected the chart name", crd, parts[1])
	}
	return nil
}

// VerifyExpectedDeployFailure pins a refused deploy or upgrade to the
// module's three-part refusal. The texts are the primitive's own (see
// pkg/kubernetes/helmcrds for the Go side and the generated helm_crds.tf for
// the Terraform side); asserting all three parts is what keeps "count
// mismatch"-class messages from ever passing again.
func (v *OtelOperatorVerifier) VerifyExpectedDeployFailure(ctx context.Context, kubeconfig, expectation string, deployErr error) error {
	// The engines wrap long diagnostics at terminal width and the Terraform
	// runner joins lines with " | ", so a sentence can be split anywhere;
	// collapse all whitespace and separators before matching phrases.
	text := strings.Join(strings.Fields(strings.ReplaceAll(deployErr.Error(), "|", " ")), " ")
	for _, part := range []string{"observed:", "meaning:", "next step:"} {
		if !strings.Contains(text, part) {
			return errors.Errorf("the refusal lacks its %q part -- every CRD-lifecycle failure explains itself in three parts; got: %s", part, firstLines(text, 12))
		}
	}
	switch expectation {
	case otelFailureChartVersionNotPublished:
		// The render or install could not locate the pinned chart; the
		// engines surface Helm's own text inside the observation.
		if !strings.Contains(text, v.ChartVersion) || !(strings.Contains(text, "not found") || strings.Contains(text, "not in the index")) {
			return errors.Errorf("expected the version-not-published refusal naming %s; got: %s", v.ChartVersion, firstLines(text, 12))
		}
		fmt.Printf("  [verify] REFUSED as expected: chart version %s is not published, and the message says what to do\n", v.ChartVersion)
	case otelFailureCrdSchemaDowngrade:
		if !strings.Contains(text, "derived from chart version") || !strings.Contains(text, "asks for chart version "+v.ChartVersion) {
			return errors.Errorf("expected the schema-downgrade refusal naming the cluster's version and %s; got: %s", v.ChartVersion, firstLines(text, 12))
		}
		if !strings.Contains(text, "kubectl delete crd") {
			return errors.Errorf("the downgrade refusal must name the deliberate remedy; got: %s", firstLines(text, 12))
		}
		// Nothing was touched: the CRDs still carry the higher version.
		for _, crd := range otelCrds {
			out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "get", "crd", crd, "-o",
				fmt.Sprintf("jsonpath={.metadata.annotations.%s}", strings.ReplaceAll(crdSourceVersionAnnotation, ".", "\\."))).CombinedOutput()
			if err != nil {
				return errors.Wrapf(err, "reading CRD %q after the refused downgrade", crd)
			}
			if strings.TrimSpace(string(out)) == v.ChartVersion {
				return errors.Errorf("CRD %q was downgraded to %s despite the refusal", crd, v.ChartVersion)
			}
		}
		fmt.Printf("  [verify] REFUSED as expected: the schema downgrade to %s was stopped before any CRD changed\n", v.ChartVersion)
	default:
		return errors.Errorf("unknown expected failure class %q for the OpenTelemetry operator", expectation)
	}
	return nil
}

// applyManifestString kubectl-applies a YAML document from a temp file
// and returns the combined output.
func applyManifestString(ctx context.Context, kubeconfig, manifest string) (string, error) {
	f, err := os.CreateTemp("", "e2e-otel-probe-*.yaml")
	if err != nil {
		return "", err
	}
	defer func() { _ = os.Remove(f.Name()) }()
	if _, err := f.WriteString(manifest); err != nil {
		_ = f.Close()
		return "", err
	}
	_ = f.Close()
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"apply", "-f", filepath.Clean(f.Name())).CombinedOutput()
	return string(out), err
}
