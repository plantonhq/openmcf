package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppEnvironmentCertificateVerifier verifies an
// AzureContainerAppEnvironmentCertificate via the generic ARM resources
// GetByID (see armResourceExists), keyed on the certificate's full ARM ID
// (a child of the managed environment).
type containerAppEnvironmentCertificateVerifier struct{}

// IDOutputKey is the certificate's full ARM ID.
func (*containerAppEnvironmentCertificateVerifier) IDOutputKey() string {
	return "certificate_id"
}

func (*containerAppEnvironmentCertificateVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentcertificate verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentcertificate %q not found after deploy", id)
	}
	return nil
}

func (*containerAppEnvironmentCertificateVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappenvironmentcertificate verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappenvironmentcertificate %q still exists after destroy", id)
	}
	return nil
}
