package module

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the opensearch-operator
// chart's values map, then merges the spec's helm_values escape hatch over
// it with Helm `-f` semantics (maps deep-merge with the later document
// winning, lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values] and
// the provider merges the documents in exactly this order. Keep every typed
// mapping below in lockstep with the Terraform module's locals.
//
// Chart-default-matching values render only on divergence (the watch
// scope, the true-defaulted toggles, the image source), so the rendered
// values stay minimal on both engines — with ONE deliberate exception:
// installCRDs is pinned false unconditionally (see below).
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	// installCRDs: false ALWAYS — never conditional, never overridable by
	// design. The chart templates its CRDs as release-owned resources with
	// NO keep-on-uninstall knob, so a Helm-owned install would
	// cascade-delete every OpenSearchCluster (and its data) on uninstall.
	// The module owns the CRD lifecycle instead (derived from the pinned
	// chart and applied kept through keptcrds). The same map is what the
	// CRD render sees (plus the switch turned on), so the derived CRDs can
	// never see different values than the install.
	values := map[string]interface{}{
		"installCRDs": false,
		// fullnameOverride pins the chart's fullname to the resource name
		// (the catalog's Helm-kind identity convention). Load-bearing here:
		// the chart's default fullname ("<release>-opensearch-operator")
		// pushes its metrics Service name past Kubernetes' 63-character
		// limit for ordinary release names — verified live (the chart
		// truncates the fullname but not the names built from it). With
		// the pin, resource names up to 27 characters stay within budget.
		"fullnameOverride": locals.ReleaseName,
	}

	// ---- manager (the operator controller) -------------------------------
	// Every manager.* key renders only when the spec sets it, so the chart
	// defaults survive an empty spec. Helm deep-merges per key, so a
	// partial manager block overrides only the halves it carries.
	manager := map[string]interface{}{}
	if spec.GetWatchNamespace() != "" {
		manager["watchNamespace"] = spec.GetWatchNamespace()
	}
	// loglevel (the chart's lowercase key) renders whenever the spec
	// carries a value; the chart default is "info".
	if spec.LogLevel != nil && spec.GetLogLevel() != "" {
		manager["loglevel"] = spec.GetLogLevel()
	}
	if spec.DnsBase != nil && spec.GetDnsBase() != "" {
		manager["dnsBase"] = spec.GetDnsBase()
	}
	// Rendered on presence — an explicit true re-states the chart default
	// harmlessly, an explicit false is the actual opt-out.
	if spec.ParallelRecoveryEnabled != nil {
		manager["parallelRecoveryEnabled"] = spec.GetParallelRecoveryEnabled()
	}
	// Plain bool (no presence): false IS the chart default, so only true
	// renders.
	if spec.GetPprofEndpointsEnabled() {
		manager["pprofEndpointsEnabled"] = true
	}
	// The chart SHIPS default requests/limits (requests 100m/350Mi,
	// limits 200m/500Mi) — the resources key renders only when the spec
	// sets them.
	if r := resourcesMap(spec.GetResources()); r != nil {
		manager["resources"] = r
	}
	// Pull secrets ride the chart's manager.imagePullSecrets list (raw
	// Kubernetes object list, piped into the pod spec).
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		manager["imagePullSecrets"] = pullSecrets
	}
	// The image override renders per-half: repository and tag deep-merge
	// over the chart defaults (opensearchproject/opensearch-operator at
	// the chart's appVersion), leaving pullPolicy intact.
	image := map[string]interface{}{}
	if spec.GetImage().GetRepository() != "" {
		image["repository"] = spec.GetImage().GetRepository()
	}
	if spec.GetImage().GetTag() != "" {
		image["tag"] = spec.GetImage().GetTag()
	}
	if len(image) > 0 {
		manager["image"] = image
	}
	if len(manager) > 0 {
		values["manager"] = manager
	}

	// ---- kube-rbac-proxy sidecar ------------------------------------------
	// The chart nests a single enable flag. Rendered on presence — an
	// explicit true re-states the chart default harmlessly, an explicit
	// false is the actual opt-out.
	// The sidecar's image repository is ALWAYS re-pointed at the
	// maintainer's own quay.io repository (same tag as the chart pins):
	// the chart's default, gcr.io/kubebuilder/kube-rbac-proxy, was
	// DELETED upstream (verified live — the registry returns not-found),
	// so a default-posture install can never pull it and times out
	// waiting on a 1/2-ready Deployment. Overridable via helm_values.
	kubeRbacProxy := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": "quay.io/brancz/kube-rbac-proxy",
		},
	}
	if spec.KubeRbacProxyEnabled != nil {
		kubeRbacProxy["enable"] = spec.GetKubeRbacProxyEnabled()
	}
	values["kubeRbacProxy"] = kubeRbacProxy

	// ---- RBAC scope ---------------------------------------------------------
	// Plain bool (no presence): false IS the chart default
	// (ClusterRoleBindings), so only true renders. Spec CEL requires
	// watch_namespace alongside it.
	if spec.GetUseRoleBindings() {
		values["useRoleBindings"] = true
	}

	// ---- scheduling ----------------------------------------------------------
	if len(spec.GetNodeSelector()) > 0 {
		values["nodeSelector"] = stringMapToInterface(spec.GetNodeSelector())
	}
	if len(spec.GetTolerations()) > 0 {
		values["tolerations"] = tolerationsSlice(spec.GetTolerations())
	}

	// ---- escape hatch (merged LAST, helm -f semantics) --------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
		}
		values = mergeMaps(values, overrides)
	}

	// installCRDs is re-pinned AFTER the escape-hatch merge — the one
	// deliberate exception to helm -f semantics. The module owns the CRD
	// lifecycle; letting an override hand them to Helm would arm the
	// cascade-delete this design exists to prevent (and would conflict
	// with the module-applied CRDs at install anyway).
	values["installCRDs"] = false

	return values, nil
}
