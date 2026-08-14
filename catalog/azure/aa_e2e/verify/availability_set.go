package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// availabilitySetAPIVersion is the stable Microsoft.Compute availability
// sets API version the generic existence probe is pinned to (the
// provider's own availabilitysets SDK pin at v5.0.0).
const availabilitySetAPIVersion = "2024-03-01"

// availabilitySetVerifier verifies an AzureAvailabilitySet via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// set's full ARM ID.
type availabilitySetVerifier struct{}

// IDOutputKey is the availability set's full ARM ID.
func (*availabilitySetVerifier) IDOutputKey() string {
	return "availability_set_id"
}

func (*availabilitySetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, availabilitySetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureavailabilityset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureavailabilityset %q not found after deploy", id)
	}
	return nil
}

func (*availabilitySetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, availabilitySetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureavailabilityset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureavailabilityset %q still exists after destroy", id)
	}
	return nil
}
