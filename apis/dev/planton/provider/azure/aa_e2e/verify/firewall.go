package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// firewallAPIVersion is the stable Microsoft.Network API version the
// generic existence probe is pinned to.
const firewallAPIVersion = "2024-05-01"

// firewallVerifier verifies an AzureFirewall via the generic ARM resources
// GetByID (see armResourceExists), keyed on the firewall's full ARM ID.
type firewallVerifier struct{}

// IDOutputKey is the firewall's full ARM ID.
func (*firewallVerifier) IDOutputKey() string {
	return "firewall_id"
}

func (*firewallVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, firewallAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefirewall verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurefirewall %q not found after deploy", id)
	}
	return nil
}

func (*firewallVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, firewallAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurefirewall verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurefirewall %q still exists after destroy", id)
	}
	return nil
}
