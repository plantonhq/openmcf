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

// HpaBehavioralVerifier asserts an HPA actually operates, not just exists —
// used for scenarios that install metrics-server (declared via the
// scenario's e2e-prerequisites annotation):
//
//  1. ScalingActive=True — the controller computes the metric (the metric
//     PIPELINE works end to end: kubelet → metrics-server → metrics API →
//     HPA controller).
//  2. desiredReplicas exceeds min_replicas — with a CPU-burning target
//     Deployment the utilization deterministically exceeds the target, so
//     the controller must scale up.
//
// Without metrics-server, HPA objects deploy fine but ScalingActive stays
// False with FailedGetResourceMetric — exactly the gap this verifier
// closes.
type HpaBehavioralVerifier struct {
	Namespace   string
	Name        string
	MinReplicas int64
}

func (v *HpaBehavioralVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] HPA %q behavioral scaling in namespace %q\n", v.Name, v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "hpa", v.Name, v.Namespace); err != nil {
		return errors.Wrapf(err, "hpa %q not found in namespace %q", v.Name, v.Namespace)
	}

	// Metric flow: the HPA controller reports ScalingActive once it can
	// read the metric and compute a desired count. Generous timeout — the
	// path spans a metrics-server scrape window plus the HPA sync period.
	if err := kubectlWait(ctx, kubeconfig, "hpa", v.Name, v.Namespace,
		"condition=ScalingActive", 3*time.Minute); err != nil {
		return errors.Wrapf(err, "hpa %q never reached ScalingActive — metric values are not flowing", v.Name)
	}

	// Scale-up: the burner target's CPU usage deterministically exceeds
	// the utilization target, so desiredReplicas must rise above min.
	deadline := time.Now().Add(3 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"-n", v.Namespace, "get", "hpa", v.Name,
			"-o", "jsonpath={.status.desiredReplicas}")
		out, err := cmd.CombinedOutput()
		last = strings.TrimSpace(string(out))
		if err == nil && last != "" {
			if desired, convErr := strconv.ParseInt(last, 10, 64); convErr == nil && desired > v.MinReplicas {
				fmt.Printf("  [verify] HPA %q scaled: desiredReplicas=%d > min=%d\n", v.Name, desired, v.MinReplicas)
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("hpa %q never scaled above min replicas (last desiredReplicas=%q) — metrics flow but the scale-up did not happen", v.Name, last)
}

func (v *HpaBehavioralVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "hpa", v.Name, v.Namespace)
}
