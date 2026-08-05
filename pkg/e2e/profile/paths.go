package profile

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
)

const (
	providerBase = "apis/dev/planton/provider"

	// ProviderProfileRelPath is the path to the provider E2E profile relative
	// to the provider directory (e.g., kubernetes/aa_e2e/profile.yaml).
	ProviderProfileRelPath = "aa_e2e/profile.yaml"
)

// ProviderDir returns the absolute path to a provider's directory.
func ProviderDir(repoRoot, provider string) string {
	return filepath.Join(repoRoot, providerBase, provider)
}

// ProviderProfilePath returns the absolute path to a provider's E2E profile.
func ProviderProfilePath(repoRoot, provider string) string {
	return filepath.Join(repoRoot, providerBase, provider, ProviderProfileRelPath)
}

// componentVersionDir resolves a component directory name to its version
// segment via the kind registry. Component E2E assets live under the kind's
// declared version directory (e.g. kubernetesvalkey/v1/e2e/...), so the
// segment follows the registry — never a literal. A directory that does not
// resolve to a registered kind is not a component.
func componentVersionDir(component string) (string, error) {
	versionDir, err := crkreflect.ComponentVersionDir(component)
	if err != nil {
		return "", errors.Wrapf(err, "cannot locate E2E assets for component %q", component)
	}
	return versionDir, nil
}

// ComponentProfilePath returns the absolute path to a component's E2E profile
// (e.g. kubernetesvalkey/v1/e2e/profile.yaml).
func ComponentProfilePath(repoRoot, provider, component string) (string, error) {
	versionDir, err := componentVersionDir(component)
	if err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, providerBase, provider, component, versionDir, "e2e", "profile.yaml"), nil
}

// ComponentScenariosDir returns the absolute path to a component's test scenarios directory.
func ComponentScenariosDir(repoRoot, provider, component string) (string, error) {
	versionDir, err := componentVersionDir(component)
	if err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, providerBase, provider, component, versionDir, "e2e", "scenarios"), nil
}

// ComponentFixturesDir returns the absolute path to a component's fixture manifests directory.
func ComponentFixturesDir(repoRoot, provider, component string) (string, error) {
	versionDir, err := componentVersionDir(component)
	if err != nil {
		return "", err
	}
	return filepath.Join(repoRoot, providerBase, provider, component, versionDir, "e2e", "fixtures"), nil
}
