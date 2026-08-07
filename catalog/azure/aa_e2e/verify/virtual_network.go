package verify

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// virtualNetworkAPIVersion is the stable Microsoft.Network API version the
// generic existence probe below is pinned to.
const virtualNetworkAPIVersion = "2024-05-01"

// virtualNetworkVerifier verifies an AzureVirtualNetwork via the generic ARM
// resources GetByID, keyed on the network's full ARM ID. A GET (not
// CheckExistenceByID's HEAD) is deliberate: not every ARM resource provider
// implements HEAD, and GET-by-ID is the uniform probe the Azure verifiers
// standardize on. A typed 404 ResponseError is the absence signal; any other
// failure (auth, network) surfaces as a real error rather than masquerading
// as "absent".
type virtualNetworkVerifier struct{}

// IDOutputKey is the virtual network's full ARM ID.
func (*virtualNetworkVerifier) IDOutputKey() string {
	return "virtual_network_id"
}

func (*virtualNetworkVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetwork verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualnetwork %q not found after deploy", id)
	}
	return nil
}

func (*virtualNetworkVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualNetworkAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualnetwork verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualnetwork %q still exists after destroy", id)
	}
	return nil
}

// armResourceExists is the shared generic ARM existence probe: GET the
// resource by its full ARM ID at the given API version and treat a typed 404
// as absence. Verifiers for plain ARM-GETtable resources share this rather
// than each importing a service-specific SDK module.
func armResourceExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, resourceID, apiVersion string) (bool, error) {
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	if _, err := client.GetByID(ctx, resourceID, apiVersion, nil); err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
