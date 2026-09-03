package component

import (
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func pendingPVC(name string, className *string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.PersistentVolumeClaimSpec{StorageClassName: className},
		Status:     corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimPending},
	}
}

func storageClass(name, provisioner string, isDefault bool) storagev1.StorageClass {
	sc := storagev1.StorageClass{
		ObjectMeta:  metav1.ObjectMeta{Name: name},
		Provisioner: provisioner,
	}
	if isDefault {
		sc.Annotations = map[string]string{"storageclass.kubernetes.io/is-default-class": "true"}
	}
	return sc
}

func csiDriver(name string) storagev1.CSIDriver {
	return storagev1.CSIDriver{ObjectMeta: metav1.ObjectMeta{Name: name}}
}

// Arm 1 (the RKE2/NetApp shape): no default class, nothing pinned -- the
// message must name the spec pin AND the kubectl default-class patch.
func TestClassifyPendingPVC_NoDefaultClass(t *testing.T) {
	msg, ok := classifyPendingPVC(
		pendingPVC("data-planton-postgresql-0", nil),
		[]storagev1.StorageClass{storageClass("trident", "csi.trident.netapp.io", false)},
		nil, nil)
	if !ok {
		t.Fatal("a claim with no resolvable class must be explained")
	}
	for _, want := range []string{
		"no default StorageClass",
		"spec.storage.storageClassName",
		"kubectl patch storageclass",
		"trident",
	} {
		if !strings.Contains(msg, want) {
			t.Errorf("message must contain %q, got: %s", want, msg)
		}
	}
}

func TestClassifyPendingPVC_NoClassesAtAll(t *testing.T) {
	msg, ok := classifyPendingPVC(pendingPVC("data-planton-postgresql-0", nil), nil, nil, nil)
	if !ok {
		t.Fatal("a cluster with no StorageClasses must be explained")
	}
	if !strings.Contains(msg, "no StorageClasses at all") {
		t.Errorf("message must say the cluster has no classes: %s", msg)
	}
}

// Arm 2 (the EKS shape): the default class's provisioner resolves via CSI
// migration to a driver that is not installed.
func TestClassifyPendingPVC_DriverNotInstalled(t *testing.T) {
	msg, ok := classifyPendingPVC(
		pendingPVC("data-planton-postgresql-0", nil),
		[]storagev1.StorageClass{storageClass("gp2", "kubernetes.io/aws-ebs", true)},
		[]storagev1.CSIDriver{csiDriver("efs.csi.aws.com")}, nil)
	if !ok {
		t.Fatal("a class whose driver is missing must be explained")
	}
	for _, want := range []string{"gp2", "ebs.csi.aws.com", "not installed", "aws-ebs-csi-driver"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message must contain %q, got: %s", want, msg)
		}
	}
}

// Non-CSI external provisioners (Kind's local-path) have no CSIDriver object
// by design -- they must never be diagnosed as a missing driver.
func TestClassifyPendingPVC_NonCSIProvisionerIsNotADriverProblem(t *testing.T) {
	msg, ok := classifyPendingPVC(
		pendingPVC("data-planton-postgresql-0", nil),
		[]storagev1.StorageClass{storageClass("standard", "rancher.io/local-path", true)},
		nil, nil)
	if ok {
		t.Fatalf("local-path has no CSIDriver by design; nothing to explain, got: %s", msg)
	}
}

// A registered driver means storage is fine -- stay with the component's own
// generic waiting message (the claim is just slow, e.g. WaitForFirstConsumer).
func TestClassifyPendingPVC_HealthyClassStaysGeneric(t *testing.T) {
	if msg, ok := classifyPendingPVC(
		pendingPVC("data-planton-postgresql-0", nil),
		[]storagev1.StorageClass{storageClass("gp3", "ebs.csi.aws.com", true)},
		[]storagev1.CSIDriver{csiDriver("ebs.csi.aws.com")}, nil); ok {
		t.Fatalf("nothing is wrong; must not explain, got: %s", msg)
	}
}

// Arm 3 (the NetApp-minimum shape): the provisioner's own rejection is
// relayed verbatim with the size-raise guidance.
func TestClassifyPendingPVC_ProvisioningFailedRelayed(t *testing.T) {
	className := "trident"
	events := []corev1.Event{
		{
			InvolvedObject: corev1.ObjectReference{Kind: "PersistentVolumeClaim", Name: "data-planton-postgresql-0"},
			Reason:         "ProvisioningFailed",
			Message:        `Size "10GB" is too small. Minimum size is "800GB".`,
			LastTimestamp:  metav1.Time{Time: time.Now()},
		},
		// Another claim's failure must not be attributed to this one.
		{
			InvolvedObject: corev1.ObjectReference{Kind: "PersistentVolumeClaim", Name: "other-claim"},
			Reason:         "ProvisioningFailed",
			Message:        "unrelated",
			LastTimestamp:  metav1.Time{Time: time.Now()},
		},
	}
	msg, ok := classifyPendingPVC(
		pendingPVC("data-planton-postgresql-0", &className),
		[]storagev1.StorageClass{storageClass("trident", "csi.trident.netapp.io", false)},
		[]storagev1.CSIDriver{csiDriver("csi.trident.netapp.io")}, events)
	if !ok {
		t.Fatal("a ProvisioningFailed event must be explained")
	}
	for _, want := range []string{`Minimum size is "800GB"`, "spec.storage.size"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message must contain %q, got: %s", want, msg)
		}
	}
	if strings.Contains(msg, "unrelated") {
		t.Errorf("another claim's event leaked into the message: %s", msg)
	}
}

// A pinned class that does not exist resolves to nothing -- the no-class arm
// covers it (the message names what IS available).
func TestClassifyPendingPVC_PinnedClassMissing(t *testing.T) {
	className := "no-such-class"
	msg, ok := classifyPendingPVC(
		pendingPVC("planton-runner-state", &className),
		[]storagev1.StorageClass{storageClass("standard", "rancher.io/local-path", false)},
		nil, nil)
	if !ok {
		t.Fatal("a pinned-but-missing class must be explained")
	}
	if !strings.Contains(msg, "standard") {
		t.Errorf("message must list the classes that exist: %s", msg)
	}
}
