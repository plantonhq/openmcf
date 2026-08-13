package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// cognitiveAPIVersion pins the Microsoft.CognitiveServices GA line the
// cognitive-family verifiers read with -- the version that carries
// accounts, deployments, projects, and the responsible-AI children.
const cognitiveAPIVersion = "2025-06-01"

// cognitiveAccountVerifier verifies an AzureCognitiveAccount via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// account's full ARM ID. Existence is the honest bar: the provider's
// own read cycle gates the deploy on the account's properties, and the
// composed responsible-AI children are ARM children under the same
// path. Absence-after-destroy is genuine absence: the module purges
// the soft-deleted ghost on destroy (the provider default), and a
// ghost is not returned by GetByID either way -- the ORPHAN sweep is
// what checks `az cognitiveservices account list-deleted`.
type cognitiveAccountVerifier struct{}

// IDOutputKey is the account's full ARM ID.
func (*cognitiveAccountVerifier) IDOutputKey() string {
	return "cognitive_account_id"
}

func (*cognitiveAccountVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cognitiveAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecognitiveaccount verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecognitiveaccount %q not found after deploy", id)
	}
	return nil
}

func (*cognitiveAccountVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, cognitiveAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecognitiveaccount verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecognitiveaccount %q still exists after destroy", id)
	}
	return nil
}
