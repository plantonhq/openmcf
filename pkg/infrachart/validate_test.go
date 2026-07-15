package infrachart

import (
	"strings"
	"testing"
)

func TestValidateChartValidFixture(t *testing.T) {
	report, err := ValidateChart("testdata/valid-chart")
	if err != nil {
		t.Fatalf("ValidateChart failed: %v", err)
	}
	if !report.Valid() {
		for _, v := range report.Variants {
			for _, e := range v.Errors {
				t.Errorf("[%s] %s: %v", v.Name, e.Template, e.Err)
			}
		}
		t.Fatal("expected the valid fixture chart to pass")
	}

	// defaults + the subnetEnabled toggle flipped.
	if len(report.Variants) != 2 {
		t.Fatalf("expected 2 variants, got %d", len(report.Variants))
	}
	if report.Variants[0].Docs != 2 {
		t.Fatalf("defaults variant: expected 2 docs (vpc + subnet), got %d", report.Variants[0].Docs)
	}
	if report.Variants[1].Name != "subnetEnabled=false" {
		t.Fatalf("unexpected variant name %q", report.Variants[1].Name)
	}
	if report.Variants[1].Docs != 1 {
		t.Fatalf("flipped variant: expected 1 doc (vpc only), got %d", report.Variants[1].Docs)
	}
}

func TestValidateChartBrokenFixture(t *testing.T) {
	report, err := ValidateChart("testdata/broken-chart")
	if err != nil {
		t.Fatalf("ValidateChart failed: %v", err)
	}
	if report.Valid() {
		t.Fatal("expected the broken fixture chart to fail")
	}
	if len(report.Variants) != 1 {
		t.Fatalf("expected 1 variant (no bool params), got %d", len(report.Variants))
	}

	errs := report.Variants[0].Errors
	if len(errs) != 4 {
		for _, e := range errs {
			t.Logf("error: %s doc %d: %v", e.Template, e.DocIndex, e.Err)
		}
		t.Fatalf("expected 4 errors (one per defect class), got %d", len(errs))
	}

	// One error per defect class: per-document errors in document order, then
	// the variant-wide intra-chart target check.
	assertErrContains(t, errs[0].Err.Error(), "cidrBlok")
	assertErrContains(t, errs[1].Err.Error(), "region")
	assertErrContains(t, errs[2].Err.Error(), "does_not_exist")
	assertErrContains(t, errs[3].Err.Error(), "phantom-vpc")
	assertErrContains(t, errs[3].Err.Error(), "does not define")
}

func assertErrContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Errorf("error %q does not mention %q", got, want)
	}
}
