// Baseline of accepted dollar tokens and the gate that compares live scan
// findings against it. Two lists with two lifetimes: `allowed` is permanent
// (dollar-typed config examples and scanner false positives, each with its
// reason recorded as a comment beside it -- hand-curated, never regenerated),
// while `gaps` is the burn-down backlog of hand-typed prices awaiting their
// driver-teaching rewrite -- it trends to 0 and never grows. The shape
// mirrors pkg/secretcoverage's and pkg/anatomy's baselines so a reader who
// knows one knows all three.
package priceprovenance

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const baselineHeader = `# Price-provenance baseline -- the dollar tokens the prose scanner accepts.
#
# ` + "`allowed`" + ` is PERMANENT and hand-curated: user-chosen dollar values that
# illustrate a dollar-typed configuration field (a budget limit, an alert
# threshold -- the user's number, never a provider rate), and non-price tokens
# the scanner cannot distinguish (regex backreferences like "/v2/$1"). Every
# entry carries its reason as a comment beside it. Regeneration preserves it.
#
# ` + "`gaps`" + ` is the BURN-DOWN backlog: hand-typed prices that have not been
# rewritten to driver-teaching copy yet. It trends to 0 and never grows.
# Prices belong in exactly one place -- catalog/_pricing/ (pinned price books
# and generated estimates, source-dated). Prose teaches cost DRIVERS.
#
# The CI guardrail (go test ./pkg/priceprovenance/...) fails when:
#   - a dollar token appears that is in neither list (a new hand-typed price
#     shipped -- rewrite it to teach the driver, or, ONLY for a user-chosen
#     dollar-typed config example, add it to allowed with its reason), or
#   - a listed entry no longer matches anything (it was fixed or the file
#     moved -- remove it here so both lists stay truthful).
#
# Entry format: <repo-relative path>:<token>
# Regenerate gaps with:  PLANTON_REGEN_PRICE_BASELINE=1 go test ./pkg/priceprovenance/
`

type baselineDoc struct {
	Allowed []string `yaml:"allowed"`
	Gaps    []string `yaml:"gaps"`
}

// Baseline is the parsed baseline: two identity sets with two lifetimes.
type Baseline struct {
	Allowed map[string]bool
	Gaps    map[string]bool
}

func LoadBaseline(path string) (Baseline, error) {
	b := Baseline{Allowed: map[string]bool{}, Gaps: map[string]bool{}}
	data, err := os.ReadFile(path)
	if err != nil {
		return b, err
	}
	var doc baselineDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return b, fmt.Errorf("parse baseline %s: %w", path, err)
	}
	for _, id := range doc.Allowed {
		b.Allowed[id] = true
	}
	for _, id := range doc.Gaps {
		b.Gaps[id] = true
	}
	return b, nil
}

// WriteBaseline regenerates the gaps list from live findings while
// PRESERVING the hand-curated allowed section byte-for-byte (including its
// reason comments): allowed entries are judgment recorded once; only the
// burn-down list is machine-refreshed. Findings already covered by an
// allowed entry are not written as gaps.
func WriteBaseline(path string, findings []Finding) error {
	allowedBlock, err := readAllowedBlock(path)
	if err != nil {
		return err
	}
	existing, _ := LoadBaseline(path)

	var ids []string
	seen := map[string]bool{}
	for _, f := range findings {
		id := f.ID()
		if existing.Allowed[id] || seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Strings(ids)

	var b strings.Builder
	b.WriteString(baselineHeader)
	b.WriteString(allowedBlock)
	if len(ids) == 0 {
		b.WriteString("gaps: []\n")
	} else {
		b.WriteString("gaps:\n")
		for _, id := range ids {
			b.WriteString("  - \"" + id + "\"\n")
		}
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// readAllowedBlock returns the existing file's `allowed:` section verbatim
// (comments included) so regeneration never destroys curated judgment. A
// missing file or missing section yields an empty allowed list.
func readAllowedBlock(path string) (string, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "allowed: []\n", nil
	}
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	start := -1
	end := len(lines)
	for i, line := range lines {
		if strings.HasPrefix(line, "allowed:") {
			start = i
			continue
		}
		if start >= 0 && strings.HasPrefix(line, "gaps:") {
			end = i
			break
		}
	}
	if start < 0 {
		return "allowed: []\n", nil
	}
	block := strings.Join(lines[start:end], "\n")
	if !strings.HasSuffix(block, "\n") {
		block += "\n"
	}
	return block, nil
}

// GateResult is the verdict of comparing a live scan to the baseline.
type GateResult struct {
	// NewTokens are dollar tokens in neither list -- a new hand-typed price
	// shipped (or a new config example missing its allowed entry).
	NewTokens []string
	// StaleGaps are gaps entries that no longer match anything -- fixed;
	// remove them so the burn-down stays truthful.
	StaleGaps []string
	// StaleAllowed are allowed entries that no longer match anything -- the
	// content or file is gone; remove them so the allowlist stays honest.
	StaleAllowed []string
}

func (g GateResult) OK() bool {
	return len(g.NewTokens) == 0 && len(g.StaleGaps) == 0 && len(g.StaleAllowed) == 0
}

// Gate is the single source of truth for the CI test, so a violation report
// and the guardrail can never disagree.
func Gate(findings []Finding, baseline Baseline) GateResult {
	var res GateResult
	live := map[string]bool{}
	for _, f := range findings {
		id := f.ID()
		live[id] = true
		if !baseline.Allowed[id] && !baseline.Gaps[id] {
			if !contains(res.NewTokens, id) {
				res.NewTokens = append(res.NewTokens, id)
			}
		}
	}
	for id := range baseline.Gaps {
		if !live[id] {
			res.StaleGaps = append(res.StaleGaps, id)
		}
	}
	for id := range baseline.Allowed {
		if !live[id] {
			res.StaleAllowed = append(res.StaleAllowed, id)
		}
	}
	sort.Strings(res.NewTokens)
	sort.Strings(res.StaleGaps)
	sort.Strings(res.StaleAllowed)
	return res
}

func contains(list []string, id string) bool {
	for _, v := range list {
		if v == id {
			return true
		}
	}
	return false
}
