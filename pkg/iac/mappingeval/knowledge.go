package mappingeval

import (
	"github.com/pkg/errors"
	componentv1 "github.com/plantonhq/planton/apis/dev/planton/iac/componentimportmap/v1"
	providerv1 "github.com/plantonhq/planton/apis/dev/planton/iac/providerimportcatalog/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/iac/importmap"
)

// ScoreOptionsFromCatalog derives the scorer's declared knowledge from the
// artifacts the import subsystem already maintains -- the provider import
// catalog and the components' import maps. Nothing is authored per suite:
// when a catalog entry gains a config-only attribute or a component's
// recipe changes, the scorer's expectations follow automatically.
func ScoreOptionsFromCatalog(repoRoot, provider string, components []string) (ScoreOptions, error) {
	catalog, err := importmap.LoadProviderCatalog(repoRoot, provider)
	if err != nil {
		return ScoreOptions{}, err
	}

	// Config-only attributes can never round-trip through a scan (they
	// exist only in IaC configuration), so their spec fields leave the spec
	// axis. The catalog declares them per resource type; the union is safe
	// because exclusion is by field NAME and the naming convention keeps
	// spec fields aligned with the module attributes they drive.
	excluded := map[string]bool{}
	for _, rt := range catalog.GetSpec().GetResourceTypes() {
		for _, attr := range rt.GetConfigOnlyAttributes() {
			excluded[attr] = true
		}
	}

	nameDerived, err := nameDerivedIdentity(repoRoot, provider, components, catalog)
	if err != nil {
		return ScoreOptions{}, err
	}

	return ScoreOptions{
		ExcludedSpecFields:  excluded,
		NameDerivedIdentity: nameDerived,
	}, nil
}

// nameDerivedIdentity finds, per component kind, the Cloud Control type
// whose claimed identifier must equal metadata.name: the component's import
// map derives a placeholder from_metadata_name, and a catalog resource type
// with a declared scan-side name imports by exactly that placeholder.
func nameDerivedIdentity(repoRoot, provider string, components []string, catalog *providerv1.ProviderImportCatalog) (map[cloudresourcekind.CloudResourceKind]string, error) {
	result := map[cloudresourcekind.CloudResourceKind]string{}
	for _, component := range components {
		if !importmap.HasComponentImportMap(repoRoot, provider, component) {
			continue
		}
		m, err := importmap.LoadComponentImportMap(repoRoot, provider, component)
		if err != nil {
			return nil, errors.Wrapf(err, "import map for %s", component)
		}
		nameDerivedValues := map[string]bool{}
		for _, v := range m.GetSpec().GetValues() {
			for _, d := range v.GetDerivations() {
				if _, ok := d.GetSource().(*componentv1.ImportValueDerivation_FromMetadataName); ok {
					nameDerivedValues[v.GetName()] = true
				}
			}
		}
		if len(nameDerivedValues) == 0 {
			continue
		}
		kind := crkreflect.KindFromString(component)
		for _, rt := range catalog.GetSpec().GetResourceTypes() {
			if rt.GetCloudControlTypeName() == "" {
				continue
			}
			placeholders := importmap.Placeholders(rt.GetIdFormat())
			if len(placeholders) == 1 && nameDerivedValues[placeholders[0]] {
				result[kind] = rt.GetCloudControlTypeName()
				break
			}
		}
	}
	return result, nil
}
