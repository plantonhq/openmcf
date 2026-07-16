package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// privateEndpointAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const privateEndpointAPIVersion = "2024-05-01"

// privateEndpointVerifier verifies an AzurePrivateEndpoint via the generic
// ARM resources GetByID (see armResourceExists), keyed on the endpoint's full
// ARM ID.
type privateEndpointVerifier struct{}

// IDOutputKey is the endpoint's full ARM ID.
func (*privateEndpointVerifier) IDOutputKey() string {
	return "private_endpoint_id"
}

func (*privateEndpointVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateEndpointAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivateendpoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivateendpoint %q not found after deploy", id)
	}
	return nil
}

func (*privateEndpointVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateEndpointAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivateendpoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivateendpoint %q still exists after destroy", id)
	}
	return nil
}
