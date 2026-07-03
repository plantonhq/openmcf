//go:build e2e

// Package gcp contains end-to-end tests that provision real GCP resources via
// Planton IaC modules and verify them through the Google Cloud APIs.
// Credentials come from the ambient ADC chain (local
// `gcloud auth application-default login` or workload identity federation in
// CI -- never a stored secret); see the aa_e2e harness package. The test
// project resolves from E2E_GCP_PROJECT / GOOGLE_PROJECT / the ADC credential.
//
// Run with: go test -tags=e2e -timeout=30m -v ./e2e/gcp/...
package gcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	gcpe2e "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/aa_e2e"
	"github.com/plantonhq/planton/e2e/framework/discovery"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/e2e/framework/runner"
)

var (
	testHarness      *gcpe2e.Harness
	repoRoot         string
	runID            string
	pulumiBackendURL string
)

func TestMain(m *testing.M) {
	var err error
	repoRoot, err = filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve repo root: %v\n", err)
		os.Exit(1)
	}

	runID = uuid.New().String()[:8]

	backendDir, err := os.MkdirTemp("", "planton-e2e-gcp-pulumi-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create temp backend dir: %v\n", err)
		os.Exit(1)
	}
	pulumiBackendURL = "file://" + backendDir
	defer os.RemoveAll(backendDir)

	if err := runner.PulumiLogin(pulumiBackendURL); err != nil {
		fmt.Fprintf(os.Stderr, "failed to login to pulumi backend: %v\n", err)
		os.Exit(1)
	}

	testHarness = gcpe2e.NewHarness()
	ctx := context.Background()
	if err := testHarness.Setup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup GCP harness: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := testHarness.Teardown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to teardown GCP harness: %v\n", err)
	}

	os.Exit(code)
}

// --- GCP Service Account (the identity leaf everything references) ---

func TestGcpServiceAccount_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpserviceaccount", "pulumi")
}
func TestGcpServiceAccount_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpserviceaccount", "terraform")
}

// --- GCP IAM Custom Role (least-privilege permission bundle) ---

func TestGcpIamCustomRole_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpiamcustomrole", "pulumi")
}
func TestGcpIamCustomRole_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpiamcustomrole", "terraform")
}

// --- GCP Project IAM Member (composed grant: deploys GcpServiceAccount + GcpIamCustomRole prerequisites) ---

func TestGcpProjectIamMember_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpprojectiammember", "pulumi")
}
func TestGcpProjectIamMember_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpprojectiammember", "terraform")
}

// --- GCP Workload Identity Pool (the keyless-auth trust boundary) ---

func TestGcpWorkloadIdentityPool_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpworkloadidentitypool", "pulumi")
}
func TestGcpWorkloadIdentityPool_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpworkloadidentitypool", "terraform")
}

// --- GCP Workload Identity Pool Provider (composed issuer: deploys a GcpWorkloadIdentityPool prerequisite) ---

func TestGcpWorkloadIdentityPoolProvider_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpworkloadidentitypoolprovider", "pulumi")
}
func TestGcpWorkloadIdentityPoolProvider_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpworkloadidentitypoolprovider", "terraform")
}

// runAllScenariosForComponent discovers and runs all E2E scenarios for a GCP component.
func runAllScenariosForComponent(t *testing.T, component, engine string) {
	t.Helper()

	var moduleDir string
	switch engine {
	case "pulumi":
		moduleDir = filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "gcp", component, "v1", "iac", "pulumi")
	case "terraform":
		moduleDir = filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "gcp", component, "v1", "iac", "tf")
	default:
		t.Fatalf("unsupported engine: %s", engine)
	}

	if !fileExists(moduleDir) {
		t.Skipf("component %s %s module not found at %s", component, engine, moduleDir)
	}

	scenarios, err := discovery.DiscoverTestScenarios(repoRoot, "gcp", component)
	if err != nil {
		t.Fatalf("failed to discover test scenarios for %s: %v", component, err)
	}

	if len(scenarios) == 0 {
		t.Skipf("no test scenarios found for %s", component)
	}

	t.Logf("Discovered %d scenarios for %s [%s]", len(scenarios), component, engine)

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.Name, func(t *testing.T) {
			runSingleScenario(t, component, moduleDir, engine, scenario)
		})
	}
}

func runSingleScenario(t *testing.T, component, moduleDir, engine string, scenario discovery.TestScenario) {
	t.Helper()

	tc := &provider.ComponentTestContext{
		Component:    component,
		Provider:     "gcp",
		Engine:       engine,
		ModuleDir:    moduleDir,
		ManifestPath: scenario.ManifestPath,
		RepoRoot:     repoRoot,
		RunID:        runID,
		T:            t,
	}

	if engine == "pulumi" {
		stackName := runner.GenerateStackName(component+"-"+scenario.Name, runID)
		if len(stackName) > 50 {
			stackName = stackName[:50]
		}
		tc.StackName = stackName
		tc.BackendURL = pulumiBackendURL
	}

	ctx := context.Background()
	result := runner.RunComponentTest(ctx, tc, testHarness)

	for _, phase := range result.Phases {
		status := "PASS"
		if !phase.Passed {
			status = "FAIL"
		}
		t.Logf("  %s: %s (%s)", phase.Phase, status, phase.Duration)
		if phase.Error != nil {
			t.Logf("    Error: %v", phase.Error)
		}
	}

	if !result.Passed {
		t.Fatalf("scenario %s/%s [%s] failed (total: %s)", component, scenario.Name, engine, result.Duration)
	}

	t.Logf("scenario %s/%s [%s] passed (total: %s)", component, scenario.Name, engine, result.Duration)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
