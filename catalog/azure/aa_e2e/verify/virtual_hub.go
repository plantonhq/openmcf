package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualHubVerifier verifies an AzureVirtualHub via the generic ARM
// resources GetByID (see armResourceExists), keyed on the hub's full
// ARM ID. Existence is the honest bar for the hub object: the router's
// routing state settles asynchronously (15-30 minutes to Provisioned),
// and the composed children (route tables, route maps, BGP
// connections, routing intent) are ARM children under the same path --
// the provider's own read cycle already gates the deploy on them. The
// Virtual WAN family shares the pinned virtualWanAPIVersion.
type virtualHubVerifier struct{}

// IDOutputKey is the hub's full ARM ID.
func (*virtualHubVerifier) IDOutputKey() string {
	return "virtual_hub_id"
}

func (*virtualHubVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualhub verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualhub %q not found after deploy", id)
	}
	return nil
}

func (*virtualHubVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualhub verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualhub %q still exists after destroy", id)
	}
	return nil
}
