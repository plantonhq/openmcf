package pulumimodule

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These tests pin the explicit-vs-implicit contract of local module
// resolution: an explicitly chosen directory that is not a usable module
// fails loudly (falling through would silently deploy the official module in
// place of the user's), while the implicit current-directory probe falls
// through quietly — the behavior every default-path deployment depends on.

func writeModuleFile(t *testing.T, dir string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "Pulumi.yaml"), []byte("name: test\nruntime: go\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestResolveLocalModuleDirExplicitValid(t *testing.T) {
	dir := t.TempDir()
	writeModuleFile(t, dir)

	got, err := resolveLocalModuleDir(dir)
	if err != nil {
		t.Fatalf("resolveLocalModuleDir() error: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Errorf("resolveLocalModuleDir() = %q, want an absolute path", got)
	}
	want, _ := filepath.Abs(dir)
	if got != want {
		t.Errorf("resolveLocalModuleDir() = %q, want %q", got, want)
	}
}

func TestResolveLocalModuleDirExplicitInvalidFailsLoudly(t *testing.T) {
	dir := t.TempDir() // exists but holds no Pulumi.yaml

	_, err := resolveLocalModuleDir(dir)
	if err == nil {
		t.Fatal("an explicitly passed non-module directory must fail, not fall through to the official module")
	}
	if !strings.Contains(err.Error(), dir) || !strings.Contains(err.Error(), "Pulumi.yaml") {
		t.Errorf("error %q must name the directory and what was missing", err.Error())
	}
}

func TestResolveLocalModuleDirImplicitPwdIsModule(t *testing.T) {
	dir := t.TempDir()
	writeModuleFile(t, dir)
	t.Chdir(dir)

	got, err := resolveLocalModuleDir("")
	if err != nil {
		t.Fatalf("resolveLocalModuleDir() error: %v", err)
	}
	// Resolve both through EvalSymlinks: on darwin, TempDir lives under a
	// symlinked /var → /private/var.
	gotReal, _ := filepath.EvalSymlinks(got)
	wantReal, _ := filepath.EvalSymlinks(dir)
	if gotReal != wantReal {
		t.Errorf("resolveLocalModuleDir() = %q, want the current directory %q", gotReal, wantReal)
	}
}

func TestResolveLocalModuleDirImplicitPwdNotModuleFallsThroughQuietly(t *testing.T) {
	t.Chdir(t.TempDir())

	got, err := resolveLocalModuleDir("")
	if err != nil {
		t.Fatalf("the implicit probe must never fail resolution, got: %v", err)
	}
	if got != "" {
		t.Errorf("resolveLocalModuleDir() = %q, want empty (fall through to binary/staging)", got)
	}
}
