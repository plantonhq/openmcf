package profile

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
)

const (
	catalogRoot = "catalog"

	// ProviderProfileRelPath is the path to the provider E2E profile relative
	// to the provider directory (e.g., kubernetes/aa_e2e/profile.yaml).
	ProviderProfileRelPath = "aa_e2e/profile.yaml"
)

// ProviderDir returns the absolute path to a provider's directory.
func ProviderDir(repoRoot, provider string) string {
	return filepath.Join(repoRoot, catalogRoot, provider)
}

// ProviderProfilePath returns the absolute path to a provider's E2E profile.
func ProviderProfilePath(repoRoot, provider string) string {
	return filepath.Join(repoRoot, catalogRoot, provider, ProviderProfileRelPath)
}

// validateComponent checks a component directory name against the kind
// registry. Component E2E assets live at the component root
// (e.g. kubernetesvalkey/e2e/...); a directory that does not resolve to a
// registered kind is not a component.
func validateComponent(component string) error {
	if _, err := crkreflect.ComponentVersionDir(component); err != nil {
		return errors.Wrapf(err, "cannot locate E2E assets for component %q", component)
	}
	return nil
}

// ComponentProfilePath returns the absolute path to a component's E2E profile
// (e.g. kubernetesvalkey/e2e/profile.yaml).
func ComponentProfilePath(repoRoot, provider, component string) (string, error) {
	if err := validateComponent(component); err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, catalogRoot, provider, component, "e2e", "profile.yaml"), nil
}

// ComponentScenariosDir returns the absolute path to a component's test scenarios directory.
func ComponentScenariosDir(repoRoot, provider, component string) (string, error) {
	if err := validateComponent(component); err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, catalogRoot, provider, component, "e2e", "scenarios"), nil
}

// ComponentFixturesDir returns the absolute path to a component's fixture manifests directory.
func ComponentFixturesDir(repoRoot, provider, component string) (string, error) {
	if err := validateComponent(component); err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, catalogRoot, provider, component, "e2e", "fixtures"), nil
}
