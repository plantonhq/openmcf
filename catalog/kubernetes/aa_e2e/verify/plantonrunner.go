package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// PlantonRunnerVerifier checks a Planton runner install to the module's
// contract: the module-created token Secret present, the chart's
// Deployment rendered pinned to exactly ONE replica with the Recreate
// strategy (two live pods under one runner name would revoke each other's
// keys), and the pod template carrying the enrollment env wiring — the
// token via a secretKeyRef (never inline) and the registration name.
//
// DELIBERATELY no pod-readiness assertion: the module installs with no
// Helm wait because the runner's readiness contract is its control-plane
// work queue — the kind lanes run with an obviously-fake token whose join
// is refused, and the pod restarting is the DESIGNED behavior there, not
// a failure. Destroy asserts the workload and the token Secret are gone.
type PlantonRunnerVerifier struct {
	Namespace string
	Name      string
}

func (v *PlantonRunnerVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] planton-runner %q in namespace %q\n", v.Name, v.Namespace)

	// The module-created token Secret (`<name>-token`) — the chart reads
	// it via its existingSecret form; the token never rides rendered
	// chart values.
	tokenSecret := v.Name + "-token"
	if err := KubectlResourceExists(ctx, kubeconfig, "secret", tokenSecret, v.Namespace); err != nil {
		return errors.Wrap(err, "the token Secret not found")
	}

	// fullnameOverride pins the Deployment name to the resource name.
	if err := KubectlResourceExists(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "the runner Deployment was not rendered")
	}

	// The singleton law and the identity-safe rollout strategy, straight
	// off the live object.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", "jsonpath={.spec.replicas}/{.spec.strategy.type}").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading deployment replicas/strategy: %s", firstLines(string(out), 3))
	}
	if got := strings.TrimSpace(string(out)); got != "1/Recreate" {
		return errors.Errorf("the runner Deployment must run exactly one replica with the Recreate strategy (a rolling surge would revoke the live pod's key); got %q", got)
	}

	// The enrollment env contract on the pod template: the token arrives
	// through a secretKeyRef naming the module's Secret, and the runner
	// registers itself under the declared name.
	out, err = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", `jsonpath={.spec.template.spec.containers[0].env[?(@.name=="PLANTON_RUNNER_TOKEN")].valueFrom.secretKeyRef.name}`).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading the token env wiring: %s", firstLines(string(out), 3))
	}
	if got := strings.TrimSpace(string(out)); got != tokenSecret {
		return errors.Errorf("PLANTON_RUNNER_TOKEN must ride a secretKeyRef naming %q (never inline values); got %q", tokenSecret, got)
	}

	out, err = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "deployment", v.Name, "-n", v.Namespace,
		"-o", `jsonpath={.spec.template.spec.containers[0].env[?(@.name=="PLANTON_RUNNER_NAME")].value}`).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "reading the runner-name env: %s", firstLines(string(out), 3))
	}
	if got := strings.TrimSpace(string(out)); got == "" {
		return errors.New("PLANTON_RUNNER_NAME is missing from the pod template -- the runner would register under the chart default instead of the declared name")
	}

	fmt.Printf("  [verify] CONTRACT: token Secret %q wired via secretKeyRef; 1 replica, Recreate strategy\n", tokenSecret)
	return nil
}

func (v *PlantonRunnerVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace); err != nil {
		return err
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "secret", v.Name+"-token", v.Namespace); err != nil {
		return err
	}
	fmt.Printf("  [verify] DESTROY: the runner Deployment and its token Secret are gone\n")
	return nil
}
