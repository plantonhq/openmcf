package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// subnetAPIVersion is the stable Microsoft.Network API version the generic
// existence probe is pinned to.
const subnetAPIVersion = "2024-05-01"

// subnetVerifier verifies an AzureSubnet via the generic ARM resources
// GetByID (see armResourceExists), keyed on the subnet's full ARM ID (a
// child of its virtual network).
type subnetVerifier struct{}

// IDOutputKey is the subnet's full ARM ID.
func (*subnetVerifier) IDOutputKey() string {
	return "subnet_id"
}

func (*subnetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, subnetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuresubnet verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuresubnet %q not found after deploy", id)
	}
	return nil
}

func (*subnetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, subnetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuresubnet verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuresubnet %q still exists after destroy", id)
	}
	return nil
}
