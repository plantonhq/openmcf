package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// monitorActionGroupAPIVersion pins the Microsoft.Insights actionGroups RP
// version the verifier reads with.
const monitorActionGroupAPIVersion = "2023-01-01"

// monitorActionGroupVerifier verifies an AzureMonitorActionGroup via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// action group's ARM ID.
type monitorActionGroupVerifier struct{}

func (*monitorActionGroupVerifier) IDOutputKey() string {
	return "action_group_id"
}

func (*monitorActionGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorActionGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitoractiongroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremonitoractiongroup %q not found after deploy", id)
	}
	return nil
}

func (*monitorActionGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, monitorActionGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremonitoractiongroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremonitoractiongroup %q still exists after destroy", id)
	}
	return nil
}
