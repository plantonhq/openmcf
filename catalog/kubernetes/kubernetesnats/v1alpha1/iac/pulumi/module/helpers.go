package module

import (
	kubernetesprovider "github.com/plantonhq/planton/catalog/kubernetes"
)

// resourcesBlock renders a ContainerResources into the chart's resources
// shape (nil when nothing is declared).
func resourcesBlock(r *kubernetesprovider.ContainerResources) map[string]interface{} {
	if r == nil {
		return nil
	}
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
	if len(resources) == 0 {
		return nil
	}
	return resources
}

// imageBlock renders a ContainerImage override into the chart's per-
// container image shape. The chart splits registry from repository
// internally; a repository value carrying a registry host works as-is
// (the chart prepends its registry value only when one is set), so the
// COMBINED repo form maps onto {repository, tag} directly. Nil when
// nothing is declared. Pull secrets are collected separately into the
// chart's global pullSecretNames list.
func imageBlock(img *kubernetesprovider.ContainerImage) map[string]interface{} {
	if img == nil || (img.GetRepo() == "" && img.GetTag() == "") {
		return nil
	}
	block := map[string]interface{}{}
	if img.GetRepo() != "" {
		block["repository"] = img.GetRepo()
	}
	if img.GetTag() != "" {
		block["tag"] = img.GetTag()
	}
	return block
}

// tolerationsSlice renders the shared WorkloadToleration list into the
// pod-spec tolerations shape.
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

// stringsToInterface converts a []string into the []interface{} YAML
// rendering expects.
func stringsToInterface(in []string) []interface{} {
	out := make([]interface{}, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}
