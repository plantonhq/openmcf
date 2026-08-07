package verify

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// CnpgOperatorInstallVerifier checks a CloudNativePG operator installation:
// the operator Deployment Available and the core CRDs Established — the
// preconditions every KubernetesPostgres apply depends on. When the Barman
// Cloud plugin is enabled, the plugin Deployment and its ObjectStore CRD
// are part of the contract too (backup blocks are dead letters without
// them).
type CnpgOperatorInstallVerifier struct {
	Namespace string
	// PluginEnabled mirrors the manifest's barman_cloud_plugin.enabled —
	// the plugin ships as a second release, so its health is asserted
	// only when the spec asked for it.
	PluginEnabled bool
}

// cnpgOperatorDeployment is the chart fullname with the module's fixed
// release name "cnpg" (`<release>-<chart>`), one per cluster.
const cnpgOperatorDeployment = "cnpg-cloudnative-pg"

// cnpgPluginDeployment is the plugin chart's fullname with the module's
// fixed release name (release == chart name, so the fullname collapses).
const cnpgPluginDeployment = "plugin-barman-cloud"

func (v *CnpgOperatorInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] cloudnative-pg operator in namespace %q (plugin=%v)\n", v.Namespace, v.PluginEnabled)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for cloudnative-pg", v.Namespace)
	}

	if err := kubectlWait(ctx, kubeconfig, "deployment", cnpgOperatorDeployment, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "cloudnative-pg operator deployment not available")
	}

	// The CRDs KubernetesPostgres renders against: the Cluster itself and
	// the two backup-facing kinds.
	for _, crd := range []string{
		"clusters.postgresql.cnpg.io",
		"scheduledbackups.postgresql.cnpg.io",
		"backups.postgresql.cnpg.io",
	} {
		if err := kubectlWait(ctx, kubeconfig, "crd", crd, "",
			"condition=Established", 2*time.Minute); err != nil {
			return errors.Wrapf(err, "CRD %s not established", crd)
		}
	}

	if !v.PluginEnabled {
		return nil
	}

	if err := kubectlWait(ctx, kubeconfig, "deployment", cnpgPluginDeployment, v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "barman-cloud plugin deployment not available (is cert-manager installed? the plugin's TLS depends on it)")
	}
	if err := kubectlWait(ctx, kubeconfig, "crd", "objectstores.barmancloud.cnpg.io", "",
		"condition=Established", 2*time.Minute); err != nil {
		return errors.Wrap(err, "ObjectStore CRD not established")
	}

	return nil
}

func (v *CnpgOperatorInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The CRDs intentionally SURVIVE uninstall (the chart stamps
	// helm.sh/resource-policy: keep on them so removing the operator never
	// cascade-deletes Cluster resources) — only the deployments' absence
	// is asserted.
	if err := KubectlResourceAbsent(ctx, kubeconfig, "deployment", cnpgOperatorDeployment, v.Namespace); err != nil {
		return err
	}
	if v.PluginEnabled {
		return KubectlResourceAbsent(ctx, kubeconfig, "deployment", cnpgPluginDeployment, v.Namespace)
	}
	return nil
}

// CnpgClusterVerifier checks a CloudNativePG-managed PostgreSQL cluster to
// the point it is actually serving: the Cluster Ready condition, every
// declared instance ready, and the -rw Service present.
//
// When Behavioral is set (the behavioral-failover scenario), the verifier
// additionally proves DATA DURABILITY through a failover: it writes a
// marker row through the current primary, DELETES the primary pod, waits
// for the operator to promote a replica, and reads the marker back through
// the new primary. The write path is `kubectl exec` + psql as the postgres
// OS user (peer auth inside the pod) — no credentials or port-forwards to
// flake on.
type CnpgClusterVerifier struct {
	Namespace   string
	ClusterName string
	// Instances is the declared instance count (read from the scenario
	// manifest) — readiness means ALL of them, not just the primary.
	Instances int64
	// Behavioral switches on the live failover proof (see above).
	Behavioral bool
	// BackupProof waits for the immediate schedule's Backup to reach
	// Completed (the with-backup scenario) — a REAL base backup landing
	// in the object store through the Barman Cloud plugin.
	BackupProof bool
}

func (v *CnpgClusterVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] cloudnative-pg cluster %q in namespace %q (instances=%d)\n",
		v.ClusterName, v.Namespace, v.Instances)

	// Ready flips once the bootstrap completed and the topology matches
	// the spec. First-boot includes an image pull plus initdb, so the
	// window is generous.
	if err := kubectlWait(ctx, kubeconfig, "cluster.postgresql.cnpg.io", v.ClusterName, v.Namespace,
		"condition=Ready", 8*time.Minute); err != nil {
		return errors.Wrapf(err, "cluster %q never became Ready", v.ClusterName)
	}

	if err := v.waitForReadyInstances(ctx, kubeconfig, 4*time.Minute); err != nil {
		return err
	}

	// The read-write Service is the application-facing contract.
	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.ClusterName+"-rw", v.Namespace); err != nil {
		return errors.Wrap(err, "read-write service not found")
	}

	if v.BackupProof {
		if err := v.proveBackupCompleted(ctx, kubeconfig); err != nil {
			return err
		}
	}

	if !v.Behavioral {
		return nil
	}
	return v.proveFailoverDurability(ctx, kubeconfig)
}

// proveBackupCompleted waits for a Backup owned by this cluster to reach
// phase Completed — the immediate ScheduledBackup fires one on creation,
// and Completed means the plugin actually wrote a base backup into the
// object store (WAL archiving is a precondition the plugin enforces).
func (v *CnpgClusterVerifier) proveBackupCompleted(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] waiting for a Completed base backup of cluster %q\n", v.ClusterName)
	deadline := time.Now().Add(6 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "backups.postgresql.cnpg.io", "-n", v.Namespace,
			"-l", "cnpg.io/cluster="+v.ClusterName,
			"-o", "jsonpath={.items[*].status.phase}").CombinedOutput()
		phases := strings.Fields(strings.TrimSpace(string(out)))
		if err == nil {
			for _, phase := range phases {
				if phase == "completed" || phase == "Completed" {
					fmt.Printf("  [verify] base backup Completed — the plugin wrote a real backup to the store\n")
					return nil
				}
				if phase == "failed" || phase == "Failed" {
					return errors.Errorf("a backup of cluster %q reached terminal phase %q", v.ClusterName, phase)
				}
			}
		}
		last = fmt.Sprintf("phases=%v err=%v", phases, err)
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("no backup of cluster %q reached Completed (last: %s)", v.ClusterName, last)
}

func (v *CnpgClusterVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "cluster.postgresql.cnpg.io", v.ClusterName, v.Namespace)
}

// proveFailoverDurability runs the write → primary-loss → promotion →
// read-back cycle. The marker table is verifier-owned in the `postgres`
// database (always present regardless of the bootstrap shape); it dies
// with the cluster, so no cleanup is needed.
func (v *CnpgClusterVerifier) proveFailoverDurability(ctx context.Context, kubeconfig string) error {
	if v.Instances < 2 {
		return errors.New("behavioral failover needs instances >= 2 — there must be a replica to promote")
	}

	oldPrimary, err := v.currentPrimary(ctx, kubeconfig)
	if err != nil {
		return err
	}
	fmt.Printf("  [verify] behavioral failover: current primary %q\n", oldPrimary)

	// 1. Write the marker through the primary. synchronous_commit rides
	//    the cluster's own settings; the row must survive whatever the
	//    spec declared.
	const markerSQL = "CREATE TABLE IF NOT EXISTS e2e_failover_proof(id int PRIMARY KEY, note text); " +
		"INSERT INTO e2e_failover_proof(id, note) VALUES (1, 'postgres-failover-round-trip') " +
		"ON CONFLICT (id) DO UPDATE SET note = EXCLUDED.note;"
	if _, err := v.psql(ctx, kubeconfig, oldPrimary, markerSQL); err != nil {
		return errors.Wrap(err, "failed to write the marker row through the primary")
	}

	// 2. The disaster: the primary pod is deleted outright. The operator
	//    detects the loss and promotes the most advanced replica.
	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", oldPrimary, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
		return errors.Errorf("failed to delete primary pod: %v: %s", err, string(out))
	}

	// 3. Promotion: status.currentPrimary must move off the deleted pod.
	//    (The deleted pod's PVC makes it eligible to REJOIN as a replica
	//    later — the check is for a DIFFERENT primary, not pod absence.)
	newPrimary := ""
	deadline := time.Now().Add(4 * time.Minute)
	for time.Now().Before(deadline) {
		current, err := v.currentPrimary(ctx, kubeconfig)
		if err == nil && current != "" && current != oldPrimary {
			newPrimary = current
			break
		}
		time.Sleep(5 * time.Second)
	}
	if newPrimary == "" {
		return errors.Errorf("operator never promoted a replica off %q", oldPrimary)
	}
	fmt.Printf("  [verify] promoted: new primary %q\n", newPrimary)

	// 4. Full recovery: every instance ready again (the old primary
	//    rejoins as a replica), then the marker read back through the NEW
	//    primary.
	if err := v.waitForReadyInstances(ctx, kubeconfig, 4*time.Minute); err != nil {
		return errors.Wrap(err, "cluster never returned to full strength after the failover")
	}
	out, err := v.psql(ctx, kubeconfig, newPrimary,
		"SELECT note FROM e2e_failover_proof WHERE id = 1;")
	if err != nil {
		return errors.Wrap(err, "failed to read the marker row from the new primary")
	}
	if strings.TrimSpace(out) != "postgres-failover-round-trip" {
		return errors.Errorf("marker row wrong after failover (got %q)", strings.TrimSpace(out))
	}

	fmt.Printf("  [verify] marker row intact on the new primary — data survived the failover\n")
	return nil
}

// currentPrimary reads status.currentPrimary — the operator maintains it
// through failovers and switchovers.
func (v *CnpgClusterVerifier) currentPrimary(ctx context.Context, kubeconfig string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "cluster.postgresql.cnpg.io", v.ClusterName, "-n", v.Namespace,
		"-o", "jsonpath={.status.currentPrimary}").CombinedOutput()
	if err != nil {
		return "", errors.Errorf("failed to read currentPrimary: %v: %s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func (v *CnpgClusterVerifier) waitForReadyInstances(ctx context.Context, kubeconfig string, timeout time.Duration) error {
	want := fmt.Sprintf("%d", v.Instances)
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "cluster.postgresql.cnpg.io", v.ClusterName, "-n", v.Namespace,
			"-o", "jsonpath={.status.readyInstances}").CombinedOutput()
		got := strings.TrimSpace(string(out))
		if err == nil && got == want {
			return nil
		}
		last = fmt.Sprintf("readyInstances=%q err=%v", got, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("cluster %q never reached %s ready instances (last: %s)", v.ClusterName, want, last)
}

// psql runs a SQL string on an instance pod as the postgres OS user (peer
// auth inside the pod — the same path the operator's own probes use).
func (v *CnpgClusterVerifier) psql(ctx context.Context, kubeconfig, podName, sql string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", podName, "-n", v.Namespace, "-c", "postgres", "--",
		"psql", "-U", "postgres", "-d", "postgres", "-tA", "-c", sql).CombinedOutput()
	if err != nil {
		return "", errors.Errorf("psql on %s: %v: %s", podName, err, string(out))
	}
	return string(out), nil
}
