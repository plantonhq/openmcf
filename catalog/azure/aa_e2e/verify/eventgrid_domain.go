package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// eventgridDomainVerifier verifies an AzureEventgridDomain via the
// generic ARM resources GetByID (see armResourceExists), keyed on the
// domain's ARM ID. Pinned to the same eventgrid 2025-02-15 line as the
// topic verifier (the version the provider vendors for the family).
type eventgridDomainVerifier struct{}

func (*eventgridDomainVerifier) IDOutputKey() string {
	return "domain_id"
}

func (*eventgridDomainVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgriddomain verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureeventgriddomain %q not found after deploy", id)
	}
	return nil
}

func (*eventgridDomainVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, eventgridAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureeventgriddomain verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureeventgriddomain %q still exists after destroy", id)
	}
	return nil
}
