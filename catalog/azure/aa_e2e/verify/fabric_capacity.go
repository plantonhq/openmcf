package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// fabricCapacityAPIVersion is the Microsoft.Fabric/capacities API
// line the pinned provider vendors (fabric/2023-11-01).
const fabricCapacityAPIVersion = "2023-11-01"

// fabricCapacityVerifier verifies an AzureFabricCapacity via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// capacity's ARM ID.
type fabricCapacityVerifier struct{}

func (*fabricCapacityVerifier) IDOutputKey() string {
	return "fabric_capacity_id"
}

func (*fabricCapacityVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, fabricCapacityAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefabriccapacity verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefabriccapacity %q not found after deploy", id)
	}
	return nil
}

func (*fabricCapacityVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, fabricCapacityAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefabriccapacity verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefabriccapacity %q still exists after destroy", id)
	}
	return nil
}
