package module

import (
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesargocdv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesargocd/v1"
)

// autoscalingBlock renders an autoscaling message into the chart's
// autoscaling shape — nil when autoscaling is not enabled, so the caller
// falls back to plain replicas (the HPA owns the count when enabled).
func autoscalingBlock(a *kubernetesargocdv1.KubernetesArgocdAutoscaling) map[string]interface{} {
	if a == nil || !a.GetEnabled() {
		return nil
	}
	minReplicas := 1
	if a.MinReplicas != nil {
		minReplicas = int(a.GetMinReplicas())
	}
	maxReplicas := 5
	if a.MaxReplicas != nil {
		maxReplicas = int(a.GetMaxReplicas())
	}
	return map[string]interface{}{
		"enabled":     true,
		"minReplicas": minReplicas,
		"maxReplicas": maxReplicas,
	}
}

// toggleableBlock renders a {enabled, resources} component. Only declared
// knobs render — the chart's own defaults (notifications/dex on, commit
// server off) stay authoritative when the block is empty.
func toggleableBlock(c *kubernetesargocdv1.KubernetesArgocdToggleableComponent) map[string]interface{} {
	block := map[string]interface{}{}
	if c == nil {
		return block
	}
	if c.Enabled != nil {
		block["enabled"] = c.GetEnabled()
	}
	if resources := resourcesBlock(c.GetResources()); resources != nil {
		block["resources"] = resources
	}
	return block
}

// mergeInto copies every entry of src into dst (shallow — used for the
// per-component metrics pair).
func mergeInto(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
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
