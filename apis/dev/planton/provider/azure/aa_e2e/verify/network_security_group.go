package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// networkSecurityGroupAPIVersion is the stable Microsoft.Network API version
// the generic existence probe is pinned to.
const networkSecurityGroupAPIVersion = "2024-05-01"

// networkSecurityGroupVerifier verifies an AzureNetworkSecurityGroup via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// group's full ARM ID.
type networkSecurityGroupVerifier struct{}

// IDOutputKey is the NSG's full ARM ID.
func (*networkSecurityGroupVerifier) IDOutputKey() string {
	return "network_security_group_id"
}

func (*networkSecurityGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, networkSecurityGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenetworksecuritygroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurenetworksecuritygroup %q not found after deploy", id)
	}
	return nil
}

func (*networkSecurityGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, networkSecurityGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenetworksecuritygroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurenetworksecuritygroup %q still exists after destroy", id)
	}
	return nil
}
