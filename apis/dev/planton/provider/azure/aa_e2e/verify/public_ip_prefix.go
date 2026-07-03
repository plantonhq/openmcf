package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// publicIpPrefixAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const publicIpPrefixAPIVersion = "2024-05-01"

// publicIpPrefixVerifier verifies an AzurePublicIpPrefix via the generic ARM
// resources GetByID (see armResourceExists), keyed on the prefix's full ARM
// ID.
type publicIpPrefixVerifier struct{}

// IDOutputKey is the prefix's full ARM ID.
func (*publicIpPrefixVerifier) IDOutputKey() string {
	return "public_ip_prefix_id"
}

func (*publicIpPrefixVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicIpPrefixAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepublicipprefix verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurepublicipprefix %q not found after deploy", id)
	}
	return nil
}

func (*publicIpPrefixVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicIpPrefixAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurepublicipprefix verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurepublicipprefix %q still exists after destroy", id)
	}
	return nil
}
