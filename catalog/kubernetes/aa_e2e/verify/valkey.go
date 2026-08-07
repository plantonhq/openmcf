package verify

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ValkeyVerifier checks a Valkey deployment (the official chart, release
// and fullname pinned to metadata.name) to the point it is actually
// serving: the workload ready (a Deployment standalone, a StatefulSet in
// replication mode) and the write Service present.
//
// When PersistenceProof is set (the behavioral-persistence scenario), the
// verifier proves DURABILITY THROUGH A POD LOSS: it writes a marker key,
// DELETES the pod outright, waits for the replacement to come up, and
// reads the marker back — real only because the scenario declares a
// persistent volume and append-only persistence; a pure in-memory cache
// would come back empty.
//
// When ReplicationProof is set (the replication scenario), the verifier
// writes through the write Service and reads the key back through the
// READ Service (`<name>-read`), which load-balances across all pods —
// proving replication actually propagates data to replicas.
//
// Commands ride `kubectl exec` + valkey-cli on the instance pods,
// authenticated with the module-materialized `<name>-auth` Secret's
// "default" user when the scenario declares auth.
type ValkeyVerifier struct {
	Namespace string
	Name      string
	// Replication mirrors the manifest's replication block presence —
	// the chart renders a StatefulSet with (replicas + 1) pods there,
	// a single-pod Deployment otherwise.
	Replication bool
	// Replicas is the declared replica count (replication mode).
	Replicas int64
	// AuthDeclared mirrors spec.auth presence — the verifier then
	// authenticates as the "default" ACL user.
	AuthDeclared     bool
	PersistenceProof bool
	ReplicationProof bool
	ServicePort      int64
}

func (v *ValkeyVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] valkey %q in namespace %q (replication=%v)\n", v.Name, v.Namespace, v.Replication)

	if v.Replication {
		total := fmt.Sprintf("%d", v.Replicas+1)
		if err := v.waitForStatefulSetReady(ctx, kubeconfig, total, 5*time.Minute); err != nil {
			return err
		}
	} else {
		if err := kubectlWait(ctx, kubeconfig, "deployment", v.Name, v.Namespace,
			"condition=Available", 5*time.Minute); err != nil {
			return errors.Wrap(err, "valkey deployment not available")
		}
	}

	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.Name, v.Namespace); err != nil {
		return errors.Wrap(err, "write service not found")
	}

	if v.PersistenceProof {
		if err := v.provePersistenceDurability(ctx, kubeconfig); err != nil {
			return err
		}
	}
	if v.ReplicationProof {
		if err := v.proveReplicationPropagation(ctx, kubeconfig); err != nil {
			return err
		}
	}
	return nil
}

func (v *ValkeyVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	if v.Replication {
		return KubectlResourceAbsent(ctx, kubeconfig, "statefulset", v.Name, v.Namespace)
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", v.Name, v.Namespace)
}

// provePersistenceDurability runs the write → pod-loss → restart →
// read-back cycle: append-only persistence on a real volume must carry
// the key across the pod's death.
func (v *ValkeyVerifier) provePersistenceDurability(ctx context.Context, kubeconfig string) error {
	pod, err := v.anyInstancePod(ctx, kubeconfig)
	if err != nil {
		return err
	}

	if _, err := v.valkeyCli(ctx, kubeconfig, pod, "", "SET", "e2e-proof", "valkey-persistence-round-trip"); err != nil {
		return errors.Wrap(err, "failed to write the marker key")
	}
	fmt.Printf("  [verify] marker written on %q — deleting the pod\n", pod)

	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", pod, "-n", v.Namespace, "--wait=true", "--timeout=120s").CombinedOutput(); err != nil {
		return errors.Errorf("failed to delete the valkey pod: %v: %s", err, string(out))
	}

	if v.Replication {
		if err := v.waitForStatefulSetReady(ctx, kubeconfig, fmt.Sprintf("%d", v.Replicas+1), 4*time.Minute); err != nil {
			return err
		}
	} else {
		if err := kubectlWait(ctx, kubeconfig, "deployment", v.Name, v.Namespace,
			"condition=Available", 4*time.Minute); err != nil {
			return errors.Wrap(err, "valkey deployment never recovered from the pod loss")
		}
	}

	newPod, err := v.anyInstancePod(ctx, kubeconfig)
	if err != nil {
		return err
	}
	// The replacement replays the append-only file on startup; give the
	// read a settle window.
	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		out, err := v.valkeyCli(ctx, kubeconfig, newPod, "", "GET", "e2e-proof")
		if err == nil && strings.Contains(out, "valkey-persistence-round-trip") {
			fmt.Printf("  [verify] marker intact on the replacement pod — the volume carried the data\n")
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("marker key missing after the pod loss — persistence did not carry the data")
}

// proveReplicationPropagation writes through the write path and reads the
// key back through the READ Service until a replica serves it.
func (v *ValkeyVerifier) proveReplicationPropagation(ctx context.Context, kubeconfig string) error {
	// Pod 0 of the StatefulSet is the primary in the chart's replication
	// topology — the write lands there.
	primary := v.Name + "-0"
	if _, err := v.valkeyCli(ctx, kubeconfig, primary, "", "SET", "e2e-replication-proof", "valkey-replica-read"); err != nil {
		return errors.Wrap(err, "failed to write through the primary")
	}

	readHost := fmt.Sprintf("%s-read.%s.svc.cluster.local", v.Name, v.Namespace)
	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		out, err := v.valkeyCli(ctx, kubeconfig, primary, readHost, "GET", "e2e-replication-proof")
		if err == nil && strings.Contains(out, "valkey-replica-read") {
			fmt.Printf("  [verify] key served through the read service — replication propagates\n")
			return nil
		}
		time.Sleep(5 * time.Second)
	}
	return errors.New("key never appeared through the read service — replication did not propagate")
}

func (v *ValkeyVerifier) waitForStatefulSetReady(ctx context.Context, kubeconfig, want string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "statefulset", v.Name, "-n", v.Namespace,
			"-o", "jsonpath={.status.readyReplicas}").CombinedOutput()
		got := strings.TrimSpace(string(out))
		if err == nil && got == want {
			return nil
		}
		last = fmt.Sprintf("readyReplicas=%q err=%v", got, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("statefulset %q never reached %s ready replicas (last: %s)", v.Name, want, last)
}

// anyInstancePod resolves one running Valkey instance pod through the
// release's instance label (Deployment pods have generated names).
func (v *ValkeyVerifier) anyInstancePod(ctx context.Context, kubeconfig string) (string, error) {
	deadline := time.Now().Add(2 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "pods", "-n", v.Namespace,
			"-l", "app.kubernetes.io/instance="+v.Name,
			"--field-selector", "status.phase=Running",
			"-o", "jsonpath={.items[0].metadata.name}").CombinedOutput()
		pod := strings.TrimSpace(string(out))
		if err == nil && pod != "" {
			return pod, nil
		}
		last = fmt.Sprintf("out=%q err=%v", pod, err)
		time.Sleep(5 * time.Second)
	}
	return "", errors.Errorf("no running valkey pod found for instance %q (last: %s)", v.Name, last)
}

// valkeyCli runs a valkey-cli command on an instance pod, optionally
// against a remote host (the read Service), authenticating as the
// "default" ACL user when the scenario declared auth.
func (v *ValkeyVerifier) valkeyCli(ctx context.Context, kubeconfig, pod, host string, cmd ...string) (string, error) {
	args := []string{"--kubeconfig", kubeconfig, "exec", pod, "-n", v.Namespace, "--", "valkey-cli"}
	if host != "" {
		args = append(args, "-h", host)
	}
	if v.ServicePort != 0 && host != "" {
		args = append(args, "-p", fmt.Sprintf("%d", v.ServicePort))
	}
	if v.AuthDeclared {
		password, err := v.defaultUserPassword(ctx, kubeconfig)
		if err != nil {
			return "", err
		}
		args = append(args, "--no-auth-warning", "-a", password)
	}
	args = append(args, cmd...)
	out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
	if err != nil {
		return "", errors.Errorf("valkey-cli %v on %s: %v: %s", cmd, pod, err, string(out))
	}
	return string(out), nil
}

// valkeyReplication reads the manifest's replication posture: block
// presence plus the declared replica count (chart default 2).
func valkeyReplication(spec map[string]interface{}) (bool, int64) {
	if spec == nil {
		return false, 0
	}
	raw, ok := spec["replication"]
	if !ok {
		return false, 0
	}
	replicas := int64(2)
	if block, ok := raw.(map[string]interface{}); ok {
		switch v := block["replicas"].(type) {
		case int:
			replicas = int64(v)
		case int64:
			replicas = v
		case float64:
			replicas = int64(v)
		}
	}
	return true, replicas
}

// valkeyServicePort reads the write Service port (chart default 6379).
func valkeyServicePort(spec map[string]interface{}) int64 {
	if spec == nil {
		return 6379
	}
	service, ok := spec["service"].(map[string]interface{})
	if !ok {
		return 6379
	}
	switch v := service["port"].(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case float64:
		return int64(v)
	default:
		return 6379
	}
}

// defaultUserPassword reads the module-materialized `<name>-auth` Secret's
// "default" key.
func (v *ValkeyVerifier) defaultUserPassword(ctx context.Context, kubeconfig string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "secret", v.Name+"-auth", "-n", v.Namespace,
		"-o", "jsonpath={.data.default}").CombinedOutput()
	if err != nil {
		return "", errors.Errorf("failed to read the valkey auth secret: %v: %s", err, string(out))
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(out)))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode the valkey default-user password")
	}
	return string(decoded), nil
}
