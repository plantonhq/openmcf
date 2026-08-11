package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualHubConnectionVerifier verifies an AzureVirtualHubConnection
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the connection's full ARM ID (an ARM child of the hub:
// .../virtualHubs/{hub}/hubVirtualNetworkConnections/{name}). ARM
// existence is the complete deployment claim -- the routing
// configuration lives on the same object and the provider's read cycle
// gates the deploy on it. The Virtual WAN family shares the pinned
// virtualWanAPIVersion.
type virtualHubConnectionVerifier struct{}

// IDOutputKey is the connection's full ARM ID.
func (*virtualHubConnectionVerifier) IDOutputKey() string {
	return "virtual_hub_connection_id"
}

func (*virtualHubConnectionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualhubconnection verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualhubconnection %q not found after deploy", id)
	}
	return nil
}

func (*virtualHubConnectionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualhubconnection verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualhubconnection %q still exists after destroy", id)
	}
	return nil
}
