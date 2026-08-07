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
	azuree2e "github.com/plantonhq/planton/catalog/azure/aa_e2e"
	"github.com/plantonhq/planton/e2e/framework/discovery"
	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/e2e/framework/runner"
	profilepkg "github.com/plantonhq/planton/pkg/e2e/profile"
	componentv1 "github.com/plantonhq/planton/qa/componente2eprofile/v1"
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

// --- Azure Application Security Group (empty micro-segmentation anchor in the fixture RG) ---

func TestAzureApplicationSecurityGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationsecuritygroup", "pulumi")
}
func TestAzureApplicationSecurityGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationsecuritygroup", "terraform")
}

// --- Azure Private Endpoint (composed: fixture subnet + scenario-local storage target + blob privatelink zone) ---

func TestAzurePrivateEndpoint_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureprivateendpoint", "pulumi")
}
func TestAzurePrivateEndpoint_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureprivateendpoint", "terraform")
}

// --- Azure Disk Encryption Set (profile-deferred: purge-protected vault cannot teardown to zero orphans) ---

func TestAzureDiskEncryptionSet_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurediskencryptionset", "pulumi")
}
func TestAzureDiskEncryptionSet_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurediskencryptionset", "terraform")
}

// --- Azure MSSQL Failover Group (composed: shared primary + scenario-local partner + database) ---

func TestAzureMssqlFailoverGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqlfailovergroup", "pulumi")
}
func TestAzureMssqlFailoverGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqlfailovergroup", "terraform")
}

// --- Azure Monitor Activity Log Alert (composed: fixture action group -> administrative alert) ---

func TestAzureMonitorActivityLogAlert_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitoractivitylogalert", "pulumi")
}
func TestAzureMonitorActivityLogAlert_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitoractivitylogalert", "terraform")
}

// --- Azure Application Insights Standard Web Test (composed: fixture App Insights -> web test) ---

func TestAzureApplicationInsightsStandardWebTest_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationinsightsstandardwebtest", "pulumi")
}
func TestAzureApplicationInsightsStandardWebTest_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationinsightsstandardwebtest", "terraform")
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

// --- Azure MySQL Flexible Server (fixture RG -> burstable public server with database, firewall rule, and server parameter) ---

func TestAzureMysqlFlexibleServer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremysqlflexibleserver", "pulumi")
}
func TestAzureMysqlFlexibleServer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremysqlflexibleserver", "terraform")
}

// --- Azure SQL logical server (fixture RG -> SQL-auth server with firewall rule and Defender alert policy) ---

func TestAzureMssqlServer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqlserver", "pulumi")
}
func TestAzureMssqlServer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqlserver", "terraform")
}

// --- Azure SQL Database (fixture RG -> fixture server -> Basic database; pool-attach joins the fixture elastic pool) ---

func TestAzureMssqlDatabase_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqldatabase", "pulumi")
}
func TestAzureMssqlDatabase_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqldatabase", "terraform")
}

// --- Azure SQL elastic pool (fixture RG -> fixture server -> BasicPool) ---

func TestAzureMssqlElasticPool_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqlelasticpool", "pulumi")
}
func TestAzureMssqlElasticPool_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremssqlelasticpool", "terraform")
}

// --- Azure storage account (fixture RG -> StorageV2 with firewall + lifecycle policy) ---

func TestAzureStorageAccount_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageaccount", "pulumi")
}
func TestAzureStorageAccount_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageaccount", "terraform")
}

// --- Azure storage container (fixture RG -> scenario-local account -> private container;
// plus the composed encryption-scope chain) ---

func TestAzureStorageContainer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragecontainer", "pulumi")
}
func TestAzureStorageContainer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragecontainer", "terraform")
}

// --- Azure storage share (fixture RG -> scenario-local account -> SMB share) ---

func TestAzureStorageShare_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageshare", "pulumi")
}
func TestAzureStorageShare_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageshare", "terraform")
}

// --- Azure storage queue (fixture RG -> scenario-local account -> work queue) ---

func TestAzureStorageQueue_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragequeue", "pulumi")
}
func TestAzureStorageQueue_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragequeue", "terraform")
}

// --- Azure storage table (fixture RG -> scenario-local account -> entities table with ACL) ---

func TestAzureStorageTable_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragetable", "pulumi")
}
func TestAzureStorageTable_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragetable", "terraform")
}

// --- Azure storage encryption scope (fixture RG -> scenario-local account ->
// platform-managed scope; destroy verified state-aware: ARM soft-disables scopes) ---

func TestAzureStorageEncryptionScope_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageencryptionscope", "pulumi")
}
func TestAzureStorageEncryptionScope_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageencryptionscope", "terraform")
}

// --- Azure Data Lake Gen2 filesystem (fixture RG -> scenario-local HNS
// account -> filesystem with a root POSIX ACL; verified via the ARM
// blob-container proxy) ---

func TestAzureStorageDataLakeGen2Filesystem_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragedatalakegen2filesystem", "pulumi")
}
func TestAzureStorageDataLakeGen2Filesystem_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragedatalakegen2filesystem", "terraform")
}

// --- Azure storage local user (fixture RG -> scenario-local SFTP account +
// container -> user with both auth methods and a container-scoped grant) ---

func TestAzureStorageLocalUser_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragelocaluser", "pulumi")
}
func TestAzureStorageLocalUser_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestoragelocaluser", "terraform")
}

// --- Azure storage object replication (fixture RG -> two scenario-local
// versioned accounts + a container on each -> the two-sided policy) ---

func TestAzureStorageObjectReplication_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageobjectreplication", "pulumi")
}
func TestAzureStorageObjectReplication_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurestorageobjectreplication", "terraform")
}

// --- Azure Redis cache (fixture RG -> Standard C0 cache with Entra auth,
// firewall rule, patch window; Redis provisioning runs 15-40 min per cache) ---

func TestAzureRedisCache_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurerediscache", "pulumi")
}
func TestAzureRedisCache_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurerediscache", "terraform")
}

// --- Azure Redis cache access policy (fixture RG -> scenario-local Entra
// cache -> prefix-scoped custom policy) ---

func TestAzureRedisCacheAccessPolicy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurerediscacheaccesspolicy", "pulumi")
}
func TestAzureRedisCacheAccessPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurerediscacheaccesspolicy", "terraform")
}

// --- Azure Redis cache access policy assignment (fixture RG + identity ->
// scenario-local Entra cache -> custom policy -> the grant; proves both FK
// seams: policy-name and identity-principal references) ---

func TestAzureRedisCacheAccessPolicyAssignment_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurerediscacheaccesspolicyassignment", "pulumi")
}
func TestAzureRedisCacheAccessPolicyAssignment_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurerediscacheaccesspolicyassignment", "terraform")
}

// --- Azure Redis linked server (fixture RG -> two scenario-local Premium P1
// caches in different regions -> the geo-replication link; the most expensive
// scenario in the suite) ---

func TestAzureRedisLinkedServer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureredislinkedserver", "pulumi")
}
func TestAzureRedisLinkedServer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureredislinkedserver", "terraform")
}

// --- Azure Managed Redis (fixture RG -> Balanced_B0 instance with access
// keys enabled; Managed Redis clusters provision and delete in tens of
// minutes) ---

func TestAzureManagedRedis_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanagedredis", "pulumi")
}
func TestAzureManagedRedis_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanagedredis", "terraform")
}

// --- Azure Managed Redis geo-replication (fixture RG -> two scenario-local
// Balanced_B3 clusters in different regions declaring the same group name ->
// the group link; the most expensive fixture pair in the suite) ---

func TestAzureManagedRedisGeoReplication_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanagedredisgeoreplication", "pulumi")
}
func TestAzureManagedRedisGeoReplication_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanagedredisgeoreplication", "terraform")
}

// --- Azure Managed Redis access policy assignment (fixture RG + identity ->
// scenario-local KEYLESS cluster -> the data-plane grant; proves the
// cluster-reference and identity-principal FK seams) ---

func TestAzureManagedRedisAccessPolicyAssignment_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanagedredisaccesspolicyassignment", "pulumi")
}
func TestAzureManagedRedisAccessPolicyAssignment_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremanagedredisaccesspolicyassignment", "terraform")
}

// --- Azure Cosmos DB account (fixture RG -> single-region SQL API account;
// Cosmos accounts provision in ~5-10 min each) ---

func TestAzureCosmosdbAccount_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbaccount", "pulumi")
}
func TestAzureCosmosdbAccount_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbaccount", "terraform")
}

// --- Azure Cosmos DB SQL database (fixture RG -> scenario-local account ->
// database with shared throughput) ---

func TestAzureCosmosdbSqlDatabase_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqldatabase", "pulumi")
}
func TestAzureCosmosdbSqlDatabase_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqldatabase", "terraform")
}

// --- Azure Cosmos DB SQL container (composed chain: scenario-local account
// -> database -> container with partition key) ---

func TestAzureCosmosdbSqlContainer_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqlcontainer", "pulumi")
}
func TestAzureCosmosdbSqlContainer_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqlcontainer", "terraform")
}

// --- Azure Cosmos DB Mongo database (fixture RG -> scenario-local Mongo
// account -> database with shared throughput) ---

func TestAzureCosmosdbMongoDatabase_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbmongodatabase", "pulumi")
}
func TestAzureCosmosdbMongoDatabase_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbmongodatabase", "terraform")
}

// --- Azure Cosmos DB Mongo collection (composed chain: scenario-local Mongo
// account -> database -> collection with shard key) ---

func TestAzureCosmosdbMongoCollection_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbmongocollection", "pulumi")
}
func TestAzureCosmosdbMongoCollection_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbmongocollection", "terraform")
}

// --- Azure Cosmos DB SQL role definition (fixture RG -> scenario-local
// SQL account -> custom data-plane role) ---

func TestAzureCosmosdbSqlRoleDefinition_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqlroledefinition", "pulumi")
}
func TestAzureCosmosdbSqlRoleDefinition_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqlroledefinition", "terraform")
}

// --- Azure Cosmos DB SQL role assignment (composed chain: scenario-local
// SQL account -> custom role definition -> grant to the fixture identity) ---

func TestAzureCosmosdbSqlRoleAssignment_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqlroleassignment", "pulumi")
}
func TestAzureCosmosdbSqlRoleAssignment_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecosmosdbsqlroleassignment", "terraform")
}

// --- Azure Front Door profile (fixture RG -> Standard profile with log
// scrubbing and tags; Front Door is global -- no region on the resource) ---

func TestAzureFrontDoorProfile_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorprofile", "pulumi")
}
func TestAzureFrontDoorProfile_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorprofile", "terraform")
}

// --- Azure Front Door endpoint (fixture RG -> fixture profile -> endpoint;
// verifies the generated *.azurefd.net hostname surfaces as an output) ---

func TestAzureFrontDoorEndpoint_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorendpoint", "pulumi")
}
func TestAzureFrontDoorEndpoint_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorendpoint", "terraform")
}

// --- Azure Front Door origin group (fixture RG -> fixture profile -> group
// with health probe and load-balancing dials) ---

func TestAzureFrontDoorOriginGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoororigingroup", "pulumi")
}
func TestAzureFrontDoorOriginGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoororigingroup", "terraform")
}

// --- Azure Front Door origin (fixture RG -> fixture profile -> fixture
// group -> origin pointing at a public backend hostname) ---

func TestAzureFrontDoorOrigin_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoororigin", "pulumi")
}
func TestAzureFrontDoorOrigin_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoororigin", "terraform")
}

// --- Azure Front Door route (the composed chain: fixture RG -> profile ->
// endpoint + origin group -> origin -> route with caching; proves the whole
// traffic-serving graph and the origin_ids ordering seam) ---

func TestAzureFrontDoorRoute_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorroute", "pulumi")
}
func TestAzureFrontDoorRoute_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorroute", "terraform")
}

// --- Azure Front Door rule set (fixture RG -> fixture profile -> a
// three-rule delivery policy: redirect, security headers, caching
// override; the route's rule-set-attach scenario proves the attach seam) ---

func TestAzureFrontDoorRuleSet_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorruleset", "pulumi")
}
func TestAzureFrontDoorRuleSet_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorruleset", "terraform")
}

// --- Azure Front Door custom domain (fixture RG -> fixture profile -> a
// managed-certificate domain in the pending-validation state; proves the
// validation_token challenge surfaces as an output) ---

func TestAzureFrontDoorCustomDomain_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorcustomdomain", "pulumi")
}
func TestAzureFrontDoorCustomDomain_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorcustomdomain", "terraform")
}

// --- Azure Front Door secret (fixture RG -> fixture profile + the Key
// Vault certificate fixture chain -> a BYO-certificate secret wrapping
// the certificate's versionless id) ---

func TestAzureFrontDoorSecret_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorsecret", "pulumi")
}
func TestAzureFrontDoorSecret_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorsecret", "terraform")
}

// --- Azure Front Door firewall policy (fixture RG -> the WAF policy;
// resource-group-scoped, NO profile fixture -- two scenarios cover the
// STANDARD custom-rules smoke and the PREMIUM managed-rules depth) ---

func TestAzureFrontDoorFirewallPolicy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorfirewallpolicy", "pulumi")
}
func TestAzureFrontDoorFirewallPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorfirewallpolicy", "terraform")
}

// --- Azure Front Door security policy (the enforcement seam: fixture RG
// -> profile -> endpoint + STANDARD WAF policy -> the association
// protecting the endpoint's default domain) ---

func TestAzureFrontDoorSecurityPolicy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorsecuritypolicy", "pulumi")
}
func TestAzureFrontDoorSecurityPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefrontdoorsecuritypolicy", "terraform")
}

// --- Azure Service Plan (fixture RG -> Basic B1 Linux plan; the compute tier the app kinds run on) ---

func TestAzureServicePlan_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureserviceplan", "pulumi")
}
func TestAzureServicePlan_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureserviceplan", "terraform")
}

// --- Azure Linux Web App (composed: fixture RG -> fixture plan -> Python app with always-on + health probe) ---

func TestAzureLinuxWebApp_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurelinuxwebapp", "pulumi")
}
func TestAzureLinuxWebApp_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurelinuxwebapp", "terraform")
}

// --- Azure Function App (composed: fixture RG -> fixture plan + scenario-local storage account -> Python function app bound by name + access key) ---

func TestAzureFunctionApp_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefunctionapp", "pulumi")
}
func TestAzureFunctionApp_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefunctionapp", "terraform")
}

// --- Azure Container App Environment (fixture RG -> consumption-only environment; the boundary the family lives in) ---

func TestAzureContainerAppEnvironment_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironment", "pulumi")
}
func TestAzureContainerAppEnvironment_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironment", "terraform")
}

// --- Azure Container App (composed: fixture RG -> fixture environment -> quickstart app with external ingress + HTTP scaling) ---

func TestAzureContainerApp_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerapp", "pulumi")
}
func TestAzureContainerApp_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerapp", "terraform")
}

// --- Azure Container App Job (composed: fixture RG -> fixture environment -> schedule-triggered job) ---

func TestAzureContainerAppJob_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappjob", "pulumi")
}
func TestAzureContainerAppJob_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappjob", "terraform")
}

// --- Azure Container App Environment Storage (composed: fixture RG -> fixture environment + scenario-local account -> share -> SMB registration) ---

func TestAzureContainerAppEnvironmentStorage_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentstorage", "pulumi")
}
func TestAzureContainerAppEnvironmentStorage_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentstorage", "terraform")
}

// --- Azure Container App Environment Dapr Component (composed: fixture RG -> fixture environment -> backendless cron-binding component) ---

func TestAzureContainerAppEnvironmentDaprComponent_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentdaprcomponent", "pulumi")
}
func TestAzureContainerAppEnvironmentDaprComponent_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentdaprcomponent", "terraform")
}

// --- Azure Log Analytics Workspace (fixture RG -> pay-as-you-go workspace with retention + quota dials) ---

func TestAzureLogAnalyticsWorkspace_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureloganalyticsworkspace", "pulumi")
}
func TestAzureLogAnalyticsWorkspace_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureloganalyticsworkspace", "terraform")
}

// --- Azure Application Insights (composed: fixture RG -> fixture workspace -> workspace-based component) ---

func TestAzureApplicationInsights_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationinsights", "pulumi")
}
func TestAzureApplicationInsights_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureapplicationinsights", "terraform")
}

// --- Azure Monitor Diagnostic Setting (composed: fixture RG -> fixture workspace routing its own audit logs into itself) ---

func TestAzureMonitorDiagnosticSetting_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitordiagnosticsetting", "pulumi")
}
func TestAzureMonitorDiagnosticSetting_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitordiagnosticsetting", "terraform")
}

// --- Azure Monitor Action Group (fixture RG -> multi-receiver global notification hub) ---

func TestAzureMonitorActionGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitoractiongroup", "pulumi")
}
func TestAzureMonitorActionGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitoractiongroup", "terraform")
}

// --- Azure Monitor Metric Alert (composed: fixture RG -> fixture action group + scenario-local storage account -> static-threshold rule) ---

func TestAzureMonitorMetricAlert_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitormetricalert", "pulumi")
}
func TestAzureMonitorMetricAlert_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitormetricalert", "terraform")
}

// --- Azure Monitor Scheduled Query Alert (composed: fixture RG -> fixture workspace + fixture action group -> row-count KQL rule) ---

func TestAzureMonitorScheduledQueryAlert_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitorscheduledqueryalert", "pulumi")
}
func TestAzureMonitorScheduledQueryAlert_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azuremonitorscheduledqueryalert", "terraform")
}

// --- Azure Service Bus family (namespace container + entities + SAS credentials + geo-DR) ---

func TestAzureServiceBusNamespace_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusnamespace", "pulumi")
}
func TestAzureServiceBusNamespace_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusnamespace", "terraform")
}

func TestAzureServiceBusQueue_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusqueue", "pulumi")
}
func TestAzureServiceBusQueue_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusqueue", "terraform")
}

func TestAzureServiceBusTopic_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebustopic", "pulumi")
}
func TestAzureServiceBusTopic_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebustopic", "terraform")
}

func TestAzureServiceBusSubscription_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebussubscription", "pulumi")
}
func TestAzureServiceBusSubscription_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebussubscription", "terraform")
}

func TestAzureServiceBusAuthorizationRule_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusauthorizationrule", "pulumi")
}
func TestAzureServiceBusAuthorizationRule_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusauthorizationrule", "terraform")
}

func TestAzureServiceBusDisasterRecoveryConfig_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusdisasterrecoveryconfig", "pulumi")
}
func TestAzureServiceBusDisasterRecoveryConfig_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureservicebusdisasterrecoveryconfig", "terraform")
}

// --- Azure Event Hubs family (namespace container + hubs + consumer groups + SAS credentials + schema registry + geo-DR + dedicated cluster + CMK) ---

func TestAzureEventHubNamespace_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubnamespace", "pulumi")
}
func TestAzureEventHubNamespace_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubnamespace", "terraform")
}

func TestAzureEventHub_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhub", "pulumi")
}
func TestAzureEventHub_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhub", "terraform")
}

func TestAzureEventHubConsumerGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubconsumergroup", "pulumi")
}
func TestAzureEventHubConsumerGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubconsumergroup", "terraform")
}

func TestAzureEventHubAuthorizationRule_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubauthorizationrule", "pulumi")
}
func TestAzureEventHubAuthorizationRule_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubauthorizationrule", "terraform")
}

func TestAzureEventHubSchemaGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubschemagroup", "pulumi")
}
func TestAzureEventHubSchemaGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubschemagroup", "terraform")
}

func TestAzureEventHubDisasterRecoveryConfig_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubdisasterrecoveryconfig", "pulumi")
}
func TestAzureEventHubDisasterRecoveryConfig_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubdisasterrecoveryconfig", "terraform")
}

func TestAzureEventHubCluster_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubcluster", "pulumi")
}
func TestAzureEventHubCluster_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubcluster", "terraform")
}

func TestAzureEventHubNamespaceCustomerManagedKey_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubnamespacecustomermanagedkey", "pulumi")
}
func TestAzureEventHubNamespaceCustomerManagedKey_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureeventhubnamespacecustomermanagedkey", "terraform")
}

// --- Azure public DNS (fixture RG -> zone with SOA customization; fixture zone -> typed record sets incl. the alias-A seam) ---

func TestAzureDnsZone_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurednszone", "pulumi")
}
func TestAzureDnsZone_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurednszone", "terraform")
}

func TestAzureDnsRecord_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurednsrecord", "pulumi")
}
func TestAzureDnsRecord_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurednsrecord", "terraform")
}

// --- Azure Container Apps TLS/domain family (fixture environment -> BYO certificate; managed certificate + custom domain are profile-deferred: both block on public-DNS domain validation) ---

func TestAzureContainerAppEnvironmentCertificate_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentcertificate", "pulumi")
}
func TestAzureContainerAppEnvironmentCertificate_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentcertificate", "terraform")
}

func TestAzureContainerAppEnvironmentManagedCertificate_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentmanagedcertificate", "pulumi")
}
func TestAzureContainerAppEnvironmentManagedCertificate_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappenvironmentmanagedcertificate", "terraform")
}

func TestAzureContainerAppCustomDomain_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappcustomdomain", "pulumi")
}
func TestAzureContainerAppCustomDomain_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurecontainerappcustomdomain", "terraform")
}

// --- Azure Firewall family (IP Group anchor; policy container; rule collection group composed on the policy; the firewall itself on a dedicated AzureFirewallSubnet -- the SLOW lane, ~10-20 min each way) ---

func TestAzureIpGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azureipgroup", "pulumi")
}
func TestAzureIpGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azureipgroup", "terraform")
}

func TestAzureFirewallPolicy_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefirewallpolicy", "pulumi")
}
func TestAzureFirewallPolicy_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefirewallpolicy", "terraform")
}

func TestAzureFirewallPolicyRuleCollectionGroup_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefirewallpolicyrulecollectiongroup", "pulumi")
}
func TestAzureFirewallPolicyRuleCollectionGroup_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefirewallpolicyrulecollectiongroup", "terraform")
}

func TestAzureFirewall_Pulumi(t *testing.T) {
	runAllScenariosForComponent(t, "azurefirewall", "pulumi")
}
func TestAzureFirewall_Terraform(t *testing.T) {
	runAllScenariosForComponent(t, "azurefirewall", "terraform")
}

// runAllScenariosForComponent discovers and runs all E2E scenarios for an Azure component.
func runAllScenariosForComponent(t *testing.T, component, engine string) {
	t.Helper()

	if cp, err := profilepkg.LoadComponentProfile(repoRoot, "azure", component); err == nil && cp.Spec != nil {
		switch cp.Spec.Status {
		case componentv1.ComponentE2EProfileSpec_deferred,
			componentv1.ComponentE2EProfileSpec_skip,
			componentv1.ComponentE2EProfileSpec_stub:
			reason := cp.Spec.DeferredReason
			if reason == "" {
				reason = cp.Spec.Status.String()
			}
			t.Skipf("component %s E2E profile status is %s: %s", component, cp.Spec.Status, reason)
		}
	}

	moduleDir, err := discovery.ModuleDir(repoRoot, "azure", component, engine)
	if err != nil {
		t.Fatalf("failed to locate %s %s module: %v", component, engine, err)
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
