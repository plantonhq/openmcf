// Package yamlemit renders YAML scalars for the repo's generators. The
// generators own their artifacts' canonical form and write it by hand (the
// way the reference generator writes markdown): no marshaler round-trip, no
// map iteration, no surprises. This package holds the one decision every
// hand-renderer must make identically -- when a scalar needs quoting -- so
// two generators can never disagree about what a value means.
//
// Style contract: strings that would parse as YAML numbers, booleans, or
// null are double-quoted (a bare 730 or 2026-08-14 must stay a string);
// prose stays plain unless YAML syntax demands otherwise. Callers force
// quoting for fields whose type is string-by-contract (money, dates).
package yamlemit

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	numericPattern = regexp.MustCompile(`^-?[0-9]+(\.[0-9]+)?$`)
	boolNullWords  = map[string]bool{
		"true": true, "false": true, "yes": true, "no": true,
		"on": true, "off": true, "null": true, "~": true,
	}
)

// NeedsQuote reports whether a scalar cannot be emitted in YAML plain
// style without changing meaning.
func NeedsQuote(s string) bool {
	if s == "" {
		return true
	}
	if numericPattern.MatchString(s) || boolNullWords[strings.ToLower(s)] {
		return true
	}
	if strings.ContainsAny(s[:1], "-?:,[]{}#&*!|>'\"%@` \t") {
		return true
	}
	last := s[len(s)-1]
	if last == ' ' || last == ':' || last == '\t' {
		return true
	}
	if strings.Contains(s, ": ") || strings.Contains(s, " #") || strings.ContainsAny(s, "\n\t") {
		return true
	}
	return false
}

// Quote renders a scalar in YAML double-quoted style.
func Quote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

// WriteKV writes one `key: value` line at the given indent. Pass a key
// like "- name" to open a block-sequence mapping entry. forceQuote is for
// fields that are strings by contract regardless of content (money
// decimals, dates, URLs carrying spaces).
func WriteKV(b *strings.Builder, indent int, key, value string, forceQuote bool) {
	if forceQuote || NeedsQuote(value) {
		value = Quote(value)
	}
	fmt.Fprintf(b, "%s%s: %s\n", strings.Repeat(" ", indent), key, value)
}

// WriteListItem writes one `- value` block-sequence line at the given
// indent.
func WriteListItem(b *strings.Builder, indent int, value string) {
	if NeedsQuote(value) {
		value = Quote(value)
	}
	fmt.Fprintf(b, "%s- %s\n", strings.Repeat(" ", indent), value)
}
