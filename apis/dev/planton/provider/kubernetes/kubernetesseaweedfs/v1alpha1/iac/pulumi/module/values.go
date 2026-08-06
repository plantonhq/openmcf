package module

import (
	"github.com/pkg/errors"
	kubernetesseaweedfsv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesseaweedfs/v1alpha1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// STORAGE POSTURE: the chart's out-of-the-box storage is hostPath under
// /ssd and /storage (bare-metal grain). This module deliberately maps every
// data volume to a PersistentVolumeClaim and every logs volume to emptyDir
// — portable across every managed cloud and kind cluster; the escape hatch
// can restore hostPath for bare-metal fleets.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// fullnameOverride pins seaweedfs.fullname to the resource name: the
	// componentName children (`<name>-master`, `-filer`, `-s3`, `-admin`,
	// the volume groups) derive deterministically and the longest suffix
	// stays far from the 63-char ceiling (the chart's componentName helper
	// truncates internally as a second guard).
	values["fullnameOverride"] = locals.ReleaseName

	// ---- global -------------------------------------------------------------
	global := map[string]interface{}{}
	// Cross-tier replication: one typed placement code flips the chart's
	// enableReplication and overrides master + filer placement together.
	if spec.GetReplication() != "" {
		global["enableReplication"] = true
		global["replicationPlacement"] = spec.GetReplication()
	}
	// ServiceMonitors are gated per component on this ONE flag (plus each
	// tier's metricsPort, which the chart defaults on).
	if spec.GetServiceMonitorEnabled() {
		global["monitoring"] = map[string]interface{}{"enabled": true}
	}
	if len(global) > 0 {
		values["global"] = map[string]interface{}{"seaweedfs": global}
	}

	// ---- image override --------------------------------------------------------
	// The chart's top-level image block wins over the global defaults:
	// image.repository REPLACES the whole image name; image.tag defaults
	// to the chart's appVersion.
	if img := spec.GetImage(); img != nil &&
		(img.GetRegistry() != "" || img.GetRepository() != "" || img.GetTag() != "") {
		image := map[string]interface{}{}
		if img.GetRegistry() != "" {
			image["registry"] = img.GetRegistry()
		}
		if img.GetRepository() != "" {
			image["repository"] = img.GetRepository()
		}
		if img.GetTag() != "" {
			image["tag"] = img.GetTag()
		}
		values["image"] = image
	}

	// ---- master tier --------------------------------------------------------------
	master := map[string]interface{}{
		"data": pvcDataMap(spec.GetMaster().GetDataVolume(), vars.DefaultMasterVolumeSize),
		"logs": map[string]interface{}{"type": "emptyDir"},
	}
	if spec.GetMaster().GetReplicas() > 0 {
		master["replicas"] = int(spec.GetMaster().GetReplicas())
	}
	if spec.GetMaster().GetVolumeSizeLimitMb() > 0 {
		master["volumeSizeLimitMB"] = int(spec.GetMaster().GetVolumeSizeLimitMb())
	}
	if r := resourcesMap(spec.GetMaster().GetResources()); r != nil {
		master["resources"] = r
	}
	values["master"] = master

	// ---- volume tier -----------------------------------------------------------------
	// dataDirs is a LIST in the chart; the typed surface models the
	// canonical single-PVC entry (named "data"), sized per pod. More
	// exotic layouts (multiple dirs, hostPath fleets) ride helm_values —
	// note lists REPLACE on merge, so an override provides the whole
	// list.
	volumeDataDir := pvcDataMap(spec.GetVolume().GetDataVolume(), vars.DefaultVolumeVolumeSize)
	volumeDataDir["name"] = "data"
	volumeDataDir["maxVolumes"] = int(spec.GetVolume().GetMaxVolumes())
	volume := map[string]interface{}{
		"dataDirs": []interface{}{volumeDataDir},
		"logs":     map[string]interface{}{"type": "emptyDir"},
	}
	if spec.GetVolume().GetReplicas() > 0 {
		volume["replicas"] = int(spec.GetVolume().GetReplicas())
	}
	if spec.GetVolume().GetIndexMode() != "" {
		volume["index"] = spec.GetVolume().GetIndexMode()
	}
	// NIL-SAFETY: presence of the optional scalar is checked on the FIELD,
	// which requires the parent message to exist — GetVolume() alone is
	// nil-safe, but `.MinFreeSpacePercent` on its nil result panics (the
	// runtime-panic class the minimal-shape preview proof exists to catch).
	if v := spec.GetVolume(); v != nil && v.MinFreeSpacePercent != nil {
		volume["minFreeSpacePercent"] = int(v.GetMinFreeSpacePercent())
	}
	if r := resourcesMap(spec.GetVolume().GetResources()); r != nil {
		volume["resources"] = r
	}
	values["volume"] = volume

	// ---- filer tier ---------------------------------------------------------------------
	// The filer's embedded leveldb metadata store lives on the data PVC
	// (the chart's WEED_LEVELDB2_ENABLED default) — external shared
	// stores (Postgres/MySQL) ride extra env vars + helm_values.
	filer := map[string]interface{}{
		"data": pvcDataMap(spec.GetFiler().GetDataVolume(), vars.DefaultFilerVolumeSize),
		"logs": map[string]interface{}{"type": "emptyDir"},
	}
	if spec.GetFiler().GetReplicas() > 0 {
		filer["replicas"] = int(spec.GetFiler().GetReplicas())
	}
	if spec.GetFiler().GetEncryptVolumeData() {
		filer["encryptVolumeData"] = true
	}
	if len(spec.GetFiler().GetExtraEnvironmentVars()) > 0 {
		filer["extraEnvironmentVars"] = stringMapToInterface(spec.GetFiler().GetExtraEnvironmentVars())
	}
	if r := resourcesMap(spec.GetFiler().GetResources()); r != nil {
		filer["resources"] = r
	}

	// ---- s3 gateway ------------------------------------------------------------------------
	// EMBEDDED (default): the gateway serves from the filer pods
	// (filer.s3.enabled). DEDICATED: its own Deployment (s3.enabled) that
	// scales independently. Both shapes expose the same `<name>-s3`
	// Service; the chart's s3-secret and bucket hook read auth and
	// existingConfigSecret from BOTH paths, so the module renders them on
	// both — only the enabled flags differ. Buckets always render under
	// s3.createBuckets (the hook's preferred path for either shape).
	s3Common := map[string]interface{}{
		"enableAuth": locals.S3Auth,
	}
	if locals.S3ConfigName != "" {
		s3Common["existingConfigSecret"] = locals.S3ConfigName
	}

	filerS3 := map[string]interface{}{"enabled": locals.S3Enabled && !locals.S3Dedicated}
	for k, v := range s3Common {
		filerS3[k] = v
	}
	if spec.GetS3().GetDomainName() != "" {
		filerS3["domainName"] = spec.GetS3().GetDomainName()
	}
	filer["s3"] = filerS3
	values["filer"] = filer

	s3 := map[string]interface{}{"enabled": locals.S3Enabled && locals.S3Dedicated}
	for k, v := range s3Common {
		s3[k] = v
	}
	if spec.GetS3().GetDomainName() != "" {
		s3["domainName"] = spec.GetS3().GetDomainName()
	}
	if buckets := bucketsSlice(spec.GetS3().GetBuckets()); len(buckets) > 0 {
		s3["createBuckets"] = buckets
	}
	if dedicated := spec.GetS3().GetDedicated(); dedicated != nil {
		if dedicated.GetReplicas() > 0 {
			s3["replicas"] = int(dedicated.GetReplicas())
		}
		if r := resourcesMap(dedicated.GetResources()); r != nil {
			s3["resources"] = r
		}
		s3["logs"] = map[string]interface{}{"type": "emptyDir"}
	}
	values["s3"] = s3

	// ---- admin console -------------------------------------------------------------------------
	// The console is never installed open: the chart requires
	// userKey/pwKey alongside existingSecret, and the module always
	// points at a credentials Secret — the module-materialized
	// `<name>-admin-auth` or the referenced existing one.
	if locals.AdminEnabled {
		admin := map[string]interface{}{
			"enabled": true,
			"secret": map[string]interface{}{
				"existingSecret": locals.AdminAuthSecretName,
				"userKey":        "user",
				"pwKey":          "password",
			},
		}
		if adminVolume := spec.GetAdmin().GetDataVolume(); adminVolume != nil {
			admin["data"] = pvcDataMap(adminVolume, vars.DefaultAdminVolumeSize)
		}
		if r := resourcesMap(spec.GetAdmin().GetResources()); r != nil {
			admin["resources"] = r
		}
		values["admin"] = admin
	}

	// ---- escape hatch (merged LAST, helm -f semantics) --------------------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). Every componentName
	// child (`-master`, `-filer`, `-s3`, `-admin`, `-volume`) and the
	// chart-generated `-s3-secret` — and the exported outputs built from
	// them — all derive from the fullname; letting an override move it
	// would break every output.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}

// pvcDataMap renders one tier's data volume as the chart's
// persistentVolumeClaim shape: size resolved to the tier default,
// storageClass only when declared (absent = the cluster's default class —
// the chart renders storageClassName empty, which Kubernetes treats as
// nil).
func pvcDataMap(volume *kubernetesseaweedfsv1alpha1.KubernetesSeaweedFsDataVolume, defaultSize string) map[string]interface{} {
	size := volume.GetSize()
	if size == "" {
		size = defaultSize
	}
	data := map[string]interface{}{
		"type": "persistentVolumeClaim",
		"size": size,
	}
	if sc := volume.GetStorageClass().GetValue(); sc != "" {
		data["storageClass"] = sc
	}
	return data
}

// bucketsSlice renders the typed bucket list into the chart's
// createBuckets shape (consumed by its post-install hook).
func bucketsSlice(buckets []*kubernetesseaweedfsv1alpha1.KubernetesSeaweedFsS3Bucket) []interface{} {
	out := make([]interface{}, 0, len(buckets))
	for _, b := range buckets {
		bucket := map[string]interface{}{"name": b.GetName()}
		if b.GetAnonymousRead() {
			bucket["anonymousRead"] = true
		}
		if b.GetTtl() != "" {
			bucket["ttl"] = b.GetTtl()
		}
		if b.GetObjectLock() {
			bucket["objectLock"] = true
		}
		if b.GetVersioning() {
			bucket["versioning"] = "Enabled"
		}
		out = append(out, bucket)
	}
	return out
}
