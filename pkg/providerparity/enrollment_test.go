//go:build !codegen
// +build !codegen

package providerparity

import (
	"testing"

	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// The enrollment table is the write-scope of the shared baseline: malformed
// entries would silently narrow every baseline write, so its shape is gated
// hermetically.
func TestEnrollmentsWellFormed(t *testing.T) {
	if len(Enrollments) == 0 {
		t.Fatal("no enrolled providers -- the gate would vacuously pass")
	}
	seen := map[cloudresourcekind.CloudResourceProvider]bool{}
	for _, e := range Enrollments {
		if e.Provider == cloudresourcekind.CloudResourceProvider_cloud_resource_provider_unspecified {
			t.Error("enrollment with unspecified provider")
		}
		if e.GASchema == "" {
			t.Errorf("%s: enrollment without a GA schema", e.Provider)
		}
		if seen[e.Provider] {
			t.Errorf("%s: enrolled twice -- one GA schema per provider", e.Provider)
		}
		seen[e.Provider] = true
	}
}

// A baseline write consumes the MERGED findings of every enrollment; losing
// any provider's findings in the merge would truncate the shared baseline
// exactly the way the single-provider write path used to.
func TestMergeFindingsPreservesEveryProvider(t *testing.T) {
	accountings := []Accounting{
		{CloudProvider: "gcp", Findings: []Finding{
			{BaselineKey: "kind:GcpGcsBucket", Detail: "a"},
			{BaselineKey: "resource:google_widget", Detail: "b"},
		}},
		{CloudProvider: "aws", Findings: []Finding{
			{BaselineKey: "kind:AwsS3Bucket", Detail: "c"},
		}},
	}
	merged := MergeFindings(accountings)
	if len(merged) != 3 {
		t.Fatalf("merged %d findings, want 3", len(merged))
	}
	keys := map[string]bool{}
	for _, f := range merged {
		keys[f.BaselineKey] = true
	}
	for _, want := range []string{"kind:GcpGcsBucket", "resource:google_widget", "kind:AwsS3Bucket"} {
		if !keys[want] {
			t.Errorf("merged findings lost %s", want)
		}
	}
}
