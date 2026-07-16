package importmap

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// placeholderPattern matches {snake_case} placeholders in an id_format
// template, with an optional trailing "?" marking the segment as optional
// (e.g. "{index_name?}"). The name charset is deliberately narrow so a
// malformed template fails the conformance guard instead of silently
// importing a literal brace.
var placeholderPattern = regexp.MustCompile(`\{([a-z0-9_]+)(\??)\}`)

// Placeholders returns the placeholder names an id_format references, in
// order of appearance -- optional ones included (they still need a
// declaration and a resolution attempt; only their absence is tolerated).
func Placeholders(idFormat string) []string {
	matches := placeholderPattern.FindAllStringSubmatch(idFormat, -1)
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, m[1])
	}
	return names
}

// RequiredPlaceholders returns only the names whose value MUST resolve for
// the ID to be importable -- the "{name?}" optional marker is for provider
// ID variants where a segment is legitimately empty (e.g. DynamoDB
// contributor insights on the table itself carry no index name:
// "name:tbl/index:/123456789012").
func RequiredPlaceholders(idFormat string) []string {
	matches := placeholderPattern.FindAllStringSubmatch(idFormat, -1)
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		if m[2] != "?" {
			names = append(names, m[1])
		}
	}
	return names
}

// RenderID fills an id_format template from resolved placeholder values.
// Every REQUIRED placeholder must resolve to a non-empty value -- a
// partially-filled import ID would import the wrong resource, so missing
// required values are an error, never an empty substitution. Optional
// ("{name?}") placeholders render as the empty string when unresolved,
// which is exactly the shape the provider documents for those variants.
func RenderID(idFormat string, values map[string]string) (string, error) {
	var missing []string
	rendered := placeholderPattern.ReplaceAllStringFunc(idFormat, func(match string) string {
		groups := placeholderPattern.FindStringSubmatch(match)
		name, optional := groups[1], groups[2] == "?"
		value := values[name]
		if value == "" && !optional {
			missing = append(missing, name)
		}
		return value
	})
	if len(missing) > 0 {
		return "", errors.Errorf("id format %q missing values for: %s", idFormat, strings.Join(missing, ", "))
	}
	return rendered, nil
}
