package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// privateDnsRecordVerifier verifies an AzurePrivateDnsRecord via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// record set's full ARM ID -- the id's TYPE segment
// (.../privateDnsZones/{zone}/{TYPE}/{name}) names whichever of the
// seven variant resources the module created, so one verifier serves
// every record type. Pinned to the private DNS family's API version
// (privateDnsAPIVersion -- the line the provider vendors).
type privateDnsRecordVerifier struct{}

// IDOutputKey is the record set's full ARM ID.
func (*privateDnsRecordVerifier) IDOutputKey() string {
	return "record_id"
}

func (*privateDnsRecordVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsrecord verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureprivatednsrecord %q not found after deploy", id)
	}
	return nil
}

func (*privateDnsRecordVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, privateDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureprivatednsrecord verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureprivatednsrecord %q still exists after destroy", id)
	}
	return nil
}
