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

// PxcClusterVerifier checks a Percona-operator-managed MySQL (XtraDB
// Cluster) to the point it is actually serving: the PerconaXtraDBCluster
// resource in state "ready" and the proxy's write Service present.
//
// When Behavioral is set (the behavioral-durability scenario), the
// verifier proves Galera's synchronous replication live: it writes a
// marker row THROUGH THE PROXY, DELETES a database node outright, waits
// for the cluster to return to full strength, and reads the marker back
// through the proxy. Galera certifies every committed transaction on
// every node, so losing one node must lose nothing — that is the claim
// this proof pins down. SQL rides `kubectl exec` + the mysql client on a
// surviving database pod, connecting to the proxy Service with the
// operator-managed root credential.
//
// When BackupProof is set (the with-backup scenario), the verifier drives
// a verifier-owned PerconaXtraDBClusterBackup to state "Succeeded" — a
// real XtraBackup landing in the declared store. Verifier-owned because
// its CRD is installed by the operator fixture; deleted in a defer.
type PxcClusterVerifier struct {
	Namespace   string
	ClusterName string
	// Size is the declared database-node count.
	Size int64
	// ProxyService is the client-facing write Service
	// ("<name>-haproxy" or "<name>-proxysql").
	ProxyService  string
	Behavioral    bool
	BackupProof   bool
	BackupStorage string
}

func (v *PxcClusterVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] percona xtradb cluster %q in namespace %q (size=%d)\n",
		v.ClusterName, v.Namespace, v.Size)

	// state "ready" flips once every component (pxc nodes + proxy)
	// reached its desired size. Galera bootstraps SERIALLY — the first
	// node initializes, then each joiner runs a full SST before the next
	// starts — so the budget scales with the node count (a single node
	// converges in minutes; three-node quorum on a constrained cluster
	// can legitimately take beyond fifteen).
	readyBudget := 10*time.Minute + time.Duration(v.Size)*5*time.Minute
	if err := kubectlWait(ctx, kubeconfig, "pxc", v.ClusterName, v.Namespace,
		"jsonpath={.status.state}=ready", readyBudget); err != nil {
		return errors.Wrapf(err, "PerconaXtraDBCluster %q never reached state ready", v.ClusterName)
	}

	if err := KubectlResourceExists(ctx, kubeconfig, "service", v.ProxyService, v.Namespace); err != nil {
		return errors.Wrap(err, "proxy write service not found")
	}

	if v.BackupProof {
		if err := v.proveBackupCompleted(ctx, kubeconfig); err != nil {
			return err
		}
	}
	if !v.Behavioral {
		return nil
	}
	return v.proveGaleraDurability(ctx, kubeconfig)
}

func (v *PxcClusterVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	return KubectlResourceAbsent(ctx, kubeconfig, "pxc", v.ClusterName, v.Namespace)
}

// proveGaleraDurability runs the write-through-proxy → node-loss →
// recovery → read-through-proxy cycle. The marker schema is
// verifier-owned; it dies with the cluster.
func (v *PxcClusterVerifier) proveGaleraDurability(ctx context.Context, kubeconfig string) error {
	if v.Size < 3 {
		return errors.New("behavioral durability needs 3 database nodes — Galera quorum must survive the loss")
	}

	rootPassword, err := v.rootPassword(ctx, kubeconfig)
	if err != nil {
		return err
	}

	// 1. Write the marker through the PROXY — the same path applications
	//    use. Galera certifies the commit on every node synchronously.
	victim := fmt.Sprintf("%s-pxc-0", v.ClusterName)
	survivor := fmt.Sprintf("%s-pxc-1", v.ClusterName)
	writeSQL := "CREATE DATABASE IF NOT EXISTS e2e; " +
		"CREATE TABLE IF NOT EXISTS e2e.proof (id INT PRIMARY KEY, note VARCHAR(64)); " +
		"REPLACE INTO e2e.proof (id, note) VALUES (1, 'mysql-galera-round-trip');"
	if _, err := v.mysql(ctx, kubeconfig, victim, rootPassword, writeSQL); err != nil {
		return errors.Wrap(err, "failed to write the marker row through the proxy")
	}
	fmt.Printf("  [verify] marker committed through %q — killing node %q\n", v.ProxyService, victim)

	// 2. The disaster: a database node is deleted outright. HAProxy
	//    detects the loss and routes to a surviving node; Galera keeps
	//    quorum with 2 of 3.
	if out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"delete", "pod", victim, "-n", v.Namespace, "--wait=false").CombinedOutput(); err != nil {
		return errors.Errorf("failed to delete database node: %v: %s", err, string(out))
	}

	// 3. The marker must be readable through the proxy WHILE the node is
	//    down (the surviving pod runs the client; the proxy routes around
	//    the loss).
	readSQL := "SELECT note FROM e2e.proof WHERE id = 1;"
	var readBack string
	deadline := time.Now().Add(3 * time.Minute)
	for time.Now().Before(deadline) {
		out, err := v.mysql(ctx, kubeconfig, survivor, rootPassword, readSQL)
		if err == nil && strings.Contains(out, "mysql-galera-round-trip") {
			readBack = out
			break
		}
		time.Sleep(5 * time.Second)
	}
	if readBack == "" {
		return errors.New("marker row not readable through the proxy while the node was down")
	}
	fmt.Printf("  [verify] marker served during the outage — waiting for full strength\n")

	// 4. Full recovery: the StatefulSet recreates the node, it rejoins
	//    via IST/SST, and the cluster returns to ready at full size.
	if err := kubectlWait(ctx, kubeconfig, "pxc", v.ClusterName, v.Namespace,
		"jsonpath={.status.state}=ready", 10*time.Minute); err != nil {
		return errors.Wrap(err, "cluster never returned to ready after the node loss")
	}
	out, err := v.mysql(ctx, kubeconfig, victim, rootPassword, readSQL)
	if err != nil {
		return errors.Wrap(err, "failed to read the marker after recovery")
	}
	if !strings.Contains(out, "mysql-galera-round-trip") {
		return errors.Errorf("marker row wrong after recovery (got %q)", strings.TrimSpace(out))
	}
	fmt.Printf("  [verify] marker intact after the rejoined node — Galera lost nothing\n")
	return nil
}

// proveBackupCompleted applies a verifier-owned PerconaXtraDBClusterBackup
// and waits for it to reach state "Succeeded" — a real XtraBackup written
// into the declared store.
func (v *PxcClusterVerifier) proveBackupCompleted(ctx context.Context, kubeconfig string) error {
	backupName := fmt.Sprintf("e2e-proof-%d", time.Now().Unix())
	manifest := fmt.Sprintf(`apiVersion: pxc.percona.com/v1
kind: PerconaXtraDBClusterBackup
metadata:
  name: %s
  namespace: %s
spec:
  pxcCluster: %s
  storageName: %s
`, backupName, v.Namespace, v.ClusterName, v.BackupStorage)

	fmt.Printf("  [verify] driving a real XtraBackup %q to storage %q\n", backupName, v.BackupStorage)
	apply := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig, "apply", "-f", "-")
	apply.Stdin = strings.NewReader(manifest)
	if out, err := apply.CombinedOutput(); err != nil {
		return errors.Errorf("failed to apply the driver backup: %v: %s", err, string(out))
	}
	defer func() {
		_ = exec.Command("kubectl", "--kubeconfig", kubeconfig,
			"delete", "pxc-backup", backupName, "-n", v.Namespace, "--ignore-not-found").Run()
	}()

	deadline := time.Now().Add(8 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "pxc-backup", backupName, "-n", v.Namespace,
			"-o", "jsonpath={.status.state}").CombinedOutput()
		state := strings.TrimSpace(string(out))
		if err == nil {
			switch state {
			case "Succeeded":
				fmt.Printf("  [verify] backup %q Succeeded — XtraBackup wrote a real backup to the store\n", backupName)
				return nil
			case "Failed":
				return errors.Errorf("backup %q failed", backupName)
			}
		}
		last = fmt.Sprintf("state=%q err=%v", state, err)
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("backup %q never Succeeded (last: %s)", backupName, last)
}

// mysqlProxyService derives the client-facing write Service from the
// manifest's proxy flavor: "<name>-proxysql" when the proxysql arm is
// declared, "<name>-haproxy" otherwise (HAProxy is the default proxy).
func mysqlProxyService(name string, spec map[string]interface{}) string {
	if spec != nil {
		if proxy, ok := spec["proxy"].(map[string]interface{}); ok {
			if _, isProxysql := proxy["proxysql"]; isProxysql {
				return name + "-proxysql"
			}
		}
	}
	return name + "-haproxy"
}

// mysqlFirstBackupStorage reads the first declared backup storage's name —
// the storage the with-backup scenario's driver Backup writes to.
func mysqlFirstBackupStorage(spec map[string]interface{}) string {
	if spec == nil {
		return ""
	}
	backup, ok := spec["backup"].(map[string]interface{})
	if !ok {
		return ""
	}
	raw, ok := backup["storages"]
	if !ok {
		return ""
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 {
		return ""
	}
	first, ok := list[0].(map[string]interface{})
	if !ok {
		return ""
	}
	storageName, _ := first["name"].(string)
	return storageName
}

// rootPassword reads the operator-managed system-users Secret
// (`<cluster>-secrets`, key "root").
func (v *PxcClusterVerifier) rootPassword(ctx context.Context, kubeconfig string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "secret", v.ClusterName+"-secrets", "-n", v.Namespace,
		"-o", "jsonpath={.data.root}").CombinedOutput()
	if err != nil {
		return "", errors.Errorf("failed to read the root password: %v: %s", err, string(out))
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(out)))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode the root password")
	}
	return string(decoded), nil
}

// mysql runs SQL on a database pod's mysql client, connecting THROUGH the
// proxy Service (the application path), authenticated as root.
func (v *PxcClusterVerifier) mysql(ctx context.Context, kubeconfig, viaPod, rootPassword, sql string) (string, error) {
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", viaPod, "-n", v.Namespace, "-c", "pxc", "--",
		"mysql", "-h", v.ProxyService, "-uroot", "-p"+rootPassword,
		"--connect-timeout=10", "-N", "-e", sql).CombinedOutput()
	if err != nil {
		return "", errors.Errorf("mysql via %s through %s: %v: %s", viaPod, v.ProxyService, err, string(out))
	}
	return string(out), nil
}
