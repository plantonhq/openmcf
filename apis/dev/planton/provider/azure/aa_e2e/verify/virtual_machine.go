package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualMachineAPIVersion is the stable Microsoft.Compute virtualMachines
// API version the generic existence probe is pinned to.
const virtualMachineAPIVersion = "2024-07-01"

// virtualMachineVerifier verifies an AzureVirtualMachine via the generic
// ARM resources GetByID (see armResourceExists), keyed on the VM's full
// ARM ID.
type virtualMachineVerifier struct{}

// IDOutputKey is the VM's full ARM ID.
func (*virtualMachineVerifier) IDOutputKey() string {
	return "vm_id"
}

func (*virtualMachineVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualMachineAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualmachine verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualmachine %q not found after deploy", id)
	}
	return nil
}

func (*virtualMachineVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualMachineAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualmachine verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualmachine %q still exists after destroy", id)
	}
	return nil
}
