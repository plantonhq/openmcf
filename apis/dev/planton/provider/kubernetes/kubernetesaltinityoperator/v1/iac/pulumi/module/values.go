package module

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the
// altinity-clickhouse-operator chart's values map, then merges the spec's
// helm_values escape hatch over it with Helm `-f` semantics (maps
// deep-merge with the later document winning, lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence, so the rendered
// values stay minimal on both engines. CRDs need no handling here: the
// chart ships them in its crds/ directory (Helm-native keep-on-uninstall)
// and its pre-install/pre-upgrade hook job carries schema upgrades —
// unlike sibling operator modules that must own CRDs because their charts
// template them release-owned.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		// fullnameOverride pins the chart's fullname to the resource name
		// (the catalog's Helm-kind identity convention). Every generated
		// child name and every exported output (Deployment, credentials
		// Secret, "<fullname>-metrics" Service) hangs off it.
		"fullnameOverride": locals.ReleaseName,
	}

	// ---- watch scope -------------------------------------------------------
	// Empty = the chart default (the operator watches only its own
	// namespace), so the key renders only when the spec widens the scope.
	if len(spec.GetWatchNamespaces()) > 0 {
		watchNamespaces := make([]interface{}, 0, len(spec.GetWatchNamespaces()))
		for _, ns := range spec.GetWatchNamespaces() {
			watchNamespaces = append(watchNamespaces, ns)
		}
		values["watchNamespaces"] = watchNamespaces
	}

	// ---- RBAC scope ----------------------------------------------------------
	// Plain bool (no presence): false IS the chart default (cluster-wide
	// RBAC), so only true renders.
	if spec.GetNamespaceScopedRbac() {
		values["rbac"] = map[string]interface{}{
			"namespaceScoped": true,
		}
	}

	// ---- operator credentials ------------------------------------------------
	// Rendered only when the message is present — an absent message keeps
	// the chart's (publicly documented) defaults, which the spec flags as
	// unsafe outside throwaway environments. The username resolves to the
	// spec default so the rendered pair is always complete.
	if creds := spec.GetOperatorCredentials(); creds != nil {
		username := creds.GetUsername()
		if username == "" {
			username = vars.DefaultCredentialsUsername
		}
		values["secret"] = map[string]interface{}{
			"username": username,
			"password": creds.GetPassword().GetValue(),
		}
	}

	// ---- metrics exporter ------------------------------------------------------
	// enabled renders on presence — an explicit true re-states the chart
	// default harmlessly, an explicit false is the actual opt-out. The
	// chart ships no default requests/limits for the sidecar, so resources
	// render only when set.
	if m := spec.GetMetrics(); m != nil {
		metrics := map[string]interface{}{}
		if m.Enabled != nil {
			metrics["enabled"] = m.GetEnabled()
		}
		if r := resourcesMap(m.GetResources()); r != nil {
			metrics["resources"] = r
		}
		if len(metrics) > 0 {
			values["metrics"] = metrics
		}
	}

	// ---- CRD hook -----------------------------------------------------------
	// enabled renders on presence (true-defaulted upstream; an explicit
	// false stops upgrades from carrying CRD schema changes). The image
	// override deep-merges per half over the chart default
	// (bitnami/kubectl:latest), leaving registry and pullPolicy intact —
	// any registry prefix belongs in repository.
	if hook := spec.GetCrdHook(); hook != nil {
		crdHook := map[string]interface{}{}
		if hook.Enabled != nil {
			crdHook["enabled"] = hook.GetEnabled()
		}
		hookImage := map[string]interface{}{}
		if hook.GetImage().GetRepo() != "" {
			hookImage["repository"] = hook.GetImage().GetRepo()
		}
		if hook.GetImage().GetTag() != "" {
			hookImage["tag"] = hook.GetImage().GetTag()
		}
		if len(hookImage) > 0 {
			crdHook["image"] = hookImage
		}
		if len(crdHook) > 0 {
			values["crdHook"] = crdHook
		}
	}

	// ---- operator container ---------------------------------------------------
	// The chart ships no default requests/limits, so resources render only
	// when set. The image override deep-merges per half over the chart
	// defaults (altinity/clickhouse-operator at the chart's appVersion),
	// leaving registry and pullPolicy intact.
	operator := map[string]interface{}{}
	if r := resourcesMap(spec.GetResources()); r != nil {
		operator["resources"] = r
	}
	operatorImage := map[string]interface{}{}
	if spec.GetImage().GetRepo() != "" {
		operatorImage["repository"] = spec.GetImage().GetRepo()
	}
	if spec.GetImage().GetTag() != "" {
		operatorImage["tag"] = spec.GetImage().GetTag()
	}
	if len(operatorImage) > 0 {
		operator["image"] = operatorImage
	}
	if len(operator) > 0 {
		values["operator"] = operator
	}

	// ---- ServiceMonitor ----------------------------------------------------------
	// Plain bool (no presence): false IS the chart default, so only true
	// renders (and requires the Prometheus Operator CRDs on the cluster).
	if spec.GetServiceMonitorEnabled() {
		values["serviceMonitor"] = map[string]interface{}{
			"enabled": true,
		}
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// Pull secrets ride the chart's top-level imagePullSecrets list (raw
	// Kubernetes object list, piped into the operator pod spec).
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// ---- escape hatch (merged LAST, helm -f semantics) --------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// fullnameOverride is re-pinned AFTER the escape-hatch merge — the one
	// deliberate exception to helm -f semantics. Overriding it would break
	// the naming budget and every fullname-derived output.
	values["fullnameOverride"] = locals.ReleaseName

	return values, nil
}
