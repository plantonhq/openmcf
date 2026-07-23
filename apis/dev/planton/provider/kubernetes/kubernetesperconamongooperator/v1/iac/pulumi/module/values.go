package module

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the psmdb-operator chart's
// values map, then merges the spec's helm_values escape hatch over it with
// Helm `-f` semantics (maps deep-merge with the later document winning,
// lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence (telemetry,
// structured logging, and the cluster-wide switch all default off
// upstream), so the rendered values stay minimal on both engines.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{}

	// ---- operator sizing -------------------------------------------------
	if spec.Replicas != nil {
		values["replicaCount"] = int(spec.GetReplicas())
	}
	if r := resourcesMap(spec.GetResources()); r != nil {
		// The chart ships no default requests/limits — this is purely
		// additive.
		values["resources"] = r
	}

	// ---- watch scope -----------------------------------------------------
	// The chart models the scope with two independent values:
	// watchAllNamespaces (cluster-wide RBAC, default false) and
	// watchNamespace (a comma-separated namespace list; blank = the
	// operator's own namespace). Spec CEL rules make the two spec arms
	// mutually exclusive, so at most one of the values renders.
	//
	// The chart's own `createNamespace` value (create the WATCHED
	// namespaces) is never rendered: this module owns only the
	// installation namespace; watched namespaces must already exist.
	if spec.GetWatch().GetClusterWide() {
		values["watchAllNamespaces"] = true
	}
	if len(spec.GetWatch().GetNamespaces()) > 0 {
		values["watchNamespace"] = strings.Join(spec.GetWatch().GetNamespaces(), ",")
	}

	// ---- reconciliation throughput ----------------------------------------
	// Rendered as a STRING to match the chart's own declared type (its
	// default is the string "1"); the deployment template quotes the value
	// into the MAX_CONCURRENT_RECONCILES environment variable either way,
	// and the string keeps both engines byte-identical with the chart's
	// values file.
	if spec.MaxConcurrentReconciles != nil {
		values["maxConcurrentReconciles"] = strconv.Itoa(int(spec.GetMaxConcurrentReconciles()))
	}

	// ---- logging -----------------------------------------------------------
	// logStructured matches the chart default (false) — rendered only when
	// on. logLevel renders whenever the spec carries a value (the chart
	// default is "INFO"; re-stating it is harmless and keeps rendering
	// presence-driven, not value-driven).
	if spec.GetLog().GetStructured() {
		values["logStructured"] = true
	}
	if spec.GetLog() != nil && spec.GetLog().Level != nil && spec.GetLog().GetLevel() != "" {
		values["logLevel"] = spec.GetLog().GetLevel()
	}

	// ---- telemetry ----------------------------------------------------------
	// Chart default false (telemetry on) — rendered only on explicit
	// opt-out.
	if spec.GetDisableTelemetry() {
		values["disableTelemetry"] = true
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// ---- image ----------------------------------------------------------------
	// Pull secrets render as the raw Kubernetes object list
	// ([{name: ...}]) — the chart pipes imagePullSecrets straight into the
	// pod spec with toYaml. The image override renders only the halves
	// that are set: the chart composes "<repository>:<tag>" itself, so an
	// unset half keeps the chart's default for it (repository
	// percona/percona-server-mongodb-operator; tag = the chart version).
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}
	image := map[string]interface{}{}
	if spec.GetImage().GetRepository() != "" {
		image["repository"] = spec.GetImage().GetRepository()
	}
	if spec.GetImage().GetTag() != "" {
		image["tag"] = spec.GetImage().GetTag()
	}
	if len(image) > 0 {
		values["image"] = image
	}

	// ---- escape hatch (merged LAST, helm -f semantics) --------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	return values, nil
}
