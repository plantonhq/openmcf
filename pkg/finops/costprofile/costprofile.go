// Package costprofile loads and validates per-component cost profiles --
// the catalog/<provider>/<kind>/cost.yaml sidecars declaring each
// component's cost anatomy (billing model, always-on baseline charges, and
// the spec fields that drive the bill). The profile carries the STABLE half
// of cost knowledge; unit prices live in the central price book. Enrollment
// is the file's presence: every discovered cost.yaml is held to this
// package's conformance gate, with no allowlist to keep in sync.
package costprofile

import (
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	costprofilev1 "github.com/plantonhq/planton/finops/componentcostprofile/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// FileName is the sidecar's name at the component root.
const FileName = "cost.yaml"

// Path is a component's cost profile location.
func Path(repoRoot, provider, component string) string {
	return filepath.Join(repoRoot, "catalog", provider, component, FileName)
}

// Discover returns provider -> components that ship a cost profile, sorted.
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

// Load reads and strictly parses a component's cost profile.
func Load(repoRoot, provider, component string) (*costprofilev1.ComponentCostProfile, error) {
	profile := &costprofilev1.ComponentCostProfile{}
	if err := protobufyaml.Load(Path(repoRoot, provider, component), profile); err != nil {
		return nil, errors.Wrapf(err, "loading cost profile for %s/%s", provider, component)
	}
	return profile, nil
}
