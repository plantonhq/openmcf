package module

import (
	"strings"

	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
)

// splitImageRepository separates the registry host from a declared image
// repository ("my.registry.com/grafana/grafana" → "my.registry.com",
// "grafana/grafana"). The chart composes its image reference as
// {registry}/{repository}:{tag} with registry defaulting to docker.io, so
// a repository that carries its own registry MUST be split — mapped
// verbatim onto image.repository it would render
// docker.io/my.registry.com/grafana/grafana and ImagePullBackOff. The
// first path segment is a registry exactly when it looks like a host (a
// dot, a port colon, or the literal localhost) — the same rule the
// container runtimes apply to bare image references.
func splitImageRepository(repo string) (registry string, repository string) {
	first, rest, found := strings.Cut(repo, "/")
	if !found {
		return "", repo
	}
	if strings.ContainsAny(first, ".:") || first == "localhost" {
		return first, rest
	}
	return "", repo
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
