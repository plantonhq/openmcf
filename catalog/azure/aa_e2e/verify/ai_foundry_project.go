package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// aiFoundryProjectVerifier verifies an AzureAiFoundryProject via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// project's full ARM ID. Projects are ML workspaces at ARM (kind
// "Project") living in the HUB's resource group, so the ML family's
// API pin reads them.
type aiFoundryProjectVerifier struct{}

// IDOutputKey is the project's full ARM ID.
func (*aiFoundryProjectVerifier) IDOutputKey() string {
	return "ai_foundry_project_id"
}

func (*aiFoundryProjectVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureaifoundryproject verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureaifoundryproject %q not found after deploy", id)
	}
	return nil
}

func (*aiFoundryProjectVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, machineLearningAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureaifoundryproject verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureaifoundryproject %q still exists after destroy", id)
	}
	return nil
}
