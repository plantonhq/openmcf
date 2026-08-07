package eject

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/cli/cliprint"
	"github.com/plantonhq/planton/pkg/fileutil"
)

// licenseFileNames are the attribution files every ejected module must
// carry: an ejected module is a standalone redistribution of Apache-2.0
// catalog code, and Apache-2.0 §4(d) requires LICENSE and NOTICE to
// accompany every redistribution.
var licenseFileNames = []string{"LICENSE", "NOTICE"}

// ensureLicenseFiles guarantees the attribution files in the output.
// Release artifacts embed them at the module root (their lanes inject them);
// a repo checkout keeps them at the repo root, so the copy is backfilled
// from there. A missing file with no repo to backfill from means the
// artifact violated its own packaging contract — surfaced loudly, since a
// silent gap here ships an incomplete redistribution.
func ensureLicenseFiles(outputDir, repoRoot string) error {
	for _, name := range licenseFileNames {
		destPath := filepath.Join(outputDir, name)
		if _, err := os.Stat(destPath); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to check for %s", destPath)
		}

		if repoRoot == "" {
			cliprint.PrintWarning(fmt.Sprintf(
				"%s is missing from the downloaded module artifact — add the file from https://github.com/plantonhq/planton before redistributing", name))
			continue
		}

		srcPath := filepath.Join(repoRoot, name)
		if err := fileutil.CopyFile(srcPath, destPath); err != nil {
			return errors.Wrapf(err, "failed to copy %s into the ejected module", name)
		}
	}
	return nil
}
