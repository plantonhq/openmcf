// Baseline of accepted logo debt and the gate that compares live findings
// against it. The baseline is a burn-down list: providers whose logo sets
// predate the logo law and have not been judged glyph by glyph yet. Each
// provider wave pays its list down as it profiles its kinds; the list never
// grows, and a fixed entry must leave it. The shape mirrors pkg/anatomy's
// baseline so a reader who knows one knows both. The malformed-svg rule is
// never baselined: a file browsers cannot draw ships nowhere.
package cataloglogo

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const baselineHeader = `# Catalog-logo baseline -- the accepted backlog of logo debt: kinds whose logo
# shares its glyph with another kind of the same provider, or names no
# provenance, because their provider's logo set predates the logo law (one
# kind, one glyph; the official mark only for the kind that is the product; a
# Planton-drawn glyph everywhere else) and has not been judged yet. Each
# provider wave pays its list down; this list trends to 0 and never grows.
#
# The CI guardrail (go test ./pkg/cataloglogo/...) fails when:
#   - a violation appears that is NOT listed here (a new shared or unlabelled
#     glyph shipped), or
#   - a listed entry is no longer a violation (it was fixed -- remove it here).
#
# Entry format: <repo-relative path>:<rule>   (rules: shared-glyph, missing-provenance)
# Regenerate with:  PLANTON_REGEN_LOGO_BASELINE=1 go test ./pkg/cataloglogo/
`

type baselineDoc struct {
	Violations []string `yaml:"violations"`
}

// LoadBaseline reads the accepted-debt set from a baseline file.
func LoadBaseline(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc baselineDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parse baseline %s: %w", path, err)
	}
	set := make(map[string]bool, len(doc.Violations))
	for _, v := range doc.Violations {
		set[v] = true
	}
	return set, nil
}

// WriteBaseline writes the baselinable violations (everything but
// malformed-svg) as the accepted baseline.
func WriteBaseline(path string, violations []Violation) error {
	ids := make([]string, 0, len(violations))
	for _, v := range violations {
		if v.Rule == RuleMalformedSVG {
			continue
		}
		ids = append(ids, v.ID())
	}
	sort.Strings(ids)
	var b strings.Builder
	b.WriteString(baselineHeader)
	if len(ids) == 0 {
		b.WriteString("violations: []\n")
	} else {
		b.WriteString("violations:\n")
		for _, id := range ids {
			b.WriteString("  - " + id + "\n")
		}
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// GateResult is the verdict of comparing live violations to the baseline.
type GateResult struct {
	Malformed     []Violation // never baselined: the file will render as a broken image
	NewViolations []Violation // debt not in the baseline -- a new shared or unlabelled glyph shipped
	StaleEntries  []string    // baseline ids no longer violated -- remove them
}

// OK reports whether the gate passes.
func (g GateResult) OK() bool {
	return len(g.Malformed) == 0 && len(g.NewViolations) == 0 && len(g.StaleEntries) == 0
}

// Gate compares violations against the accepted baseline. It is the single
// source of truth for the CI test, so the test and any future doctor command
// can never disagree.
func Gate(violations []Violation, baseline map[string]bool) GateResult {
	var res GateResult
	current := map[string]bool{}
	for _, v := range violations {
		if v.Rule == RuleMalformedSVG {
			res.Malformed = append(res.Malformed, v)
			continue
		}
		id := v.ID()
		current[id] = true
		if !baseline[id] {
			res.NewViolations = append(res.NewViolations, v)
		}
	}
	for id := range baseline {
		if !current[id] {
			res.StaleEntries = append(res.StaleEntries, id)
		}
	}
	sort.Slice(res.Malformed, func(i, j int) bool { return res.Malformed[i].ID() < res.Malformed[j].ID() })
	sort.Slice(res.NewViolations, func(i, j int) bool { return res.NewViolations[i].ID() < res.NewViolations[j].ID() })
	sort.Strings(res.StaleEntries)
	return res
}
