// The catalog-schemas artifact CLI: `package` zips the buf-generated JSON
// tree into catalog-schemas.zip (sorted entries -- deterministic output for
// identical inputs), and `verify` proves the artifact serves the whole
// user-facing registry before it may ship. Mirrors the catalog-bundle CLI's
// build/verify shape: a schema artifact that fails verification must never
// ship.
package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/plantonhq/planton/pkg/catalogschema"
)

func main() {
	if len(os.Args) < 2 {
		fatal("usage: catalogschema-cli <package|verify> [flags]")
	}

	switch os.Args[1] {
	case "package":
		fs := flag.NewFlagSet("package", flag.ExitOnError)
		dir := fs.String("dir", "", "directory holding the buf-generated *.proto.json tree")
		out := fs.String("out", "", "output zip path")
		_ = fs.Parse(os.Args[2:])
		if *dir == "" || *out == "" {
			fatal("package requires --dir and --out")
		}
		if err := packageSchemas(*dir, *out); err != nil {
			fatal(err.Error())
		}
	case "verify":
		fs := flag.NewFlagSet("verify", flag.ExitOnError)
		zipPath := fs.String("zip", "", "catalog-schemas.zip to verify")
		_ = fs.Parse(os.Args[2:])
		if *zipPath == "" {
			fatal("verify requires --zip")
		}
		if err := catalogschema.VerifyArtifact(*zipPath); err != nil {
			fatal(err.Error())
		}
		fmt.Println("catalog-schemas artifact verified")
	default:
		fatal("unknown subcommand " + os.Args[1])
	}
}

// packageSchemas zips every *.proto.json under dir, preserving dir-relative
// paths, in sorted order with metadata-free entries -- identical inputs
// produce an identical zip.
func packageSchemas(dir, out string) error {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".proto.json") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("refusing to package an empty schema tree at %s -- the generation step produced nothing", dir)
	}
	sort.Strings(files)

	outFile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outFile.Close()
	zw := zip.NewWriter(outFile)
	for _, file := range files {
		rel, err := filepath.Rel(dir, file)
		if err != nil {
			return err
		}
		w, err := zw.Create(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		content, err := os.Open(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, content); err != nil {
			content.Close()
			return err
		}
		content.Close()
	}
	if err := zw.Close(); err != nil {
		return err
	}
	fmt.Printf("packaged %d schema documents into %s\n", len(files), out)
	return nil
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "error: "+msg)
	os.Exit(1)
}
