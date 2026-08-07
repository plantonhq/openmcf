//go:build !codegen
// +build !codegen

// Baseline of accepted provider-parity gaps and the gate that compares live
// findings against it -- the pkg/anatomy / pkg/secretcoverage burn-down
// shape (a reader who knows one knows all three). Entries are baseline KEYS,
// one per kind or resource, never per argument: the file reads as the depth
// and breadth work lists and burns down one line per closed kind or judged
// resource. The permanent record lives in the manifests and the dispositions
// ledger; this file only holds what has not been judged YET.

package providerparity

import (
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// DefaultBaselinePath is repo-root-relative (where the CLI and CI run).
const DefaultBaselinePath = "pkg/providerparity/baseline.yaml"

const baselineHeader = `# Provider-parity baseline -- the accepted backlog of parity gaps: kinds not
# yet at total accounting against the pinned provider, and GA resources whose
# breadth disposition is not recorded yet. This list trends to 0 and never
# grows; the permanent record lives in each kind's iac/provider-parity.yaml
# and in pkg/providerparity/dispositions/.
#
# The CI guardrail (go test ./pkg/providerparity/...) fails when:
#   - a finding appears for a kind/resource NOT listed here (new gap shipped,
#     or a pin bump surfaced migration work), or
#   - a listed entry has no findings left (it was closed -- remove it here).
#
# Entry format: kind:<Kind> | resource:<terraform resource type>
# Regenerate with:  PLANTON_REGEN_PROVIDERPARITY_BASELINE=1 go test ./pkg/providerparity/
`

type baselineDoc struct {
	Entries []string `yaml:"entries"`
}

// LoadBaseline reads the accepted-gap set from a baseline file.
func LoadBaseline(path string) (map[string]bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc baselineDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, errors.Wrapf(err, "parse baseline %s", path)
	}
	set := make(map[string]bool, len(doc.Entries))
	for _, e := range doc.Entries {
		set[e] = true
	}
	return set, nil
}

// WriteBaseline writes the current findings' baseline keys as the accepted
// baseline, one line per kind/resource, sorted.
func WriteBaseline(path string, findings []Finding) error {
	keys := map[string]bool{}
	for _, f := range findings {
		keys[f.BaselineKey] = true
	}
	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)
	var b strings.Builder
	b.WriteString(baselineHeader)
	if len(sorted) == 0 {
		b.WriteString("entries: []\n")
	} else {
		b.WriteString("entries:\n")
		for _, k := range sorted {
			b.WriteString("  - " + k + "\n")
		}
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// GateResult is the verdict of comparing live findings to the baseline.
type GateResult struct {
	// NewFindings hit kinds/resources the baseline does not list -- a new
	// gap shipped, or a pin bump surfaced migration work.
	NewFindings []Finding
	// StaleEntries are baseline keys with no findings left -- closed;
	// remove them so the ledger only ever shrinks truthfully.
	StaleEntries []string
}

// OK reports whether the gate passes.
func (g GateResult) OK() bool {
	return len(g.NewFindings) == 0 && len(g.StaleEntries) == 0
}

// Gate compares findings against the accepted baseline. It is the single
// source of truth for the CI test and the CLI --check, so the two can never
// disagree.
func Gate(findings []Finding, baseline map[string]bool) GateResult {
	var res GateResult
	current := map[string]bool{}
	for _, f := range findings {
		current[f.BaselineKey] = true
		if !baseline[f.BaselineKey] {
			res.NewFindings = append(res.NewFindings, f)
		}
	}
	for key := range baseline {
		if !current[key] {
			res.StaleEntries = append(res.StaleEntries, key)
		}
	}
	sort.Slice(res.NewFindings, func(i, j int) bool {
		if res.NewFindings[i].BaselineKey != res.NewFindings[j].BaselineKey {
			return res.NewFindings[i].BaselineKey < res.NewFindings[j].BaselineKey
		}
		return res.NewFindings[i].Detail < res.NewFindings[j].Detail
	})
	sort.Strings(res.StaleEntries)
	return res
}
