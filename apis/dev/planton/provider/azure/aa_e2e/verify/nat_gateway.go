package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// natGatewayAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const natGatewayAPIVersion = "2024-05-01"

// natGatewayVerifier verifies an AzureNatGateway via the generic ARM
// resources GetByID (see armResourceExists), keyed on the gateway's full ARM
// ID.
type natGatewayVerifier struct{}

// IDOutputKey is the NAT gateway's full ARM ID.
func (*natGatewayVerifier) IDOutputKey() string {
	return "nat_gateway_id"
}

func (*natGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, natGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenatgateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurenatgateway %q not found after deploy", id)
	}
	return nil
}

func (*natGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, natGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenatgateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurenatgateway %q still exists after destroy", id)
	}
	return nil
}
