package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// virtualMachineScaleSetAPIVersion is the stable Microsoft.Compute API
// version the generic existence probe is pinned to.
const virtualMachineScaleSetAPIVersion = "2024-07-01"

// virtualMachineScaleSetVerifier verifies an AzureVirtualMachineScaleSet
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the scale set's full ARM ID. The probe is orchestration-mode
// agnostic: UNIFORM and FLEXIBLE sets share the one ARM resource type.
type virtualMachineScaleSetVerifier struct{}

// IDOutputKey is the scale set's full ARM ID.
func (*virtualMachineScaleSetVerifier) IDOutputKey() string {
	return "scale_set_id"
}

func (*virtualMachineScaleSetVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualMachineScaleSetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualmachinescaleset verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevirtualmachinescaleset %q not found after deploy", id)
	}
	return nil
}

func (*virtualMachineScaleSetVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualMachineScaleSetAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevirtualmachinescaleset verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevirtualmachinescaleset %q still exists after destroy", id)
	}
	return nil
}
