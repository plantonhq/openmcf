package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualWanAPIVersion is the stable Microsoft.Network API version the
// existence probes for the Virtual WAN family are pinned to (the WAN
// and the hub resources created under it share the namespace and GA
// API line).
const virtualWanAPIVersion = "2024-05-01"

// virtualWanVerifier verifies an AzureVirtualWan via the generic ARM
// resources GetByID (see armResourceExists), keyed on the WAN's full
// ARM ID. The WAN is a pure policy object (its hubs are separate
// resources), so ARM existence IS the complete deployment claim.
type virtualWanVerifier struct{}

// IDOutputKey is the WAN's full ARM ID.
func (*virtualWanVerifier) IDOutputKey() string {
	return "virtual_wan_id"
}

func (*virtualWanVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualwan verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualwan %q not found after deploy", id)
	}
	return nil
}

func (*virtualWanVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualwan verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualwan %q still exists after destroy", id)
	}
	return nil
}
