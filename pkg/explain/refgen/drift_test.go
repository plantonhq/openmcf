package main

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// TestReferenceDrift is the freshness gate for committed reference.md files:
// for every provider that has committed pages, a fresh generation must be
// byte-identical to what is committed. Scoping by committed providers keeps
// the gate exactly as wide as the rollout -- a provider's first generated
// commit automatically brings it under guard.
//
// On drift: PLANTON_REGEN_REFERENCE=1 go test ./pkg/explain/refgen -run TestReferenceDrift
// (or `make generate-reference`).
func TestReferenceDrift(t *testing.T) {
	root := repoRoot(t)
	providers := committedProviders(t, root)
	if len(providers) == 0 {
		t.Skip("no committed reference.md files yet -- the gate arms with the first generated commit")
	}

	regen := os.Getenv("PLANTON_REGEN_REFERENCE") != ""
	for _, provider := range providers {
		summary, err := Generate(root, provider)
		if err != nil {
			t.Fatalf("generate %s: %v", provider, err)
		}
		// A broken hack manifest silently thins the committed page (its
		// Example disappears on the next regeneration) -- fail loudly instead.
		for _, invalid := range summary.InvalidManifests {
			t.Errorf("%s: hack manifest fails validation (a catalog bug): %v", invalid.Kind, invalid.Err)
		}
		for _, path := range summary.SortedPaths() {
			want := summary.Files[path]
			target := filepath.Join(root, path)
			if regen {
				if err := os.WriteFile(target, []byte(want), 0644); err != nil {
					t.Fatalf("write %s: %v", target, err)
				}
				continue
			}
			got, err := os.ReadFile(target)
			if err != nil {
				t.Errorf("%s is generated but not committed -- run `make generate-reference provider=%s`: %v", path, provider, err)
				continue
			}
			if string(got) != want {
				t.Errorf("%s is out of sync with the schema -- run `make generate-reference provider=%s`", path, provider)
			}
		}
	}
}

// committedProviders lists providers that already have committed reference
// pages, discovered from the filesystem so no version segment is assumed.
func committedProviders(t *testing.T, root string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(root, "apis", providerPathPrefix, "*", "*", "*", "reference.md"))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, match := range matches {
		rel, err := filepath.Rel(filepath.Join(root, "apis", providerPathPrefix), match)
		if err != nil {
			t.Fatal(err)
		}
		// rel is <provider>/<kind>/<version>/reference.md.
		seen[strings.Split(filepath.ToSlash(rel), "/")[0]] = true
	}
	providers := make([]string, 0, len(seen))
	for provider := range seen {
		providers = append(providers, provider)
	}
	sort.Strings(providers)
	return providers
}

// repoRoot walks up from this test file to the directory containing go.mod.
// Under the Bazel sandbox the checkout (and its committed reference.md
// files) is not available, so the gate runs via plain `go test` and skips
// explicitly under Bazel.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			if os.Getenv("TEST_WORKSPACE") != "" {
				t.Skip("skipping drift gate under the Bazel sandbox: the repo checkout is not available; the gate runs via `go test`")
			}
			t.Fatal("could not locate repo root (go.mod)")
		}
		dir = parent
	}
}
