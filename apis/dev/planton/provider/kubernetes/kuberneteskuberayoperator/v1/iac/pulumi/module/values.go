package module

import (
	"github.com/pkg/errors"
	kuberneteskuberayoperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskuberayoperator/v1"
	"sigs.k8s.io/yaml"
)

// chartDefaultFeatureGates is the featureGates list the pinned 1.6.2 chart
// ships in its values.yaml. Helm LISTS REPLACE, never merge: rendering only
// the spec's entries would silently DROP every chart-default gate, so any
// spec override renders the FULL list seeded from these. Keep in lockstep
// with the Terraform module's feature_gate_defaults AND re-verify on every
// chart bump.
var chartDefaultFeatureGates = []struct {
	name    string
	enabled bool
}{
	{"RayClusterStatusConditions", true},
	{"RayJobDeletionPolicy", true},
	{"RayMultiHostIndexing", true},
	{"RayServiceIncrementalUpgrade", false},
	{"RayCronJob", false},
}

// buildHelmValues renders the typed spec into the kuberay-operator chart's
// values map, then merges the spec's helm_values escape hatch over it with
// Helm `-f` semantics (maps deep-merge with the later document winning,
// lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every
// typed mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence — with ONE
// deliberate always-rendered set: the NAME PINS. The chart hardcodes
// nameOverride, fullnameOverride, AND serviceAccount.name to
// "kuberay-operator" in its values, so every install collapses onto the
// same child names (and the same ServiceAccount) by construction.
// Pinning all three to metadata.name keeps instances distinguishable.
// This chart needs no post-escape-hatch re-pin document: it has no
// release-owned CRDs and no webhook machinery whose keys an escape-hatch
// value could weaponize — an operator deliberately overriding the name
// pins owns the collision consciously.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		"nameOverride":     locals.ReleaseName,
		"fullnameOverride": locals.ReleaseName,
		"serviceAccount":   map[string]interface{}{"name": locals.ReleaseName},
	}

	// Chart key is singular ("watchNamespace") but takes a LIST — it
	// feeds the operator's --watch-namespace flag and scopes the
	// per-namespace reconcile RBAC. Unrendered = cluster-wide (the chart
	// default).
	if len(locals.WatchNamespaces) > 0 {
		watch := make([]interface{}, 0, len(locals.WatchNamespaces))
		for _, ns := range locals.WatchNamespaces {
			watch = append(watch, ns)
		}
		values["watchNamespace"] = watch
	}

	// Explicit false only — an unset optional keeps the chart default
	// (true: safe for single replicas, required for standbys).
	if spec.LeaderElectionEnabled != nil && !spec.GetLeaderElectionEnabled() {
		values["leaderElectionEnabled"] = false
	}

	// batchScheduler.name is the standard knob; batchScheduler.enabled is
	// the deprecated legacy flag and MUTUALLY EXCLUSIVE with name (the
	// chart errors when both are set) — never render it.
	if spec.GetBatchScheduler() != "" {
		values["batchScheduler"] = map[string]interface{}{"name": spec.GetBatchScheduler()}
	}

	if len(spec.GetFeatureGates()) > 0 {
		values["featureGates"] = featureGatesList(spec.GetFeatureGates())
	}

	// metrics.enabled renders only on an EXPLICIT false (proto optional
	// bool; chart default true). The ServiceMonitor requires the
	// monitoring.coreos.com CRDs on the cluster
	// (KubernetesKubePrometheusStack) — the install FAILS without them.
	metrics := map[string]interface{}{}
	if spec.MetricsEnabled != nil && !spec.GetMetricsEnabled() {
		metrics["enabled"] = false
	}
	if spec.GetServiceMonitorEnabled() {
		metrics["serviceMonitor"] = map[string]interface{}{"enabled": true}
	}
	if len(metrics) > 0 {
		values["metrics"] = metrics
	}

	// The chart ships REAL defaults here (100m CPU / 512Mi limits —
	// upstream sizes ~500MB per 500 managed Ray pods) — the resources key
	// renders only when the spec sets them, so the upstream-tested sizing
	// stands otherwise. Helm deep-merges per key: a partial block
	// overrides only the halves it carries.
	if r := resourcesMap(spec.GetResources()); r != nil {
		values["resources"] = r
	}

	// image_registry replaces ONLY the registry part of the operator
	// image (chart default quay.io/kuberay/operator — the swap drops
	// quay.io); the tag stays the chart's appVersion-locked default. Ray
	// CLUSTER images ride each KubernetesRayCluster's own image field —
	// this never rewrites those.
	if spec.GetImageRegistry() != "" {
		values["image"] = map[string]interface{}{
			"repository": spec.GetImageRegistry() + "/" + vars.OperatorImagePath,
		}
	}

	// ---- escape hatch (merged LAST, helm -f semantics) -----------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}

// featureGatesList renders the FULL featureGates list: the chart defaults
// overridden by name from the spec entries, then spec gates the defaults
// don't know appended in spec order. Twin of the Terraform module's
// feature_gates local.
func featureGatesList(gates []*kuberneteskuberayoperatorv1.KubernetesKubeRayOperatorFeatureGate) []interface{} {
	specGateByName := make(map[string]bool, len(gates))
	for _, g := range gates {
		specGateByName[g.GetName()] = g.GetEnabled()
	}

	defaultNames := make(map[string]bool, len(chartDefaultFeatureGates))
	out := make([]interface{}, 0, len(chartDefaultFeatureGates)+len(gates))
	for _, d := range chartDefaultFeatureGates {
		defaultNames[d.name] = true
		enabled := d.enabled
		if specEnabled, ok := specGateByName[d.name]; ok {
			enabled = specEnabled
		}
		out = append(out, map[string]interface{}{"name": d.name, "enabled": enabled})
	}
	for _, g := range gates {
		if defaultNames[g.GetName()] {
			continue
		}
		out = append(out, map[string]interface{}{"name": g.GetName(), "enabled": g.GetEnabled()})
	}
	return out
}
