// Package costderivation loads and validates per-component cost
// derivations -- the catalog/_pricing/derivations/<component>.yaml
// documents carrying the machine-executable rules that turn ANY
// manifest's spec values into metered quantities and price choices. A
// derivation replaces the hand-authored estimate model for its component
// (a component carries exactly one of the two): the estimate generator
// replays every catalog preset through the rules to produce the
// committed estimates, and the same rules are what a server-side
// estimator evaluates against a live manifest. Enrollment is the file's
// presence: every discovered derivation is held to this package's
// conformance gate, with no allowlist to keep in sync.
package costderivation

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	derivationv1 "github.com/plantonhq/planton/finops/componentcostderivation/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// Dir is the derivations' home, relative to the repo root. Derivations
// live beside the price book and the generated estimates -- the estimate
// pipeline's data in one tree -- while the component's cost.yaml stays
// beside the component as its authored fact sheet.
const Dir = "catalog/_pricing/derivations"

// Path is a component's cost derivation location. The filename is the
// component's identity (the same convention the models and estimates use).
func Path(repoRoot, component string) string {
	return filepath.Join(repoRoot, Dir, component+".yaml")
}

// Discover returns the component names that ship a cost derivation, sorted.
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

// Load reads and strictly parses a component's cost derivation.
func Load(repoRoot, component string) (*derivationv1.ComponentCostDerivation, error) {
	derivation := &derivationv1.ComponentCostDerivation{}
	if err := protobufyaml.Load(Path(repoRoot, component), derivation); err != nil {
		return nil, errors.Wrapf(err, "loading cost derivation for %s", component)
	}
	return derivation, nil
}
