package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// trafficManagerEndpointVerifier verifies an AzureTrafficManagerEndpoint
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the endpoint's full ARM ID -- the id's TYPE segment
// ({profile_id}/{TYPE}/{name} with AzureEndpoints, ExternalEndpoints,
// or NestedEndpoints) names whichever of the three variant resources
// the module created, so one verifier serves every endpoint type.
// Pinned to the Traffic Manager family's API version
// (trafficManagerAPIVersion).
type trafficManagerEndpointVerifier struct{}

// IDOutputKey is the endpoint's full ARM ID.
func (*trafficManagerEndpointVerifier) IDOutputKey() string {
	return "endpoint_id"
}

func (*trafficManagerEndpointVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, trafficManagerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuretrafficmanagerendpoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuretrafficmanagerendpoint %q not found after deploy", id)
	}
	return nil
}

func (*trafficManagerEndpointVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, trafficManagerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuretrafficmanagerendpoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuretrafficmanagerendpoint %q still exists after destroy", id)
	}
	return nil
}
