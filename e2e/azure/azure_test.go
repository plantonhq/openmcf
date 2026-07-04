//go:build e2e

// Package azure contains end-to-end tests that provision real Azure resources via
// Planton IaC modules and verify them through the Azure SDK. Credentials come from
// the ambient chain (local `az` login or GitHub Actions OIDC federation -- never a
// stored secret); see the aa_e2e harness package.
//
// Run with: go test -tags=e2e -timeout=30m -v ./e2e/azure/...
package azure

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	azuree2e "github.com/plantonhq/planton/apis/dev/planton/provider/azure/aa_e2e"
	"github.com/plantonhq/planton/e2e/framework/discovery"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/e2e/framework/runner"
)

var (
	testHarness      *azuree2e.Harness
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

	backendDir, err := os.MkdirTemp("", "planton-e2e-azure-pulumi-*")
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

	testHarness = azuree2e.NewHarness()
	ctx := context.Background()
	if err := testHarness.Setup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup Azure harness: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	if err := testHarness.Teardown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to teardown Azure harness: %v\n", err)
	}

	os.Exit(code)
}

// --- Azure Resource Group (walking skeleton: the Layer-0 container every other kind references) ---

func TestAzureResourceGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureresourcegroup", "pulumi")
}
func TestAzureResourceGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureresourcegroup", "terraform")
}

// --- Azure Role Assignment (composed: fixture RG + fixture identity -> Reader grant at RG scope) ---

func TestAzureRoleAssignment_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureroleassignment", "pulumi")
}
func TestAzureRoleAssignment_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureroleassignment", "terraform")
}

// --- Azure Role Definition (composed: custom role scoped at the fixture RG) ---

func TestAzureRoleDefinition_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureroledefinition", "pulumi")
}
func TestAzureRoleDefinition_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureroledefinition", "terraform")
}

// --- Azure Federated Identity Credential (composed: GitHub-shaped trust rule on the fixture identity) ---

func TestAzureFederatedIdentityCredential_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefederatedidentitycredential", "pulumi")
}
func TestAzureFederatedIdentityCredential_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefederatedidentitycredential", "terraform")
}

// --- Azure User Assigned Identity (the identity workloads act as; tags exercised) ---

func TestAzureUserAssignedIdentity_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureuserassignedidentity", "pulumi")
}
func TestAzureUserAssignedIdentity_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureuserassignedidentity", "terraform")
}

// --- Azure Virtual Network (composed: multi-CIDR network in the fixture RG) ---

func TestAzureVirtualNetwork_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualnetwork", "pulumi")
}
func TestAzureVirtualNetwork_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualnetwork", "terraform")
}

// --- Azure Route Table (composed: black-hole + service-tag routes in the fixture RG) ---

func TestAzureRouteTable_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureroutetable", "pulumi")
}
func TestAzureRouteTable_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureroutetable", "terraform")
}

// --- Azure Private DNS Zone (composed: custom internal zone in the fixture RG) ---

func TestAzurePrivateDnsZone_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureprivatednszone", "pulumi")
}
func TestAzurePrivateDnsZone_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureprivatednszone", "terraform")
}

// --- Azure Private DNS Zone Virtual Network Link (composed two-parent chain: fixture zone + fixture network -> link) ---

func TestAzurePrivateDnsZoneVirtualNetworkLink_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureprivatednszonevirtualnetworklink", "pulumi")
}
func TestAzurePrivateDnsZoneVirtualNetworkLink_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureprivatednszonevirtualnetworklink", "terraform")
}

// --- Azure Subnet (the attach-model showcase: fixture network + extra route-table/NSG/NAT fixtures -> subnet with all three seams) ---

func TestAzureSubnet_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuresubnet", "pulumi")
}
func TestAzureSubnet_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuresubnet", "terraform")
}

// --- Azure Network Security Group (singular + plural rule forms in the fixture RG) ---

func TestAzureNetworkSecurityGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurenetworksecuritygroup", "pulumi")
}
func TestAzureNetworkSecurityGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurenetworksecuritygroup", "terraform")
}

// --- Azure Public IP (zone-redundant static address with a scope-hashed DNS label) ---

func TestAzurePublicIp_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurepublicip", "pulumi")
}
func TestAzurePublicIp_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurepublicip", "terraform")
}

// --- Azure Public IP Prefix (smallest reservable range, /31) ---

func TestAzurePublicIpPrefix_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurepublicipprefix", "pulumi")
}
func TestAzurePublicIpPrefix_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurepublicipprefix", "terraform")
}

// --- Azure NAT Gateway (composed: extra public-IP + prefix fixtures -> gateway with both association forms) ---

func TestAzureNatGateway_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurenatgateway", "pulumi")
}
func TestAzureNatGateway_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurenatgateway", "terraform")
}

// --- Azure Virtual Network Peering (composed: fixture network + path-declared second network -> one peering direction) ---

func TestAzureVirtualNetworkPeering_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualnetworkpeering", "pulumi")
}
func TestAzureVirtualNetworkPeering_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualnetworkpeering", "terraform")
}

// --- Azure AKS Cluster (composed: fixture RG -> managed-networking cluster with a single-node default pool) ---

func TestAzureAksCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureakscluster", "pulumi")
}
func TestAzureAksCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureakscluster", "terraform")
}

// --- Azure AKS Node Pool (composed: fixture RG -> cluster prerequisite -> zero-node user pool) ---

func TestAzureAksNodePool_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureaksnodepool", "pulumi")
}
func TestAzureAksNodePool_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureaksnodepool", "terraform")
}

// --- Azure Managed Disk (fixture RG -> empty zonal Standard SSD data disk) ---

func TestAzureManagedDisk_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanageddisk", "pulumi")
}
func TestAzureManagedDisk_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanageddisk", "terraform")
}

// --- Azure Network Interface (composed: fixture RG -> VNet -> subnet -> NIC with an extra-fixture NSG attached) ---

func TestAzureNetworkInterface_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurenetworkinterface", "pulumi")
}
func TestAzureNetworkInterface_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurenetworkinterface", "terraform")
}

// --- Azure Virtual Machine (composed: fixture RG -> VNet -> subnet -> NIC -> Linux VM with an extra-fixture disk at LUN 0) ---

func TestAzureVirtualMachine_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualmachine", "pulumi")
}
func TestAzureVirtualMachine_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualmachine", "terraform")
}

// --- Azure Container Registry (fixture RG -> Standard registry with the admin account on) ---

func TestAzureContainerRegistry_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerregistry", "pulumi")
}
func TestAzureContainerRegistry_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerregistry", "terraform")
}

// --- Azure Load Balancer (composed: fixture RG + extra-fixture public IP -> public LB with pool, probe, rule, NAT rule, and outbound rule) ---

func TestAzureLoadBalancer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureloadbalancer", "pulumi")
}
func TestAzureLoadBalancer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureloadbalancer", "terraform")
}

// --- Azure Virtual Machine Scale Set (composed: fixture RG -> VNet -> subnet -> one-instance fleet, BOTH orchestration modes) ---

func TestAzureVirtualMachineScaleSet_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualmachinescaleset", "pulumi")
}
func TestAzureVirtualMachineScaleSet_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurevirtualmachinescaleset", "terraform")
}

// --- Azure Key Vault (fixture RG -> Standard RBAC vault with network rules; purge-on-destroy frees the global name) ---

func TestAzureKeyVault_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurekeyvault", "pulumi")
}
func TestAzureKeyVault_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurekeyvault", "terraform")
}

// --- Azure Key Vault Key (composed: fixture RG -> vault -> RSA CMK with a rotation policy; data-plane create) ---

func TestAzureKeyVaultKey_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurekeyvaultkey", "pulumi")
}
func TestAzureKeyVaultKey_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurekeyvaultkey", "terraform")
}

// --- Azure Key Vault Certificate (composed: fixture RG -> vault -> self-signed auto-renewing certificate; data-plane create + verify) ---

func TestAzureKeyVaultCertificate_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurekeyvaultcertificate", "pulumi")
}
func TestAzureKeyVaultCertificate_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurekeyvaultcertificate", "terraform")
}

// --- Azure Application Gateway (composed: fixture RG -> VNet -> dedicated subnet + extra-fixture public IP -> Standard_v2 gateway; the waf-attach scenario adds the fixture WAF policy) ---

func TestAzureApplicationGateway_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationgateway", "pulumi")
}
func TestAzureApplicationGateway_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationgateway", "terraform")
}

// --- Azure Web Application Firewall Policy (fixture RG -> OWASP 3.2 policy with a rate-limit rule, override, and log scrubbing) ---

func TestAzureWebApplicationFirewallPolicy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurewebapplicationfirewallpolicy", "pulumi")
}
func TestAzureWebApplicationFirewallPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurewebapplicationfirewallpolicy", "terraform")
}

// --- Azure PostgreSQL Flexible Server (fixture RG -> burstable public server with database, firewall rule, and server parameter) ---

func TestAzurePostgresqlFlexibleServer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurepostgresqlflexibleserver", "pulumi")
}
func TestAzurePostgresqlFlexibleServer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurepostgresqlflexibleserver", "terraform")
}

// runAllScenariosForComponent discovers and runs all E2E scenarios for an Azure component.
func runAllScenariosForComponent(t *testing.T, component, engine string) {
	t.Helper()

	var moduleDir string
	switch engine {
	case "pulumi":
		moduleDir = filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "azure", component, "v1", "iac", "pulumi")
	case "terraform":
		moduleDir = filepath.Join(repoRoot, "apis", "dev", "planton", "provider", "azure", component, "v1", "iac", "tf")
	default:
		t.Fatalf("unsupported engine: %s", engine)
	}

	if !fileExists(moduleDir) {
		t.Skipf("component %s %s module not found at %s", component, engine, moduleDir)
	}

	scenarios, err := discovery.DiscoverTestScenarios(repoRoot, "azure", component)
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
		Provider:     "azure",
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
