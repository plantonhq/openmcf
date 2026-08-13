package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cognitiveAccountProjectVerifier verifies an
// AzureCognitiveAccountProject via the generic ARM resources GetByID
// (see armResourceExists), keyed on the project's full ARM ID (an
// account child: .../accounts/{account}/projects/{name}). Existence is
// the honest bar: the object is a workspace container -- its agents
// and evaluations are data-plane and never the provisioning contract.
// The cognitive family shares the pinned cognitiveAPIVersion.
type cognitiveAccountProjectVerifier struct{}

// IDOutputKey is the project's full ARM ID.
func (*cognitiveAccountProjectVerifier) IDOutputKey() string {
	return "project_id"
}

func (*cognitiveAccountProjectVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cognitiveAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecognitiveaccountproject verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecognitiveaccountproject %q not found after deploy", id)
	}
	return nil
}

func (*cognitiveAccountProjectVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cognitiveAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecognitiveaccountproject verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecognitiveaccountproject %q still exists after destroy", id)
	}
	return nil
}
