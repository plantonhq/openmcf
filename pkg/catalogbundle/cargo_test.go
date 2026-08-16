package catalogbundle

import (
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"

	costestimatev1 "github.com/plantonhq/planton/finops/componentcostestimate/v1"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// The fact-sheet cargo end to end over the REAL tree: a covered exemplar's
// four documents ride the bundle byte-identical to their tree sources, the
// central documents ride beside them, and the exemplar's entry summaries
// agree with an INDEPENDENT read of the tree -- computed here from the
// estimate document itself, so a price refresh can never rot this test and
// a projection bug can never hide behind its own arithmetic.
func TestFactSheetCargoRoundTrip(t *testing.T) {
	dir := t.TempDir()
	descriptorsPath := filepath.Join(dir, "descriptors.binpb")
	writeLinkedDescriptorSet(t, descriptorsPath)

	bundlePath := filepath.Join(dir, "catalog-bundle.zip")
	if _, err := Build(BuildInput{
		DescriptorSetPath: descriptorsPath,
		CatalogDir:        catalogDir(t),
		OutputPath:        bundlePath,
	}); err != nil {
		t.Fatalf("build failed: %v", err)
	}
	bundle, err := Load(bundlePath)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// The exemplar's cargo is aboard, byte-identical to the tree.
	for entryName, treePath := range map[string]string{
		"costs/aws/awsalb.yaml":       filepath.Join(catalogDir(t), "aws", "awsalb", "cost.yaml"),
		"controls/aws/awsalb.yaml":    filepath.Join(catalogDir(t), "aws", "awsalb", "controls.yaml"),
		"permissions/aws/awsalb.yaml": filepath.Join(catalogDir(t), "aws", "awsalb", "iac", "permissions.yaml"),
		"estimates/aws/awsalb.yaml":   filepath.Join(catalogDir(t), "_pricing", "estimates", "awsalb.yaml"),
		"derivations/aws/awsalb.yaml": filepath.Join(catalogDir(t), "_pricing", "derivations", "awsalb.yaml"),
	} {
		want, err := os.ReadFile(treePath)
		if err != nil {
			t.Fatalf("reading %s: %v", treePath, err)
		}
		got, aboardOK := bundle.entries[entryName]
		if !aboardOK {
			t.Errorf("%s is missing from the bundle", entryName)
			continue
		}
		if string(got) != string(want) {
			t.Errorf("%s is not byte-identical to its tree source %s", entryName, treePath)
		}
	}

	// The central documents ride with the component cargo.
	if _, aboardOK := bundle.Compliance()[controlsCatalogEntryName]; !aboardOK {
		t.Error("the control catalog is missing from the bundle")
	}
	if len(bundle.PriceBooks()) == 0 {
		t.Error("no price book is aboard")
	}
	crosswalks := 0
	for name := range bundle.Compliance() {
		if strings.HasPrefix(name, frameworksEntryPrefix) {
			crosswalks++
		}
	}
	if crosswalks == 0 {
		t.Error("no framework crosswalk is aboard")
	}

	// The exemplar's entry summaries agree with an independent read of the
	// tree's own estimate document.
	var exemplar *CatalogEntry
	for _, entry := range bundle.CatalogEntries() {
		if entry.Kind == "AwsAlb" {
			e := entry
			exemplar = &e
			break
		}
	}
	if exemplar == nil {
		t.Fatal("AwsAlb has no catalog entry")
	}
	if exemplar.CostSummary == nil || exemplar.ControlSummary == nil || exemplar.PermissionsProvenance == "" {
		t.Fatalf("AwsAlb is covered but its entry carries incomplete summaries: %+v", exemplar)
	}

	estimateRaw, err := os.ReadFile(filepath.Join(catalogDir(t), "_pricing", "estimates", "awsalb.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	estimate := &costestimatev1.ComponentCostEstimate{}
	if err := protobufyaml.LoadYamlBytes(estimateRaw, estimate); err != nil {
		t.Fatal(err)
	}
	var wantMin, wantMax *big.Rat
	var wantMinStr, wantMaxStr string
	for _, preset := range estimate.GetSpec().GetPresets() {
		total := preset.GetTotalListCost()
		if total == "" {
			continue
		}
		rat, ratOK := new(big.Rat).SetString(total)
		if !ratOK {
			t.Fatalf("preset %s total %q is not a decimal", preset.GetPreset(), total)
		}
		if wantMin == nil || rat.Cmp(wantMin) < 0 {
			wantMin, wantMinStr = rat, total
		}
		if wantMax == nil || rat.Cmp(wantMax) > 0 {
			wantMax, wantMaxStr = rat, total
		}
	}
	if exemplar.CostSummary.MonthlyMin != wantMinStr || exemplar.CostSummary.MonthlyMax != wantMaxStr {
		t.Errorf("AwsAlb cost range = [%s, %s], independent read of its estimate says [%s, %s]",
			exemplar.CostSummary.MonthlyMin, exemplar.CostSummary.MonthlyMax, wantMinStr, wantMaxStr)
	}
	if exemplar.CostSummary.Currency != "USD" {
		t.Errorf("AwsAlb currency = %q, want USD", exemplar.CostSummary.Currency)
	}
	if exemplar.CostSummary.BillingModel != "hybrid" {
		t.Errorf("AwsAlb billing model = %q, want hybrid", exemplar.CostSummary.BillingModel)
	}
	if got := exemplar.ControlSummary.EnforcedByDefault + exemplar.ControlSummary.Configurable +
		exemplar.ControlSummary.NotApplicable; got == 0 {
		t.Error("AwsAlb control summary counts nothing -- the profile examines every catalog control")
	}
	if exemplar.PermissionsProvenance != "derived" {
		t.Errorf("AwsAlb permissions provenance = %q, want derived", exemplar.PermissionsProvenance)
	}

	// An uncovered component's entry document carries no summary keys at
	// all (omitempty holds) -- absence is the honest state. The exemplar is
	// chosen DYNAMICALLY (the first entry document, in name order, whose
	// kind ships no cost cargo) so growing coverage can never rot this
	// test; when every component is covered, there is nothing left to
	// assert and the honest-absence property holds vacuously.
	uncoveredName := ""
	for name := range bundle.entries {
		if !strings.HasPrefix(name, "entries/") {
			continue
		}
		costName := "costs/" + strings.TrimPrefix(name, "entries/")
		if _, covered := bundle.entries[costName]; covered {
			continue
		}
		if uncoveredName == "" || name < uncoveredName {
			uncoveredName = name
		}
	}
	if uncoveredName == "" {
		t.Log("every component ships fact-sheets -- no uncovered entry left to assert honest absence on")
		return
	}
	var uncovered map[string]any
	if err := yaml.Unmarshal(bundle.entries[uncoveredName], &uncovered); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"costSummary", "controlSummary", "permissionsProvenance"} {
		if _, present := uncovered[key]; present {
			t.Errorf("%s ships no fact-sheets but its entry document carries %q", uncoveredName, key)
		}
	}
}

// Fact-sheet cargo violations are conformance findings with the offenders
// named: a document keyed by no user-facing kind, partial coverage, an
// entry claiming summaries it has no cargo for, and a malformed document.
func TestConformanceRefusesIncoherentCargo(t *testing.T) {
	descriptors := linkedDescriptorSetBytes(t)
	controlsDoc := []byte("apiVersion: compliance.planton.dev/v1\nkind: ComponentControlProfile\n")

	cases := []struct {
		name        string
		contents    map[string][]byte
		wantFinding string
	}{
		{
			name: "cargo keyed by no user-facing kind",
			contents: map[string][]byte{
				"costs/aws/awsimaginarything.yaml": []byte("apiVersion: finops.planton.dev/v1\nkind: ComponentCostProfile\n"),
			},
			wantFinding: "not keyed by a user-facing registry kind",
		},
		{
			name: "partial coverage",
			contents: map[string][]byte{
				"controls/aws/awsvpc.yaml": controlsDoc,
			},
			wantFinding: "whole-or-not-at-all",
		},
		{
			name: "entry claims summaries without cargo",
			contents: map[string][]byte{
				"entries/aws/awsvpc.yaml": mustMarshalEntry(t, CatalogEntry{
					Kind: "AwsVpc", Title: "AWS VPC", Description: "d", Slug: "aws-vpc",
					IacModules:            CatalogEntryIacModules{TerraformModuleDir: "catalog/aws/awsvpc/iac/tf"},
					PermissionsProvenance: "derived",
				}),
			},
			wantFinding: "no complete cargo",
		},
		{
			name: "malformed cargo document",
			contents: map[string][]byte{
				"costs/aws/awsvpc.yaml": []byte("kind: ComponentCostProfile\nnot_a_field: true\n"),
			},
			wantFinding: "does not parse against its schema",
		},
		{
			name: "malformed derivation document",
			contents: map[string][]byte{
				"derivations/aws/awsvpc.yaml": []byte("kind: ComponentCostDerivation\nnot_a_field: true\n"),
			},
			wantFinding: "does not parse against its schema",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.contents["descriptors.binpb"] = descriptors
			bundle, err := Load(writeRawBundle(t, c.contents))
			if err != nil {
				t.Fatalf("structurally sound bundles must load: %v", err)
			}
			err = CheckConformance(bundle)
			if err == nil {
				t.Fatal("incoherent cargo must fail conformance")
			}
			if !strings.Contains(err.Error(), c.wantFinding) {
				t.Fatalf("the refusal must contain %q, got: %v", c.wantFinding, err)
			}
		})
	}
}

// A tree whose fact-sheet coverage is partial fails the BUILD with the
// component named -- the whole-or-not-at-all standard holds at packing
// time, not just at conformance time.
func TestCollectCargoRefusesPartialCoverage(t *testing.T) {
	dir := t.TempDir()
	componentDir := filepath.Join(dir, "aws", "awswidget")
	if err := os.MkdirAll(componentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	costDoc := "apiVersion: finops.planton.dev/v1\nkind: ComponentCostProfile\nmetadata:\n  name: awswidget\nspec:\n  billingModel: usage_based\n"
	if err := os.WriteFile(filepath.Join(componentDir, "cost.yaml"), []byte(costDoc), 0o644); err != nil {
		t.Fatal(err)
	}

	entries := map[string][]byte{}
	if _, err := collectCargo(dir, entries); err == nil {
		t.Fatal("a cost profile without controls and permissions must fail cargo collection")
	} else if !strings.Contains(err.Error(), "whole-or-not-at-all") {
		t.Fatalf("the refusal must name the standard, got: %v", err)
	}
}

// A derivation whose component ships no cost profile fails the BUILD with
// the file named -- derivations price covered components only.
func TestCollectCargoRefusesOrphanDerivation(t *testing.T) {
	dir := t.TempDir()
	derivationsDir := filepath.Join(dir, "_pricing", "derivations")
	if err := os.MkdirAll(derivationsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	derivationDoc := "apiVersion: finops.planton.dev/v1\nkind: ComponentCostDerivation\nmetadata:\n  name: awswidget\nspec:\n  currency: USD\n"
	if err := os.WriteFile(filepath.Join(derivationsDir, "awswidget.yaml"), []byte(derivationDoc), 0o644); err != nil {
		t.Fatal(err)
	}

	entries := map[string][]byte{}
	if _, err := collectCargo(dir, entries); err == nil {
		t.Fatal("a derivation without a covered component must fail cargo collection")
	} else if !strings.Contains(err.Error(), "derivations price covered components only") {
		t.Fatalf("the refusal must name the rule, got: %v", err)
	}
}
