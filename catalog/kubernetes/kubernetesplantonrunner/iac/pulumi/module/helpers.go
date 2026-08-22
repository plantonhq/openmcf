package module

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// mergeMaps deep-merges override into base with Helm `-f` semantics: maps
// merge recursively with the override winning, everything else (scalars
// and LISTS) replaces. Returns a new map; neither input is mutated.
func mergeMaps(base, override map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(base)+len(override))
	for key, value := range base {
		merged[key] = value
	}
	for key, overrideValue := range override {
		if baseMap, baseIsMap := merged[key].(map[string]interface{}); baseIsMap {
			if overrideMap, overrideIsMap := overrideValue.(map[string]interface{}); overrideIsMap {
				merged[key] = mergeMaps(baseMap, overrideMap)
				continue
			}
		}
		merged[key] = overrideValue
	}
	return merged
}

// chartVersionAtLeast compares two exact semver strings ("0.4.0")
// numerically, part by part. The spec's CEL already guarantees the exact
// three-part shape; the parse error arm exists for the default-version
// constant's own sake.
func chartVersionAtLeast(version, floor string) (bool, error) {
	versionParts, err := semverParts(version)
	if err != nil {
		return false, err
	}
	floorParts, err := semverParts(floor)
	if err != nil {
		return false, err
	}
	for i := 0; i < 3; i++ {
		if versionParts[i] != floorParts[i] {
			return versionParts[i] > floorParts[i], nil
		}
	}
	return true, nil
}

func semverParts(version string) ([3]int, error) {
	var parts [3]int
	segments := strings.Split(version, ".")
	if len(segments) != 3 {
		return parts, errors.Errorf("%q is not an exact major.minor.patch version", version)
	}
	for i, segment := range segments {
		number, err := strconv.Atoi(segment)
		if err != nil {
			return parts, errors.Errorf("%q is not an exact major.minor.patch version", version)
		}
		parts[i] = number
	}
	return parts, nil
}
