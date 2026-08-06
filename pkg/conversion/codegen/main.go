//go:build codegen

// Regenerates pkg/conversion/embedded/specs -- the byte-for-byte mirror of
// every co-located conversion spec (apis/.../<kind>/conversions/*.yaml) that
// ships inside standalone binaries. Run via:
//
//	make generate-conversion-registry
//
// The mirror is committed; the drift gate in pkg/conversion keeps it honest.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	providerBase = "apis/dev/planton/provider"
	mirrorRoot   = "pkg/conversion/embedded/specs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	specs, err := filepath.Glob(filepath.Join(providerBase, "*", "*", "conversions", "*.yaml"))
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no conversion specs found under %s -- refusing to write an empty mirror", providerBase)
	}

	if err := os.RemoveAll(mirrorRoot); err != nil {
		return err
	}
	for _, spec := range specs {
		rel, err := filepath.Rel(providerBase, spec)
		if err != nil {
			return err
		}
		destination := filepath.Join(mirrorRoot, rel)
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}
		content, err := os.ReadFile(spec)
		if err != nil {
			return err
		}
		if err := os.WriteFile(destination, content, 0o644); err != nil {
			return err
		}
	}
	fmt.Printf("mirrored %d conversion specs into %s\n", len(specs), mirrorRoot)
	return nil
}
