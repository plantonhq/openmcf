package outputs

import "strings"

// preprocessKeys normalizes IaC output map keys to match proto field paths.
//
// IaC engines (Terraform/Pulumi) may emit keys with conventions that differ
// from protobuf field naming:
//
//   - ".\[" -> "[" : Remove the dot before square brackets. Terraform uses
//     "subnets.[0].id" but the dot-path walker expects "subnets[0].id".
//
// Hyphenated FIELD names ("load-balancer-arn") are deliberately NOT rewritten
// here: a whole-key rewrite would also corrupt proto map KEYS, which are data
// that legitimately contains hyphens (subnet IDs like "subnet-0abc" keying a
// map<string,string> output). Field-name segments are instead normalized per
// segment at lookup time in setFieldRecursively, which knows which segments
// are field names and which are map keys.
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
