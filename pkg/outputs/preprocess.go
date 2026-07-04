package outputs

import "strings"

// preprocessKeys normalizes IaC output map keys' PATH SYNTAX to match the
// dot-path walker:
//
//   - ".\[" -> "[" : Remove the dot before square brackets. Terraform uses
//     "subnets.[0].id" but the dot-path walker expects "subnets[0].id".
//
// Field-NAME normalization (hyphens to underscores, for IaC outputs whose
// names use hyphenated conventions) deliberately does NOT happen here: a
// dotted key may traverse a map field, and everything after the map field's
// segment is a map KEY -- user data such as a backend pool's name -- that
// must be preserved verbatim ("ssh-admin" must not become "ssh_admin").
// populateMessage normalizes each segment at field-descriptor lookup time
// instead, where it knows which segments are field names.
//
// The original map is not modified; a new map is returned.
func preprocessKeys(outputs map[string]string) map[string]string {
	result := make(map[string]string, len(outputs))
	for key, val := range outputs {
		normalized := strings.ReplaceAll(key, ".[", "[")
		result[normalized] = val
	}
	return result
}
