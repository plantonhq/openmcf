package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// trafficManagerAPIVersion pins the Traffic Manager family's ARM API
// version -- the line the provider vendors (trafficmanager/2022-04-01).
const trafficManagerAPIVersion = "2022-04-01"

// trafficManagerProfileVerifier verifies an AzureTrafficManagerProfile
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the profile's full ARM ID. Traffic Manager is a global service --
// the profile lives at location "global", which GetByID handles like
// any other resource.
type trafficManagerProfileVerifier struct{}

// IDOutputKey is the profile's full ARM ID.
func (*trafficManagerProfileVerifier) IDOutputKey() string {
	return "traffic_manager_profile_id"
}

func (*trafficManagerProfileVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, trafficManagerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuretrafficmanagerprofile verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuretrafficmanagerprofile %q not found after deploy", id)
	}
	return nil
}

func (*trafficManagerProfileVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, trafficManagerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuretrafficmanagerprofile verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuretrafficmanagerprofile %q still exists after destroy", id)
	}
	return nil
}
