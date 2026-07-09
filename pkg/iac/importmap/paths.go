package importmap

import "path/filepath"

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
