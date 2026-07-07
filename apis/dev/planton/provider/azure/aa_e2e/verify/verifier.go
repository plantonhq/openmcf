// Package verify checks that Azure resources created by an E2E scenario exist
// after DEPLOY and are gone after DESTROY. Each component family has its own
// verifier because Azure verification is service-specific (CheckExistence for a
// resource group, a GET for a VNet, ...). All verifiers run against the same
// ambient credential chain the deploy used, so a verification failure reflects
// real cloud state.
package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/pkg/errors"
)

// Verifier checks a single component's Azure resource for existence/absence.
// Azure resource verification is subscription-scoped, so verifiers take the
// subscription id and the ambient token credential -- not a region, unlike AWS.
type Verifier interface {
	// IDOutputKey is the stack-output key carrying the identifier used to verify
	// the resource (e.g. "resource_group_name").
	IDOutputKey() string
	// VerifyExists returns an error unless the resource exists.
	VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error
	// VerifyAbsent returns an error unless the resource is gone.
	VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error
}

// verifiers maps a component name to its verifier. New Azure components register
// here as they are forged.
var verifiers = map[string]Verifier{
	"azureakscluster":                       &aksClusterVerifier{},
	"azureaksnodepool":                      &aksNodePoolVerifier{},
	"azureapplicationgateway":               &applicationGatewayVerifier{},
	"azurecontainerregistry":                &containerRegistryVerifier{},
	"azurecosmosdbaccount":                  &cosmosdbAccountVerifier{},
	"azurecosmosdbmongocollection":          &cosmosdbMongoCollectionVerifier{},
	"azurecosmosdbmongodatabase":            &cosmosdbMongoDatabaseVerifier{},
	"azurecosmosdbsqlcontainer":             &cosmosdbSqlContainerVerifier{},
	"azurecosmosdbsqldatabase":              &cosmosdbSqlDatabaseVerifier{},
	"azurefederatedidentitycredential":      &federatedIdentityCredentialVerifier{},
	"azurekeyvault":                         &keyVaultVerifier{},
	"azurekeyvaultcertificate":              &keyVaultCertificateVerifier{},
	"azurekeyvaultkey":                      &keyVaultKeyVerifier{},
	"azureloadbalancer":                     &loadBalancerVerifier{},
	"azuremanageddisk":                      &managedDiskVerifier{},
	"azuremssqldatabase":                    &mssqlDatabaseVerifier{},
	"azuremssqlelasticpool":                 &mssqlElasticPoolVerifier{},
	"azuremssqlserver":                      &mssqlServerVerifier{},
	"azuremysqlflexibleserver":              &mysqlFlexibleServerVerifier{},
	"azurenatgateway":                       &natGatewayVerifier{},
	"azurenetworkinterface":                 &networkInterfaceVerifier{},
	"azurenetworksecuritygroup":             &networkSecurityGroupVerifier{},
	"azurepostgresqlflexibleserver":         &postgresqlFlexibleServerVerifier{},
	"azureprivatednszone":                   &privateDnsZoneVerifier{},
	"azureprivatednszonevirtualnetworklink": &privateDnsZoneVirtualNetworkLinkVerifier{},
	"azurepublicip":                         &publicIpVerifier{},
	"azurepublicipprefix":                   &publicIpPrefixVerifier{},
	"azurerediscache":                       &redisCacheVerifier{},
	"azurerediscacheaccesspolicy":           &redisCacheAccessPolicyVerifier{},
	"azurerediscacheaccesspolicyassignment": &redisCacheAccessPolicyAssignmentVerifier{},
	"azureredislinkedserver":                &redisLinkedServerVerifier{},
	"azureresourcegroup":                    &resourceGroupVerifier{},
	"azureroleassignment":                   &roleAssignmentVerifier{},
	"azureroledefinition":                   &roleDefinitionVerifier{},
	"azureroutetable":                       &routeTableVerifier{},
	"azurestorageaccount":                   &storageAccountVerifier{},
	"azurestoragecontainer":                 &storageContainerVerifier{},
	"azurestorageencryptionscope":           &storageEncryptionScopeVerifier{},
	"azurestoragequeue":                     &storageQueueVerifier{},
	"azurestorageshare":                     &storageShareVerifier{},
	"azurestoragetable":                     &storageTableVerifier{},
	"azuresubnet":                           &subnetVerifier{},
	"azureuserassignedidentity":             &userAssignedIdentityVerifier{},
	"azurevirtualmachine":                   &virtualMachineVerifier{},
	"azurevirtualmachinescaleset":           &virtualMachineScaleSetVerifier{},
	"azurevirtualnetwork":                   &virtualNetworkVerifier{},
	"azurevirtualnetworkpeering":            &virtualNetworkPeeringVerifier{},
	"azurewebapplicationfirewallpolicy":     &webApplicationFirewallPolicyVerifier{},
}

// GetVerifier returns the verifier for a component, or an error if none is registered.
func GetVerifier(component string) (Verifier, error) {
	v, ok := verifiers[component]
	if !ok {
		return nil, errors.Errorf("no Azure verifier registered for component %q", component)
	}
	return v, nil
}
