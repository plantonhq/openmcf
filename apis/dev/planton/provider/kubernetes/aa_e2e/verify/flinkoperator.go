package verify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// flinkOperatorCrds are the four CRDs the chart ships from its crds/
// directory: Helm installs them once and KEEPS them on uninstall, so
// destroy asserts their survival, not their deletion.
var flinkOperatorCrds = []string{
	"flinkdeployments.flink.apache.org",
	"flinksessionjobs.flink.apache.org",
	"flinkstatesnapshots.flink.apache.org",
	"flinkbluegreendeployments.flink.apache.org",
}

// flinkWebhookService is the chart-fixed name of the operator's webhook
// Service (no fullname derivation exists for it upstream).
const flinkWebhookService = "flink-operator-webhook-service"

// FlinkOperatorVerifier checks an Apache Flink Kubernetes Operator
// install to the kind's actual definition of working. On webhook lanes
// that definition is THE ADMISSION GATE: the fail-closed validating
// webhook must REJECT an invalid FlinkDeployment (standby JobManagers
// without HA — the operator's own validator rule) at apply time. A
// webhook that admits garbage is worse than none: it costs the
// cert-machinery and blocks nothing. The fenced posture also proves the
// gate's SCOPING (the webhook's namespaceSelector excludes namespaces
// outside the watch set — the fence is real in both directions), and the
// webhook-less posture proves THE POSTURE CONTRAST (the same invalid CR
// is ACCEPTED at admission — validation deferred to reconcile, the
// honest trade that arm buys).
type FlinkOperatorVerifier struct {
	// Namespace is the release namespace; the module pins the chart
	// nameOverride (and fullnameOverride) to the resource name, so the
	// Deployment is deployment/<Name> here — the chart's
	// `flink-operator.name` helper honors nameOverride, not fullname.
	Namespace string
	Name      string
	// WebhookEnabled mirrors the spec: it selects between the admission
	// gate proof and the posture-contrast proof.
	WebhookEnabled bool
	// WatchNamespaces is the watch fence from the spec. Non-empty means
	// RBAC AND the webhook namespaceSelector are scoped to these
	// namespaces (and the module owns those namespaces); empty means
	// cluster-wide.
	WatchNamespaces []string
}

func (v *FlinkOperatorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] flink-operator %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the operator deployment never rolled out")
	}
	if err := waitForCrdsEstablished(ctx, kubeconfig, flinkOperatorCrds); err != nil {
		return err
	}
	fmt.Printf("  [verify] operator rolled out and all 4 flink.apache.org CRDs Established\n")

	if v.WebhookEnabled {
		if err := KubectlResourceExists(ctx, kubeconfig, "service", flinkWebhookService, v.Namespace); err != nil {
			return errors.Wrap(err, "the webhook service is missing on a webhook-enabled lane")
		}
		// The keystore Secret is MODULE-generated — a random credential
		// replacing the chart's hardcoded default keystore password. Its
		// presence is the proof the webhook's TLS keystore is not
		// protected by a password every reader of the chart knows.
		keystoreSecret := v.Name + "-webhook-keystore"
		if err := KubectlResourceExists(ctx, kubeconfig, "secret", keystoreSecret, v.Namespace); err != nil {
			return errors.Wrapf(err, "the module-generated webhook keystore secret %q is missing", keystoreSecret)
		}
		fmt.Printf("  [verify] webhook service present and secret %q present — the module-generated keystore credential replaces the chart's hardcoded default\n", keystoreSecret)
	}

	// THE DESIGN INVARIANT: the operator install deploys no Flink
	// cluster — a FlinkDeployment CR appearing here would mean the chart
	// grew an auto-provision path the two-kind grain does not expect.
	deployments, err := listCustomResourcesAllNamespaces(ctx, kubeconfig, "flinkdeployments.flink.apache.org", "")
	if err != nil {
		return errors.Wrap(err, "listing FlinkDeployment CRs")
	}
	if deployments != "" {
		return errors.Errorf("a FlinkDeployment CR exists after installing the operator alone (found: %s) — the declaration kind owns every Flink cluster", deployments)
	}
	fmt.Printf("  [verify] INVARIANT: NO FlinkDeployment deployed by the operator install — the declaration kind owns clusters\n")

	if v.WebhookEnabled {
		if err := v.proveAdmissionGate(ctx, kubeconfig); err != nil {
			return err
		}
		if len(v.WatchNamespaces) > 0 {
			return v.proveGateScoping(ctx, kubeconfig)
		}
		return nil
	}
	return v.provePostureContrast(ctx, kubeconfig)
}

// invalidFlinkDeployment renders the CR every proof arm shares: standby
// JobManagers (replicas: 2) WITHOUT any HA configuration — the operator's
// validator rejects standby replicas when HA is not enabled, so this CR
// is invalid by the operator's own rules while being schema-valid (it
// passes CRD validation and only the webhook / reconcile validator can
// object).
func (v *FlinkOperatorVerifier) invalidFlinkDeployment(name, namespace string) string {
	return fmt.Sprintf(`apiVersion: flink.apache.org/v1beta1
kind: FlinkDeployment
metadata:
  name: %s
  namespace: %s
spec:
  flinkVersion: v2_1
  image: flink:2.1
  serviceAccount: flink
  jobManager:
    replicas: 2
    resource:
      cpu: 0.5
      memory: "1Gi"
  taskManager:
    resource:
      cpu: 0.5
      memory: "1Gi"
`, name, namespace)
}

// proveAdmissionGate is THE ADMISSION GATE: the invalid FlinkDeployment
// must be REJECTED at apply time by the fail-closed validating webhook —
// the behavioral definition of "the webhook works" on this lane. On the
// fenced posture the probe goes INSIDE the first watch namespace, where
// the webhook's namespaceSelector must match.
func (v *FlinkOperatorVerifier) proveAdmissionGate(ctx context.Context, kubeconfig string) error {
	gateNamespace := v.Namespace
	if len(v.WatchNamespaces) > 0 {
		gateNamespace = v.WatchNamespaces[0]
		if err := ensureNamespace(ctx, kubeconfig, gateNamespace); err != nil {
			return err
		}
	}
	proofName := v.Name + "-gate-proof"

	out, err := applyManifestString(ctx, kubeconfig, v.invalidFlinkDeployment(proofName, gateNamespace))
	if err == nil {
		// It went through — remove it (confirmed: the operator watches
		// here and a finalizer-held leftover would wedge the destroy
		// phase that still runs after this failure) and fail loudly.
		_ = kubectlDeleteResource(ctx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, gateNamespace)
		_ = KubectlResourceAbsent(ctx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, gateNamespace)
		return errors.New("ADMISSION GATE FAILED: a FlinkDeployment with standby JobManagers and no HA was ADMITTED — the validating webhook is not enforcing")
	}
	// The denial must come FROM the webhook — any other apply failure
	// (connection refused, missing CRD) is a different broken thing.
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "webhook") && !strings.Contains(lower, "denied") && !strings.Contains(lower, "validat") {
		return errors.Errorf("the invalid FlinkDeployment failed to apply, but NOT with a webhook denial: %s", firstLines(out, 3))
	}
	fmt.Printf("  [verify] THE ADMISSION GATE: invalid FlinkDeployment REJECTED at admission (%s)\n", firstLines(out, 1))
	return nil
}

// proveGateScoping is THE SCOPING ARM: the SAME invalid CR applied in a
// namespace OUTSIDE the watch fence must be ADMITTED — the webhook's
// namespaceSelector excludes it, so admission never evaluates it. This
// proves the fence in the OTHER direction: the gate polices exactly the
// namespaces the spec fenced, no more. The CR is deleted immediately
// after admission (the operator does not watch there anyway, so nothing
// ever acts on it).
func (v *FlinkOperatorVerifier) proveGateScoping(ctx context.Context, kubeconfig string) error {
	outsideNamespace := v.Name + "-outside"
	proofName := v.Name + "-gate-proof"

	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		_ = kubectlDeleteResource(cleanupCtx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, outsideNamespace)
		_ = KubectlResourceAbsent(cleanupCtx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, outsideNamespace)
		_ = kubectlDeleteResource(cleanupCtx, kubeconfig, "namespace", outsideNamespace, "")
	}()

	if err := ensureNamespace(ctx, kubeconfig, outsideNamespace); err != nil {
		return err
	}
	out, err := applyManifestString(ctx, kubeconfig, v.invalidFlinkDeployment(proofName, outsideNamespace))
	if err != nil {
		return errors.Errorf("the invalid FlinkDeployment was REJECTED outside the watch fence — the webhook namespaceSelector is not scoping (%s)", firstLines(out, 3))
	}
	// Admitted, as the scoping predicts. Remove it right away.
	if err := kubectlDeleteResource(ctx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, outsideNamespace); err != nil {
		return errors.Wrap(err, "sweeping the scoping-probe FlinkDeployment")
	}
	fmt.Printf("  [verify] THE SCOPING ARM: the same invalid FlinkDeployment ADMITTED outside the watch fence — the webhook namespaceSelector scopes the gate exactly to the fence\n")
	return nil
}

// provePostureContrast is THE POSTURE CONTRAST on the webhook-less lane:
// the invalid CR the webhook lanes reject must be ACCEPTED here — no
// admission gate exists, so validation is deferred to the operator's
// reconcile loop. Asserting the acceptance (then sweeping before the
// operator acts) teaches the honest trade of the webhook-less arm with a
// live assertion instead of a comment.
func (v *FlinkOperatorVerifier) provePostureContrast(ctx context.Context, kubeconfig string) error {
	proofName := v.Name + "-gate-proof"

	out, err := applyManifestString(ctx, kubeconfig, v.invalidFlinkDeployment(proofName, v.Namespace))
	if err != nil {
		return errors.Errorf("the invalid FlinkDeployment was REJECTED on the webhook-less lane — something is still gating admission (%s)", firstLines(out, 3))
	}
	// Unlike the scoping arm's outside namespace, the operator WATCHES
	// this namespace and may already hold a finalizer on the probe —
	// deletion must be CONFIRMED while the operator is alive, or the
	// leftover wedges the namespace teardown on the finalizer.
	if err := kubectlDeleteResource(ctx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, v.Namespace); err != nil {
		return errors.Wrap(err, "sweeping the contrast-probe FlinkDeployment")
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "flinkdeployments.flink.apache.org", proofName, v.Namespace); err != nil {
		return errors.Wrap(err, "the contrast-probe FlinkDeployment never finished deleting — its finalizer must clear while the operator lives")
	}
	fmt.Printf("  [verify] THE POSTURE CONTRAST: invalid FlinkDeployment ACCEPTED at admission on the webhook-less lane — validation deferred to reconcile, the trade this arm buys\n")
	return nil
}

// VerifyAbsent asserts the destroy posture: the operator Deployment gone,
// all four CRDs SURVIVING (the crds/-directory keep — a designed outcome,
// asserted, not tolerated), the webhook Service gone when it existed,
// module-owned watch namespaces gone, and zero FlinkDeployment CRs
// (which covers any gate-proof leftovers).
func (v *FlinkOperatorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	for _, crd := range flinkOperatorCrds {
		if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %q was DELETED on destroy — the crds/-directory keep posture regressed", crd)
		}
	}
	if v.WebhookEnabled {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "service", flinkWebhookService, v.Namespace); err != nil {
			return errors.Wrap(err, "the webhook service survived destroy")
		}
	}
	for _, ns := range v.WatchNamespaces {
		if err := KubectlResourceAbsent(ctx, kubeconfig, "namespace", ns, ""); err != nil {
			return errors.Wrapf(err, "module-owned watch namespace %q survived destroy", ns)
		}
	}
	// The CRDs survive, so the API is still queryable: any
	// FlinkDeployment remaining would be verifier residue (a leaked
	// gate proof) or an unexpected tenant.
	leftovers, err := listCustomResourcesAllNamespaces(ctx, kubeconfig, "flinkdeployments.flink.apache.org", "")
	if err != nil {
		return errors.Wrap(err, "listing leftover FlinkDeployment CRs")
	}
	if leftovers != "" {
		return errors.Errorf("FlinkDeployment CRs survived the destroy: %s", leftovers)
	}
	fmt.Printf("  [verify] DESTROY: operator deployment gone, all 4 CRDs RETAINED by design, webhook service gone, module-owned watch namespaces gone, no FlinkDeployment CRs anywhere\n")
	return nil
}

// flinkOperatorWebhookEnabled reads spec.webhook_enabled — the webhook
// is ON by default (the chart posture this component keeps), so only an
// explicit false selects the webhook-less arm.
func flinkOperatorWebhookEnabled(spec map[string]interface{}) bool {
	for _, key := range []string{"webhook_enabled", "webhookEnabled"} {
		if enabled, ok := spec[key].(bool); ok {
			return enabled
		}
	}
	return true
}
