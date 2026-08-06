package module

import (
	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
)

// resourcesMap renders the shared ContainerResources into the chart's
// requests/limits shape; nil when nothing is declared so the chart keeps
// its own defaults.
func resourcesMap(resources *kubernetesprovider.ContainerResources) map[string]interface{} {
	if resources == nil {
		return nil
	}
	out := map[string]interface{}{}
	if req := resources.GetRequests(); req != nil {
		requests := map[string]interface{}{}
		if req.GetCpu() != "" {
			requests["cpu"] = req.GetCpu()
		}
		if req.GetMemory() != "" {
			requests["memory"] = req.GetMemory()
		}
		if len(requests) > 0 {
			out["requests"] = requests
		}
	}
	if lim := resources.GetLimits(); lim != nil {
		limits := map[string]interface{}{}
		if lim.GetCpu() != "" {
			limits["cpu"] = lim.GetCpu()
		}
		if lim.GetMemory() != "" {
			limits["memory"] = lim.GetMemory()
		}
		if len(limits) > 0 {
			out["limits"] = limits
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
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

// stringMapToInterface converts a map[string]string into the
// map[string]interface{} the YAML rendering expects.
func stringMapToInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
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
