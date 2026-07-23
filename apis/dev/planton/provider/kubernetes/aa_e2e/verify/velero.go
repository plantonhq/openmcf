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

// VeleroInstallVerifier checks a Velero installation to the point the DR
// engine is actually serving: the server deployment Available and the
// default BackupStorageLocation phase Available — the store handshake
// (plugin loaded, bucket reachable, credentials valid) actually happened,
// not merely pods running.
//
// When Behavioral is set (the behavioral-backup-restore scenario), the
// verifier additionally proves the DR LOOP: it creates a subject namespace
// with a marker ConfigMap, drives a Backup CR to Completed, DELETES the
// subject namespace, drives a Restore CR to Completed, and asserts the
// marker ConfigMap came back. Backup and Restore CRs are verifier-owned on
// purpose: their CRDs are installed by the component under test, so
// scenario fixtures (which deploy first) can never carry them.
type VeleroInstallVerifier struct {
	Namespace string
	// Behavioral switches on the live backup/restore proof (see above).
	Behavioral bool
}

// veleroSubject* identify the verifier-owned backup subject. Fixed names so
// the Backup's includedNamespaces, the assertions, and the cleanup cannot
// drift apart. The subject lives in its OWN namespace: deleting and
// restoring it exercises exactly the disaster path (namespace gone,
// resources recovered from the store).
const (
	veleroSubjectNamespace = "e2e-velero-subject"
	veleroSubjectConfigMap = "e2e-velero-subject-marker"
	veleroBackupName       = "e2e-velero-backup"
	veleroRestoreName      = "e2e-velero-restore"
)

func (v *VeleroInstallVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] velero installation in namespace %q\n", v.Namespace)

	if err := KubectlResourceExists(ctx, kubeconfig, "namespace", v.Namespace, ""); err != nil {
		return errors.Wrapf(err, "namespace %q not found for velero", v.Namespace)
	}

	// The chart's fullname with the fixed release name "velero" is
	// "velero" — one server deployment per cluster.
	if err := kubectlWait(ctx, kubeconfig, "deployment", "velero", v.Namespace,
		"condition=Available", 3*time.Minute); err != nil {
		return errors.Wrap(err, "velero server deployment not available")
	}

	// BSL Available is the real install contract: the object-store plugin
	// loaded and the bucket answered. Velero revalidates on a cadence, so
	// poll rather than expect an instant phase.
	if err := v.waitForJSONPath(ctx, kubeconfig, 3*time.Minute,
		"backupstoragelocation", "default", v.Namespace,
		"{.status.phase}", "Available"); err != nil {
		return errors.Wrap(err, "default BackupStorageLocation never became Available — the store handshake failed")
	}

	if !v.Behavioral {
		return nil
	}
	return v.proveBackupRestore(ctx, kubeconfig)
}

func (v *VeleroInstallVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The CRDs intentionally SURVIVE uninstall (the chart serves them from
	// its crds/ directory, which Helm never deletes — backup records must
	// outlive the component), so only the server's presence is asserted.
	return KubectlResourceAbsent(ctx, kubeconfig, "deployment", "velero", v.Namespace)
}

// proveBackupRestore runs the full DR cycle. Every resource it creates is
// verifier-owned and removed in defers, so a failed assertion cannot leak
// state that would block later lanes; the backed-up object data dies with
// the scenario's MinIO fixture.
func (v *VeleroInstallVerifier) proveBackupRestore(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] behavioral backup/restore: %q must survive namespace deletion\n", veleroSubjectConfigMap)

	// 1. The backup subject: a namespace holding a marker ConfigMap whose
	//    contents the restore assertion checks byte-for-byte.
	if err := v.kubectl(ctx, kubeconfig, "create", "namespace", veleroSubjectNamespace); err != nil {
		return errors.Wrap(err, "failed to create backup-subject namespace")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "namespace", veleroSubjectNamespace, "--ignore-not-found")
	}()
	if err := v.kubectl(ctx, kubeconfig, "create", "configmap", veleroSubjectConfigMap,
		"-n", veleroSubjectNamespace, "--from-literal=proof=velero-round-trip"); err != nil {
		return errors.Wrap(err, "failed to create marker configmap")
	}

	// 2. Back the subject up. The Backup CR is served by the component
	//    under test — applied here, after the install assertions.
	backup := fmt.Sprintf(`apiVersion: velero.io/v1
kind: Backup
metadata:
  name: %s
  namespace: %s
spec:
  includedNamespaces:
    - %s
  snapshotVolumes: false
  ttl: 1h0m0s
`, veleroBackupName, v.Namespace, veleroSubjectNamespace)
	backupFile, err := writeTempManifest(backup)
	if err != nil {
		return err
	}
	defer os.Remove(backupFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", backupFile); err != nil {
		return errors.Wrap(err, "failed to apply Backup CR")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "backup.velero.io", veleroBackupName,
			"-n", v.Namespace, "--ignore-not-found")
	}()
	if err := v.waitForJSONPath(ctx, kubeconfig, 4*time.Minute,
		"backup.velero.io", veleroBackupName, v.Namespace,
		"{.status.phase}", "Completed"); err != nil {
		return errors.Wrap(err, "Backup never reached phase Completed")
	}

	// 3. The disaster: the subject namespace is deleted outright. Wait for
	//    it to be fully gone — a restore into a Terminating namespace
	//    would flake.
	if err := v.kubectl(ctx, kubeconfig, "delete", "namespace", veleroSubjectNamespace, "--wait=true", "--timeout=2m"); err != nil {
		return errors.Wrap(err, "failed to delete the backup-subject namespace")
	}

	// 4. Restore from the store and assert the marker returned.
	restore := fmt.Sprintf(`apiVersion: velero.io/v1
kind: Restore
metadata:
  name: %s
  namespace: %s
spec:
  backupName: %s
`, veleroRestoreName, v.Namespace, veleroBackupName)
	restoreFile, err := writeTempManifest(restore)
	if err != nil {
		return err
	}
	defer os.Remove(restoreFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", restoreFile); err != nil {
		return errors.Wrap(err, "failed to apply Restore CR")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "restore.velero.io", veleroRestoreName,
			"-n", v.Namespace, "--ignore-not-found")
	}()
	if err := v.waitForJSONPath(ctx, kubeconfig, 4*time.Minute,
		"restore.velero.io", veleroRestoreName, v.Namespace,
		"{.status.phase}", "Completed"); err != nil {
		return errors.Wrap(err, "Restore never reached phase Completed")
	}

	cmd := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"get", "configmap", veleroSubjectConfigMap, "-n", veleroSubjectNamespace,
		"-o", "jsonpath={.data.proof}")
	out, err := cmd.CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) != "velero-round-trip" {
		return errors.Errorf("restored marker configmap wrong or missing (data=%q err=%v)", strings.TrimSpace(string(out)), err)
	}
	fmt.Printf("  [verify] marker configmap restored with intact data — Velero completed a real DR loop\n")
	return nil
}

// waitForJSONPath polls a resource's jsonpath until it equals want. Velero
// phases move through queueing states (New → InProgress → Completed), so a
// plain `kubectl wait --for=condition` cannot express them.
func (v *VeleroInstallVerifier) waitForJSONPath(ctx context.Context, kubeconfig string, timeout time.Duration,
	kind, name, namespace, jsonPath, want string) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		args := []string{"--kubeconfig", kubeconfig, "get", kind, name, "-o", "jsonpath=" + jsonPath}
		if namespace != "" {
			args = append(args, "-n", namespace)
		}
		out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
		got := strings.TrimSpace(string(out))
		if err == nil && got == want {
			return nil
		}
		// A terminal failure phase will never progress — fail fast with
		// the phase in hand instead of burning the whole timeout.
		if err == nil && (got == "Failed" || got == "PartiallyFailed" || got == "FailedValidation") {
			return errors.Errorf("%s %q reached terminal phase %q", kind, name, got)
		}
		last = fmt.Sprintf("value=%q err=%v", got, err)
		time.Sleep(5 * time.Second)
	}
	return errors.Errorf("%s %q never reached %q (last: %s)", kind, name, want, last)
}

func (v *VeleroInstallVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}
