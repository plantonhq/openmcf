// Package discovery scans the planton repository to find testable components
// and their associated IaC modules and test manifests.
package discovery

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
)

// Component represents a discovered Planton component that can be E2E tested.
type Component struct {
	// Name is the component name in lowercase (e.g., "kubernetesnamespace").
	Name string

	// Provider is the cloud provider (e.g., "kubernetes", "aws", "gcp").
	Provider string

	// ManifestPath is the absolute path to the component's base test manifest
	// (e2e/manifest.yaml at the component root).
	ManifestPath string

	// PulumiDir is the absolute path to iac/pulumi/ (empty if not present).
	PulumiDir string

	// TerraformDir is the absolute path to iac/tf/ (empty if not present).
	TerraformDir string
}

// DiscoverComponents scans the catalog tree to find all components
// that have an e2e/manifest.yaml file (meaning they're testable).
func DiscoverComponents(repoRoot string) ([]Component, error) {
	catalogDir := filepath.Join(repoRoot, "catalog")

	var components []Component

	providerDirs, err := os.ReadDir(catalogDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read provider directory %s", catalogDir)
	}

	for _, providerEntry := range providerDirs {
		if !providerEntry.IsDir() {
			continue
		}
		providerName := providerEntry.Name()
		providerPath := filepath.Join(catalogDir, providerName)

		componentDirs, err := os.ReadDir(providerPath)
		if err != nil {
			continue
		}

		for _, componentEntry := range componentDirs {
			if !componentEntry.IsDir() {
				continue
			}
			componentName := componentEntry.Name()

			// A directory that does not resolve to a registered kind is not
			// a component (e.g. a provider's aa_e2e/ folder) — skip it,
			// exactly as a directory without a test manifest is skipped.
			if _, err := crkreflect.ComponentVersionDir(componentName); err != nil {
				continue
			}

			// The living component sits at the component root:
			// iac/ modules beside the e2e/ assets, no version segment.
			iacBase := filepath.Join(providerPath, componentName, "iac")
			manifestPath := filepath.Join(providerPath, componentName, "e2e", "manifest.yaml")

			if _, err := os.Stat(manifestPath); err != nil {
				continue
			}

			comp := Component{
				Name:         componentName,
				Provider:     providerName,
				ManifestPath: manifestPath,
			}

			pulumiDir := filepath.Join(iacBase, "pulumi")
			if info, err := os.Stat(pulumiDir); err == nil && info.IsDir() {
				comp.PulumiDir = pulumiDir
			}

			tfDir := filepath.Join(iacBase, "tf")
			if info, err := os.Stat(tfDir); err == nil && info.IsDir() {
				comp.TerraformDir = tfDir
			}

			components = append(components, comp)
		}
	}

	return components, nil
}

// DiscoverByProvider filters discovered components to a single provider.
func DiscoverByProvider(repoRoot, providerName string) ([]Component, error) {
	all, err := DiscoverComponents(repoRoot)
	if err != nil {
		return nil, err
	}

	var filtered []Component
	for _, c := range all {
		if strings.EqualFold(c.Provider, providerName) {
			filtered = append(filtered, c)
		}
	}
	return filtered, nil
}

// DiscoverByName finds a single component by name (case-insensitive).
func DiscoverByName(repoRoot, componentName string) (*Component, error) {
	all, err := DiscoverComponents(repoRoot)
	if err != nil {
		return nil, err
	}

	for _, c := range all {
		if strings.EqualFold(c.Name, componentName) {
			return &c, nil
		}
	}
	return nil, errors.Errorf("component %q not found", componentName)
}

// TestScenario represents one test manifest for a component under e2e/testdata/.
type TestScenario struct {
	// Name is the scenario name derived from the filename (e.g., "minimal", "with-hpa").
	Name string

	// ManifestPath is the absolute path to the test manifest YAML.
	ManifestPath string

	// Component is the component name (e.g., "kubernetesnamespace").
	Component string

	// Provider is the provider name (e.g., "kubernetes").
	Provider string
}

// ModuleDir returns a component's IaC module directory for the given engine
// ("pulumi" or "terraform"). Modules live at the component root — the path is
// fully derivable from provider and component; the registry check only guards
// against unregistered names.
func ModuleDir(repoRoot, provider, component, engine string) (string, error) {
	if _, err := crkreflect.ComponentVersionDir(component); err != nil {
		return "", err
	}
	var engineDir string
	switch engine {
	case "pulumi":
		engineDir = "pulumi"
	case "terraform":
		engineDir = "tf"
	default:
		return "", errors.Errorf("unsupported engine %q: want pulumi or terraform", engine)
	}
	return filepath.Join(repoRoot, "catalog", provider, component, "iac", engineDir), nil
}

// SecondActAnnotation marks a manifest in e2e/scenarios/ that is NOT a
// scenario of its own but the second act of one: the manifest a lifecycle
// annotation on another scenario (planton.dev/e2e-upgrade-manifest) deploys
// against that scenario's stack. Its value names the first-act scenario file.
// Discovery skips such manifests, so a second act never runs as a lane; the
// runner reaches it only through the first act's annotation.
const SecondActAnnotation = "planton.dev/e2e-second-act"

// isSecondActManifest reports whether the manifest carries SecondActAnnotation.
// Read errors count as "not a second act" so a malformed scenario still
// surfaces as a lane failure rather than vanishing silently.
func isSecondActManifest(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	// A structural read would be more exact, but the marker is a fixed key
	// under metadata.annotations and scenarios are hand-authored YAML; the
	// key's presence anywhere in the file is the author's declaration.
	return strings.Contains(string(data), SecondActAnnotation+":")
}

// DiscoverTestScenarios scans the component's colocated e2e/scenarios/ directory for YAML manifests.
// Path: catalog/{provider}/{component}/e2e/scenarios/
func DiscoverTestScenarios(repoRoot, provider, component string) ([]TestScenario, error) {
	if _, err := crkreflect.ComponentVersionDir(component); err != nil {
		return nil, err
	}
	scenarioDir := filepath.Join(repoRoot, "catalog", provider, component, "e2e", "scenarios")

	entries, err := os.ReadDir(scenarioDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to read test scenario directory %s", scenarioDir)
	}

	var scenarios []TestScenario
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}

		if isSecondActManifest(filepath.Join(scenarioDir, name)) {
			continue
		}

		scenarioName := strings.TrimSuffix(strings.TrimSuffix(name, ".yaml"), ".yml")
		scenarios = append(scenarios, TestScenario{
			Name:         scenarioName,
			ManifestPath: filepath.Join(scenarioDir, name),
			Component:    component,
			Provider:     provider,
		})
	}

	return scenarios, nil
}

// DiscoverAllTestScenarios scans all components under a provider for colocated e2e/ directories.
func DiscoverAllTestScenarios(repoRoot, provider string) (map[string][]TestScenario, error) {
	providerDir := filepath.Join(repoRoot, "catalog", provider)

	entries, err := os.ReadDir(providerDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to read provider testdata directory %s", providerDir)
	}

	result := make(map[string][]TestScenario)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		component := entry.Name()
		scenarios, err := DiscoverTestScenarios(repoRoot, provider, component)
		if err != nil {
			return nil, err
		}
		if len(scenarios) > 0 {
			result[component] = scenarios
		}
	}

	return result, nil
}
