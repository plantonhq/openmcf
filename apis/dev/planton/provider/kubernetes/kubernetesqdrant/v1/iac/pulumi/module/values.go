package module

import (
	"fmt"

	"github.com/pkg/errors"
	kubernetesqdrantv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesqdrant/v1"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// fullnameOverride pins qdrant.fullname to the resource name: child
	// names (the Service, the headless Service `-headless`, the chart's
	// `-apikey` Secret, the ConfigMap) derive deterministically and the
	// longest suffix stays far from the 63-char ceiling regardless of how
	// the release name composes with the chart name.
	values["fullnameOverride"] = locals.ReleaseName

	// ---- cluster size -----------------------------------------------------
	// Distributed mode is the chart default (config.cluster.enabled:
	// true); pod 0 bootstraps consensus and later pods join over p2p —
	// scaling is only this count.
	replicas := int32(1)
	if spec.Replicas != nil {
		replicas = spec.GetReplicas()
	}
	values["replicaCount"] = int(replicas)

	// ---- container resources ----------------------------------------------
	// The chart's resources value is the standard PodSpec shape.
	if r := spec.GetResources(); r != nil {
		resources := map[string]interface{}{}
		if q := r.GetRequests(); q != nil && (q.GetCpu() != "" || q.GetMemory() != "") {
			requests := map[string]interface{}{}
			if q.GetCpu() != "" {
				requests["cpu"] = q.GetCpu()
			}
			if q.GetMemory() != "" {
				requests["memory"] = q.GetMemory()
			}
			resources["requests"] = requests
		}
		if l := r.GetLimits(); l != nil && (l.GetCpu() != "" || l.GetMemory() != "") {
			limits := map[string]interface{}{}
			if l.GetCpu() != "" {
				limits["cpu"] = l.GetCpu()
			}
			if l.GetMemory() != "" {
				limits["memory"] = l.GetMemory()
			}
			resources["limits"] = limits
		}
		if len(resources) > 0 {
			values["resources"] = resources
		}
	}

	// ---- storage ------------------------------------------------------------
	// The chart's persistence block always renders a PVC; the module
	// resolves the spec default (10Gi) and sets storageClassName only
	// when declared (absent = the cluster's default class).
	storageSize := spec.GetStorage().GetSize()
	if storageSize == "" {
		storageSize = vars.DefaultStorageSize
	}
	persistence := map[string]interface{}{"size": storageSize}
	if sc := spec.GetStorage().GetStorageClass().GetValue(); sc != "" {
		persistence["storageClassName"] = sc
	}
	values["persistence"] = persistence

	// ---- snapshots volume -----------------------------------------------------
	// Rendered only when declared: a separate volume for snapshots and
	// snapshot shard transfers (upstream sizing guidance: like the main
	// volume).
	if snap := spec.GetSnapshots(); snap != nil {
		snapshotsSize := snap.GetSize()
		if snapshotsSize == "" {
			snapshotsSize = vars.DefaultSnapshotsSize
		}
		snapshotPersistence := map[string]interface{}{
			"enabled": true,
			"size":    snapshotsSize,
		}
		if sc := snap.GetStorageClass().GetValue(); sc != "" {
			snapshotPersistence["storageClassName"] = sc
		}
		values["snapshotPersistence"] = snapshotPersistence
	}

	// ---- api keys ---------------------------------------------------------------
	// The chart owns key materialization either way: the generate arm
	// renders `apiKey: true` (the chart creates a random key ONCE and
	// keeps it stable across upgrades via its lookup); the existing arm
	// renders the chart's valueFrom shape — the chart reads the
	// referenced Secret AT TEMPLATE TIME (it must exist before the
	// install) and copies the key into its own `<name>-apikey` Secret.
	// Key MATERIAL never lands in these values.
	if apiKeyValue := apiKeyValue(spec.GetApiKey()); apiKeyValue != nil {
		values["apiKey"] = apiKeyValue
	}
	if roKeyValue := apiKeyValue(spec.GetReadOnlyApiKey()); roKeyValue != nil {
		values["readOnlyApiKey"] = roKeyValue
	}

	// ---- engine config (TLS) ------------------------------------------------------
	// The chart renders .Values.config verbatim into production.yaml.
	// The typed tls block turns on service.enable_tls and points the
	// engine at the mounted certificate Secret; the chart's own probes
	// switch to HTTPS off the same flag. Inter-node p2p TLS is a separate
	// surface (config.cluster.p2p) and rides helm_values.
	if tlsSecret := spec.GetTls().GetSecret().GetValue(); tlsSecret != "" {
		values["config"] = map[string]interface{}{
			"service": map[string]interface{}{"enable_tls": true},
			"tls": map[string]interface{}{
				"cert": fmt.Sprintf("%s/tls.crt", vars.TlsMountPath),
				"key":  fmt.Sprintf("%s/tls.key", vars.TlsMountPath),
			},
		}
		values["additionalVolumes"] = []interface{}{
			map[string]interface{}{
				"name":   "qdrant-tls",
				"secret": map[string]interface{}{"secretName": tlsSecret},
			},
		}
		values["additionalVolumeMounts"] = []interface{}{
			map[string]interface{}{
				"name":      "qdrant-tls",
				"mountPath": vars.TlsMountPath,
				"readOnly":  true,
			},
		}
	}

	// ---- scheduling -----------------------------------------------------------------
	if sched := spec.GetScheduling(); sched != nil {
		if len(sched.GetNodeSelector()) > 0 {
			values["nodeSelector"] = stringMapToInterface(sched.GetNodeSelector())
		}
		if len(sched.GetTolerations()) > 0 {
			values["tolerations"] = tolerationsSlice(sched.GetTolerations())
		}
		// The chart's affinity value is a raw PodSpec affinity object;
		// the typed flag renders the chart's own documented
		// anti-affinity recipe (spread members across nodes so one node
		// loss takes one member, not the quorum).
		if sched.GetPodAntiAffinity() {
			values["affinity"] = map[string]interface{}{
				"podAntiAffinity": map[string]interface{}{
					"requiredDuringSchedulingIgnoredDuringExecution": []interface{}{
						map[string]interface{}{
							"labelSelector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app.kubernetes.io/name":     vars.HelmChartName,
									"app.kubernetes.io/instance": locals.ReleaseName,
								},
							},
							"topologyKey": "kubernetes.io/hostname",
						},
					},
				},
			}
		}
		if sched.GetPriorityClassName() != "" {
			values["priorityClassName"] = sched.GetPriorityClassName()
		}
	}

	// ---- observability ------------------------------------------------------------------
	if spec.GetServiceMonitorEnabled() {
		values["metrics"] = map[string]interface{}{
			"serviceMonitor": map[string]interface{}{"enabled": true},
		}
	}

	// ---- image --------------------------------------------------------------------------
	// The chart's image.repository carries the registry
	// (docker.io/qdrant/qdrant); useUnprivilegedImage switches to the
	// `qdrant-unprivileged` variant for restricted PSS environments.
	if img := spec.GetImage(); img != nil &&
		(img.GetRepository() != "" || img.GetTag() != "" || img.GetUseUnprivileged()) {
		image := map[string]interface{}{}
		if img.GetRepository() != "" {
			image["repository"] = img.GetRepository()
		}
		if img.GetTag() != "" {
			image["tag"] = img.GetTag()
		}
		if img.GetUseUnprivileged() {
			image["useUnprivilegedImage"] = true
		}
		values["image"] = image
	}

	// ---- escape hatch (merged LAST, helm -f semantics) -----------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). The Service /
	// `-headless` / `-apikey` Secret names — and the exported outputs
	// built from them — all derive from the fullname; letting an override
	// move it would break every output.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}

// apiKeyValue renders one KubernetesQdrantApiKey arm into the chart's
// apiKey value shape: bool true for the generate arm, the valueFrom
// secretKeyRef map for the existing-secret arm, nil when the key is not
// declared.
func apiKeyValue(key *kubernetesqdrantv1.KubernetesQdrantApiKey) interface{} {
	if key == nil {
		return nil
	}
	if existing := key.GetExistingSecret(); existing != nil {
		return map[string]interface{}{
			"valueFrom": map[string]interface{}{
				"secretKeyRef": map[string]interface{}{
					"name": existing.GetName(),
					"key":  existing.GetKey(),
				},
			},
		}
	}
	if key.GetGenerate() {
		return true
	}
	return nil
}
