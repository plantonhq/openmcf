package actioninventory

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/plantonhq/planton/pkg/iac/permissions"
)

// repoRoot resolves the repository root from this file's location so the
// gate always reads the live tree it ships in.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller path")
	}
	return filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(thisFile))))
}

func packageDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot resolve caller path")
	}
	return filepath.Dir(thisFile)
}

// TestAwsActionsExist is the gate: every AWS action every committed runner
// permissions manifest names must exist in the provider's own inventory
// snapshot. Exact names must match a published spelling (case-insensitive,
// IAM's evaluation semantics); wildcard patterns must match at least one
// published action -- a pattern matching nothing is a fabricated
// permission wearing a wildcard. Providers without an inventory arm (gcp,
// azure, kubernetes) are exempt HERE deliberately: their structural
// validation lives in pkg/iac/permissions, and their existence arms join
// this gate when their machine-readable inventories are proven.
func TestAwsActionsExist(t *testing.T) {
	root := repoRoot(t)
	inv, err := LoadAws(packageDir(t))
	if err != nil {
		t.Fatalf("loading inventory: %v", err)
	}

	discovered, err := permissions.Discover(root)
	if err != nil {
		t.Fatalf("discovering permissions manifests: %v", err)
	}

	referenced := map[string]bool{}
	for provider, components := range discovered {
		for _, component := range components {
			manifest, err := permissions.Load(root, provider, component)
			if err != nil {
				t.Fatalf("loading %s/%s: %v", provider, component, err)
			}
			for _, statement := range manifest.GetSpec().GetAws().GetStatements() {
				for _, action := range statement.GetActions() {
					prefix, name, found := strings.Cut(action, ":")
					if !found {
						t.Errorf("%s/%s: action %q has no service prefix", provider, component, action)
						continue
					}
					referenced[prefix] = true
					published := inv.ServiceActions(prefix)
					if published == nil {
						t.Errorf("%s/%s: action %q names service %q which the inventory snapshot does not cover -- run `make generate-action-inventory`", provider, component, action, prefix)
						continue
					}
					if MatchAction(published, name) == 0 {
						t.Errorf("%s/%s: action %q does not exist in AWS's service reference for %q -- the name is invented or misspelled", provider, component, action, prefix)
					}
				}
			}
		}
	}

	// Dead weight is rejected like the price books' dead prices: a
	// snapshot service no manifest references bloats every refresh for
	// nothing -- the fetcher scopes to referenced services, so a stale
	// extra means the manifests moved and the snapshot did not.
	for _, svc := range inv.Services {
		if !referenced[svc.Prefix] {
			t.Errorf("inventory covers service %q which no permissions manifest references -- run `make generate-action-inventory`", svc.Prefix)
		}
	}
}

// TestMatchAction pins the matching semantics the gate stands on.
func TestMatchAction(t *testing.T) {
	published := []string{"CreateFunction", "DeleteFunction", "GetFunction", "GetFunctionConfiguration", "ListFunctions", "TagResource"}
	cases := []struct {
		name    string
		matches int
	}{
		{"CreateFunction", 1},
		{"createfunction", 1}, // IAM evaluates action names case-insensitively
		{"Get*", 2},           // wildcard expands against published names
		{"GetFunction?onfiguration", 1},
		{"*", 6},
		{"List*", 1},
		{"DeleteFunctionScalingConfig", 0}, // the invented-name class the gate exists for
		{"Create*Url", 0},                  // a wildcard matching nothing is still a fabrication
	}
	for _, c := range cases {
		if got := MatchAction(published, c.name); got != c.matches {
			t.Errorf("MatchAction(%q) = %d, want %d", c.name, got, c.matches)
		}
	}
}

// TestRenderLoadRoundTrip proves the renderer and the strict loader agree:
// what Render writes, LoadAws accepts, byte-stably.
func TestRenderLoadRoundTrip(t *testing.T) {
	inv := &Inventory{
		Provider: "aws",
		Services: []Service{
			{
				Prefix:         "lambda",
				SourceURL:      "https://servicereference.us-east-1.amazonaws.com/v1/lambda/lambda.json",
				SourceModified: "2026-08-01",
				RetrievedOn:    "2026-08-15",
				Actions:        []string{"CreateFunction", "DeleteFunction"},
			},
			{
				Prefix:         "s3",
				SourceURL:      "https://servicereference.us-east-1.amazonaws.com/v1/s3/s3.json",
				SourceModified: "2026-08-02",
				RetrievedOn:    "2026-08-15",
				Actions:        []string{"GetObject", "PutObject"},
			},
		},
	}
	dir := t.TempDir()
	rendered := Render(inv)
	if err := writeFile(t, filepath.Join(dir, AwsFileName), rendered); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadAws(dir)
	if err != nil {
		t.Fatalf("round-trip load: %v", err)
	}
	if Render(loaded) != rendered {
		t.Error("render -> load -> render is not byte-stable")
	}
}

// TestLoadAwsRefusals pins the loader's structural invariants: unsorted or
// duplicated content refuses loudly rather than letting gate results
// depend on snapshot ordering accidents.
func TestLoadAwsRefusals(t *testing.T) {
	base := func() *Inventory {
		return &Inventory{
			Provider: "aws",
			Services: []Service{{
				Prefix:         "lambda",
				SourceURL:      "https://example.invalid/lambda.json",
				SourceModified: "2026-08-01",
				RetrievedOn:    "2026-08-15",
				Actions:        []string{"CreateFunction", "DeleteFunction"},
			}},
		}
	}
	cases := []struct {
		name    string
		mutate  func(*Inventory)
		wantErr string
	}{
		{"wrong provider", func(i *Inventory) { i.Provider = "gcp" }, "provider"},
		{"unsorted actions", func(i *Inventory) { i.Services[0].Actions = []string{"DeleteFunction", "CreateFunction"} }, "not sorted"},
		{"duplicate action", func(i *Inventory) { i.Services[0].Actions = []string{"CreateFunction", "CreateFunction"} }, "duplicates"},
		{"empty actions", func(i *Inventory) { i.Services[0].Actions = nil }, "no actions"},
		{"missing provenance", func(i *Inventory) { i.Services[0].RetrievedOn = "" }, "required"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inv := base()
			c.mutate(inv)
			dir := t.TempDir()
			// Render sorts, so unsorted/duplicate cases write raw YAML to
			// exercise the loader against genuinely malformed bytes.
			raw := renderRaw(inv)
			if err := writeFile(t, filepath.Join(dir, AwsFileName), raw); err != nil {
				t.Fatal(err)
			}
			_, err := LoadAws(dir)
			if err == nil || !strings.Contains(err.Error(), c.wantErr) {
				t.Errorf("LoadAws error = %v, want it to contain %q", err, c.wantErr)
			}
		})
	}
}

// renderRaw writes an inventory verbatim, WITHOUT Render's sorting, so
// refusal tests can present genuinely malformed snapshots.
func renderRaw(inv *Inventory) string {
	var b strings.Builder
	b.WriteString("provider: " + inv.Provider + "\n")
	b.WriteString("services:\n")
	for _, svc := range inv.Services {
		b.WriteString("  - prefix: " + svc.Prefix + "\n")
		b.WriteString("    source_url: " + svc.SourceURL + "\n")
		b.WriteString("    source_modified: \"" + svc.SourceModified + "\"\n")
		b.WriteString("    retrieved_on: \"" + svc.RetrievedOn + "\"\n")
		if len(svc.Actions) == 0 {
			b.WriteString("    actions: []\n")
			continue
		}
		b.WriteString("    actions:\n")
		for _, action := range svc.Actions {
			b.WriteString("      - " + action + "\n")
		}
	}
	return b.String()
}

func writeFile(t *testing.T, path, content string) error {
	t.Helper()
	return os.WriteFile(path, []byte(content), 0o644)
}
