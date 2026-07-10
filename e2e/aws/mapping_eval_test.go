//go:build e2e

package aws

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	evalrunner "github.com/plantonhq/planton/e2e/framework/mappingeval"
	"github.com/plantonhq/planton/pkg/iac/mappingeval"
	"github.com/plantonhq/planton/pkg/iac/mappingeval/baseline"
	"github.com/plantonhq/planton/pkg/iac/mappingeval/inventory"
	"google.golang.org/protobuf/encoding/protojson"
)

// MappingEvalEnvVar opts the mapping-eval lane in. Opt-in like the import
// round-trip: it deploys a whole multi-component suite (create-and-destroy)
// and scans the account, which is a deliberate, scheduled spend -- not part
// of every e2e invocation.
const MappingEvalEnvVar = "PLANTON_E2E_MAPPING_EVAL"

// TestMappingEval_NetworkStaples is the live proof of the mapping eval
// harness: seed the account from the network-staples suite (the answer
// key), scan it back blind through the read-only scanner, have the
// deterministic baseline propose the mapping, and machine-score the
// proposal. The baseline is PINNED to a perfect score on this suite -- any
// drop is a harness, recipe, or scanner regression, never model variance.
//
// Artifacts (scan, proposal, report) are written to
// PLANTON_E2E_MAPPING_EVAL_ARTIFACTS when set, so runs leave durable
// evidence.
func TestMappingEval_NetworkStaples(t *testing.T) {
	if os.Getenv(MappingEvalEnvVar) != "1" {
		t.Skipf("mapping eval lane is opt-in; set %s=1", MappingEvalEnvVar)
	}
	ctx := context.Background()

	suitePath := filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "aws", "aa_eval", "suites", "network-staples.yaml")
	suite, err := mappingeval.LoadSuite(repoRoot, suitePath)
	if err != nil {
		t.Fatalf("loading suite: %v", err)
	}
	region := suite.Suite.GetSpec().GetScanScope().GetRegion()

	// Seed: deploy the answer key. Teardown ALWAYS runs, and its failure
	// fails the test -- leaked fixtures are a real cost, not a footnote.
	deployed, groundTruth, deployErr := evalrunner.DeploySuite(t, repoRoot, "aws", suite)
	defer func() {
		if err := evalrunner.TeardownSuite(t, deployed); err != nil {
			t.Errorf("suite teardown failed -- fixtures may be leaking: %v", err)
		}
	}()
	if deployErr != nil {
		t.Fatalf("deploying suite: %v", deployErr)
	}

	// Blind side: scan the region over the declared type allowlist.
	scanTypes, err := evalrunner.ScanTypeNames(repoRoot, "aws", suite)
	if err != nil {
		t.Fatalf("deriving scan type allowlist: %v", err)
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("loading AWS config: %v", err)
	}
	scan, err := inventory.NewScanner(cfg).Scan(ctx, region, scanTypes)
	if err != nil {
		t.Fatalf("scanning account: %v", err)
	}
	writeArtifact(t, "scan.json", func() ([]byte, error) { return json.MarshalIndent(scan, "", "  ") })

	// Propose + validate against the contract.
	proposal, err := baseline.Propose(scan)
	if err != nil {
		t.Fatalf("baseline proposer: %v", err)
	}
	writeArtifact(t, "proposal.json", func() ([]byte, error) {
		return protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(proposal)
	})
	loaded, err := mappingeval.ParseProposal(proposal)
	if err != nil {
		t.Fatalf("baseline proposal violates the contract: %v", err)
	}

	// Score against the answer key.
	components := make([]string, 0, len(suite.Members))
	for _, member := range suite.Members {
		components = append(components, member.Component)
	}
	scoreOptions, err := mappingeval.ScoreOptionsFromCatalog(repoRoot, "aws", components)
	if err != nil {
		t.Fatalf("deriving score options: %v", err)
	}
	report := mappingeval.Score(groundTruth, loaded, scoreOptions)
	writeArtifact(t, "report.json", func() ([]byte, error) { return json.MarshalIndent(report, "", "  ") })
	t.Logf("mapping eval report:\n%s", report.Summary())

	if !report.Perfect() {
		reportJSON, _ := json.MarshalIndent(report, "", "  ")
		t.Fatalf("baseline must score PERFECT on the seeded suite; it did not.\n%s\n\nfull report:\n%s",
			report.Summary(), string(reportJSON))
	}
}

// writeArtifact persists one run artifact when an artifacts dir is set.
func writeArtifact(t *testing.T, name string, render func() ([]byte, error)) {
	dir := os.Getenv("PLANTON_E2E_MAPPING_EVAL_ARTIFACTS")
	if dir == "" {
		return
	}
	payload, err := render()
	if err != nil {
		t.Logf("rendering artifact %s failed: %v", name, err)
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Logf("creating artifacts dir failed: %v", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, name), payload, 0o644); err != nil {
		t.Logf("writing artifact %s failed: %v", name, err)
	}
}
