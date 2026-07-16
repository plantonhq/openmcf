package importmap

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
)

// DiscoverComponentImportMaps walks the provider tree and returns every
// component that ships an import map, keyed by provider directory name
// (e.g. "aws" -> ["awss3bucket", "awsvpc", ...]), with both providers and
// components sorted for deterministic iteration.
//
// File presence IS enrollment: a component with a
// `v1/iac/import-map.yaml` is validated by the offline conformance guard,
// picked up by the E2E import round-trip gate, and bundled into the
// platform's catalog data for the import wizard — all from this one signal.
// There is deliberately no separate allowlist to keep in sync; an authored
// map that would not be checked, or a checked map that would not ship, is a
// divergence this single signal makes impossible.
func DiscoverComponentImportMaps(repoRoot string) (map[string][]string, error) {
	providersDir := filepath.Join(repoRoot, providerBase)
	providerEntries, err := os.ReadDir(providersDir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading providers directory %s", providersDir)
	}

	discovered := map[string][]string{}
	for _, providerEntry := range providerEntries {
		if !providerEntry.IsDir() {
			continue
		}
		provider := providerEntry.Name()
		componentEntries, err := os.ReadDir(filepath.Join(providersDir, provider))
		if err != nil {
			return nil, errors.Wrapf(err, "reading provider directory %s", provider)
		}
		for _, componentEntry := range componentEntries {
			if !componentEntry.IsDir() {
				continue
			}
			component := componentEntry.Name()
			if HasComponentImportMap(repoRoot, provider, component) {
				discovered[provider] = append(discovered[provider], component)
			}
		}
	}
	for provider := range discovered {
		sort.Strings(discovered[provider])
	}
	return discovered, nil
}
