package fileutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// CopyDir recursively copies the directory tree rooted at src into dst,
// preserving file modes. dst is created if it does not exist.
func CopyDir(src, dst string) error {
	return CopyDirFiltered(src, dst, nil)
}

// CopyDirFiltered recursively copies the directory tree rooted at src into
// dst, preserving file modes. When include is non-nil it is consulted with
// each entry's path relative to src; entries it rejects are skipped (for a
// directory, the whole subtree).
//
// Only directories and regular files are copied. Anything else (symlinks,
// devices) fails loudly: silently dropping or dereferencing an entry would
// produce a copy that looks complete but is not.
func CopyDirFiltered(src, dst string, include func(relPath string, d os.DirEntry) bool) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return os.MkdirAll(dst, 0755)
		}

		if include != nil && !include(relPath, d) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		destPath := filepath.Join(dst, relPath)

		info, err := d.Info()
		if err != nil {
			return errors.Wrapf(err, "failed to stat %s", path)
		}

		switch {
		case d.IsDir():
			return os.MkdirAll(destPath, info.Mode().Perm())
		case info.Mode().IsRegular():
			return CopyFile(path, destPath)
		default:
			return errors.Errorf("cannot copy %s: unsupported file type %s", path, info.Mode().Type())
		}
	})
}

// CopyFile copies a single regular file from src to dst, preserving the file
// mode. The destination's parent directory must already exist.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", src)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return errors.Wrapf(err, "failed to stat %s", src)
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", dst)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return errors.Wrapf(err, "failed to copy %s to %s", src, dst)
	}
	return nil
}
