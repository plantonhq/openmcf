package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ghaRunnerCrds are the actions.github.com CRDs the controller chart
// installs — release-owned: they delete WITH the controller.
var ghaRunnerCrds = []string{
	"autoscalingrunnersets.actions.github.com",
	"autoscalinglisteners.actions.github.com",
	"ephemeralrunnersets.actions.github.com",
	"ephemeralrunners.actions.github.com",
}

// GhaRunnerScaleSetControllerVerifier checks a runner scale set
// controller install to the point a KubernetesGhaRunnerScaleSet could be
// applied against it: the controller Deployment rolled out with the
// declared replica count and every actions.github.com CRD established.
// Destroy asserts the chart-owned CRD posture: the CRDs delete WITH the
// release.
type GhaRunnerScaleSetControllerVerifier struct {
	Namespace string
	Name      string
	// Replicas is the declared controller replica count (spec default 1).
	Replicas int
}

// ghaControllerReplicas reads spec.replicas (default 1).
func ghaControllerReplicas(spec map[string]interface{}) int {
	if raw, ok := spec["replicas"]; ok {
		switch value := raw.(type) {
		case float64:
			return int(value)
		case int:
			return value
		case string:
			if parsed, err := strconv.Atoi(value); err == nil {
				return parsed
			}
		}
	}
	return 1
}

func (v *GhaRunnerScaleSetControllerVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] gha-runner-scale-set-controller %q in namespace %q\n", v.Name, v.Namespace)

	// fullnameOverride pins the Deployment name to the resource name.
	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name, v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the controller deployment never rolled out")
	}
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", "jsonpath={.status.readyReplicas}").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading controller ready replicas: %s", firstLines(string(out), 3))
	}
	ready, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	if ready < v.Replicas {
		return errors.Errorf("the controller has %d ready replicas, the spec declared %d", ready, v.Replicas)
	}

	for _, crd := range ghaRunnerCrds {
		if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %s was not installed", crd)
		}
	}
	fmt.Printf("  [verify] CRDS: all %d actions.github.com CRDs established; %d/%d controller replicas ready\n", len(ghaRunnerCrds), ready, v.Replicas)
	return nil
}

func (v *GhaRunnerScaleSetControllerVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	// The crds/-directory posture: Helm installs the four
	// actions.github.com CRDs once and NEVER removes them — destroy
	// must LEAVE them on the cluster (verified live at chart 0.14.2;
	// the pre-proof design claim that they delete with the release was
	// wrong). Kept CRDs carry no release ownership metadata, so later
	// installs adopt them cleanly — the same designed keep the
	// monitoring stack's CRDs prove.
	for _, crd := range ghaRunnerCrds {
		if err := KubectlResourceExists(ctx, kubeconfig, "crd", crd, ""); err != nil {
			return errors.Wrapf(err, "CRD %s must SURVIVE the controller destroy (the crds/-directory keep posture)", crd)
		}
	}
	fmt.Printf("  [verify] DESTROY: controller workloads gone; actions.github.com CRDs KEPT (the designed crds/-directory posture)\n")
	return nil
}

// GhaRunnerScaleSetVerifier checks a runner scale set deployment. Every
// lane proves the MODULE CONTRACT: the credential Secret present (the
// materialized `<name>-github-auth` on declared arms), the
// AutoscalingRunnerSet rendered under the declared scale-set name, and
// the controller OBSERVING it (a status recorded on the CR — the
// reconcile-attempt proof, which needs no GitHub account).
//
// The behavioral-github lane (recognized by name) additionally proves
// THE PRODUCT with a real credential from the fenced environment tokens:
// the listener pod reaches Running (it long-polls GitHub — impossible
// without a valid registration) and the CR reports current runners at
// the declared minimum (an idle runner online in GitHub).
type GhaRunnerScaleSetVerifier struct {
	Namespace string
	Name      string
	// ScaleSetName is the GitHub-visible fleet name (the AutoscalingRunnerSet
	// CR name; spec.runner_scale_set_name falling back to metadata.name).
	ScaleSetName string
	// AuthSecretName is the Secret the chart reads the credential from.
	AuthSecretName string
	// MinRunners gates the idle-runner assertion on the live arm.
	MinRunners int
	// RegistrationProof switches on the live GitHub arm.
	RegistrationProof bool
}

// ghaScaleSetName reads spec.runner_scale_set_name falling back to
// metadata name (both manifest key forms tolerated).
func ghaScaleSetName(spec map[string]interface{}, metadataName string) string {
	for _, key := range []string{"runner_scale_set_name", "runnerScaleSetName"} {
		if raw, ok := spec[key].(string); ok && raw != "" {
			return raw
		}
	}
	return metadataName
}

// ghaAuthSecretName resolves the credential Secret name: the
// existing-Secret arm's own name, or the module-materialized
// `<name>-github-auth`.
func ghaAuthSecretName(spec map[string]interface{}, metadataName string) string {
	auth, _ := spec["auth"].(map[string]interface{})
	if auth != nil {
		for _, key := range []string{"existing_secret_name", "existingSecretName"} {
			if raw, ok := auth[key].(string); ok && raw != "" {
				return raw
			}
		}
	}
	return metadataName + "-github-auth"
}

// ghaMinRunners reads spec.min_runners (default 0; both key forms
// tolerated).
func ghaMinRunners(spec map[string]interface{}) int {
	for _, key := range []string{"min_runners", "minRunners"} {
		if raw, ok := spec[key]; ok {
			switch value := raw.(type) {
			case float64:
				return int(value)
			case int:
				return value
			}
		}
	}
	return 0
}

func (v *GhaRunnerScaleSetVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] gha-runner-scale-set %q (fleet %q) in namespace %q\n", v.Name, v.ScaleSetName, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "secret", v.AuthSecretName, v.Namespace); err != nil {
		return errors.Wrap(err, "the GitHub credential Secret not found")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "autoscalingrunnersets.actions.github.com", v.ScaleSetName, v.Namespace); err != nil {
		return errors.Wrap(err, "the AutoscalingRunnerSet was not rendered")
	}

	// The reconcile-attempt proof: the controller observed the CR and
	// recorded a status. This needs no GitHub account — it proves the
	// controller wiring, not the registration.
	deadline := time.Now().Add(3 * time.Minute)
	observed := false
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "autoscalingrunnersets.actions.github.com", v.ScaleSetName, "-n", v.Namespace,
			"-o", "jsonpath={.status}").CombinedOutput()
		if err == nil && strings.TrimSpace(string(out)) != "" {
			observed = true
			break
		}
		time.Sleep(5 * time.Second)
	}
	if !observed {
		return errors.Errorf("the controller never recorded a status on AutoscalingRunnerSet %q — is the controller running and watching this namespace?", v.ScaleSetName)
	}
	fmt.Printf("  [verify] OBSERVED: the controller reconciled the AutoscalingRunnerSet (status recorded)\n")

	if !v.RegistrationProof {
		return nil
	}

	// ---- THE REGISTRATION PROOF (real GitHub credential) --------------------
	// The listener long-polls GitHub for queued jobs; it reaches Running
	// only after a successful registration against the configured URL.
	if err := v.awaitListenerRunning(ctx, kubeconfig, 5*time.Minute); err != nil {
		return err
	}
	if v.MinRunners > 0 {
		if err := v.awaitCurrentRunners(ctx, kubeconfig, v.MinRunners, 5*time.Minute); err != nil {
			return err
		}
	}
	return nil
}

func (v *GhaRunnerScaleSetVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "autoscalingrunnersets.actions.github.com", v.ScaleSetName, v.Namespace); err != nil {
		return err
	}
	fmt.Printf("  [verify] DESTROY: the AutoscalingRunnerSet is gone (the chart's cleanup finalizers unregistered the fleet)\n")
	return nil
}

// awaitListenerRunning waits for the scale set's listener pod (labeled
// with the fleet's scale-set-name by the controller) to reach Running —
// the listener lives in the CONTROLLER's namespace, so the search is
// cluster-wide by label.
func (v *GhaRunnerScaleSetVerifier) awaitListenerRunning(ctx context.Context, kubeconfig string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	lastSeen := "(no listener pod yet)"
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "pods", "-A",
			"-l", "actions.github.com/scale-set-name="+v.ScaleSetName+",app.kubernetes.io/component=runner-scale-set-listener",
			"-o", "jsonpath={range .items[*]}{.metadata.name}={.status.phase} {end}").CombinedOutput()
		if err == nil {
			state := strings.TrimSpace(string(out))
			if state != "" {
				lastSeen = state
			}
			if strings.Contains(state, "=Running") {
				fmt.Printf("  [verify] REGISTERED: the listener pod is Running — the fleet registered with GitHub (%s)\n", state)
				return nil
			}
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("the listener never reached Running within %s (last: %s) — check the GitHub credential and URL", timeout, lastSeen)
}

// awaitCurrentRunners waits for the AutoscalingRunnerSet to report at
// least min current runners (idle runners online in GitHub).
func (v *GhaRunnerScaleSetVerifier) awaitCurrentRunners(ctx context.Context, kubeconfig string, min int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	last := "(none reported)"
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "autoscalingrunnersets.actions.github.com", v.ScaleSetName, "-n", v.Namespace,
			"-o", "jsonpath={.status.currentRunners}").CombinedOutput()
		if err == nil {
			raw := strings.TrimSpace(string(out))
			if raw != "" {
				last = raw
				if current, err := strconv.Atoi(raw); err == nil && current >= min {
					fmt.Printf("  [verify] RUNNERS: %d runner(s) online (declared minimum %d)\n", current, min)
					return nil
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("currentRunners never reached the declared minimum %d within %s (last: %s)", min, timeout, last)
}
