package module

import (
	"fmt"

	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesvelerov1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesvelero/v1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), credentials doc,
// helm_values] and the provider merges the documents in exactly this
// order. Keep every typed mapping below in lockstep with the Terraform
// module's locals.
//
// No fullnameOverride: with the release named "velero" the chart's
// velero.fullname already collapses to "velero" (the release name contains
// the chart name) — every derived name, including the "velero-server"
// ServiceAccount, is deterministic without one.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// The chart's `configuration` block collects the BackupStorageLocation
	// and VolumeSnapshotLocation lists AND the velero-server CLI flags —
	// several spec sections contribute to it, so it is assembled across
	// this function and attached once at the end.
	configuration := map[string]interface{}{}

	// ---- CRD lifecycle ------------------------------------------------------
	// The chart ships its CRDs in the crds/ directory: Helm installs them
	// on first install and — by Helm's own contract for crds/-dir CRDs —
	// NEVER upgrades or deletes them on its own. The chart compensates
	// with two jobs: upgradeCRDs (re-applies the pinned CRDs on
	// install/upgrade; chart default true) and cleanUpCRDs (DESTRUCTIVE
	// CI-oriented delete on uninstall; chart default false). Both flags
	// are rendered only when they differ from the chart default, so the
	// values stay minimal — and backup records survive uninstall unless
	// the user explicitly opted into cleanup.
	if crds := spec.GetCrds(); crds != nil {
		if crds.UpgradeOnInstall != nil && !crds.GetUpgradeOnInstall() {
			values["upgradeCRDs"] = false
		}
		if crds.GetCleanupOnUninstall() {
			values["cleanUpCRDs"] = true
		}
	}

	// ---- backup storage backend ------------------------------------------------
	// Exactly one arm (spec-enforced oneof). Each arm renders the plugin
	// initContainer, the default BackupStorageLocation and the credential
	// posture — nothing from inactive arms leaks into the values.
	backupStorage := spec.GetBackupStorage()

	// BSL entry shared shape (chart values.yaml
	// configuration.backupStorageLocation[0]): name/provider/bucket/
	// default/prefix/caCert are item-level keys — the template renders
	// bucket, prefix and caCert under the BSL's spec.objectStorage — and
	// everything provider-specific goes into the `config` map, which the
	// template renders with every value QUOTED (config values are
	// strings; booleans must be rendered as "true"/"false").
	bsl := map[string]interface{}{
		"name":    vars.BackupStorageLocationName,
		"default": true,
	}
	if backupStorage.GetPrefix() != "" {
		bsl["prefix"] = backupStorage.GetPrefix()
	}

	var vslProvider string
	vslConfig := map[string]interface{}{}

	switch {
	case backupStorage.GetS3() != nil:
		s3 := backupStorage.GetS3()
		bsl["provider"] = "aws"
		bsl["bucket"] = s3.GetBucket()
		// caCert is an item-level BSL key (rendered under
		// objectStorage.caCert by templates/backupstoragelocation.yaml),
		// NOT a config entry.
		if s3.GetCaCert() != "" {
			bsl["caCert"] = s3.GetCaCert()
		}
		config := map[string]interface{}{
			"region": s3.GetRegion(),
		}
		if s3.GetS3Url() != "" {
			config["s3Url"] = s3.GetS3Url()
		}
		if s3.GetForcePathStyle() {
			// String, not bool: the chart template quotes every config
			// value ({{ $value | quote }}) — keep the rendered document
			// honest about the type the BSL API receives.
			config["s3ForcePathStyle"] = "true"
		}
		if s3.GetKmsKeyId() != "" {
			config["kmsKeyId"] = s3.GetKmsKeyId()
		}
		bsl["config"] = config

		values["initContainers"] = []interface{}{
			pluginInitContainer(vars.AwsPluginName, backupStorage.GetPluginImage(), vars.AwsPluginImage),
		}

		vslProvider = "aws"
		vslConfig["region"] = s3.GetRegion()

		// Credential posture (spec-enforced: at most one of IRSA /
		// access keys; an s3_url endpoint never rides IRSA).
		switch {
		case s3.GetIrsaRoleArn() != "":
			// Keyless EKS posture: the role-arn annotation on the server
			// ServiceAccount makes the pod exchange its projected token
			// for role credentials — no Secret exists at all.
			values["serviceAccount"] = serverServiceAccountAnnotations(map[string]interface{}{
				"eks.amazonaws.com/role-arn": s3.GetIrsaRoleArn(),
			})
			values["credentials"] = map[string]interface{}{"useSecret": false}
		case s3.GetAccessKeys() != nil:
			// Declared-credential posture: the AWS plugin reads the
			// shared-credentials `cloud` file (chart values.yaml
			// documents the exact format under
			// credentials.secretContents).
			values["credentials"] = map[string]interface{}{
				"useSecret": true,
				"secretContents": map[string]interface{}{
					"cloud": fmt.Sprintf("[default]\naws_access_key_id=%s\naws_secret_access_key=%s\n",
						s3.GetAccessKeys().GetAccessKeyId(),
						s3.GetAccessKeys().GetSecretAccessKey()),
				},
			}
		default:
			// Ambient node credentials (instance profile, kube2iam, ...).
			values["credentials"] = map[string]interface{}{"useSecret": false}
		}

	case backupStorage.GetGcs() != nil:
		gcs := backupStorage.GetGcs()
		bsl["provider"] = "gcp"
		bsl["bucket"] = gcs.GetBucket()

		values["initContainers"] = []interface{}{
			pluginInitContainer(vars.GcpPluginName, backupStorage.GetPluginImage(), vars.GcpPluginImage),
		}

		vslProvider = "gcp"
		// GCP needs no VolumeSnapshotLocation config (values.yaml lists
		// only optional keys — project, snapshotLocation).

		switch {
		case gcs.GetWorkloadIdentityServiceAccountEmail() != "":
			// Keyless GKE posture: the WI annotation binds the pod to the
			// GCP service account, and the BSL's config.serviceAccount
			// (values.yaml: "Specify the service account here if you want
			// to use workload identity instead of providing the key
			// file") points the plugin at the same identity.
			bsl["config"] = map[string]interface{}{
				"serviceAccount": gcs.GetWorkloadIdentityServiceAccountEmail(),
			}
			values["serviceAccount"] = serverServiceAccountAnnotations(map[string]interface{}{
				"iam.gke.io/gcp-service-account": gcs.GetWorkloadIdentityServiceAccountEmail(),
			})
			values["credentials"] = map[string]interface{}{"useSecret": false}
		case gcs.GetServiceAccountKeyJson() != "":
			// Declared-credential posture: the GCP plugin reads the JSON
			// key file verbatim as the `cloud` credentials file.
			values["credentials"] = map[string]interface{}{
				"useSecret": true,
				"secretContents": map[string]interface{}{
					"cloud": gcs.GetServiceAccountKeyJson(),
				},
			}
		default:
			// Ambient node credentials (GCE default service account).
			values["credentials"] = map[string]interface{}{"useSecret": false}
		}

	case backupStorage.GetAzureBlob() != nil:
		azure := backupStorage.GetAzureBlob()
		bsl["provider"] = "azure"
		// Velero's generic "bucket" is the blob CONTAINER on Azure.
		bsl["bucket"] = azure.GetContainer()
		bsl["config"] = map[string]interface{}{
			"resourceGroup":  azure.GetResourceGroup(),
			"storageAccount": azure.GetStorageAccount(),
			"subscriptionId": azure.GetSubscriptionId(),
		}

		values["initContainers"] = []interface{}{
			pluginInitContainer(vars.AzurePluginName, backupStorage.GetPluginImage(), vars.AzurePluginImage),
		}

		vslProvider = "azure"
		vslConfig["resourceGroup"] = azure.GetResourceGroup()
		vslConfig["subscriptionId"] = azure.GetSubscriptionId()

		switch {
		case azure.GetUseWorkloadIdentity():
			// Keyless AKS posture. Unlike AWS/GCP, the Azure plugin STILL
			// reads a `cloud` file for the non-credential parameters
			// (subscription, resource group, cloud name) — the chart
			// README defers the file's contents to the Azure plugin's own
			// README, which prescribes exactly these three lines for
			// workload identity. The federated token then rides the
			// azure.workload.identity/client-id annotation plus the
			// azure.workload.identity/use pod label (the AKS webhook only
			// injects the token into labeled pods).
			values["credentials"] = map[string]interface{}{
				"useSecret": true,
				"secretContents": map[string]interface{}{
					"cloud": fmt.Sprintf("AZURE_SUBSCRIPTION_ID=%s\nAZURE_RESOURCE_GROUP=%s\nAZURE_CLOUD_NAME=AzurePublicCloud\n",
						azure.GetSubscriptionId(),
						azure.GetResourceGroup()),
				},
			}
			values["serviceAccount"] = serverServiceAccountAnnotations(map[string]interface{}{
				"azure.workload.identity/client-id": azure.GetWorkloadIdentityClientId(),
			})
			values["podLabels"] = map[string]interface{}{
				"azure.workload.identity/use": "true",
			}
		case azure.GetServicePrincipal() != nil:
			// Declared-credential posture: the plugin's environment-file
			// format (AZURE_* lines), per the Azure plugin README the
			// chart values.yaml points to.
			sp := azure.GetServicePrincipal()
			values["credentials"] = map[string]interface{}{
				"useSecret": true,
				"secretContents": map[string]interface{}{
					"cloud": fmt.Sprintf("AZURE_SUBSCRIPTION_ID=%s\nAZURE_TENANT_ID=%s\nAZURE_CLIENT_ID=%s\nAZURE_CLIENT_SECRET=%s\nAZURE_RESOURCE_GROUP=%s\nAZURE_CLOUD_NAME=AzurePublicCloud\n",
						azure.GetSubscriptionId(),
						sp.GetTenantId(),
						sp.GetClientId(),
						sp.GetClientSecret(),
						azure.GetResourceGroup()),
				},
			}
		default:
			// Ambient credentials (AAD pod identity / managed identity on
			// the node pool).
			values["credentials"] = map[string]interface{}{"useSecret": false}
		}
	}

	configuration["backupStorageLocation"] = []interface{}{bsl}

	// ---- volume snapshots -------------------------------------------------------
	// snapshotsEnabled matches the chart's own default (true) — rendered
	// only on explicit opt-out. While snapshots are on, the default
	// VolumeSnapshotLocation rides the active arm's provider with the
	// provider-required config keys (values.yaml VSL comments: region for
	// aws; resourceGroup + subscriptionId for azure; nothing required for
	// gcp).
	snapshotsEnabled := true
	if vs := spec.GetVolumeSnapshots(); vs != nil && vs.Enabled != nil && !vs.GetEnabled() {
		snapshotsEnabled = false
		values["snapshotsEnabled"] = false
	}
	if snapshotsEnabled {
		vsl := map[string]interface{}{
			"name":     vars.BackupStorageLocationName,
			"provider": vslProvider,
		}
		if len(vslConfig) > 0 {
			vsl["config"] = vslConfig
		}
		configuration["volumeSnapshotLocation"] = []interface{}{vsl}
	}
	if vs := spec.GetVolumeSnapshots(); vs != nil {
		if vs.GetEnableCsi() {
			// The velero server's feature-flag list; EnableCSI switches PVC
			// snapshotting to the CSI snapshot API.
			configuration["features"] = "EnableCSI"
		}
		if vs.GetDefaultSnapshotMoveData() {
			configuration["defaultSnapshotMoveData"] = true
		}
	}

	// ---- file-system backup -------------------------------------------------------
	// deployNodeAgent matches the chart's own default (false) — rendered
	// only when the DaemonSet is wanted. defaultVolumesToFsBackup is a
	// velero-server flag, so it lives under `configuration`, not under
	// `nodeAgent`.
	if fs := spec.GetFsBackup(); fs != nil {
		if fs.GetDeployNodeAgent() {
			values["deployNodeAgent"] = true
		}
		if fs.GetDefaultVolumesToFsBackup() {
			configuration["defaultVolumesToFsBackup"] = true
		}
		nodeAgent := map[string]interface{}{}
		if r := resourcesMap(fs.GetNodeAgentResources()); r != nil {
			nodeAgent["resources"] = r
		}
		if len(fs.GetNodeAgentTolerations()) > 0 {
			nodeAgent["tolerations"] = tolerationsSlice(fs.GetNodeAgentTolerations())
		}
		if len(nodeAgent) > 0 {
			values["nodeAgent"] = nodeAgent
		}
	}

	// ---- scheduled backups ----------------------------------------------------------
	// The chart's `schedules` value is a MAP keyed by schedule name (the
	// rendered Schedule object is named "velero-<key>" — velero.fullname
	// plus the key). The spec's list converts to that map; name
	// uniqueness is spec-enforced.
	if len(spec.GetSchedules()) > 0 {
		schedules := map[string]interface{}{}
		for _, s := range spec.GetSchedules() {
			schedules[s.GetName()] = scheduleMap(s)
		}
		values["schedules"] = schedules
	}

	// ---- server tuning ---------------------------------------------------------------
	// All velero-server CLI flags rendered through the chart's
	// `configuration` passthrough — only what the spec sets, so the
	// server's own defaults stay in force otherwise.
	if server := spec.GetServer(); server != nil {
		if server.GetDefaultBackupTtl() != "" {
			configuration["defaultBackupTTL"] = server.GetDefaultBackupTtl()
		}
		if server.GetDefaultItemOperationTimeout() != "" {
			configuration["defaultItemOperationTimeout"] = server.GetDefaultItemOperationTimeout()
		}
		if server.GetGarbageCollectionFrequency() != "" {
			configuration["garbageCollectionFrequency"] = server.GetGarbageCollectionFrequency()
		}
		if server.GetRestoreOnlyMode() {
			configuration["restoreOnlyMode"] = true
		}
		if server.GetLogLevel() != "" {
			configuration["logLevel"] = server.GetLogLevel()
		}
		if server.GetLogFormat() != "" {
			configuration["logFormat"] = server.GetLogFormat()
		}
	}

	values["configuration"] = configuration

	// ---- deployment sizing / scheduling --------------------------------------------------
	// These are top-level chart keys applying to the Velero server
	// Deployment (the node-agent has its own block, handled above).
	if deployment := spec.GetDeployment(); deployment != nil {
		if r := resourcesMap(deployment.GetResources()); r != nil {
			values["resources"] = r
		}
		if deployment.GetPriorityClassName() != "" {
			values["priorityClassName"] = deployment.GetPriorityClassName()
		}
		if len(deployment.GetNodeSelector()) > 0 {
			values["nodeSelector"] = stringMapToInterface(deployment.GetNodeSelector())
		}
		if len(deployment.GetTolerations()) > 0 {
			values["tolerations"] = tolerationsSlice(deployment.GetTolerations())
		}
	}

	// ---- own telemetry ------------------------------------------------------------------
	// metrics.enabled matches the chart's own default (true) — rendered
	// only on explicit opt-out. The ServiceMonitor is opt-in (it needs the
	// Prometheus operator CRDs on the cluster or the release fails).
	if p := spec.GetPrometheus(); p != nil {
		metrics := map[string]interface{}{}
		if p.Enabled != nil && !p.GetEnabled() {
			metrics["enabled"] = false
		}
		if p.GetServiceMonitor() {
			metrics["serviceMonitor"] = map[string]interface{}{"enabled": true}
		}
		if len(metrics) > 0 {
			values["metrics"] = metrics
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) ------------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}

// hasSecretMaterial reports whether the rendered values carry actual
// secret content (an S3 secret access key, a GCP service-account key, an
// Azure client secret). Used to mask the whole values map in state —
// coarser than field-level masking, but nothing secret can leak through a
// missed path. The Azure workload-identity `cloud` file carries only
// non-secret identifiers and does not trigger masking.
func hasSecretMaterial(spec *kubernetesvelerov1.KubernetesVeleroSpec) bool {
	backupStorage := spec.GetBackupStorage()
	switch {
	case backupStorage.GetS3() != nil:
		return backupStorage.GetS3().GetAccessKeys() != nil
	case backupStorage.GetGcs() != nil:
		return backupStorage.GetGcs().GetServiceAccountKeyJson() != ""
	case backupStorage.GetAzureBlob() != nil:
		return backupStorage.GetAzureBlob().GetServicePrincipal() != nil
	}
	return false
}

// pluginInitContainer renders the provider plugin initContainer the chart
// expects verbatim under `initContainers` (values.yaml documents the exact
// shape: image + volumeMount of the shared `plugins` dir at /target, where
// the Velero server discovers plugin binaries).
func pluginInitContainer(name, imageOverride, defaultImage string) map[string]interface{} {
	image := defaultImage
	if imageOverride != "" {
		image = imageOverride
	}
	return map[string]interface{}{
		"name":  name,
		"image": image,
		"volumeMounts": []interface{}{
			map[string]interface{}{
				"mountPath": "/target",
				"name":      "plugins",
			},
		},
	}
}

// serverServiceAccountAnnotations wraps identity annotations in the
// chart's serviceAccount.server block (the chart still creates the
// ServiceAccount; only annotations are injected).
func serverServiceAccountAnnotations(annotations map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"server": map[string]interface{}{
			"annotations": annotations,
		},
	}
}

// scheduleMap renders one spec schedule into the chart's schedules-map
// entry: {schedule, paused, template{...}} with empties omitted and the
// three optional booleans rendered presence-aware (unset means "Velero
// decides", which is different from false).
func scheduleMap(s *kubernetesvelerov1.KubernetesVeleroSchedule) map[string]interface{} {
	entry := map[string]interface{}{
		"schedule": s.GetSchedule(),
	}
	if s.GetPaused() {
		entry["paused"] = true
	}

	template := map[string]interface{}{}
	if s.GetTtl() != "" {
		template["ttl"] = s.GetTtl()
	}
	if len(s.GetIncludedNamespaces()) > 0 {
		template["includedNamespaces"] = stringSliceToInterface(s.GetIncludedNamespaces())
	}
	if len(s.GetExcludedNamespaces()) > 0 {
		template["excludedNamespaces"] = stringSliceToInterface(s.GetExcludedNamespaces())
	}
	if len(s.GetIncludedResources()) > 0 {
		template["includedResources"] = stringSliceToInterface(s.GetIncludedResources())
	}
	if len(s.GetExcludedResources()) > 0 {
		template["excludedResources"] = stringSliceToInterface(s.GetExcludedResources())
	}
	if len(s.GetLabelSelector()) > 0 {
		template["labelSelector"] = map[string]interface{}{
			"matchLabels": stringMapToInterface(s.GetLabelSelector()),
		}
	}
	if s.IncludeClusterResources != nil {
		template["includeClusterResources"] = s.GetIncludeClusterResources()
	}
	if s.SnapshotVolumes != nil {
		template["snapshotVolumes"] = s.GetSnapshotVolumes()
	}
	if s.DefaultVolumesToFsBackup != nil {
		template["defaultVolumesToFsBackup"] = s.GetDefaultVolumesToFsBackup()
	}
	if s.GetStorageLocation() != "" {
		template["storageLocation"] = s.GetStorageLocation()
	}
	if len(template) > 0 {
		entry["template"] = template
	}

	return entry
}

// resourcesMap renders the shared ContainerResources message into the
// chart's resources shape. Returns nil when nothing is set.
func resourcesMap(r *kubernetesprovider.ContainerResources) map[string]interface{} {
	if r == nil {
		return nil
	}
	out := map[string]interface{}{}
	if l := r.GetLimits(); l != nil && (l.GetCpu() != "" || l.GetMemory() != "") {
		limits := map[string]interface{}{}
		if l.GetCpu() != "" {
			limits["cpu"] = l.GetCpu()
		}
		if l.GetMemory() != "" {
			limits["memory"] = l.GetMemory()
		}
		out["limits"] = limits
	}
	if q := r.GetRequests(); q != nil && (q.GetCpu() != "" || q.GetMemory() != "") {
		requests := map[string]interface{}{}
		if q.GetCpu() != "" {
			requests["cpu"] = q.GetCpu()
		}
		if q.GetMemory() != "" {
			requests["memory"] = q.GetMemory()
		}
		out["requests"] = requests
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// tolerationsSlice renders the shared WorkloadToleration list into the
// chart's tolerations shape.
func tolerationsSlice(tolerations []*kubernetesprovider.WorkloadToleration) []interface{} {
	out := make([]interface{}, 0, len(tolerations))
	for _, t := range tolerations {
		tol := map[string]interface{}{}
		if t.GetKey() != "" {
			tol["key"] = t.GetKey()
		}
		if t.GetOperator() != "" {
			tol["operator"] = t.GetOperator()
		}
		if t.GetValue() != "" {
			tol["value"] = t.GetValue()
		}
		if t.GetEffect() != "" {
			tol["effect"] = t.GetEffect()
		}
		if t.TolerationSeconds != nil {
			tol["tolerationSeconds"] = t.GetTolerationSeconds()
		}
		out = append(out, tol)
	}
	return out
}

// mergeMaps deep-merges b over a with Helm's `-f` semantics: nested maps
// merge recursively with b winning per key; everything else (scalars,
// lists) is replaced by b's value.
func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if bChild, ok := v.(map[string]interface{}); ok {
			if aChild, ok := out[k].(map[string]interface{}); ok {
				out[k] = mergeMaps(aChild, bChild)
				continue
			}
		}
		out[k] = v
	}
	return out
}

// stringMapToInterface converts a map[string]string into the
// map[string]interface{} YAML rendering expects.
func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// stringSliceToInterface converts a []string into the []interface{} YAML
// rendering expects.
func stringSliceToInterface(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}
