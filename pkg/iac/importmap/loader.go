package importmap

import (
	"os"

	"github.com/pkg/errors"

	componentv1 "github.com/plantonhq/planton/apis/dev/planton/iac/componentimportmap/v1"
	providerv1 "github.com/plantonhq/planton/apis/dev/planton/iac/providerimportcatalog/v1"
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
	m := &componentv1.ComponentImportMap{}
	if err := protobufyaml.Load(ComponentImportMapPath(repoRoot, provider, component), m); err != nil {
		return nil, errors.Wrapf(err, "loading component import map for %s/%s", provider, component)
	}
	return m, nil
}

// HasComponentImportMap reports whether a component ships an import map --
// the "is this kind mapped?" check callers use before offering derived import.
func HasComponentImportMap(repoRoot, provider, component string) bool {
	_, err := os.Stat(ComponentImportMapPath(repoRoot, provider, component))
	return err == nil
}
