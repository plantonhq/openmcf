package mappingeval

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// tfResourcePattern matches top-level terraform resource declarations --
// the same shape the import-map conformance guard scans.
var tfResourcePattern = regexp.MustCompile(`(?m)^resource\s+"([a-z0-9_]+)"\s+"`)

// moduleResourceTypes scans a module directory's .tf files and returns the
// distinct resource types it declares.
func moduleResourceTypes(moduleDir string) ([]string, error) {
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading module dir %s", moduleDir)
	}
	seen := map[string]bool{}
	var types []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tf") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(moduleDir, entry.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "reading %s", entry.Name())
		}
		for _, match := range tfResourcePattern.FindAllStringSubmatch(string(content), -1) {
			if !seen[match[1]] {
				seen[match[1]] = true
				types = append(types, match[1])
			}
		}
	}
	return types, nil
}
