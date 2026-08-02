package importmap

import (
	"path/filepath"
	"strings"
)

const (
	providerBase = "apis/dev/planton/provider"

	// ProviderCatalogRelPath is the provider import catalog's path relative to
	// the provider directory (e.g. aws/aa_import/catalog.yaml).
	ProviderCatalogRelPath = "aa_import/catalog.yaml"

	// ComponentImportMapRelPath is the component import map's path relative to
	// the component directory (e.g. awss3bucket/v1/iac/import-map.yaml). It
	// sits next to the module it maps.
	ComponentImportMapRelPath = "v1/iac/import-map.yaml"
)

// ProviderCatalogPath returns the absolute path to a provider's import catalog.
func ProviderCatalogPath(repoRoot, provider string) string {
	return filepath.Join(repoRoot, providerBase, provider, ProviderCatalogRelPath)
}

// ComponentImportMapPath returns the absolute path to a component's import map.
func ComponentImportMapPath(repoRoot, provider, component string) string {
	return filepath.Join(repoRoot, providerBase, provider, component, ComponentImportMapRelPath)
}

// SplitAttributePath splits a declared attribute sub-path (a provider
// catalog config-only sub-path or a component map import-normalized
// sub-path) into its segments. Plain segments are dot-separated:
// "spec.update_strategy" → ["spec", "update_strategy"]. A segment may
// also be written bracket-quoted — `data["password.db"]` → ["data",
// "password.db"] — for map keys that themselves contain dots
// (Kubernetes Secret data keys are the canonical case), which a plain
// dotted path cannot express. An unterminated bracket segment is kept
// verbatim so the mismatch fails loud downstream instead of silently
// tolerating the wrong key.
func SplitAttributePath(path string) []string {
	var segments []string
	var current strings.Builder
	flush := func() {
		if current.Len() > 0 {
			segments = append(segments, current.String())
			current.Reset()
		}
	}
	for i := 0; i < len(path); {
		switch {
		case strings.HasPrefix(path[i:], `["`):
			flush()
			end := strings.Index(path[i+2:], `"]`)
			if end < 0 {
				current.WriteString(path[i:])
				i = len(path)
				continue
			}
			segments = append(segments, path[i+2:i+2+end])
			i += 2 + end + 2
		case path[i] == '.':
			flush()
			i++
		default:
			current.WriteByte(path[i])
			i++
		}
	}
	flush()
	return segments
}
