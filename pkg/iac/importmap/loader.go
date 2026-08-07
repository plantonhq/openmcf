package importmap

import (
	"os"

	"github.com/pkg/errors"

	componentv1 "github.com/plantonhq/planton/iac/componentimportmap/v1"
	providerv1 "github.com/plantonhq/planton/iac/providerimportcatalog/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// LoadProviderCatalog reads and parses a provider's import catalog from disk.
func LoadProviderCatalog(repoRoot, provider string) (*providerv1.ProviderImportCatalog, error) {
	catalog := &providerv1.ProviderImportCatalog{}
	if err := protobufyaml.Load(ProviderCatalogPath(repoRoot, provider), catalog); err != nil {
		return nil, errors.Wrapf(err, "loading provider import catalog for %s", provider)
	}
	return catalog, nil
}

// LoadComponentImportMap reads and parses a component's import map from disk.
func LoadComponentImportMap(repoRoot, provider, component string) (*componentv1.ComponentImportMap, error) {
	path, err := ComponentImportMapPath(repoRoot, provider, component)
	if err != nil {
		return nil, err
	}
	m := &componentv1.ComponentImportMap{}
	if err := protobufyaml.Load(path, m); err != nil {
		return nil, errors.Wrapf(err, "loading component import map for %s/%s", provider, component)
	}
	return m, nil
}

// HasComponentImportMap reports whether a component ships an import map --
// the "is this kind mapped?" check callers use before offering derived import.
// A name that does not resolve to a registered kind has no import map.
func HasComponentImportMap(repoRoot, provider, component string) bool {
	path, err := ComponentImportMapPath(repoRoot, provider, component)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}
