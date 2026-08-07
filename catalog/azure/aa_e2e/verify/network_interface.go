package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// networkInterfaceAPIVersion is the stable Microsoft.Network API version
// the generic existence probe is pinned to.
const networkInterfaceAPIVersion = "2024-05-01"

// networkInterfaceVerifier verifies an AzureNetworkInterface via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// NIC's full ARM ID.
type networkInterfaceVerifier struct{}

// IDOutputKey is the NIC's full ARM ID.
func (*networkInterfaceVerifier) IDOutputKey() string {
	return "network_interface_id"
}

func (*networkInterfaceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, networkInterfaceAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenetworkinterface verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurenetworkinterface %q not found after deploy", id)
	}
	return nil
}

func (*networkInterfaceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, networkInterfaceAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurenetworkinterface verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurenetworkinterface %q still exists after destroy", id)
	}
	return nil
}
