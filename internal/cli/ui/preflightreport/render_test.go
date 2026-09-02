package preflightreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plantonhq/planton/pkg/setdeploy"
)

// The rendered report is pinned by golden files — the refusal sentences and
// the report's shape are the offline lane's product surface, so a wording
// change must be a visible diff someone chose, never an accident.
//
// Regenerate with:  PLANTON_REGEN_PREFLIGHT_GOLDENS=1 go test ./internal/cli/ui/preflightreport/
const regenEnvVar = "PLANTON_REGEN_PREFLIGHT_GOLDENS"

// refusedReport is a hand-built report exercising every severity and every
// line shape the renderer knows: verified facts, refusals with and without
// sources, warnings, and assumptions.
func refusedReport() *setdeploy.Report {
	return &setdeploy.Report{Checks: []setdeploy.Check{
		{
			Name: "load-and-schema", Title: "Manifests load and validate",
			Verified: []string{"2 of 3 documents load as known kinds and pass schema validation"},
			Entries: []setdeploy.Entry{{
				Severity: setdeploy.SeverityRefusal, Source: "manifests/cache.yaml",
				Message: "spec.memory_size_gb: value is required",
			}},
		},
		{
			Name: "references", Title: "References resolve inside this set",
			Entries: []setdeploy.Entry{
				{
					Severity: setdeploy.SeverityRefusal, Source: "manifests/service.yaml", FieldPath: "spec.redis_host",
					Message: "spec.redis_host: references GcpMemorystoreRedis \"cache\" outside this set; the set does not deploy it — the value must come from a resource that already exists — no backend exists here to discover it; add its manifest to this set, or deploy connected",
				},
				{
					Severity: setdeploy.SeverityAssumption, Source: "manifests/service.yaml",
					Message: "relationship contains GcpProject \"acme-prod\" is outside this set — its existence is assumed and verified by the module at apply",
				},
			},
		},
		{
			Name: "state-backend", Title: "State backends are configured and reachable",
			Verified: []string{"s3 state bucket \"acme-state\" reachable"},
			Entries: []setdeploy.Entry{{
				Severity: setdeploy.SeverityWarning,
				Message:  "1 node(s) keep state on this machine only — re-runs must happen here; name a remote backend (s3/gcs/azurerm) for CI",
			}},
		},
	}}
}

func passingReport() *setdeploy.Report {
	return &setdeploy.Report{Checks: []setdeploy.Check{
		{
			Name: "load-and-schema", Title: "Manifests load and validate",
			Verified: []string{"2 of 2 documents load as known kinds and pass schema validation"},
		},
		{
			Name: "cycles", Title: "Dependencies form a deployable order",
			Verified: []string{"deploy order: GcpMemorystoreRedis/cache@prod -> GcpCloudRun/storefront@prod"},
		},
		{
			Name: "provider-credentials", Title: "Provider credentials authenticate",
			Verified: []string{"gcp credentials authenticate"},
			Entries: []setdeploy.Entry{{
				Severity: setdeploy.SeverityAssumption,
				Message:  "kubernetes access is verified per kube context",
			}},
		},
	}}
}

func TestReportGoldens(t *testing.T) {
	cases := []struct {
		name   string
		report *setdeploy.Report
	}{
		{"refused", refusedReport()},
		{"passing", passingReport()},
	}
	regen := os.Getenv(regenEnvVar) != ""
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var rendered strings.Builder
			for _, line := range FormatReportLines(tc.report) {
				rendered.WriteString(line + "\n")
			}
			rendered.WriteString("\n" + FormatVerdict(tc.report) + "\n")

			goldenPath := filepath.Join("testdata", tc.name+".golden.txt")
			if regen {
				if err := os.MkdirAll("testdata", 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(goldenPath, []byte(rendered.String()), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("missing golden (regenerate with %s=1): %v", regenEnvVar, err)
			}
			if rendered.String() != string(want) {
				t.Fatalf("rendered report drifted from its golden — a wording change must be deliberate.\ngot:\n%s\nwant:\n%s", rendered.String(), want)
			}
		})
	}
}

// The verdict's refusal arm counts and pluralizes correctly.
func TestVerdictRefusalCount(t *testing.T) {
	report := refusedReport()
	verdict := FormatVerdict(report)
	if !strings.Contains(verdict, "2 problems named above") {
		t.Fatalf("verdict must count refusals: %s", verdict)
	}
	if !strings.Contains(verdict, "nothing was handed to an IaC engine") {
		t.Fatalf("verdict must state the guarantee: %s", verdict)
	}
}
