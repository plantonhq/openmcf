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

// PvcDataSourceVerifier proves data-source PROVISIONING for a claim whose
// spec carries the typed data_source arm (clone or snapshot-restore): a
// reader pod scheduled against the claim under test triggers the
// WaitForFirstConsumer provisioning, and the marker the fixture's writer
// stamped onto the SOURCE volume must be readable on the NEW volume — data
// actually traveled through the CSI clone/snapshot path.
//
// The VolumeSnapshot for the restore variant is verifier-owned: it must be
// cut only after the writer Job completes (a fixture-applied snapshot could
// freeze the volume before the marker lands), and cutting it is itself part
// of the behavioral surface under proof.
type PvcDataSourceVerifier struct {
	// Name of the claim under test.
	Name string
	// Snapshot switches the restore variant on: cut a VolumeSnapshot from
	// the source PVC before the read (the claim references it by name).
	Snapshot bool
}

// pvcDs* mirror the fixture manifest (fixture-data-source.yaml) and the
// scenarios' data_source references — fixed so they cannot drift apart.
const (
	pvcDsNamespace    = "e2e-pvc-ds"
	pvcDsSourceClaim  = "e2e-pvc-ds-source"
	pvcDsWriterJob    = "e2e-pvc-ds-writer"
	pvcDsSnapshotName = "e2e-pvc-ds-snapshot"
	pvcDsMarker       = "pvc-data-source-proof"
)

func (v *PvcDataSourceVerifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	fmt.Printf("  [verify] data-source provisioning for claim %q (snapshot=%v)\n", v.Name, v.Snapshot)

	if err := KubectlResourceExists(ctx, kubeconfig, "pvc", v.Name, pvcDsNamespace); err != nil {
		return err
	}

	// The source volume's marker must exist before anything is cloned or
	// snapshotted from it.
	if err := kubectlWait(ctx, kubeconfig, "job", pvcDsWriterJob, pvcDsNamespace,
		"condition=Complete", 4*time.Minute); err != nil {
		return errors.Wrap(err, "writer job never completed — the source volume has no marker to prove with")
	}

	if v.Snapshot {
		if err := v.cutSnapshot(ctx, kubeconfig); err != nil {
			return err
		}
	}

	// The reader triggers WaitForFirstConsumer provisioning of the claim
	// under test and asserts the marker traveled with the data.
	return v.readMarker(ctx, kubeconfig)
}

func (v *PvcDataSourceVerifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	// The verifier-owned snapshot is removed here rather than in a defer:
	// the claim under test references it, so it must outlive the whole
	// deploy-verify phase and go away with the destroy assertions.
	if v.Snapshot {
		_ = exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"delete", "volumesnapshot", pvcDsSnapshotName, "-n", pvcDsNamespace, "--ignore-not-found").Run()
	}
	return KubectlResourceAbsent(ctx, kubeconfig, "pvc", v.Name, pvcDsNamespace)
}

// cutSnapshot creates the VolumeSnapshot the claim under test restores from
// and waits for readyToUse (the CSI sidecar finished cutting it).
func (v *PvcDataSourceVerifier) cutSnapshot(ctx context.Context, kubeconfig string) error {
	snapshot := fmt.Sprintf(`apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshot
metadata:
  name: %s
  namespace: %s
spec:
  volumeSnapshotClassName: e2e-ebs-snapclass
  source:
    persistentVolumeClaimName: %s
`, pvcDsSnapshotName, pvcDsNamespace, pvcDsSourceClaim)
	snapFile, err := writeTempManifest(snapshot)
	if err != nil {
		return err
	}
	defer os.Remove(snapFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", snapFile); err != nil {
		return errors.Wrap(err, "failed to apply VolumeSnapshot")
	}

	deadline := time.Now().Add(5 * time.Minute)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
			"get", "volumesnapshot", pvcDsSnapshotName, "-n", pvcDsNamespace,
			"-o", "jsonpath={.status.readyToUse}").CombinedOutput()
		ready := strings.TrimSpace(string(out))
		if err == nil && ready == "true" {
			return nil
		}
		last = fmt.Sprintf("readyToUse=%q err=%v", ready, err)
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("VolumeSnapshot never became readyToUse (last: %s)", last)
}

// readMarker schedules a verifier-owned reader pod on the claim under test
// and asserts the source's marker is present on the provisioned volume.
func (v *PvcDataSourceVerifier) readMarker(ctx context.Context, kubeconfig string) error {
	readerName := "e2e-pvc-ds-reader"
	reader := fmt.Sprintf(`apiVersion: v1
kind: Pod
metadata:
  name: %s
  namespace: %s
spec:
  restartPolicy: Never
  containers:
    - name: reader
      image: busybox:1.36
      command:
        - sleep
        - "600"
      volumeMounts:
        - name: data
          mountPath: /data
      resources:
        requests:
          cpu: 25m
          memory: 16Mi
        limits:
          cpu: 100m
          memory: 32Mi
  volumes:
    - name: data
      persistentVolumeClaim:
        claimName: %s
`, readerName, pvcDsNamespace, v.Name)
	readerFile, err := writeTempManifest(reader)
	if err != nil {
		return err
	}
	defer os.Remove(readerFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", readerFile); err != nil {
		return errors.Wrap(err, "failed to apply reader pod")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "pod", readerName,
			"-n", pvcDsNamespace, "--ignore-not-found", "--wait=false")
	}()

	// Provisioning timing differs by path: restoring a pre-cut, readyToUse
	// snapshot lands in ~1 minute, but the EBS driver implements PVC
	// CLONING as an internal snapshot + restore (verified in the driver
	// source — EBS has no native volume clone), so the clone path routinely
	// takes several minutes for even a small volume.
	if err := kubectlWait(ctx, kubeconfig, "pod", readerName, pvcDsNamespace,
		"condition=Ready", 12*time.Minute); err != nil {
		return errors.Wrap(err, "reader pod never became Ready — data-source provisioning failed")
	}
	out, err := exec.CommandContext(ctx, "kubectl", "--kubeconfig", kubeconfig,
		"exec", "-n", pvcDsNamespace, readerName, "--", "cat", "/data/marker").CombinedOutput()
	if err != nil || strings.TrimSpace(string(out)) != pvcDsMarker {
		return errors.Errorf("marker missing on the provisioned volume (data=%q err=%v)", strings.TrimSpace(string(out)), err)
	}
	fmt.Printf("  [verify] source marker present on the provisioned volume — data traveled the CSI %s path\n",
		map[bool]string{true: "snapshot-restore", false: "clone"}[v.Snapshot])
	return nil
}

func (v *PvcDataSourceVerifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}
