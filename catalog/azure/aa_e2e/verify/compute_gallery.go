package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// computeGalleryAPIVersion is the stable Microsoft.Compute galleries API
// version the generic existence probe is pinned to (the provider's own
// galleries SDK pin at v5.0.0).
const computeGalleryAPIVersion = "2022-03-03"

// computeGalleryVerifier verifies an AzureComputeGallery via the generic
// ARM resources GetByID (see armResourceExists), keyed on the gallery's
// full ARM ID.
type computeGalleryVerifier struct{}

// IDOutputKey is the gallery's full ARM ID.
func (*computeGalleryVerifier) IDOutputKey() string {
	return "gallery_id"
}

func (*computeGalleryVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, computeGalleryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecomputegallery verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecomputegallery %q not found after deploy", id)
	}
	return nil
}

func (*computeGalleryVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, computeGalleryAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecomputegallery verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecomputegallery %q still exists after destroy", id)
	}
	return nil
}
