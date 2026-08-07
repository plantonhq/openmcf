// Command catalogbundle builds and verifies catalog bundles. It is the thin
// operational face of pkg/catalogbundle, used by the Makefile targets and
// the release/CI lanes:
//
//	go run ./pkg/catalogbundle/cli build  --descriptors <fds> --catalog-dir <dir> --out <zip> [--tag <tag>]
//	go run ./pkg/catalogbundle/cli verify --bundle <zip>
//
// verify loads the bundle (checksum self-verification included) and runs the
// registry conformance gate; a bundle that fails must never ship.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/plantonhq/planton/pkg/catalogbundle"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: catalogbundle build|verify [flags]")
	}
	switch os.Args[1] {
	case "build":
		fs := flag.NewFlagSet("build", flag.ExitOnError)
		descriptors := fs.String("descriptors", "", "buf-built FileDescriptorSet path")
		catalogDir := fs.String("catalog-dir", "catalog", "catalog tree holding conversions/ and presets/")
		out := fs.String("out", "", "bundle zip to write")
		tag := fs.String("tag", "", "release tag stamped into the manifest")
		if err := fs.Parse(os.Args[2:]); err != nil {
			return err
		}
		if *descriptors == "" || *out == "" {
			return fmt.Errorf("build requires --descriptors and --out")
		}
		manifest, err := catalogbundle.Build(catalogbundle.BuildInput{
			DescriptorSetPath: *descriptors,
			CatalogDir:        *catalogDir,
			ReleaseTag:        *tag,
			OutputPath:        *out,
		})
		if err != nil {
			return err
		}
		fmt.Printf("built %s (%d entries)\n", *out, len(manifest.Checksums))
		return nil

	case "verify":
		fs := flag.NewFlagSet("verify", flag.ExitOnError)
		bundlePath := fs.String("bundle", "", "bundle zip to verify")
		if err := fs.Parse(os.Args[2:]); err != nil {
			return err
		}
		if *bundlePath == "" {
			return fmt.Errorf("verify requires --bundle")
		}
		bundle, err := catalogbundle.Load(*bundlePath)
		if err != nil {
			return err
		}
		if err := catalogbundle.CheckConformance(bundle); err != nil {
			return err
		}
		fmt.Printf("bundle %s conforms to the compiled registry\n", *bundlePath)
		return nil

	default:
		return fmt.Errorf("unknown subcommand %q (want build or verify)", os.Args[1])
	}
}
