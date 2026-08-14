package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// computeGalleryImageAPIVersion is the stable Microsoft.Compute gallery
// images API version the generic existence probe is pinned to (the
// provider's own galleryimages SDK pin at v5.0.0).
const computeGalleryImageAPIVersion = "2022-03-03"

// computeGalleryImageVerifier verifies an AzureComputeGalleryImage via
// the generic ARM resources GetByID (see armResourceExists), keyed on
// the image definition's full ARM ID. Composed versions live under the
// definition's ARM path ({image_id}/versions/{name}) and are proven by
// the import round-trip's address-key derivations.
type computeGalleryImageVerifier struct{}

// IDOutputKey is the image definition's full ARM ID.
func (*computeGalleryImageVerifier) IDOutputKey() string {
	return "image_id"
}

func (*computeGalleryImageVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, computeGalleryImageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecomputegalleryimage verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecomputegalleryimage %q not found after deploy", id)
	}
	return nil
}

func (*computeGalleryImageVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, computeGalleryImageAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecomputegalleryimage verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecomputegalleryimage %q still exists after destroy", id)
	}
	return nil
}
