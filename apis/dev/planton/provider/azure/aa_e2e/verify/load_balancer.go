package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// loadBalancerAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const loadBalancerAPIVersion = "2024-05-01"

// loadBalancerVerifier verifies an AzureLoadBalancer via the generic ARM
// resources GetByID (see armResourceExists), keyed on the load
// balancer's full ARM ID.
type loadBalancerVerifier struct{}

// IDOutputKey is the load balancer's full ARM ID.
func (*loadBalancerVerifier) IDOutputKey() string {
	return "load_balancer_id"
}

func (*loadBalancerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, loadBalancerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureloadbalancer verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureloadbalancer %q not found after deploy", id)
	}
	return nil
}

func (*loadBalancerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, loadBalancerAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureloadbalancer verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureloadbalancer %q still exists after destroy", id)
	}
	return nil
}
