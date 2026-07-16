package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// publicIpAPIVersion is the stable Microsoft.Network API version the generic
// existence probe is pinned to.
const publicIpAPIVersion = "2024-05-01"

// publicIpVerifier verifies an AzurePublicIp via the generic ARM resources
// GetByID (see armResourceExists), keyed on the address's full ARM ID.
type publicIpVerifier struct{}

// IDOutputKey is the public IP's full ARM ID.
func (*publicIpVerifier) IDOutputKey() string {
	return "public_ip_id"
}

func (*publicIpVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicIpAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepublicip verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurepublicip %q not found after deploy", id)
	}
	return nil
}

func (*publicIpVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicIpAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepublicip verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurepublicip %q still exists after destroy", id)
	}
	return nil
}
