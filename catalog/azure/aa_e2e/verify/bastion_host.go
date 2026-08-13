package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// bastionHostAPIVersion is the stable Microsoft.Network API version the
// existence probe is pinned to -- the line the pinned provider vendors
// for bastionHosts.
const bastionHostAPIVersion = "2025-01-01"

// bastionHostVerifier verifies an AzureBastionHost via the generic ARM
// resources GetByID (see armResourceExists), keyed on the host's full
// ARM ID.
type bastionHostVerifier struct{}

// IDOutputKey is the host's full ARM ID.
func (*bastionHostVerifier) IDOutputKey() string {
	return "bastion_host_id"
}

func (*bastionHostVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, bastionHostAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebastionhost verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurebastionhost %q not found after deploy", id)
	}
	return nil
}

func (*bastionHostVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, bastionHostAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurebastionhost verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurebastionhost %q still exists after destroy", id)
	}
	return nil
}
