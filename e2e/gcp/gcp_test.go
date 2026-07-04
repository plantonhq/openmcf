//go:build e2e

// Package gcp contains end-to-end tests that provision real GCP resources via
// Planton IaC modules and verify them through the Google Cloud APIs.
// Credentials come from the ambient ADC chain (local
// `gcloud auth application-default login` or workload identity federation in
// CI -- never a stored secret); see the aa_e2e harness package. The test
// project resolves from E2E_GCP_PROJECT / GOOGLE_PROJECT / the ADC credential.
//
// Run with: go test -tags=e2e -timeout=120m -v ./e2e/gcp/...
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

// --- GCP Health Check (the load-balancing family's leaf; global + regional scenarios) ---

func TestGcpHealthCheck_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcphealthcheck", "pulumi")
}
func TestGcpHealthCheck_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcphealthcheck", "terraform")
}

// --- GCP Backend Bucket (composed static origin: deploys a GcpGcsBucket prerequisite) ---

func TestGcpBackendBucket_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpbackendbucket", "pulumi")
}
func TestGcpBackendBucket_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpbackendbucket", "terraform")
}

// --- GCP Backend Service (composed LB hub: deploys a GcpHealthCheck prerequisite) ---

func TestGcpBackendService_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpbackendservice", "pulumi")
}
func TestGcpBackendService_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpbackendservice", "terraform")
}

// --- GCP Subnetwork (composed address plan: deploys a GcpVpc prerequisite) ---

func TestGcpSubnetwork_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpsubnetwork", "pulumi")
}
func TestGcpSubnetwork_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpsubnetwork", "terraform")
}

// --- GCP Region Network Endpoint Group (serverless NEG bridge) ---

func TestGcpRegionNetworkEndpointGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpregionnetworkendpointgroup", "pulumi")
}
func TestGcpRegionNetworkEndpointGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpregionnetworkendpointgroup", "terraform")
}

// --- GCP URL Map (composed routing: deploys GcpBackendService + GcpHealthCheck prerequisites) ---

func TestGcpUrlMap_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpurlmap", "pulumi")
}
func TestGcpUrlMap_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpurlmap", "terraform")
}

// --- GCP Managed SSL Certificate (global managed cert leaf) ---

func TestGcpManagedSslCertificate_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpmanagedsslcertificate", "pulumi")
}
func TestGcpManagedSslCertificate_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpmanagedsslcertificate", "terraform")
}

// --- GCP Target HTTP Proxy (composed frontend adapter: deploys the GcpUrlMap prerequisite chain) ---

func TestGcpTargetHttpProxy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcptargethttpproxy", "pulumi")
}
func TestGcpTargetHttpProxy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcptargethttpproxy", "terraform")
}

// --- GCP Target HTTPS Proxy (composed TLS frontend: deploys GcpUrlMap + GcpManagedSslCertificate prerequisites) ---

func TestGcpTargetHttpsProxy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcptargethttpsproxy", "pulumi")
}
func TestGcpTargetHttpsProxy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcptargethttpsproxy", "terraform")
}

// --- GCP Global Forwarding Rule (the VIP node; the deepest composed chain in the GCP harness) ---

func TestGcpGlobalForwardingRule_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpglobalforwardingrule", "pulumi")
}
func TestGcpGlobalForwardingRule_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpglobalforwardingrule", "terraform")
}

// --- GCP SSL Policy (TLS version/cipher hardening; leaf with global + regional scenarios) ---

func TestGcpSslPolicy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpsslpolicy", "pulumi")
}
func TestGcpSslPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpsslpolicy", "terraform")
}

// --- GCP SSL Certificate (self-managed cert upload; leaf with a checked-in throwaway keypair) ---

func TestGcpSslCertificate_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpsslcertificate", "pulumi")
}
func TestGcpSslCertificate_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpsslcertificate", "terraform")
}

// --- GCP Service Networking Connection (private services access peering) ---

func TestGcpServiceNetworkingConnection_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpservicenetworkingconnection", "pulumi")
}
func TestGcpServiceNetworkingConnection_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpservicenetworkingconnection", "terraform")
}

// --- GCP Address (regional static IP reservation) ---

func TestGcpAddress_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpaddress", "pulumi")
}
func TestGcpAddress_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpaddress", "terraform")
}

// --- GCP VPC (deep-rebuilt network; leaf scenario) ---

func TestGcpVpc_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpvpc", "pulumi")
}
func TestGcpVpc_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpvpc", "terraform")
}

// --- GCP Cloud SQL (composed PSA chain + public/private instance scenarios) ---

func TestGcpCloudSql_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpcloudsql", "pulumi")
}
func TestGcpCloudSql_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpcloudsql", "terraform")
}

// --- GCP Cloud SQL Database (composed instance prerequisite) ---

func TestGcpCloudSqlDatabase_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpcloudsqldatabase", "pulumi")
}
func TestGcpCloudSqlDatabase_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpcloudsqldatabase", "terraform")
}

// --- GCP Cloud SQL User (composed instance prerequisite) ---

func TestGcpCloudSqlUser_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpcloudsqluser", "pulumi")
}
func TestGcpCloudSqlUser_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpcloudsqluser", "terraform")
}

// GcpRedisInstance scenarios (Memorystore for Redis; the PSA scenario rides
// the same service networking chain Cloud SQL private IP uses).
func TestGcpRedisInstance_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcpredisinstance", "pulumi")
}
func TestGcpRedisInstance_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcpredisinstance", "terraform")
}

// GcpRouterNat scenarios (Cloud Router + NAT gateway; the manual scenario
// proves NAT IPs referencing GcpAddress reservations).
func TestGcpRouterNat_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "gcprouternat", "pulumi")
}
func TestGcpRouterNat_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "gcprouternat", "terraform")
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
