package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plantonhq/planton/e2e/framework/provider"
)

func writeSetupScript(t *testing.T, dir, body string) string {
	t.Helper()
	rel := "e2e-setup-test.sh"
	if err := os.WriteFile(filepath.Join(dir, rel), []byte("#!/usr/bin/env bash\nset -euo pipefail\n"+body+"\n"), 0o644); err != nil {
		t.Fatalf("writing script: %v", err)
	}
	return rel
}

func TestRunSetupScript_PassesEnvAndSucceeds(t *testing.T) {
	repoRoot := t.TempDir()
	outPath := filepath.Join(repoRoot, "captured.txt")
	rel := writeSetupScript(t, repoRoot, `echo "$E2E_RUN_ID $E2E_SCENARIO $PWD" > captured.txt`)

	tc := &provider.ComponentTestContext{
		RepoRoot:     repoRoot,
		ManifestPath: filepath.Join(repoRoot, "minimal.yaml"),
	}
	if err := runSetupScript(tc, rel, "run42-p"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	captured, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("script did not write its capture file: %v", err)
	}
	got := strings.TrimSpace(string(captured))
	if !strings.HasPrefix(got, "run42-p minimal ") {
		t.Fatalf("script env/cwd wrong: %q", got)
	}
	// The script must run from the repo root -- that is the path contract
	// every annotation value is written against. macOS reports TempDir
	// through /private symlinks, so compare resolved paths.
	wantDir, _ := filepath.EvalSymlinks(repoRoot)
	gotDir, _ := filepath.EvalSymlinks(strings.Fields(got)[2])
	if wantDir != gotDir {
		t.Fatalf("script ran from %q, want repo root %q", gotDir, wantDir)
	}
}

func TestRunSetupScript_FailureIsAnError(t *testing.T) {
	repoRoot := t.TempDir()
	rel := writeSetupScript(t, repoRoot, `echo "seeding failed" >&2; exit 3`)

	tc := &provider.ComponentTestContext{RepoRoot: repoRoot, ManifestPath: "minimal.yaml"}
	if err := runSetupScript(tc, rel, "run42-t"); err == nil {
		t.Fatal("expected a non-zero script exit to error the phase")
	}
}

func TestRunSetupScript_MissingScriptIsAnError(t *testing.T) {
	tc := &provider.ComponentTestContext{RepoRoot: t.TempDir(), ManifestPath: "minimal.yaml"}
	if err := runSetupScript(tc, "does/not/exist.sh", "run42-p"); err == nil {
		t.Fatal("expected a missing script to error the phase")
	}
}
