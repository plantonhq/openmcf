package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDirFiltered(t *testing.T) {
	src := t.TempDir()
	if err := os.MkdirAll(filepath.Join(src, "keep", "nested"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(src, "skipdir"), 0755); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		"root.txt":             "root",
		"keep/nested/deep.txt": "deep",
		"skipdir/inside.txt":   "never copied",
		"skip.txt":             "never copied",
	} {
		if err := os.WriteFile(filepath.Join(src, path), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	dst := filepath.Join(t.TempDir(), "copy")
	err := CopyDirFiltered(src, dst, func(relPath string, d os.DirEntry) bool {
		return d.Name() != "skipdir" && d.Name() != "skip.txt"
	})
	if err != nil {
		t.Fatalf("CopyDirFiltered() error: %v", err)
	}

	for _, path := range []string{"root.txt", "keep/nested/deep.txt"} {
		if _, err := os.Stat(filepath.Join(dst, path)); err != nil {
			t.Errorf("expected %s in the copy: %v", path, err)
		}
	}
	// A rejected directory skips its whole subtree.
	for _, path := range []string{"skip.txt", "skipdir", "skipdir/inside.txt"} {
		if _, err := os.Stat(filepath.Join(dst, path)); !os.IsNotExist(err) {
			t.Errorf("%s must not be in the copy", path)
		}
	}
}

func TestCopyDir_RejectsSymlinks(t *testing.T) {
	// Silently dereferencing or dropping a symlink would produce a copy
	// that looks complete but is not — the copy must fail loudly instead.
	src := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "real.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(src, "real.txt"), filepath.Join(src, "link.txt")); err != nil {
		t.Skipf("symlinks unsupported here: %v", err)
	}

	if err := CopyDir(src, filepath.Join(t.TempDir(), "copy")); err == nil {
		t.Fatal("expected an error copying a tree containing a symlink")
	}
}

func TestCopyFile_PreservesMode(t *testing.T) {
	src := filepath.Join(t.TempDir(), "script.sh")
	if err := os.WriteFile(src, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(t.TempDir(), "copy.sh")
	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile() error: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0755 {
		t.Errorf("mode = %v, want 0755", info.Mode().Perm())
	}
}
