package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// dnsRecordVerifier verifies an AzureDnsRecord via the generic ARM
// resources GetByID (see armResourceExists), keyed on the record set's
// full ARM ID -- which embeds the record type as its own path segment
// (.../dnsZones/{zone}/{TYPE}/{name}), so one verifier covers all nine
// record types.
type dnsRecordVerifier struct{}

// IDOutputKey is the record set's full ARM ID.
func (*dnsRecordVerifier) IDOutputKey() string {
	return "record_id"
}

func (*dnsRecordVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurednsrecord verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurednsrecord %q not found after deploy", id)
	}
	return nil
}

func (*dnsRecordVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, publicDnsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurednsrecord verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurednsrecord %q still exists after destroy", id)
	}
	return nil
}
