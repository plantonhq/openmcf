package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// resourceGroupVerifier verifies an AzureResourceGroup via
// ResourceGroupsClient.CheckExistence -- a purpose-built boolean HEAD that reads
// no resource body and needs only Reader on the resource group. The typed
// Success flag is the existence signal, so absence never masquerades as a real
// error (auth/network failures surface instead).
type resourceGroupVerifier struct{}

// IDOutputKey is the resource-group NAME (not the full ARM id): CheckExistence is
// keyed by name within the subscription.
func (*resourceGroupVerifier) IDOutputKey() string { return "resource_group_name" }

func (*resourceGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := resourceGroupExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureresourcegroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureresourcegroup %q not found after deploy", id)
	}
	return nil
}

func (*resourceGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := resourceGroupExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureresourcegroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureresourcegroup %q still exists after destroy", id)
	}
	return nil
}

func resourceGroupExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, name string) (bool, error) {
	client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	resp, err := client.CheckExistence(ctx, name, nil)
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
