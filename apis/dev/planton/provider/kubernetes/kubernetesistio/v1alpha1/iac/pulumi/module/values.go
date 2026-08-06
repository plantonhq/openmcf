package module

import (
	"github.com/pkg/errors"
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesistiov1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesistio/v1alpha1"
	"sigs.k8s.io/yaml"
)

// This file renders the typed spec into per-release Helm values. Every typed
// mapping is verified against the pinned upstream chart values
// (manifests/charts/{base,istio-control/istio-discovery,istio-cni,ztunnel}).
//
// PARITY: the Terraform module reaches the same result natively — each
// helm_release passes values = [yamlencode(typed values), helm_values.<release>]
// and the provider merges the documents in that order (Helm -f semantics).
// Keep every typed mapping below in lockstep with the Terraform module's
// locals.

// buildBaseValues renders values for the `base` release (validation webhook
// plumbing; the CRDs are module-owned and excluded from the release — see
// main.go).
func buildBaseValues(locals *Locals) (map[string]interface{}, error) {
	values := map[string]interface{}{
		// The default-revision validating webhook must point at THIS control
		// plane's revision ("default" for the unnamed revision).
		"defaultRevision": locals.Revision,
		// Exclude the ENTIRE CRD bundle from the release: the module applies
		// the CRDs itself via server-side apply so they are never Helm-owned
		// (Helm cannot adopt pre-existing CRDs — the KubernetesIstioBaseCrds
		// upgrade path).
		"base": map[string]interface{}{
			"excludedCRDs": toInterfaceSlice(vars.CrdNames),
		},
	}
	if locals.Namespace != "istio-system" {
		setNested(values, map[string]interface{}{"istioNamespace": locals.Namespace}, "global")
	}
	return mergeEscapeHatch(values, locals.Spec.GetHelmValues().GetBase())
}

// buildIstiodValues renders values for the `istiod` release — the main chart.
func buildIstiodValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	// Ambient mode rides the chart's own profile overlay (the upstream
	// install path: --set profile=ambient), which enables HBONE, the ambient
	// pilot env, and the distroless variant defaults in one place.
	if locals.Ambient {
		values["profile"] = "ambient"
	}

	if spec.GetRevision() != "" {
		values["revision"] = spec.GetRevision()
	}

	if locals.Namespace != "istio-system" {
		setNested(values, map[string]interface{}{"istioNamespace": locals.Namespace}, "global")
	}

	// ---- istiod sizing / scheduling -----------------------------------
	if istiod := spec.GetIstiod(); istiod != nil {
		if istiod.Replicas != nil {
			values["replicaCount"] = int(istiod.GetReplicas())
			// A fixed replica count only holds if the chart's HPA is off
			// (the spec forbids replicas+autoscale, so this is safe).
			values["autoscaleEnabled"] = false
		}
		if a := istiod.GetAutoscale(); a != nil {
			if a.Enabled != nil {
				values["autoscaleEnabled"] = a.GetEnabled()
			}
			if a.MinReplicas != nil {
				values["autoscaleMin"] = int(a.GetMinReplicas())
			}
			if a.MaxReplicas != nil {
				values["autoscaleMax"] = int(a.GetMaxReplicas())
			}
			if a.TargetCpuUtilizationPercent != nil {
				setNested(values, map[string]interface{}{"targetAverageUtilization": int(a.GetTargetCpuUtilizationPercent())}, "cpu")
			}
		}
		if r := resourcesMap(istiod.GetResources()); r != nil {
			values["resources"] = r
		}
		if istiod.LogLevel != nil && istiod.GetLogLevel() != "default:info" {
			setNested(values, map[string]interface{}{"level": istiod.GetLogLevel()}, "global", "logging")
		}
		if istiod.PodDisruptionBudget != nil {
			setNested(values, map[string]interface{}{"enabled": istiod.GetPodDisruptionBudget()}, "global", "defaultPodDisruptionBudget")
		}
		if istiod.GetPriorityClassName() != "" {
			setNested(values, map[string]interface{}{"priorityClassName": istiod.GetPriorityClassName()}, "global")
		}
		if len(istiod.GetNodeSelector()) > 0 {
			values["nodeSelector"] = stringMapToInterface(istiod.GetNodeSelector())
		}
		if tolerations := tolerationsList(istiod.GetTolerations()); tolerations != nil {
			values["tolerations"] = tolerations
		}
	}

	// ---- mesh config ----------------------------------------------------
	if mc := spec.GetMeshConfig(); mc != nil {
		meshConfig := map[string]interface{}{}
		if mc.TrustDomain != nil && mc.GetTrustDomain() != "cluster.local" {
			meshConfig["trustDomain"] = mc.GetTrustDomain()
		}
		if mc.OutboundTrafficPolicyMode != nil {
			meshConfig["outboundTrafficPolicy"] = map[string]interface{}{"mode": mc.GetOutboundTrafficPolicyMode()}
		}
		if mc.GetAccessLogFile() != "" {
			meshConfig["accessLogFile"] = mc.GetAccessLogFile()
		}
		if mc.EnablePrometheusMerge != nil {
			meshConfig["enablePrometheusMerge"] = mc.GetEnablePrometheusMerge()
		}
		if len(meshConfig) > 0 {
			values["meshConfig"] = meshConfig
		}
		if mc.GetClusterName() != "" {
			setNested(values, map[string]interface{}{"clusterName": mc.GetClusterName()}, "global", "multiCluster")
		}
		if mc.GetNetwork() != "" {
			setNested(values, map[string]interface{}{"network": mc.GetNetwork()}, "global")
		}
		if mc.GetMeshId() != "" {
			setNested(values, map[string]interface{}{"meshID": mc.GetMeshId()}, "global")
		}
	}

	// ---- proxy defaults --------------------------------------------------
	if proxy := spec.GetProxy(); proxy != nil {
		proxyValues := map[string]interface{}{}
		if r := resourcesMap(proxy.GetResources()); r != nil {
			proxyValues["resources"] = r
		}
		if proxy.LogLevel != nil {
			proxyValues["logLevel"] = proxy.GetLogLevel()
		}
		if proxy.AutoInject != nil {
			proxyValues["autoInject"] = proxy.GetAutoInject()
		}
		if proxy.ClusterDomain != nil && proxy.GetClusterDomain() != "cluster.local" {
			proxyValues["clusterDomain"] = proxy.GetClusterDomain()
		}
		if len(proxyValues) > 0 {
			setNested(values, proxyValues, "global", "proxy")
		}
	}

	// ---- sidecar injection ------------------------------------------------
	if si := spec.GetSidecarInjector(); si != nil {
		injector := map[string]interface{}{}
		if si.GetEnableNamespacesByDefault() {
			injector["enableNamespacesByDefault"] = true
		}
		if si.RewriteAppHttpProbe != nil {
			injector["rewriteAppHTTPProbe"] = si.GetRewriteAppHttpProbe()
		}
		if len(injector) > 0 {
			values["sidecarInjectorWebhook"] = injector
		}
	}

	// In sidecar mode with the node-level CNI agent installed, istiod's
	// injector must emit pods that rely on it instead of the privileged
	// istio-init container. (Ambient's profile overlay handles this itself.)
	if !locals.Ambient && locals.InstallCni {
		values["cni"] = map[string]interface{}{"enabled": true}
	}

	// ---- gateway-class defaults --------------------------------------------
	// Renders the per-GatewayClass deployment overlay ConfigMap the chart
	// ships for istiod's Gateway API auto-provisioning.
	if gd := spec.GetGatewayDefaults(); gd != nil && gd.ServiceType != nil {
		values["gatewayClasses"] = map[string]interface{}{
			vars.GatewayClassName: map[string]interface{}{
				"service": map[string]interface{}{
					"spec": map[string]interface{}{"type": gd.GetServiceType()},
				},
			},
		}
	}

	applyImageValues(values, spec.GetImages())

	return mergeEscapeHatch(values, spec.GetHelmValues().GetIstiod())
}

// buildCniValues renders values for the `cni` release (the istio-cni node
// agent DaemonSet).
func buildCniValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	if locals.Ambient {
		values["profile"] = "ambient"
	}
	if spec.GetRevision() != "" {
		values["revision"] = spec.GetRevision()
	}

	if cni := spec.GetCni(); cni != nil {
		if len(cni.GetExcludeNamespaces()) > 0 {
			values["excludeNamespaces"] = toInterfaceSlice(cni.GetExcludeNamespaces())
		}
		if cni.CniBinDir != nil {
			values["cniBinDir"] = cni.GetCniBinDir()
		}
		if cni.CniConfDir != nil {
			values["cniConfDir"] = cni.GetCniConfDir()
		}
		if cni.Chained != nil {
			values["chained"] = cni.GetChained()
		}
	}

	applyImageValues(values, spec.GetImages())

	return mergeEscapeHatch(values, spec.GetHelmValues().GetCni())
}

// buildZtunnelValues renders values for the `ztunnel` release (the ambient
// per-node L4 proxy DaemonSet).
func buildZtunnelValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{}

	if spec.GetRevision() != "" {
		values["revision"] = spec.GetRevision()
	}
	if locals.Namespace != "istio-system" {
		values["istioNamespace"] = locals.Namespace
	}

	if zt := spec.GetZtunnel(); zt != nil {
		if r := resourcesMap(zt.GetResources()); r != nil {
			values["resources"] = r
		}
		if zt.LogLevel != nil {
			values["logLevel"] = zt.GetLogLevel()
		}
	}

	// The ztunnel chart reads hub/variant/imagePullSecrets at the top level,
	// unlike the global.* convention of the other three charts.
	if images := spec.GetImages(); images != nil {
		if images.GetHub() != "" {
			values["hub"] = images.GetHub()
		}
		if images.Variant != nil {
			values["variant"] = images.GetVariant()
		}
		if len(images.GetImagePullSecrets()) > 0 {
			values["imagePullSecrets"] = toInterfaceSlice(images.GetImagePullSecrets())
		}
	}

	return mergeEscapeHatch(values, spec.GetHelmValues().GetZtunnel())
}

// applyImageValues writes the shared image-source knobs into a chart's
// global.* values (base/istiod/cni convention; ztunnel reads them top-level).
func applyImageValues(values map[string]interface{}, images *kubernetesistiov1alpha1.KubernetesIstioImages) {
	if images == nil {
		return
	}
	globalValues := map[string]interface{}{}
	if images.GetHub() != "" {
		globalValues["hub"] = images.GetHub()
	}
	if images.Variant != nil {
		globalValues["variant"] = images.GetVariant()
	}
	if len(images.GetImagePullSecrets()) > 0 {
		globalValues["imagePullSecrets"] = toInterfaceSlice(images.GetImagePullSecrets())
	}
	if len(globalValues) > 0 {
		setNested(values, globalValues, "global")
	}
}

// mergeEscapeHatch parses a per-release helm_values YAML document and merges
// it LAST over the typed values (Helm -f semantics: nested maps deep-merge
// with the override winning, lists and scalars replace).
func mergeEscapeHatch(values map[string]interface{}, escapeHatch string) (map[string]interface{}, error) {
	if escapeHatch == "" {
		return values, nil
	}
	overrides := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(escapeHatch), &overrides); err != nil {
		return nil, errors.Wrap(err, "failed to parse helm_values as a YAML document")
	}
	return mergeMaps(values, overrides), nil
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

// tolerationsList renders the shared WorkloadToleration list into the chart's
// tolerations shape. Returns nil when empty.
func tolerationsList(tolerations []*kubernetesprovider.WorkloadToleration) []interface{} {
	if len(tolerations) == 0 {
		return nil
	}
	out := make([]interface{}, 0, len(tolerations))
	for _, t := range tolerations {
		entry := map[string]interface{}{}
		if t.GetKey() != "" {
			entry["key"] = t.GetKey()
		}
		if t.GetOperator() != "" {
			entry["operator"] = t.GetOperator()
		}
		if t.GetValue() != "" {
			entry["value"] = t.GetValue()
		}
		if t.GetEffect() != "" {
			entry["effect"] = t.GetEffect()
		}
		if t.TolerationSeconds != nil {
			entry["tolerationSeconds"] = t.GetTolerationSeconds()
		}
		out = append(out, entry)
	}
	return out
}

// setNested writes entries into a nested map path, creating intermediate
// maps as needed and preserving siblings already set.
func setNested(root map[string]interface{}, entries map[string]interface{}, path ...string) {
	node := root
	for _, key := range path {
		child, ok := node[key].(map[string]interface{})
		if !ok {
			child = map[string]interface{}{}
			node[key] = child
		}
		node = child
	}
	for k, v := range entries {
		node[k] = v
	}
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

// toInterfaceSlice converts a []string into the []interface{} YAML rendering
// expects.
func toInterfaceSlice(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}
