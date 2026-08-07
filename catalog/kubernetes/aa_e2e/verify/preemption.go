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

// PriorityClassPreemptionVerifier proves preemption under REAL scheduling
// pressure: the scenario's fixture saturates the batch's system nodes with
// low-priority pad pods, and a verifier-owned pod at the class-under-test's
// priority must (a) schedule anyway — impossible without evicting a pad —
// and (b) leave a scheduler Preempted event on a victim. The driver pod is
// verifier-owned because it must reference the PriorityClass the component
// under test creates.
type PriorityClassPreemptionVerifier struct {
	// Name of the PriorityClass under test (the driver's priorityClassName).
	Name string
}

// preemption* mirror the pressure fixture (fixture-pressure.yaml) — fixed so
// the fixture and the proof cannot drift apart.
const (
	preemptionNamespace = "e2e-pc-pressure"
	preemptionPadsName  = "e2e-pc-pads"
	preemptionDriverPod = "e2e-pc-preemptor-pod"
)

func (v *PriorityClassPreemptionVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] PriorityClass %q must preempt under real pressure\n", v.Name)

	if err := KubectlResourceExists(ctx, kubeconfig, "priorityclass", v.Name, ""); err != nil {
		return err
	}

	// Pressure must be real before the proof: every pad replica that CAN
	// run is running (some may stay Pending if the nodes filled early —
	// they are fodder either way; what matters is zero remaining headroom
	// for the driver's request).
	if err := v.waitForSaturation(ctx, kubeconfig, 3*time.Minute); err != nil {
		return err
	}

	driver := fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  priorityClassName: %s
  nodeSelector:
    planton.dev/e2e-node-role: system
  containers:
    - name: preemptor
      image: registry.k8s.io/pause:3.10
      resources:
        requests:
          cpu: 500m
          memory: 64Mi
`, preemptionDriverPod, preemptionNamespace, v.Name)
	driverFile, err := writeTempManifest(driver)
	if err != nil {
		return err
	}
	defer os.Remove(driverFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", driverFile); err != nil {
		return errors.Wrap(err, "failed to apply preemptor pod")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "pod", preemptionDriverPod,
			"-n", preemptionNamespace, "--ignore-not-found", "--wait=false")
	}()

	// Scheduling requires an eviction first (pad termination grace is the
	// default 30s), so allow a few minutes.
	if err := kubectlWait(ctx, kubeconfig, "pod", preemptionDriverPod, preemptionNamespace,
		"condition=Ready", 4*time.Minute); err != nil {
		return errors.Wrap(err, "preemptor pod never scheduled — preemption did not happen")
	}

	// The victim's eviction is the scheduler's own testimony: a Preempted
	// event in the pressure namespace.
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "events", "-n", preemptionNamespace,
		"--field-selector", "reason=Preempted",
		"-o", "jsonpath={.items[*].involvedObject.name}").CombinedOutput()
	victims := strings.TrimSpace(string(out))
	if err != nil || victims == "" {
		return errors.Errorf("preemptor scheduled but no Preempted event found (out=%q err=%v) — capacity may not have been saturated", victims, err)
	}
	fmt.Printf("  [verify] preemptor scheduled and scheduler recorded Preempted victim(s): %s\n", victims)
	return nil
}

func (v *PriorityClassPreemptionVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "priorityclass", v.Name, "")
}

// waitForSaturation waits until the pad Deployment stops making scheduling
// progress with the nodes full: at least one pad is running (pressure is
// real) and the driver-sized headroom is gone, approximated by the pads'
// own Pending tail or full availability.
func (v *PriorityClassPreemptionVerifier) waitForSaturation(ctx context.Context, kubeconfig string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "deployment", preemptionPadsName, "-n", preemptionNamespace,
			"-o", "jsonpath={.status.readyReplicas}/{.spec.replicas}").CombinedOutput()
		state := strings.TrimSpace(string(out))
		if err == nil {
			parts := strings.Split(state, "/")
			if len(parts) == 2 && parts[0] != "" && parts[0] != "0" {
				// Ready pads exist; give the scheduler a settling beat so
				// remaining pads reach their terminal Pending state, then
				// treat the cluster as saturated.
				time.Sleep(15 * time.Second)
				return nil
			}
		}
		last = fmt.Sprintf("readiness=%q err=%v", state, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("pressure fixture never saturated the nodes (last: %s)", last)
}

func (v *PriorityClassPreemptionVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}
