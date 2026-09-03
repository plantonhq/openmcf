package module

import (
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the opentelemetry-operator
// chart's values map, then merges the spec's helm_values escape hatch over
// it with Helm `-f` semantics (maps deep-merge with the later document
// winning, lists replace).
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(re-pins)] and the provider merges the documents in exactly
// this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
//
// Chart-default-matching values render only on divergence, so the
// rendered values stay minimal on both engines — with TWO deliberate
// exceptions re-pinned AFTER the escape-hatch merge (see the bottom of
// this function).
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		// crds.create: false ALWAYS — never conditional, never a spec
		// knob. The chart templates its CRDs as release-owned resources
		// (a Helm uninstall would cascade-delete every collector in the
		// cluster); the module owns the CRDs instead (main.go derives
		// them from the pinned chart and applies them kept through
		// keptcrds).
		"crds": map[string]interface{}{"create": false},
		// fullnameOverride pins the chart's fullname to the resource
		// name (the catalog's Helm-kind identity convention).
		// Load-bearing here: the derived CRDs' conversion webhook and
		// inject-ca-from annotation point at names derived from it, and
		// the CRD render runs with these exact values.
		"fullnameOverride": locals.ReleaseName,
	}

	// Rendered on presence — an explicit 1 re-states the chart default
	// harmlessly; 2 gives a warm standby behind leader election.
	if spec.Replicas != nil {
		values["replicaCount"] = spec.GetReplicas()
	}

	// ---- manager (the operator controller) -------------------------------
	// Every manager.* key renders only when the spec sets it, so the chart
	// defaults survive an empty spec. Helm deep-merges per key.
	manager := map[string]interface{}{}
	// The chart ships NO default requests/limits for the manager — the
	// resources key renders only when the spec sets them.
	if r := resourcesMap(spec.GetResources()); r != nil {
		manager["resources"] = r
	}
	// image_registry replaces ONLY the registry part of the manager image
	// (the one image this component's pods pull); the path stays the
	// upstream one. The default collector image the operator INJECTS into
	// CRs is mirrored via default_collector_image instead — collector
	// pods pull that one, not this component.
	if spec.GetImageRegistry() != "" {
		manager["image"] = map[string]interface{}{
			"repository": spec.GetImageRegistry() + "/" + vars.ManagerImagePath,
		}
	}
	// The spec carries ONE image string; the chart takes repository and
	// tag separately (and renders --collector-image only when BOTH are
	// present — a repository-only override deep-merges with the chart's
	// default tag, so the flag still renders).
	if spec.GetDefaultCollectorImage() != "" {
		repo, tag := splitImageRef(spec.GetDefaultCollectorImage())
		collectorImage := map[string]interface{}{"repository": repo}
		if tag != "" {
			collectorImage["tag"] = tag
		}
		manager["collectorImage"] = collectorImage
	}
	// Plain bool (no presence): false IS the chart default, so only true
	// renders. Requires the monitoring.coreos.com CRDs on the cluster
	// (KubernetesKubePrometheusStack).
	if spec.GetServiceMonitorEnabled() {
		manager["serviceMonitor"] = map[string]interface{}{"enabled": true}
	}
	if len(manager) > 0 {
		values["manager"] = manager
	}

	// ---- cert-manager issuer reference ------------------------------------
	// Rendered only when the spec names an issuer; empty means the chart
	// creates its own self-signed Issuer (the default posture).
	// cert-manager itself is NOT optional — see the re-pin below.
	if ref := spec.GetWebhook().GetIssuerRef(); ref != nil && (ref.GetKind() != "" || ref.GetName() != "") {
		issuerRef := map[string]interface{}{}
		if ref.GetKind() != "" {
			issuerRef["kind"] = ref.GetKind()
		}
		if ref.GetName() != "" {
			issuerRef["name"] = ref.GetName()
		}
		values["admissionWebhooks"] = map[string]interface{}{
			"certManager": map[string]interface{}{"issuerRef": issuerRef},
		}
	}

	// ---- pull secrets (TOP-LEVEL in this chart) -----------------------------
	// Raw Kubernetes object list, piped into the pod spec — verified in
	// the deployment template; this chart does NOT nest them under
	// manager.
	if len(spec.GetImagePullSecrets()) > 0 {
		pullSecrets := make([]interface{}, 0, len(spec.GetImagePullSecrets()))
		for _, name := range spec.GetImagePullSecrets() {
			pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
		}
		values["imagePullSecrets"] = pullSecrets
	}

	// ---- scheduling (TOP-LEVEL keys in this chart) ---------------------------
	// nodeSelector deep-merges over the chart's default
	// {kubernetes.io/os: linux}, so the OS pin survives a spec selector.
	if sched := spec.GetScheduling(); sched != nil {
		if len(sched.GetNodeSelector()) > 0 {
			values["nodeSelector"] = stringMapToInterface(sched.GetNodeSelector())
		}
		if len(sched.GetTolerations()) > 0 {
			values["tolerations"] = tolerationsSlice(sched.GetTolerations())
		}
		if sched.GetPriorityClassName() != "" {
			values["priorityClassName"] = sched.GetPriorityClassName()
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

	// The TWO design-load-bearing keys are re-pinned AFTER the
	// escape-hatch merge — the deliberate exceptions to helm -f
	// semantics (twin of the Terraform module's third values document):
	//   - crds.create=false: the module owns the CRD lifecycle; handing
	//     them to Helm would arm the uninstall cascade-delete this
	//     design exists to prevent.
	//   - admissionWebhooks.certManager.enabled=true: the kept CRDs'
	//     conversion trust rides cert-manager's CA injector; disabling
	//     it would leave module-owned CRDs pointing at a Certificate
	//     that no longer exists and silently break collector-CR
	//     conversion.
	if crds, ok := values["crds"].(map[string]interface{}); ok {
		crds["create"] = false
	} else {
		values["crds"] = map[string]interface{}{"create": false}
	}
	admissionWebhooks, ok := values["admissionWebhooks"].(map[string]interface{})
	if !ok {
		admissionWebhooks = map[string]interface{}{}
		values["admissionWebhooks"] = admissionWebhooks
	}
	certManager, ok := admissionWebhooks["certManager"].(map[string]interface{})
	if !ok {
		certManager = map[string]interface{}{}
		admissionWebhooks["certManager"] = certManager
	}
	certManager["enabled"] = true

	return values, nil
}

// splitImageRef splits an image reference into repository and tag. A tag
// exists when the LAST ":" comes after the last "/" (registry ports carry
// ":" too — "reg.example.com:5000/x" has no tag). Twin of the Terraform
// module's collector_image locals — keep the semantics identical.
func splitImageRef(ref string) (repository string, tag string) {
	lastColon := strings.LastIndex(ref, ":")
	lastSlash := strings.LastIndex(ref, "/")
	if lastColon > lastSlash {
		return ref[:lastColon], ref[lastColon+1:]
	}
	return ref, ""
}
