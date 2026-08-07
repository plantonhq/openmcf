package verify

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// OtelCollectorVerifier checks an operator-managed OpenTelemetry
// Collector to the product bar: the mode's workload rolled out, the
// receiver-derived Service present — and THE PIPELINE PROOF on every
// lane: a real OTLP/HTTP trace push whose marker lands in the debug
// exporter's output (telemetry THROUGH the declared pipeline, not just
// a running pod). The daemonset lane additionally proves THE FILELOG
// PROOF (node log files actually ingested — the pattern needs the
// spec's run-as-root security context on containerd nodes). The
// behavioral-pipeline lane adds THE RECONCILE PROOF: the verifier
// patches the CR's config with a resource-attribute marker and asserts
// the operator rolls the collector onto the NEW pipeline (a fresh push
// carries the new attribute) — declared config change becoming live
// behavior is this kind's whole contract.
//
// Every push runs through a FRESH port-forward (the dead-tunnel class:
// a tunnel opened before a pod replacement silently blackholes).
type OtelCollectorVerifier struct {
	// Namespace is the deploy namespace from the spec.
	Namespace string
	// Name is metadata.name — the operator derives every child name
	// from it ("<name>-collector" for the workload and Service).
	Name string
	// Daemonset marks the daemonset-mode lane (rollout target + the
	// filelog proof).
	Daemonset bool
	// Behavioral marks the behavioral-pipeline lane (the reconcile
	// proof).
	Behavioral bool
}

func (v *OtelCollectorVerifier) workload() string {
	if v.Daemonset {
		return "daemonset/" + v.Name + "-collector"
	}
	return "deployment/" + v.Name + "-collector"
}

func (v *OtelCollectorVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] otel collector %q (%s) in namespace %q\n", v.Name, v.workload(), v.Namespace)

	if err := kubectlRolloutStatus(ctx, kubeconfig, v.workload(), v.Namespace, 5*time.Minute); err != nil {
		return errors.Wrap(err, "the collector workload never rolled out")
	}
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name+"-collector", v.Namespace); err != nil {
		return errors.Wrap(err, "the collector service not found")
	}

	marker := fmt.Sprintf("e2e-%s-%d", v.Name, time.Now().Unix())
	if err := v.pushAndAssert(ctx, kubeconfig, marker, 3*time.Minute); err != nil {
		return errors.Wrap(err, "THE PIPELINE PROOF failed")
	}
	fmt.Printf("  [verify] PIPELINE PROOF: OTLP push %q observed in the debug exporter output\n", marker)

	if v.Daemonset {
		if err := v.filelogProof(ctx, kubeconfig); err != nil {
			return errors.Wrap(err, "THE FILELOG PROOF failed")
		}
	}

	if v.Behavioral {
		if err := v.reconcileProof(ctx, kubeconfig); err != nil {
			return errors.Wrap(err, "THE RECONCILE PROOF failed")
		}
	}
	return nil
}

// pushAndAssert opens a FRESH port-forward, POSTs one OTLP/JSON trace
// carrying the marker as service.name, and polls the workload logs until
// the debug exporter prints it.
func (v *OtelCollectorVerifier) pushAndAssert(ctx context.Context, kubeconfig, marker string, budget time.Duration) error {
	deadline := time.Now().Add(budget)
	var lastErr error
	for time.Now().Before(deadline) {
		if err := v.pushOnce(ctx, kubeconfig, marker); err != nil {
			lastErr = err
			time.Sleep(10 * time.Second)
			continue
		}
		// Poll the logs for the debug exporter's record of the marker.
		logDeadline := time.Now().Add(45 * time.Second)
		for time.Now().Before(logDeadline) {
			out, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
				"logs", v.workload(), "-n", v.Namespace, "--since=5m").CombinedOutput()
			if strings.Contains(string(out), marker) {
				return nil
			}
			time.Sleep(5 * time.Second)
		}
		lastErr = errors.Errorf("the pushed marker %q never appeared in the collector logs", marker)
	}
	return lastErr
}

// pushOnce POSTs a single OTLP/JSON trace through a fresh tunnel.
func (v *OtelCollectorVerifier) pushOnce(ctx context.Context, kubeconfig, marker string) error {
	const localPort = "43180"
	cancel, err := startPortForward(ctx, kubeconfig, "svc/"+v.Name+"-collector", v.Namespace, localPort+":4318")
	if err != nil {
		return errors.Wrap(err, "starting the OTLP port-forward")
	}
	defer cancel()

	now := time.Now().UnixNano()
	payload := fmt.Sprintf(`{
  "resourceSpans": [{
    "resource": {"attributes": [{"key": "service.name", "value": {"stringValue": "%s"}}]},
    "scopeSpans": [{"spans": [{
      "traceId": "5b8efff798038103d269b633813fc60c",
      "spanId": "eee19b7ec3c1b174",
      "name": "e2e-pipeline-probe",
      "kind": 1,
      "startTimeUnixNano": "%d",
      "endTimeUnixNano": "%d"
    }]}]
  }]
}`, marker, now, now+1000000)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"http://127.0.0.1:"+localPort+"/v1/traces", bytes.NewBufferString(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "posting the OTLP trace")
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("OTLP push answered HTTP %d", resp.StatusCode)
	}
	return nil
}

// filelogProof asserts the daemonset lane's node-log ingestion: the
// debug exporter must have printed LogRecords sourced from log files
// (the filelog receiver stamps log.file.path). This is the proof the
// run-as-root security context exists for — without it the receiver
// reports permission errors and ships nothing.
func (v *OtelCollectorVerifier) filelogProof(ctx context.Context, kubeconfig string) error {
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		out, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"logs", v.workload(), "-n", v.Namespace, "--since=5m").CombinedOutput()
		if strings.Contains(string(out), "log.file.path") {
			fmt.Printf("  [verify] FILELOG PROOF: node log files ingested (log.file.path records in the debug output)\n")
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return errors.New("no filelog records appeared in the debug output — check the run-as-root security context and the hostPath mounts")
}

// reconcileProof patches the CR's config with a resource-attribute
// marker, waits for the operator to roll the collector, and asserts a
// fresh push carries the new attribute — the declared-config-change →
// live-pipeline-change loop.
func (v *OtelCollectorVerifier) reconcileProof(ctx context.Context, kubeconfig string) error {
	marker := fmt.Sprintf("e2e-reconcile-%d", time.Now().Unix())
	patch := fmt.Sprintf(`{"spec":{"config":{"processors":{"resource":{"attributes":[{"key":"planton.e2e.reconcile","value":"%s","action":"insert"}]}},"service":{"pipelines":{"traces":{"receivers":["otlp"],"processors":["resource"],"exporters":["debug"]}}}}}}`, marker)

	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"patch", "opentelemetrycollectors.opentelemetry.io", v.Name, "-n", v.Namespace,
		"--type=merge", "-p", patch).CombinedOutput(); err != nil {
		return errors.Wrapf(err, "patching the collector config: %s", firstLines(string(out), 3))
	}
	fmt.Printf("  [verify] RECONCILE: config patched with resource attribute %q — awaiting the operator's rollout\n", marker)

	// The operator rolls the Deployment onto the new ConfigMap; the
	// rollout wait absorbs the transition.
	time.Sleep(10 * time.Second)
	if err := kubectlRolloutStatus(ctx, kubeconfig, v.workload(), v.Namespace, 4*time.Minute); err != nil {
		return errors.Wrap(err, "the collector never rolled onto the patched config")
	}

	// The new attribute must appear on freshly pushed telemetry (each
	// push attempt opens a fresh tunnel — the old pod is gone).
	if err := v.pushAndAssert(ctx, kubeconfig, marker+"-svc", 4*time.Minute); err != nil {
		return err
	}
	// The debug output must ALSO carry the injected resource attribute.
	out, _ := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"logs", v.workload(), "-n", v.Namespace, "--since=5m").CombinedOutput()
	if !strings.Contains(string(out), marker) {
		return errors.Errorf("the patched resource attribute %q never appeared on pushed telemetry", marker)
	}
	fmt.Printf("  [verify] RECONCILE PROOF: pushed telemetry carries the patched resource attribute — declared config change became live pipeline behavior\n")
	return nil
}

func (v *OtelCollectorVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	kind := strings.Split(v.workload(), "/")[0]
	if err := KubectlResourceAbsent(ctx, kubeconfig, kind, v.Name+"-collector", v.Namespace); err != nil {
		return err
	}
	if err := KubectlResourceAbsent(ctx, kubeconfig, "opentelemetrycollectors.opentelemetry.io", v.Name, v.Namespace); err != nil {
		return err
	}
	fmt.Printf("  [verify] DESTROY: the collector CR and its workload are gone\n")
	return nil
}
