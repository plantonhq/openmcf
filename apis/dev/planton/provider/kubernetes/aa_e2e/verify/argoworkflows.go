package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ArgoWorkflowsVerifier checks an Argo Workflows install to the point a
// customer could run their pipelines on it: the workflow controller
// rolled out, the Argo server rolled out (when enabled) with an
// AUTHENTICATED version-API round-trip (a TokenRequest token for the
// runner ServiceAccount — the default `client` auth mode forwards the
// caller's Kubernetes identity), the runner ServiceAccount present, and
// THE ENGINE PROOF on every lane: a verifier-owned Workflow submitted as
// a CR runs to Succeeded under the runner identity (a workflow engine
// that cannot run a workflow is not a workflow engine).
//
// The behavioral-resilience scenario (recognized by name) additionally
// DELETES the controller pod, waits for a REPLACEMENT (a new UID), and
// submits a SECOND workflow that must also succeed — the engine recovers
// its reconciliation loop from cluster state alone.
//
// Destroy asserts the designed CRD keep posture: the release's workloads
// are gone while workflows.argoproj.io REMAINS — deleting it would
// cascade to every Workflow in the cluster.
type ArgoWorkflowsVerifier struct {
	Namespace string
	Name      string
	// ServerEnabled gates the server rollout + authenticated API proof.
	ServerEnabled bool
	// WorkflowServiceAccount is the runner identity workflows execute
	// under (the module always creates it).
	WorkflowServiceAccount string
	// InstanceId labels the proof Workflows when the controller claims an
	// instance ID (it reconciles ONLY matching workflows).
	InstanceId string
	// ResilienceProof switches on the controller-pod-replacement arm.
	ResilienceProof bool
}

const argoWorkflowsCrd = "workflows.argoproj.io"

// argoWorkflowsServerEnabled reads spec.server.enabled (default true —
// the proto default; both manifest key forms tolerated).
func argoWorkflowsServerEnabled(spec map[string]interface{}) bool {
	server, _ := spec["server"].(map[string]interface{})
	if server == nil {
		return true
	}
	if raw, ok := server["enabled"]; ok {
		if enabled, ok := raw.(bool); ok {
			return enabled
		}
	}
	return true
}

// argoWorkflowsRunnerSA reads spec.workflow_service_account (default
// "argo-workflow" — the chart's own runner name).
func argoWorkflowsRunnerSA(spec map[string]interface{}) string {
	for _, key := range []string{"workflow_service_account", "workflowServiceAccount"} {
		if raw, ok := spec[key].(string); ok && raw != "" {
			return raw
		}
	}
	return "argo-workflow"
}

// argoWorkflowsInstanceId reads spec.controller.instance_id ("" when the
// controller claims no instance ID).
func argoWorkflowsInstanceId(spec map[string]interface{}) string {
	controller, _ := spec["controller"].(map[string]interface{})
	if controller == nil {
		return ""
	}
	for _, key := range []string{"instance_id", "instanceId"} {
		if raw, ok := controller[key].(string); ok && raw != "" {
			return raw
		}
	}
	return ""
}

func (v *ArgoWorkflowsVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] argo-workflows %q in namespace %q\n", v.Name, v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-workflow-controller", v.Namespace, 10*time.Minute); err != nil {
		return errors.Wrap(err, "the workflow controller deployment never rolled out")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "serviceaccount", v.WorkflowServiceAccount, v.Namespace); err != nil {
		return errors.Wrap(err, "the workflow runner ServiceAccount not found (the module renders workflow.serviceAccount.create=true)")
	}

	if v.ServerEnabled {
		if err := kubectlRolloutStatus(ctx, kubeconfig, "deployment/"+v.Name+"-server", v.Namespace, 10*time.Minute); err != nil {
			return errors.Wrap(err, "the argo server deployment never rolled out")
		}
		if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-server", v.Namespace); err != nil {
			return errors.Wrap(err, "argo server service not found")
		}
		if err := v.proveServerApi(ctx, kubeconfig); err != nil {
			return err
		}
	}

	// THE ENGINE PROOF — on every lane.
	if err := v.proveWorkflowRun(ctx, kubeconfig, "e2e-proof-run"); err != nil {
		return err
	}

	if !v.ResilienceProof {
		return nil
	}

	// ---- the resilience proof: controller pod replacement ------------------
	if err := deletePodAwaitReplacement(ctx, kubeconfig, v.Namespace,
		"app.kubernetes.io/instance="+v.Name+",app.kubernetes.io/component=workflow-controller", 10*time.Minute); err != nil {
		return errors.Wrap(err, "the workflow controller pod did not recover after deletion")
	}
	if err := v.proveWorkflowRun(ctx, kubeconfig, "e2e-proof-run-recovered"); err != nil {
		return errors.Wrap(err, "a workflow submitted AFTER the controller pod replacement should still run to completion")
	}
	fmt.Printf("  [verify] RESILIENCE: a fresh workflow ran to Succeeded AFTER controller pod replacement\n")
	return nil
}

func (v *ArgoWorkflowsVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name+"-workflow-controller", v.Namespace); err != nil {
		return err
	}
	// The designed keep posture: destroying the release must LEAVE the
	// Workflow CRD behind (crds.keep defaults true).
	if err := KubectlResourceExists(ctx, kubeconfig, "crd", argoWorkflowsCrd, ""); err != nil {
		return errors.Wrap(err, "the argoproj CRDs must SURVIVE destroy (the crds.keep posture) but workflows.argoproj.io is gone")
	}
	fmt.Printf("  [verify] DESTROY: release workloads gone; workflows.argoproj.io KEPT (the designed crds.keep posture)\n")
	return nil
}

// proveServerApi mints a TokenRequest token for the runner ServiceAccount
// and drives an authenticated GET /api/v1/version through the server (the
// default `client` auth mode acts with the presented Kubernetes identity;
// route verified in the app source at v4.0.8).
func (v *ArgoWorkflowsVerifier) proveServerApi(ctx context.Context, kubeconfig string) error {
	tokenCmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"create", "token", v.WorkflowServiceAccount, "-n", v.Namespace)
	rawToken, err := tokenCmd.Output()
	if err != nil {
		return errors.Wrap(err, "minting a TokenRequest token for the runner ServiceAccount")
	}
	token := strings.TrimSpace(string(rawToken))

	const apiPort = "12746"
	cancel, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name+"-server", v.Namespace, apiPort+":2746")
	if err != nil {
		return errors.Wrap(err, "starting port-forward to the argo server")
	}
	defer cancel()

	// The server speaks plain HTTP at the chart default (secure=false);
	// the argo API expects the token with the "Bearer v2" style plain
	// bearer header.
	body, err := httpBearerRoundTrip(ctx, "GET", "http://127.0.0.1:"+apiPort+"/api/v1/version", token, 3*time.Minute)
	if err != nil {
		return errors.Wrap(err, "the authenticated version-API round-trip failed")
	}
	if !strings.Contains(body, "version") {
		return errors.Errorf("the version API answered without a version: %s", firstLines(body, 2))
	}
	fmt.Printf("  [verify] API: authenticated /api/v1/version round-trip OK as the runner identity\n")
	return nil
}

// proveWorkflowRun submits a verifier-owned Workflow CR and waits for
// phase Succeeded. The CR carries the controller's instance-ID label when
// one is claimed (an instanced controller reconciles ONLY matching
// workflows — the label key is the app's own contract).
func (v *ArgoWorkflowsVerifier) proveWorkflowRun(ctx context.Context, kubeconfig, runName string) error {
	instanceLabel := ""
	if v.InstanceId != "" {
		instanceLabel = fmt.Sprintf("\n  labels:\n    workflows.argoproj.io/controller-instanceid: %s", v.InstanceId)
	}
	manifest := fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  name: %s
  namespace: %s%s
spec:
  entrypoint: main
  serviceAccountName: %s
  templates:
    - name: main
      container:
        image: busybox:1.37
        command: ["sh", "-c", "echo e2e-proof-ok"]
`, runName, v.Namespace, instanceLabel, v.WorkflowServiceAccount)

	path, err := writeVerifierTempManifest("argo-workflow-proof", manifest)
	if err != nil {
		return errors.Wrap(err, "writing the proof Workflow manifest")
	}
	defer func() { _ = os.Remove(path) }()

	if err := kubectlApplyFile(ctx, kubeconfig, path); err != nil {
		return errors.Wrap(err, "submitting the proof Workflow")
	}
	fmt.Printf("  [verify] SUBMIT: Workflow %q submitted under the runner identity\n", runName)

	// Zero-orphan duty: the Workflow CR (and its pods, which it owns)
	// leave with the verifier.
	defer func() {
		delCmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"delete", "workflow", runName, "-n", v.Namespace, "--ignore-not-found", "--wait=false")
		_, _ = delCmd.CombinedOutput()
	}()

	deadline := time.Now().Add(6 * time.Minute)
	lastPhase := ""
	for time.Now().Before(deadline) {
		phase, err := kubectlGetJSONPath(ctx, kubeconfig, "workflow", runName, v.Namespace, "{.status.phase}")
		if err == nil {
			phase = strings.TrimSpace(phase)
			if phase != "" {
				lastPhase = phase
			}
			switch phase {
			case "Succeeded":
				fmt.Printf("  [verify] RUN: Workflow %q ran to Succeeded — the engine executes pipelines\n", runName)
				return nil
			case "Failed", "Error":
				return errors.Errorf("the proof Workflow %q ended %s", runName, phase)
			}
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("the proof Workflow %q never completed (last phase %q)", runName, lastPhase)
}
