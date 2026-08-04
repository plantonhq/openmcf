// Command refgen generates the committed per-kind markdown reference files
// (reference.md, co-located with each kind's protos) from the descriptors
// compiled into the binary, via the explain engine's markdown renderer.
//
// Run through `make generate-reference` (optionally scoped with
// provider=<name>). Output is deterministic: identical descriptors produce
// byte-identical files, and the drift test in this package holds committed
// files to that promise.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	provider := flag.String("provider", "", "generate only kinds under this provider (e.g. aws); empty generates every kind")
	flag.Parse()

	repoRoot, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	summary, err := Generate(repoRoot, *provider)
	if err != nil {
		fatal(err)
	}

	for _, path := range summary.SortedPaths() {
		if err := os.WriteFile(filepath.Join(repoRoot, path), []byte(summary.Files[path]), 0644); err != nil {
			fatal(err)
		}
	}
	fmt.Printf("generated %d reference file(s)\n", len(summary.Files))

	if len(summary.MissingManifests) > 0 {
		fmt.Printf("\n%d kind(s) have no iac/hack/manifest.yaml (page rendered without an Example):\n", len(summary.MissingManifests))
		for _, kind := range summary.MissingManifests {
			fmt.Printf("  - %s\n", kind)
		}
	}
	if len(summary.InvalidManifests) > 0 {
		fmt.Printf("\n%d hack manifest(s) FAILED validation -- these are catalog bugs, fix the manifests:\n", len(summary.InvalidManifests))
		for _, invalid := range summary.InvalidManifests {
			fmt.Printf("  - %s: %v\n", invalid.Kind, invalid.Err)
		}
		os.Exit(1)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "refgen:", err)
	os.Exit(1)
}
