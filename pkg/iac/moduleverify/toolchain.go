//go:build !codegen
// +build !codegen

package moduleverify

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/plantonhq/planton/pkg/fileutil"
)

// maxToolOutputLines caps how much toolchain output a violation carries —
// enough to see the real failure, not the whole scroll.
const maxToolOutputLines = 25

// runTofuValidate runs the engine's own validation (init without a backend,
// then validate) against a TEMPORARY copy of the module — init writes into
// the working directory, and verify never mutates the user's module.
// Skipped with a notice when neither tofu nor terraform is on PATH.
func runTofuValidate(moduleDir string, result *Result) {
	binary, err := exec.LookPath("tofu")
	if err != nil {
		binary, err = exec.LookPath("terraform")
	}
	if err != nil {
		result.addNotice("engine validation skipped: neither tofu nor terraform is on PATH")
		return
	}

	tempDir, err := os.MkdirTemp("", "module-verify-*")
	if err != nil {
		result.addNotice(fmt.Sprintf("engine validation skipped: %v", err))
		return
	}
	defer os.RemoveAll(tempDir)

	if err := fileutil.CopyDir(moduleDir, tempDir); err != nil {
		result.addNotice(fmt.Sprintf("engine validation skipped: %v", err))
		return
	}

	// -backend=false skips state backend setup; providers still download on
	// a cold cache, which is where most of this check's time goes.
	initCmd := exec.Command(binary, "init", "-backend=false", "-input=false", "-no-color")
	initCmd.Dir = tempDir
	if out, err := initCmd.CombinedOutput(); err != nil {
		result.addError("", fmt.Sprintf("%s init failed:\n%s", binary, tailLines(string(out))))
		return
	}

	validateCmd := exec.Command(binary, "validate", "-no-color")
	validateCmd.Dir = tempDir
	if out, err := validateCmd.CombinedOutput(); err != nil {
		result.addError("", fmt.Sprintf("%s validate failed:\n%s", binary, tailLines(string(out))))
	}
}

// runGoBuild compiles the pulumi module — the strongest proof its source,
// imports, and go.mod actually agree. The workspace overlay is disabled so
// the module builds from its own go.mod graph, exactly as deployments build
// it. Skipped with a notice when no go toolchain is on PATH.
func runGoBuild(moduleDir string, result *Result) {
	goBinary, err := exec.LookPath("go")
	if err != nil {
		result.addNotice("compile check skipped: no go toolchain on PATH")
		return
	}

	cmd := exec.Command(goBinary, "build", "./...")
	cmd.Dir = moduleDir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		result.addError("", fmt.Sprintf("go build failed:\n%s", tailLines(string(out))))
	}
}

// tailLines keeps the last maxToolOutputLines lines of toolchain output.
func tailLines(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) > maxToolOutputLines {
		lines = lines[len(lines)-maxToolOutputLines:]
	}
	return strings.Join(lines, "\n")
}
