package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// ipGroupAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const ipGroupAPIVersion = "2024-05-01"

// ipGroupVerifier verifies an AzureIpGroup via the generic ARM resources
// GetByID (see armResourceExists), keyed on the group's full ARM ID.
type ipGroupVerifier struct{}

// IDOutputKey is the IP Group's full ARM ID.
func (*ipGroupVerifier) IDOutputKey() string {
	return "ip_group_id"
}

func (*ipGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, ipGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureipgroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureipgroup %q not found after deploy", id)
	}
	return nil
}

func (*ipGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, ipGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureipgroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureipgroup %q still exists after destroy", id)
	}
	return nil
}
