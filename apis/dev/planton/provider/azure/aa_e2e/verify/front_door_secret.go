package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// frontDoorSecretVerifier verifies an AzureFrontDoorSecret via the
// generic ARM resources GetByID, keyed on the secret's full ARM ID.
// The secret is ARM metadata wrapping a Key Vault certificate -- the
// control plane answers for it; no data-plane read is needed.
type frontDoorSecretVerifier struct{}

// IDOutputKey is the secret's full ARM ID.
func (*frontDoorSecretVerifier) IDOutputKey() string {
	return "secret_id"
}

func (*frontDoorSecretVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorsecret verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefrontdoorsecret %q not found after deploy", id)
	}
	return nil
}

func (*frontDoorSecretVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cdnFrontDoorAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefrontdoorsecret verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefrontdoorsecret %q still exists after destroy", id)
	}
	return nil
}
