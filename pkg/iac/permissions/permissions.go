// Package permissions loads and validates per-component runner permission
// manifests -- the catalog/<provider>/<kind>/iac/permissions.yaml sidecars
// declaring the least-privilege permissions the IaC runner's cloud
// principal needs for the OFFICIAL modules. Every entry carries provenance
// (derived from module static analysis, or proven by live-run capture);
// the distinction is a trust feature consumers display. Enrollment is the
// file's presence: every discovered permissions.yaml is held to this
// package's conformance gate, with no allowlist to keep in sync.
package permissions

import (
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	permissionsv1 "github.com/plantonhq/planton/iac/componentpermissions/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// FileName is the sidecar's name inside a component's iac/ directory.
const FileName = "permissions.yaml"

// Path is a component's permissions manifest location.
func Path(repoRoot, provider, component string) string {
	return filepath.Join(repoRoot, "catalog", provider, component, "iac", FileName)
}

// Discover returns provider -> components that ship a permissions manifest,
// sorted.
func Discover(repoRoot string) (map[string][]string, error) {
	matches, err := filepath.Glob(filepath.Join(repoRoot, "catalog", "*", "*", "iac", FileName))
	if err != nil {
		return nil, err
	}
	discovered := map[string][]string{}
	for _, m := range matches {
		component := filepath.Base(filepath.Dir(filepath.Dir(m)))
		provider := filepath.Base(filepath.Dir(filepath.Dir(filepath.Dir(m))))
		discovered[provider] = append(discovered[provider], component)
	}
	for _, components := range discovered {
		sort.Strings(components)
	}
	return discovered, nil
}

// Load reads and strictly parses a component's permissions manifest.
func Load(repoRoot, provider, component string) (*permissionsv1.ComponentPermissions, error) {
	manifest := &permissionsv1.ComponentPermissions{}
	if err := protobufyaml.Load(Path(repoRoot, provider, component), manifest); err != nil {
		return nil, errors.Wrapf(err, "loading permissions for %s/%s", provider, component)
	}
	return manifest, nil
}
