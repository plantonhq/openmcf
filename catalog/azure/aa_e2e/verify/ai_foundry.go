package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// aiFoundryVerifier verifies an AzureAiFoundry hub via the generic ARM
// resources GetByID (see armResourceExists), keyed on the hub's full
// ARM ID. Hubs are ML workspaces at ARM (kind "Hub"), so the ML
// family's API pin reads them; existence is the honest bar -- the
// provider's own read cycle gates the deploy on the hub's properties
// and kind.
// Absence-after-destroy is genuine absence: a soft-deleted hub ghost
// is not returned by GetByID -- the ORPHAN sweep is what checks
// `az ml workspace list --archived` (ghosts hold the hub name).
type aiFoundryVerifier struct{}

// IDOutputKey is the hub's full ARM ID.
func (*aiFoundryVerifier) IDOutputKey() string {
	return "ai_foundry_id"
}

func (*aiFoundryVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureaifoundry verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureaifoundry %q not found after deploy", id)
	}
	return nil
}

func (*aiFoundryVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureaifoundry verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureaifoundry %q still exists after destroy", id)
	}
	return nil
}
