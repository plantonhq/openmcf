package infrachart

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// DiscoverCharts expands the given directories into chart directories. When all is
// false, each entry is taken as a chart directory itself. When all is true, each entry
// is a root walked recursively for every chart under it (a chart directory is marked by
// its Chart.yaml -- see IsChartDir); the walk does not descend into a chart, so nested
// template trees are never mistaken for charts. Results are sorted for deterministic
// iteration and reporting.
func DiscoverCharts(roots []string, all bool) ([]string, error) {
	if !all {
		return roots, nil
	}
	var dirs []string
	for _, rootDir := range roots {
		err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return err
			}
			if IsChartDir(path) {
				dirs = append(dirs, path)
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking %s: %w", rootDir, err)
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}
