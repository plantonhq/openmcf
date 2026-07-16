package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppEnvironmentManagedCertificateVerifier verifies an
// AzureContainerAppEnvironmentManagedCertificate via the generic ARM
// resources GetByID (see armResourceExists), keyed on the managed
// certificate's full ARM ID (a child of the managed environment).
type containerAppEnvironmentManagedCertificateVerifier struct{}

// IDOutputKey is the managed certificate's full ARM ID.
func (*containerAppEnvironmentManagedCertificateVerifier) IDOutputKey() string {
	return "certificate_id"
}

func (*containerAppEnvironmentManagedCertificateVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentmanagedcertificate verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentmanagedcertificate %q not found after deploy", id)
	}
	return nil
}

func (*containerAppEnvironmentManagedCertificateVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentmanagedcertificate verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentmanagedcertificate %q still exists after destroy", id)
	}
	return nil
}
