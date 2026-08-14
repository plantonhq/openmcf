// Command defspack validates the skills/ and agents/ definition trees and
// packages a definitions release: one deterministic zip per skill, each
// agent's instructions file, and a definitions-manifest.json carrying the
// release version, per-artifact SHA-256 checksums, and the compatibility
// floor from skills/compat.yaml.
//
// Run `go run ./pkg/skills/defspack` to validate only (the lint gate), or
// add -version and -out to package a release (the release workflow):
//
//	go run ./pkg/skills/defspack -version v0.4.0 -out build/definitions
//
// Validation is deliberately strict in both directions: every reference
// file SKILL.md cites must exist and be non-empty, and every file present
// under references/ must be cited. An empty or orphaned reference is
// content rot an agent would silently load at runtime, so it fails here
// instead.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	version := flag.String("version", "", "release version (vX.Y.Z) to stamp into the manifest; required with -out")
	outDir := flag.String("out", "", "directory to write release artifacts into; omit to validate only")
	root := flag.String("root", "", "repository root (defaults to the working directory)")
	// The catalog skill's assembled pack is ~18MB, and Stigmer engines cap
	// gRPC messages at 10MB in both directions (hardcoded in the engine's
	// server library, contradicting the skill layer's own 100MB acceptance
	// -- reported upstream: https://github.com/stigmer/stigmer/issues/675).
	// Until that lifts, PACKAGING excludes the pack by default -- released
	// artifacts stay pushable and mountable -- while loading and VALIDATION
	// always assemble it, so a pull request that breaks the pack fails the
	// lint gate today, and activation is this one flag on the day the cap
	// lifts.
	embedPack := flag.Bool("embed-catalog-pack", false, "package the catalog skill self-contained (its components/ reference pack inside the archive)")
	flag.Parse()

	repoRoot := *root
	if repoRoot == "" {
		wd, err := os.Getwd()
		if err != nil {
			fatal(err)
		}
		repoRoot = wd
	}

	tree, err := LoadTree(repoRoot)
	if err != nil {
		fatal(err)
	}
	if errs := Validate(tree); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintln(os.Stderr, "defspack: INVALID:", e)
		}
		os.Exit(1)
	}
	fmt.Printf("validated %d skill(s) and %d agent(s)\n", len(tree.Skills), len(tree.Agents))

	if *outDir == "" {
		return
	}
	if *version == "" {
		fatal(fmt.Errorf("-version is required when packaging with -out"))
	}
	if !*embedPack {
		StripPackFiles(tree)
	}
	manifest, err := PackageRelease(tree, *version, *outDir)
	if err != nil {
		fatal(err)
	}
	fmt.Printf("packaged definitions release %s into %s (%d skill archive(s), %d agent file(s))\n",
		manifest.Version, filepath.Clean(*outDir), len(manifest.Skills), len(manifest.Agents))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "defspack:", err)
	os.Exit(1)
}
