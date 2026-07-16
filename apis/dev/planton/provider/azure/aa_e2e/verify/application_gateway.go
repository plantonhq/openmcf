package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// applicationGatewayAPIVersion is the stable Microsoft.Network API version
// the generic existence probe is pinned to.
const applicationGatewayAPIVersion = "2024-05-01"

// applicationGatewayVerifier verifies an AzureApplicationGateway via the
// generic ARM resources GetByID, keyed on the gateway's full ARM ID.
type applicationGatewayVerifier struct{}

// IDOutputKey is the gateway's full ARM ID.
func (*applicationGatewayVerifier) IDOutputKey() string {
	return "application_gateway_id"
}

func (*applicationGatewayVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, applicationGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationgateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureapplicationgateway %q not found after deploy", id)
	}
	return nil
}

func (*applicationGatewayVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, applicationGatewayAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationgateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureapplicationgateway %q still exists after destroy", id)
	}
	return nil
}
