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

// segmentGroupPattern matches an OPTIONAL SEGMENT GROUP: literal text plus
// exactly one placeholder wrapped in square brackets, e.g. "[//{namespace}]".
// When the placeholder resolves empty, the WHOLE group -- its literal
// delimiters included -- is dropped from the rendered ID.
//
// This is a distinct syntax from the "{name?}" optional marker, deliberately:
// "{name?}" keeps its surrounding literals when empty (the DynamoDB
// contributor-insights format "name:{table_name}/index:{index_name?}/..."
// legitimately renders "index:/" for the table-level variant), while a
// bracketed group drops them (the kubectl_manifest composed ID
// "{api_version}//{kind}//{name}[//{namespace}]" must render a 3-part ID with
// NO trailing "//" for cluster-scoped resources -- the importer rejects a
// trailing delimiter). Changing "{name?}" semantics instead would silently
// alter proven formats; a new bracket syntax is purely additive.
var segmentGroupPattern = regexp.MustCompile(`\[([^\[\]{}]*)\{([a-z0-9_]+)\}([^\[\]{}]*)\]`)

// Placeholders returns the placeholder names an id_format references, in
// order of appearance -- optional ones (both "{name?}" and bracketed-group
// placeholders) included: they still need a declaration and a resolution
// attempt; only their absence is tolerated.
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
// "name:tbl/index:/123456789012"), and a bracketed segment group
// ("[//{namespace}]") is dropped wholesale when its placeholder is empty.
func RequiredPlaceholders(idFormat string) []string {
	// Strip bracketed groups first: their placeholders are optional by
	// construction (the group's presence depends on resolution).
	withoutGroups := segmentGroupPattern.ReplaceAllString(idFormat, "")
	matches := placeholderPattern.FindAllStringSubmatch(withoutGroups, -1)
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
// ("{name?}") placeholders render as the empty string when unresolved (their
// surrounding literals stay), and bracketed segment groups ("[//{namespace}]")
// render with their literals when the placeholder resolves and disappear
// entirely when it does not.
func RenderID(idFormat string, values map[string]string) (string, error) {
	// Pass 1: bracketed segment groups (all-or-nothing with their literals).
	rendered := segmentGroupPattern.ReplaceAllStringFunc(idFormat, func(match string) string {
		groups := segmentGroupPattern.FindStringSubmatch(match)
		prefix, name, suffix := groups[1], groups[2], groups[3]
		value := values[name]
		if value == "" {
			return ""
		}
		return prefix + value + suffix
	})

	// Pass 2: plain placeholders.
	var missing []string
	rendered = placeholderPattern.ReplaceAllStringFunc(rendered, func(match string) string {
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
