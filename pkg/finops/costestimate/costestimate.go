// Package costestimate loads and validates per-component preset cost
// estimates -- the catalog/_pricing/estimates/<component>.yaml documents
// stating what each preset costs per month at published list prices. The
// estimates carry the VOLATILE half of cost knowledge (today's prices,
// each with source URL and retrieval date); the component's cost.yaml
// carries the stable half (what bills at all), and the conformance gate
// binds the two: an estimate can only price meters its component's cost
// profile declares. Enrollment is the file's presence: every discovered
// estimate is held to this package's conformance gate, with no allowlist
// to keep in sync.
package costestimate

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	costestimatev1 "github.com/plantonhq/planton/finops/componentcostestimate/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// Dir is the estimates' home, relative to the repo root. Estimates live
// centrally rather than beside their components because prices churn on
// their own cadence -- one tree to refresh, no touch on 676 components.
const Dir = "catalog/_pricing/estimates"

// Path is a component's estimate document location. The filename is the
// component's identity (the same convention the framework crosswalks use).
func Path(repoRoot, component string) string {
	return filepath.Join(repoRoot, Dir, component+".yaml")
}

// Discover returns the component names that ship an estimate, sorted.
func Discover(repoRoot string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(repoRoot, Dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	var components []string
	for _, m := range matches {
		components = append(components, strings.TrimSuffix(filepath.Base(m), ".yaml"))
	}
	sort.Strings(components)
	return components, nil
}

// Load reads and strictly parses a component's estimate document.
func Load(repoRoot, component string) (*costestimatev1.ComponentCostEstimate, error) {
	estimate := &costestimatev1.ComponentCostEstimate{}
	if err := protobufyaml.Load(Path(repoRoot, component), estimate); err != nil {
		return nil, errors.Wrapf(err, "loading cost estimate for %s", component)
	}
	return estimate, nil
}
