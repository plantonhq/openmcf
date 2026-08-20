package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// plantonPlatformCrd is the CRD the KubernetesPlantonPlatform declaration
// renders against. It is MODULE-OWNED by the operator kind with
// keep-on-uninstall semantics — destroy asserts it SURVIVES.
const plantonPlatformCrd = "plantonplatforms.planton.ai"

// plantonOperatorDeployment is the chart fullname with the module's fixed
// release name "planton-operator" (release name == chart name, so the
// fullname collapses to the release name).
const plantonOperatorDeployment = "planton-operator"

// PlantonOperatorInstallVerifier checks a Planton operator installation to
// the point a KubernetesPlantonPlatform declaration could be applied
// against it: the manager Deployment Available and the PlantonPlatform CRD
// Established — and THE DESIGN INVARIANT proven on every lane: NO
// PlantonPlatform exists after installing the operator alone. The operator
// never auto-creates a platform; every platform is a deliberate
// declaration, and an auto-created one here would mean the two-kind grain
// regressed into an SSA field-manager fight.
type PlantonOperatorInstallVerifier struct {
	Namespace string
}

func (v *PlantonOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] planton-operator in namespace %q\n", v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for planton-operator", v.Namespace)
	}
	if err := kubectlWait(ctx, kubeconfig, "deployment", plantonOperatorDeployment, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "planton-operator deployment not available (a sibling operator on the cluster makes the startup guard refuse — read the pod log)")
	}
	if err := kubectlWait(ctx, kubeconfig, "crd", plantonPlatformCrd, "",
		"condition=Established", 2*time.Minute); err != nil {
		return errors.Wrapf(err, "CRD %s not established (module-owned — the staged copy applies before the release)", plantonPlatformCrd)
	}

	// THE DESIGN INVARIANT: installing the operator alone deploys no
	// platform. Give the manager a settle window first, so a regression
	// toward startup auto-creation cannot slip under the check.
	time.Sleep(15 * time.Second)
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", plantonPlatformCrd, "-A", "-o", "name").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "listing PlantonPlatforms: %s", firstLines(string(out), 3))
	}
	if strings.TrimSpace(string(out)) != "" {
		return errors.Errorf("a PlantonPlatform exists after installing the operator alone — the operator must never auto-create a platform (found: %s)", strings.TrimSpace(string(out)))
	}
	fmt.Printf("  [verify] INVARIANT: no PlantonPlatform after install — platforms are always deliberate declarations\n")
	return nil
}

func (v *PlantonOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", plantonOperatorDeployment, v.Namespace); err != nil {
		return err
	}
	// THE KEEP POSTURE, asserted positively: the module-owned CRD
	// SURVIVES destroy (retainOnDelete / apply_only), so removing the
	// operator can never cascade-delete platform declarations. A missing
	// CRD here means the keep semantics regressed.
	if err := KubectlResourceExists(ctx, kubeconfig, "crd", plantonPlatformCrd, ""); err != nil {
		return errors.Wrapf(err, "the %s CRD must SURVIVE the operator's destroy (module-owned, keep-on-uninstall) — its absence means the retain semantics regressed", plantonPlatformCrd)
	}
	fmt.Printf("  [verify] KEEP POSTURE: the PlantonPlatform CRD survived the operator's destroy\n")
	return nil
}

// PlantonPlatformVerifier checks a declared platform to the point a person
// could sign in: the PlantonPlatform reaches phase Ready (the operator's
// own per-component gates — databases, identity, control plane, console,
// vault, runner — all pass inside it), and the two first-visit handles the
// module exports actually exist (the gateway Service and the setup-code
// Secret).
//
// Destroy relies on Kubernetes GARBAGE COLLECTION (every operator-created
// object is owner-referenced to the CR — the operator has no finalizers),
// so absence is polled: the CR's deletion returns quickly and the children
// drain asynchronously.
type PlantonPlatformVerifier struct {
	Namespace string
	Name      string
}

func (v *PlantonPlatformVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] planton platform %q in namespace %q\n", v.Name, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, plantonPlatformCrd, v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "the PlantonPlatform declaration was not created")
	}

	// The whole platform boots inside this one wait — databases, identity
	// server, control plane, console, secrets manager, runner. On an
	// emulated-amd64 kind cluster the full boot runs 10-15 minutes; the
	// budget leaves headroom for cold image pulls.
	fmt.Printf("  [verify] waiting for phase Ready (a full platform boot — expect ~10-15 minutes on a kind cluster)\n")
	if err := kubectlWait(ctx, kubeconfig, plantonPlatformCrd, v.Name, v.Namespace,
		"jsonpath={.status.phase}=Ready", 30*time.Minute); err != nil {
		// The operator explains itself: surface the per-component status
		// alongside the timeout so the lane's failure names the culprit.
		out, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", plantonPlatformCrd, v.Name, "-n", v.Namespace,
			"-o", "jsonpath={.status.components}").CombinedOutput()
		return errors.Wrapf(err, "the platform never reached Ready; component status: %s", firstLines(string(out), 6))
	}

	// The first-visit handles the module exports — a Ready platform must
	// actually serve them.
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-gateway", v.Namespace); err != nil {
		return errors.Wrap(err, "the front-door gateway Service is missing on a Ready platform")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "secret", v.Name+"-identity-setup-code", v.Namespace); err != nil {
		return errors.Wrap(err, "the first-run setup-code Secret is missing on a Ready platform")
	}
	fmt.Printf("  [verify] platform Ready — gateway Service and setup-code Secret present\n")
	return nil
}

func (v *PlantonPlatformVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The CR's own deletion returns quickly (no operator finalizer);
	// the owner-referenced children drain through garbage collection
	// asynchronously — poll both with a bounded budget.
	deadline := time.Now().Add(5 * time.Minute)
	for {
		crErr := KubectlResourceAbsent(ctx, kubeconfig, plantonPlatformCrd, v.Name, v.Namespace)
		gwErr := KubectlResourceAbsent(ctx, kubeconfig, "service", v.Name+"-gateway", v.Namespace)
		cpErr := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name+"-control-plane", v.Namespace)
		if crErr == nil && gwErr == nil && cpErr == nil {
			fmt.Printf("  [verify] platform gone — CR deleted, children garbage-collected\n")
			return nil
		}
		if time.Now().After(deadline) {
			for _, err := range []error{crErr, gwErr, cpErr} {
				if err != nil {
					return errors.Wrap(err, "platform teardown did not complete within the garbage-collection budget")
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}
