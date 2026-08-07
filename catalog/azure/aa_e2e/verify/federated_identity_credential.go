package verify

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// federatedIdentityCredentialAPIVersion is the stable Microsoft.ManagedIdentity
// API version the generic existence probe below is pinned to. Federated
// identity credentials are child resources of user-assigned identities and are
// served by the same resource provider (and API version) as the identities
// themselves.
const federatedIdentityCredentialAPIVersion = "2023-01-31"

// federatedIdentityCredentialVerifier verifies an
// AzureFederatedIdentityCredential via the generic ARM resources GetByID,
// keyed on the credential's full ARM ID
// ({identity-id}/federatedIdentityCredentials/{name}). A GET (not
// CheckExistenceByID's HEAD) is deliberate: Microsoft.ManagedIdentity does not
// implement HEAD and answers it with 405 Method Not Allowed (verified live on
// the parent identities). A typed 404 ResponseError is the absence signal; any
// other failure (auth, network) surfaces as a real error rather than
// masquerading as "absent".
type federatedIdentityCredentialVerifier struct{}

// IDOutputKey is the credential's full ARM ID.
func (*federatedIdentityCredentialVerifier) IDOutputKey() string {
	return "federated_identity_credential_id"
}

func (*federatedIdentityCredentialVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := federatedIdentityCredentialExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefederatedidentitycredential verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefederatedidentitycredential %q not found after deploy", id)
	}
	return nil
}

func (*federatedIdentityCredentialVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := federatedIdentityCredentialExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefederatedidentitycredential verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefederatedidentitycredential %q still exists after destroy", id)
	}
	return nil
}

func federatedIdentityCredentialExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, credentialID string) (bool, error) {
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	if _, err := client.GetByID(ctx, credentialID, federatedIdentityCredentialAPIVersion, nil); err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
