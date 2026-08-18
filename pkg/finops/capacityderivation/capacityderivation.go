// Package capacityderivation loads and validates per-component capacity
// derivations -- the catalog/_pricing/capacity/<component>.yaml documents
// carrying the machine-executable rules that turn a cluster-capacity
// manifest's spec values into the capacity footprint it reserves from its
// target cluster (CPU/memory requests and limits, persistent volume
// storage). A capacity derivation replaces the hand-authored estimate
// model for its component (a component carries exactly one of the two):
// the estimate generator replays every catalog preset through the rules
// to produce the committed footprint estimates, and the same rules can
// compute a live manifest's footprint server-side. Enrollment is the
// file's presence: every discovered derivation is held to this package's
// conformance gate, with no allowlist to keep in sync.
package capacityderivation

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	capacityv1 "github.com/plantonhq/planton/finops/componentcapacityderivation/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// Dir is the capacity derivations' home, relative to the repo root --
// beside the cost derivations, models, price book, and generated
// estimates: the estimate pipeline's data in one tree.
const Dir = "catalog/_pricing/capacity"

// Path is a component's capacity derivation location. The filename is the
// component's identity (the same convention the derivations, models, and
// estimates use).
func Path(repoRoot, component string) string {
	return filepath.Join(repoRoot, Dir, component+".yaml")
}

// Discover returns the component names that ship a capacity derivation,
// sorted.
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

// Load reads and strictly parses a component's capacity derivation.
func Load(repoRoot, component string) (*capacityv1.ComponentCapacityDerivation, error) {
	derivation := &capacityv1.ComponentCapacityDerivation{}
	if err := protobufyaml.Load(Path(repoRoot, component), derivation); err != nil {
		return nil, errors.Wrapf(err, "loading capacity derivation for %s", component)
	}
	return derivation, nil
}
