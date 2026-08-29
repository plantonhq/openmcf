// Package presetvalidity is the machine-enforced truth of the catalog's
// presets: every preset YAML under catalog/<provider>/<kind>/presets/ must
// load as a manifest of its own kind and pass that kind's complete
// validation rules, so a preset that dead-ends the user who copies it is
// unshippable instead of auditable.
//
// A preset is a promise -- the catalog's own authored starting point for a
// kind. Two placeholder grammars keep that promise while staying visibly
// replace-me:
//
//   - Free-string fields may carry angle-bracket tokens (`<aws-region>`):
//     nothing fences them, consoles recognize and clear them, and they read
//     as obviously-not-real.
//   - Pattern- or CEL-fenced fields need pattern-VALID placeholders (a
//     documentation-account id like 123456789012, `cdn.replaceme.example.com`,
//     `arn:aws:iam::123456789012:role/replace-me`), and proto-sensitive
//     fields carry reference-shaped placeholders (`$secret/replace-with-...`)
//     that teach the managed-secret contract instead of plaintext habits.
//
// This gate is what makes the second grammar mandatory: an angle-bracket
// token seeded into a fenced field fails the kind's own validation right
// here, at authoring time, instead of in the user's first apply.
//
// The division of labor with the sibling catalog gates is deliberate and
// closed:
//
//   - pkg/anatomy proves the FILES EXIST (a preset and its .md sidecar are
//     anatomy).
//   - pkg/catalogpage proves the catalog PAGE has the right shape and its
//     embedded manifests validate.
//   - pkg/presetvalidity (this package) proves the PRESETS THEMSELVES
//     validate.
//
// The descriptor of accepted gaps lives in baseline.yaml -- a burn-down
// list mirroring pkg/anatomy's and pkg/catalogpage's baselines (a reader
// who knows one knows all three). The CI lane is
// .github/workflows/lint.preset-validity.yaml.
package presetvalidity

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"buf.build/go/protovalidate"

	"github.com/plantonhq/planton/internal/manifest"
)

// Violation is one preset-validity rule broken at one preset.
type Violation struct {
	// Path is the repo-root-relative preset YAML path.
	Path string
	// Rule is the stable rule identifier (baseline entries key on it).
	Rule string
	// Detail says what the fix is, in plain language.
	Detail string
}

// ID is the stable baseline identifier: "<path>:<rule>".
func (v Violation) ID() string { return v.Path + ":" + v.Rule }

// Rule identifiers. Stable: baseline.yaml entries reference them.
const (
	// RuleInvalidPreset fires when a preset fails to load as a manifest of
	// its declared kind or fails that kind's validation rules -- the exact
	// rejection the user who copies the preset would hit.
	RuleInvalidPreset = "invalid-preset"
)

// walkValidator is one lazily-compiling validator shared across the whole
// walk. The manifest package's ValidateLoaded builds (and rule-compiles) a
// fresh validator per call -- the right trade for validating one user
// manifest in the CLI, and the wrong one for 1,800 presets in a gate, where
// it recompiles the same kinds' rules over and over. Sharing one validator
// caches each kind's compiled rules on first encounter; the rules evaluated
// are byte-identical to the CLI path's.
var walkValidator = sync.OnceValues(func() (protovalidate.Validator, error) {
	return protovalidate.New()
})

// CheckPreset validates one preset's bytes as a manifest of its own kind.
// path is the repo-root-relative preset path used in violation IDs.
func CheckPreset(path string, content []byte) []Violation {
	loaded, err := manifest.LoadManifestBytes(content, path)
	if err == nil {
		var v protovalidate.Validator
		if v, err = walkValidator(); err == nil {
			err = v.Validate(loaded)
		}
	}
	if err == nil {
		return nil
	}
	return []Violation{{
		Path:   path,
		Rule:   RuleInvalidPreset,
		Detail: fmt.Sprintf("the preset does not validate against its own kind's rules: %s", firstMeaningfulLine(err.Error())),
	}}
}

// Check walks every preset YAML at repoRoot and returns every violation,
// sorted by ID. It never consults the baseline -- Gate does the comparison.
// Preset PRESENCE (and the .md sidecar) is pkg/anatomy's rule, so only
// existing presets are checked; the _test provider's fixtures are not
// product presets.
func Check(repoRoot string) ([]Violation, error) {
	presets, err := filepath.Glob(filepath.Join(repoRoot, "catalog", "*", "*", "presets", "*.yaml"))
	if err != nil {
		return nil, err
	}
	sort.Strings(presets)
	var vs []Violation
	for _, preset := range presets {
		rel, err := filepath.Rel(repoRoot, preset)
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(rel, filepath.Join("catalog", "_test")+string(filepath.Separator)) {
			continue
		}
		content, err := os.ReadFile(preset)
		if err != nil {
			return nil, err
		}
		vs = append(vs, CheckPreset(filepath.ToSlash(rel), content)...)
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].ID() < vs[j].ID() })
	return vs, nil
}

// firstMeaningfulLine collapses the manifest validator's multi-line, colored
// error rendering into the first line that carries an actual violation, so a
// gate failure reads as one actionable sentence per preset.
func firstMeaningfulLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(stripANSI(line))
		if line == "" {
			continue
		}
		// Skip the validator's banner/framing lines; keep the first line
		// that names a field path or a plain-language violation.
		if line == "validation error:" || line == "validation errors:" ||
			strings.ContainsAny(line, "╔╚║═") ||
			strings.HasPrefix(line, "❌") || strings.HasPrefix(line, "⚠️") ||
			strings.HasPrefix(line, "💡") || strings.HasPrefix(line, "📋") ||
			strings.HasPrefix(line, "📚") || strings.HasPrefix(line, "•") ||
			strings.HasPrefix(line, "Please review") ||
			strings.HasPrefix(line, "in your manifest") {
			continue
		}
		return strings.TrimPrefix(strings.TrimPrefix(line, "validation error: "), "- ")
	}
	return strings.TrimSpace(s)
}

// stripANSI removes terminal color escape sequences from the validator's
// user-facing rendering so baseline details stay plain text.
func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inEscape {
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				inEscape = false
			}
			continue
		}
		if c == 0x1b {
			inEscape = true
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}
