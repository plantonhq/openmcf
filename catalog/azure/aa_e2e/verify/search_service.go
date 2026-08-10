package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// searchAPIVersion pins the Microsoft.Search GA line the search
// verifier reads with -- the version azurerm v5's own search clients
// vend (services and sharedPrivateLinkResources share it).
const searchAPIVersion = "2025-05-01"

// searchServiceVerifier verifies an AzureSearchService via the generic
// ARM resources GetByID (see armResourceExists), keyed on the
// service's full ARM ID. Existence is the honest bar: the composed
// shared private links are ARM children under the same path, and a
// link's PENDING approval state is expected (nothing in the lane
// approves the target side) -- the verifier asserts ARM state, never
// connection health. Search services have no soft-delete class:
// absence after destroy is genuine absence.
type searchServiceVerifier struct{}

// IDOutputKey is the service's full ARM ID.
func (*searchServiceVerifier) IDOutputKey() string {
	return "search_service_id"
}

func (*searchServiceVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, searchAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuresearchservice verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuresearchservice %q not found after deploy", id)
	}
	return nil
}

func (*searchServiceVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, searchAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuresearchservice verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuresearchservice %q still exists after destroy", id)
	}
	return nil
}
