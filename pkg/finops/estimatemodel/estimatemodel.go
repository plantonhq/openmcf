// Package estimatemodel loads and validates per-component estimate models
// -- the catalog/_pricing/models/<component>.yaml documents stating each
// preset's quantity assumptions (how much of which declared meter, defended
// in prose). Models carry the AUTHORED half of an estimate; unit prices
// live in the provider's price book, and the estimate generator joins the
// two into the generated catalog/_pricing/estimates/ documents. Enrollment
// is the file's presence: every discovered model is held to this package's
// conformance gate, with no allowlist to keep in sync.
package estimatemodel

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	estimatemodelv1 "github.com/plantonhq/planton/finops/componentcostestimatemodel/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// Dir is the models' home, relative to the repo root. Models live beside
// the price book and the generated estimates -- the estimate pipeline's
// data in one tree -- while the component's cost.yaml stays beside the
// component as its authored fact sheet.
const Dir = "catalog/_pricing/models"

// Path is a component's estimate model location. The filename is the
// component's identity (the same convention the estimates use).
func Path(repoRoot, component string) string {
	return filepath.Join(repoRoot, Dir, component+".yaml")
}

// Discover returns the component names that ship an estimate model, sorted.
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

// Load reads and strictly parses a component's estimate model.
func Load(repoRoot, component string) (*estimatemodelv1.ComponentCostEstimateModel, error) {
	model := &estimatemodelv1.ComponentCostEstimateModel{}
	if err := protobufyaml.Load(Path(repoRoot, component), model); err != nil {
		return nil, errors.Wrapf(err, "loading estimate model for %s", component)
	}
	return model, nil
}
