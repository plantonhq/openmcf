package manifestgraph

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/plantonhq/planton/internal/manifest"
	goyaml "gopkg.in/yaml.v3"
)

// The corpus under testdata/corpus is the behavior-parity contract: one
// scenario per semantics point, each with a committed golden holding the
// dependency order, derived targets, and classification verdicts. These
// goldens encode the TARGET semantics for every lane that orders manifests —
// a change that shifts a golden is a semantics change and must be made
// deliberately, in the open.
//
// Regenerate with:  PLANTON_REGEN_MANIFESTGRAPH_GOLDENS=1 go test ./pkg/manifestgraph/
const regenEnvVar = "PLANTON_REGEN_MANIFESTGRAPH_GOLDENS"

// golden is the committed shape of one scenario's expected outcome.
type golden struct {
	// Order is the deployment order as identity strings. Empty when the
	// scenario expects no order (a cycle).
	Order []string `yaml:"order"`
	// Derived are implied placement targets not deployed by the set.
	Derived []string `yaml:"derived,omitempty"`
	// Findings are "<class> <node-or-empty> <fieldPath-or-empty>" triples —
	// class and location, deliberately not the message text, so wording can
	// improve without a semantics-change ceremony.
	Findings []string `yaml:"findings,omitempty"`
}

func TestCorpusGoldens(t *testing.T) {
	corpusDir := filepath.Join("testdata", "corpus")
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		t.Fatalf("reading corpus dir: %v", err)
	}

	regen := os.Getenv(regenEnvVar) != ""

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		scenario := entry.Name()
		t.Run(scenario, func(t *testing.T) {
			dir := filepath.Join(corpusDir, scenario)
			actual := runScenario(t, dir)

			goldenPath := filepath.Join(dir, "golden.yaml")
			if regen {
				out, err := goyaml.Marshal(actual)
				if err != nil {
					t.Fatalf("marshaling golden: %v", err)
				}
				if err := os.WriteFile(goldenPath, out, 0o644); err != nil {
					t.Fatalf("writing golden: %v", err)
				}
				return
			}

			raw, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("reading golden (%s to regenerate): %v", regenEnvVar, err)
			}
			var want golden
			if err := goyaml.Unmarshal(raw, &want); err != nil {
				t.Fatalf("parsing golden: %v", err)
			}

			assertSlicesEqual(t, "order", want.Order, actual.Order)
			assertSlicesEqual(t, "derived", want.Derived, actual.Derived)
			assertSlicesEqual(t, "findings", want.Findings, actual.Findings)
		})
	}
}

// runScenario loads every manifest in the scenario dir (sorted filename
// order — the authored order), builds the set and graph, and renders the
// comparable outcome.
func runScenario(t *testing.T, dir string) golden {
	t.Helper()

	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		t.Fatalf("globbing scenario: %v", err)
	}
	sort.Strings(files)

	var items []Item
	for _, f := range files {
		if filepath.Base(f) == "golden.yaml" {
			continue
		}
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("reading %s: %v", f, err)
		}
		msg, err := manifest.LoadManifestBytes(raw, filepath.Base(f))
		if err != nil {
			t.Fatalf("loading %s: %v", f, err)
		}
		items = append(items, Item{Msg: msg, Source: filepath.Base(f)})
	}

	set, setFindings := NewSet(items)
	graph := BuildGraph(set)

	var out golden
	order, cycleFinding := graph.TopoOrder()
	for _, idx := range order {
		out.Order = append(out.Order, set.Nodes[idx].Identity.String())
	}
	for _, d := range graph.Derived {
		out.Derived = append(out.Derived, d.String())
	}

	all := append(append([]Finding{}, setFindings...), graph.Findings...)
	if cycleFinding != nil {
		all = append(all, *cycleFinding)
	}
	for _, f := range all {
		node := ""
		if f.Node != nil {
			node = f.Node.String()
		}
		out.Findings = append(out.Findings, strings.TrimSpace(fmt.Sprintf("%s %s %s", f.Class, node, f.FieldPath)))
	}
	sort.Strings(out.Findings)
	return out
}

func assertSlicesEqual(t *testing.T, label string, want, got []string) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("%s mismatch:\n want: %v\n got:  %v", label, want, got)
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("%s[%d] mismatch:\n want: %q\n got:  %q", label, i, want[i], got[i])
		}
	}
}
