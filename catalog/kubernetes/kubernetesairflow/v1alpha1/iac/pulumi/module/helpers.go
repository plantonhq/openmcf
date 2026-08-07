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

// tolerationsBlock renders WorkloadTolerations into the chart's
// tolerations list shape.
func tolerationsBlock(tolerations []*kubernetesprovider.WorkloadToleration) []interface{} {
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
			entry["tolerationSeconds"] = int(t.GetTolerationSeconds())
		}
		out = append(out, entry)
	}
	return out
}

// toInterfaceMap converts a proto string map into the generic map shape
// the values document uses.
func toInterfaceMap(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// mergeMaps deep-merges b over a with Helm `-f` semantics: nested maps
// merge recursively with b winning per key, everything else (lists,
// scalars) replaces wholesale.
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
