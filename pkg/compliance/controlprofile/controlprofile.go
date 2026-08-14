// Package controlprofile loads and validates per-component control profiles
// -- the catalog/<provider>/<kind>/controls.yaml sidecars declaring each
// component's posture against the central control catalog
// (pkg/compliance/controlcatalog). Enrollment is the file's presence: every
// discovered controls.yaml is held to this package's conformance gate --
// including COMPLETE examination of the catalog's controls -- with no
// allowlist to keep in sync.
package controlprofile

import (
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	controlprofilev1 "github.com/plantonhq/planton/compliance/componentcontrolprofile/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// FileName is the sidecar's name at the component root.
const FileName = "controls.yaml"

// Path is a component's control profile location.
func Path(repoRoot, provider, component string) string {
	return filepath.Join(repoRoot, "catalog", provider, component, FileName)
}

// Discover returns provider -> components that ship a control profile,
// sorted.
func Discover(repoRoot string) (map[string][]string, error) {
	matches, err := filepath.Glob(filepath.Join(repoRoot, "catalog", "*", "*", FileName))
	if err != nil {
		return nil, err
	}
	discovered := map[string][]string{}
	for _, m := range matches {
		component := filepath.Base(filepath.Dir(m))
		provider := filepath.Base(filepath.Dir(filepath.Dir(m)))
		discovered[provider] = append(discovered[provider], component)
	}
	for _, components := range discovered {
		sort.Strings(components)
	}
	return discovered, nil
}

// Load reads and strictly parses a component's control profile.
func Load(repoRoot, provider, component string) (*controlprofilev1.ComponentControlProfile, error) {
	profile := &controlprofilev1.ComponentControlProfile{}
	if err := protobufyaml.Load(Path(repoRoot, provider, component), profile); err != nil {
		return nil, errors.Wrapf(err, "loading control profile for %s/%s", provider, component)
	}
	return profile, nil
}
