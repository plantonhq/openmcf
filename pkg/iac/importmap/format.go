package importmap

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// placeholderPattern matches {snake_case} placeholders in an id_format
// template. The name charset is deliberately narrow so a malformed template
// fails the conformance guard instead of silently importing a literal brace.
var placeholderPattern = regexp.MustCompile(`\{([a-z0-9_]+)\}`)

// Placeholders returns the placeholder names an id_format references, in
// order of appearance.
func Placeholders(idFormat string) []string {
	matches := placeholderPattern.FindAllStringSubmatch(idFormat, -1)
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, m[1])
	}
	return names
}

// RenderID fills an id_format template from resolved placeholder values.
// Every placeholder must resolve to a non-empty value -- a partially-filled
// import ID would import the wrong resource, so missing values are an error,
// never an empty substitution.
func RenderID(idFormat string, values map[string]string) (string, error) {
	var missing []string
	rendered := placeholderPattern.ReplaceAllStringFunc(idFormat, func(match string) string {
		name := placeholderPattern.FindStringSubmatch(match)[1]
		value := values[name]
		if value == "" {
			missing = append(missing, name)
		}
		return value
	})
	if len(missing) > 0 {
		return "", errors.Errorf("id format %q missing values for: %s", idFormat, strings.Join(missing, ", "))
	}
	return rendered, nil
}
