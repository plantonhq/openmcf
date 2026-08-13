package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dataProtectionResourceGuardVerifier verifies an
// AzureDataProtectionResourceGuard via the generic ARM resources
// GetByID (see armResourceExists), keyed on the guard's full ARM ID.
// The guard is a pure governance object with no soft-delete semantics
// -- absence after destroy is immediate (a guard deletes cleanly even
// while vaults reference it).
type dataProtectionResourceGuardVerifier struct{}

// IDOutputKey is the guard's full ARM ID.
func (*dataProtectionResourceGuardVerifier) IDOutputKey() string {
	return "resource_guard_id"
}

func (*dataProtectionResourceGuardVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionresourceguard verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuredataprotectionresourceguard %q not found after deploy", id)
	}
	return nil
}

func (*dataProtectionResourceGuardVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, dataProtectionAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuredataprotectionresourceguard verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuredataprotectionresourceguard %q still exists after destroy", id)
	}
	return nil
}
