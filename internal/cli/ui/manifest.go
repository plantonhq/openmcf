package ui

import (
	"fmt"
	"regexp"
	"strings"
)

// ManifestLoadError displays a beautiful error when proto unmarshaling fails.
// This handles errors like "unknown field 'spc'" or wrong types.
func ManifestLoadError(manifestPath string, err error) {
	sep := separator(errorIcon)

	// Parse the error to extract useful information
	errMsg := err.Error()
	fieldName := extractUnknownField(errMsg)
	lineInfo := extractLineInfo(errMsg)

	fmt.Println()
	fmt.Println(sep)
	fmt.Printf("%s  %s\n", errorIcon.Render(iconError), errorTitle.Render("Manifest Load Error"))
	fmt.Println(sep)

	// File path
	fmt.Printf("Failed to load manifest from: %s\n", Path(manifestPath))
	fmt.Println()

	// Error details
	if fieldName != "" {
		if lineInfo != "" {
			fmt.Printf("Error: Unknown field %s %s\n", errorTitle.Render("\""+fieldName+"\""), Dim(lineInfo))
		} else {
			fmt.Printf("Error: Unknown field %s\n", errorTitle.Render("\""+fieldName+"\""))
		}
	} else {
		// Generic error display
		fmt.Printf("Error: %s\n", errorMessage.Render(errMsg))
	}
	fmt.Println()

	// Explanation
	fmt.Println(warningTitle.Render("This typically happens when:"))
	fmt.Printf("  %s %s\n", Dim("-"), errorMessage.Render("There's a typo in a field name"))
	fmt.Printf("  %s %s\n", Dim("-"), errorMessage.Render("The field doesn't exist in the resource schema"))
	fmt.Printf("  %s %s\n", Dim("-"), errorMessage.Render("The manifest is from a different API version"))
	fmt.Println()

	// Suggest similar field if we found an unknown field
	if fieldName != "" {
		suggestion := suggestSimilarField(fieldName)
		if suggestion != "" {
			fmt.Printf("Did you mean: %s\n", Cmd(suggestion))
			fmt.Println()
		}
	}

	// Helpful commands
	fmt.Printf("%s %s\n", infoIcon.Render(iconTip),
		infoMessage.Render("Tip: Check your manifest YAML for typos in field names"))

	fmt.Println(sep)
}

// extractUnknownField extracts the field name from proto error messages.
// Examples:
//   - 'proto: (line 1:927): unknown field "spc"' -> "spc"
//   - 'unknown field "spc"' -> "spc"
func extractUnknownField(errMsg string) string {
	// Pattern: unknown field "fieldname"
	re := regexp.MustCompile(`unknown field "([^"]+)"`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// extractLineInfo extracts line/column info from proto error messages.
// Example: 'proto: (line 1:927): unknown field "spc"' -> "at line 1, column 927"
func extractLineInfo(errMsg string) string {
	// Pattern: (line N:M)
	re := regexp.MustCompile(`\(line (\d+):(\d+)\)`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) >= 3 {
		return fmt.Sprintf("at line %s, column %s", matches[1], matches[2])
	}
	return ""
}

// suggestSimilarField suggests common manifest fields that might be what the user meant.
func suggestSimilarField(fieldName string) string {
	// Common field name mappings based on typos
	suggestions := map[string]string{
		"spc":         "spec",
		"sepc":        "spec",
		"spce":        "spec",
		"spec ":       "spec",
		" spec":       "spec",
		"meta":        "metadata",
		"metdata":     "metadata",
		"meatadata":   "metadata",
		"metadta":     "metadata",
		"metadat":     "metadata",
		"apiversion":  "apiVersion",
		"api_version": "apiVersion",
		"ApiVersion":  "apiVersion",
		"Apiversion":  "apiVersion",
		"king":        "kind",
		"knd":         "kind",
		"knid":        "kind",
		"nam":         "name",
		"nmae":        "name",
		"naem":        "name",
		"lables":      "labels",
		"labesl":      "labels",
		"lable":       "labels",
	}

	// Check for exact match
	lower := strings.ToLower(fieldName)
	if suggestion, ok := suggestions[lower]; ok {
		return suggestion
	}
	if suggestion, ok := suggestions[fieldName]; ok {
		return suggestion
	}

	// Check for close matches using simple heuristics
	commonFields := []string{"spec", "metadata", "apiVersion", "kind", "name", "labels", "annotations"}
	for _, common := range commonFields {
		// If the field is very similar (off by 1-2 chars), suggest it
		if levenshteinDistance(lower, strings.ToLower(common)) <= 2 {
			return common
		}
	}

	return ""
}

// levenshteinDistance calculates the edit distance between two strings.
func levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Create matrix
	matrix := make([][]int, len(a)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(b)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(a)][len(b)]
}

func min(nums ...int) int {
	m := nums[0]
	for _, n := range nums[1:] {
		if n < m {
			m = n
		}
	}
	return m
}
