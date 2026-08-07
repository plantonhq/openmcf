package fileutil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// ExtractZip extracts a zip archive into destDir, creating directories as
// needed and preserving file modes. Entries whose paths would escape destDir
// (zip slip) are rejected.
func ExtractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open zip file %s", zipPath)
	}
	defer r.Close()

	for _, f := range r.File {
		destPath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return errors.Errorf("illegal file path in zip: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return errors.Wrapf(err, "failed to create directory %s", destPath)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return errors.Wrapf(err, "failed to create parent directory for %s", destPath)
		}

		if err := extractZipFile(f, destPath); err != nil {
			return err
		}
	}

	return nil
}

func extractZipFile(f *zip.File, destPath string) error {
	rc, err := f.Open()
	if err != nil {
		return errors.Wrapf(err, "failed to open file in zip: %s", f.Name)
	}
	defer rc.Close()

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", destPath)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, rc); err != nil {
		return errors.Wrapf(err, "failed to write file %s", destPath)
	}

	return nil
}
