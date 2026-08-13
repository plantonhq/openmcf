package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// vpnServerConfigurationVerifier verifies an AzureVpnServerConfiguration
// via the generic ARM resources GetByID (see armResourceExists), keyed
// on the configuration's full ARM ID. Existence is the honest bar: the
// object is pure ARM metadata (authentication policy) and its composed
// policy groups are ARM children under the same path -- the provider's
// own read cycle gates the deploy on them. The Virtual WAN family
// shares the pinned virtualWanAPIVersion.
type vpnServerConfigurationVerifier struct{}

// IDOutputKey is the configuration's full ARM ID.
func (*vpnServerConfigurationVerifier) IDOutputKey() string {
	return "vpn_server_configuration_id"
}

func (*vpnServerConfigurationVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpnserverconfiguration verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurevpnserverconfiguration %q not found after deploy", id)
	}
	return nil
}

func (*vpnServerConfigurationVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, virtualWanAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurevpnserverconfiguration verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurevpnserverconfiguration %q still exists after destroy", id)
	}
	return nil
}
