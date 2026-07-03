package verify

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// userAssignedIdentityAPIVersion is the stable Microsoft.ManagedIdentity API
// version the generic existence probe below is pinned to. The generic ARM
// resources client requires an explicit api-version because it is not generated
// per-service; pinning a stable version here beats pulling in the dedicated
// armmsi SDK module for what is only an existence check.
const userAssignedIdentityAPIVersion = "2023-01-31"

// userAssignedIdentityVerifier verifies an AzureUserAssignedIdentity via the
// generic ARM resources GetByID, keyed on the identity's fully-scoped ARM ID.
// A GET (not CheckExistenceByID's HEAD) is deliberate: Microsoft.ManagedIdentity
// does not implement HEAD and answers it with 405 Method Not Allowed (verified
// live). A typed 404 ResponseError is the absence signal; any other failure
// (auth, network) surfaces as a real error rather than masquerading as "absent".
type userAssignedIdentityVerifier struct{}

// IDOutputKey is the identity's fully-scoped ARM ID.
func (*userAssignedIdentityVerifier) IDOutputKey() string { return "identity_id" }

func (*userAssignedIdentityVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := userAssignedIdentityExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureuserassignedidentity verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureuserassignedidentity %q not found after deploy", id)
	}
	return nil
}

func (*userAssignedIdentityVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := userAssignedIdentityExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureuserassignedidentity verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureuserassignedidentity %q still exists after destroy", id)
	}
	return nil
}

func userAssignedIdentityExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, identityID string) (bool, error) {
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	if _, err := client.GetByID(ctx, identityID, userAssignedIdentityAPIVersion, nil); err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
