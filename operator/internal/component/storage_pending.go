package component

import (
	"context"
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// A stuck volume is the first wall real clusters hit (no default class, a
// default class whose driver was never installed, a backend rejecting our
// sizes) -- and without this file it surfaces as a generic "Waiting for X"
// forever. Each data component's not-ready branch asks ExplainPendingStorage
// for a better answer; the classification itself is pure so the wording and
// arms are pinned by offline tests.

// The operator only READS these to explain a stuck volume; it never installs
// storage drivers or mutates classes (that would require cloud credentials
// the platform deliberately never holds). Events are read so the status can
// relay the storage provisioner's own ProvisioningFailed message (a backend's
// minimum-size rejection, for instance) instead of a generic wait.
// +kubebuilder:rbac:groups=storage.k8s.io,resources=storageclasses;csidrivers,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list

// inTreeCSITranslations maps legacy in-tree provisioner names to the CSI
// driver that serves them after CSI migration (GA everywhere the chart's
// kubeVersion floor allows). A default class carrying an in-tree name is
// only provisionable when the translated driver is registered.
//
// Deliberate twin of the table in the CLI's self-hosted install preflight --
// the installer and
// the operator must diagnose identically; unification rides the shared
// boot-contract extraction.
var inTreeCSITranslations = map[string]string{
	"kubernetes.io/aws-ebs":        "ebs.csi.aws.com",
	"kubernetes.io/gce-pd":         "pd.csi.storage.gke.io",
	"kubernetes.io/azure-disk":     "disk.csi.azure.com",
	"kubernetes.io/azure-file":     "file.csi.azure.com",
	"kubernetes.io/vsphere-volume": "csi.vsphere.vmware.com",
	"kubernetes.io/cinder":         "cinder.csi.openstack.org",
}

// driverInstallHints name the substrate-specific way out when a class's CSI
// driver is missing -- the exact trap a fresh EKS cluster ships in (gp2 is
// default but the EBS CSI driver is an optional addon).
var driverInstallHints = map[string]string{
	"ebs.csi.aws.com":       "on EKS, install the aws-ebs-csi-driver addon and grant it an IAM role via Pod Identity or IRSA",
	"pd.csi.storage.gke.io": "on GKE, enable the Compute Engine persistent disk CSI driver addon",
	"disk.csi.azure.com":    "on AKS, the disk CSI driver ships enabled -- re-enable it if it was removed",
}

// ExplainPendingStorage looks for a Pending volume belonging to a component
// (claims whose name contains releaseName) and, when it can tell WHY the
// volume is stuck, returns a plain-language message naming the fix. Returns
// ("", false) when there is nothing better to say than the component's own
// generic waiting message. Read failures also return ("", false): explaining
// is best-effort and must never fail a reconcile.
func (b *Base) ExplainPendingStorage(ctx context.Context, c client.Client, namespace, releaseName string) (string, bool) {
	var pvcs corev1.PersistentVolumeClaimList
	if err := c.List(ctx, &pvcs, client.InNamespace(namespace)); err != nil {
		return "", false
	}
	var pending *corev1.PersistentVolumeClaim
	for i := range pvcs.Items {
		pvc := &pvcs.Items[i]
		if pvc.Status.Phase == corev1.ClaimPending && strings.Contains(pvc.Name, releaseName) {
			pending = pvc
			break
		}
	}
	if pending == nil {
		return "", false
	}

	var classes storagev1.StorageClassList
	if err := c.List(ctx, &classes); err != nil {
		return "", false
	}
	var drivers storagev1.CSIDriverList
	// Best-effort: a cluster without the CSIDriver API (or RBAC lag on an
	// upgrade) still gets the arms that need no driver facts.
	_ = c.List(ctx, &drivers)
	var events corev1.EventList
	_ = c.List(ctx, &events, client.InNamespace(namespace))

	return classifyPendingPVC(pending, classes.Items, drivers.Items, events.Items)
}

// classifyPendingPVC is the pure classifier: given one Pending claim and the
// cluster's storage facts, name the cause and the fix. Arms are ordered from
// most to least deterministic.
func classifyPendingPVC(pvc *corev1.PersistentVolumeClaim, classes []storagev1.StorageClass, drivers []storagev1.CSIDriver, events []corev1.Event) (string, bool) {
	// Arm 3 first when the provisioner itself already said WHY (e.g. "size
	// below the backend minimum") -- its message beats our inference.
	if msg := provisioningFailureMessage(pvc.Name, events); msg != "" {
		return fmt.Sprintf(
			"the storage backend rejected volume %s: %s -- if the backend enforces a minimum volume size, "+
				"raise spec.storage.size (or the component's storageSize) to at least that minimum",
			pvc.Name, msg), true
	}

	class := resolvedStorageClass(pvc, classes)
	if class == nil {
		return fmt.Sprintf(
			"volume %s cannot be provisioned: this cluster has no default StorageClass and the platform does not pin one. "+
				"Set spec.storage.storageClassName%s, or mark a class as default "+
				`(kubectl patch storageclass <name> -p '{"metadata": {"annotations": {"storageclass.kubernetes.io/is-default-class": "true"}}}')`,
			pvc.Name, availableClassesHint(classes)), true
	}

	if provisioner, ok := unservedCSIProvisioner(class, drivers); ok {
		hint := driverInstallHints[provisioner]
		if hint == "" {
			hint = "install the driver, or pin a provisionable class via spec.storage.storageClassName"
		}
		return fmt.Sprintf(
			"volume %s cannot be provisioned: StorageClass %q uses provisioner %q but that storage driver is not installed on this cluster -- %s",
			pvc.Name, class.Name, provisioner, hint), true
	}

	return "", false
}

// resolvedStorageClass returns the class this claim provisions through: its
// own spec.storageClassName when set, else the cluster's default class, else
// nil (nothing will ever provision it).
func resolvedStorageClass(pvc *corev1.PersistentVolumeClaim, classes []storagev1.StorageClass) *storagev1.StorageClass {
	if pvc.Spec.StorageClassName != nil && *pvc.Spec.StorageClassName != "" {
		for i := range classes {
			if classes[i].Name == *pvc.Spec.StorageClassName {
				return &classes[i]
			}
		}
		// A pinned class that does not exist is its own dead end; surface it
		// through the no-class arm by returning nil.
		return nil
	}
	for i := range classes {
		if classes[i].Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			return &classes[i]
		}
	}
	return nil
}

// unservedCSIProvisioner reports a provisioner that NEEDS a registered CSI
// driver but has none. Non-CSI external provisioners (local-path, NFS) have
// no CSIDriver object by design and must never trip this -- the rule is:
// translate in-tree names, then require a driver only for CSI-shaped names.
func unservedCSIProvisioner(class *storagev1.StorageClass, drivers []storagev1.CSIDriver) (string, bool) {
	provisioner := class.Provisioner
	if translated, ok := inTreeCSITranslations[provisioner]; ok {
		provisioner = translated
	} else if !strings.Contains(provisioner, "csi") {
		return "", false
	}
	for i := range drivers {
		if drivers[i].Name == provisioner {
			return "", false
		}
	}
	return provisioner, true
}

// provisioningFailureMessage relays the newest ProvisioningFailed event for
// the claim -- the provisioner's own words (e.g. NetApp's "Minimum size is
// 800GB") are more precise than anything we could infer.
func provisioningFailureMessage(pvcName string, events []corev1.Event) string {
	var newest *corev1.Event
	for i := range events {
		ev := &events[i]
		if ev.InvolvedObject.Kind != "PersistentVolumeClaim" || ev.InvolvedObject.Name != pvcName {
			continue
		}
		if ev.Reason != "ProvisioningFailed" {
			continue
		}
		if newest == nil || ev.LastTimestamp.After(newest.LastTimestamp.Time) {
			newest = ev
		}
	}
	if newest == nil {
		return ""
	}
	return strings.TrimSpace(newest.Message)
}

func availableClassesHint(classes []storagev1.StorageClass) string {
	if len(classes) == 0 {
		return " after installing a storage provisioner (this cluster has no StorageClasses at all)"
	}
	names := make([]string, 0, len(classes))
	for i := range classes {
		names = append(names, classes[i].Name)
	}
	sort.Strings(names)
	return fmt.Sprintf(" (available: %s)", strings.Join(names, ", "))
}
