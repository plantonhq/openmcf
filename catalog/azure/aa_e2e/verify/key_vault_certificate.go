package verify

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	pkgerrors "github.com/pkg/errors"
)

// keyVaultCertificateVerifier verifies an AzureKeyVaultCertificate through
// the Key Vault DATA PLANE (azcertificates). Unlike keys, ARM exposes no
// read proxy for certificates -- azurerm synthesizes the certificate's
// "resource manager id" client-side -- so the data-plane GET is the only
// honest probe. It runs under the same ambient credential the deploy used,
// which holds the test subscription's Key Vault data-plane bootstrap grant
// (see the harness Setup).
//
// Keyed on the certificate's VERSIONLESS data-plane id
// (https://{vault}.vault.azure.net/certificates/{name}) so the probe is
// renewal-agnostic.
type keyVaultCertificateVerifier struct{}

// IDOutputKey is the certificate's versionless data-plane ID.
func (*keyVaultCertificateVerifier) IDOutputKey() string {
	return "versionless_id"
}

func (*keyVaultCertificateVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := keyVaultCertificateExists(ctx, cred, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvaultcertificate verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurekeyvaultcertificate %q not found after deploy", id)
	}
	return nil
}

func (*keyVaultCertificateVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := keyVaultCertificateExists(ctx, cred, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurekeyvaultcertificate verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurekeyvaultcertificate %q still exists after destroy", id)
	}
	return nil
}

// keyVaultCertificateExists GETs the certificate (latest version) through the
// vault's data plane and treats a typed 404 as absence. The providers purge
// soft-deleted certificates on destroy by default, and a soft-deleted
// certificate 404s the live GET regardless, so the absence probe holds in
// both cases.
func keyVaultCertificateExists(ctx context.Context, cred azcore.TokenCredential, versionlessID string) (bool, error) {
	vaultURL, name, err := splitCertificateID(versionlessID)
	if err != nil {
		return false, err
	}
	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		return false, err
	}
	// An empty version resolves to the latest.
	if _, err := client.GetCertificate(ctx, name, "", nil); err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// splitCertificateID parses a versionless certificate data-plane id
// (https://{vault}.vault.azure.net/certificates/{name}) into the vault base
// URL and the certificate name, failing loudly on any other shape.
func splitCertificateID(id string) (vaultURL string, name string, err error) {
	u, err := url.Parse(id)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", "", pkgerrors.Errorf("certificate id %q is not a vault data-plane URL", id)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "certificates" || parts[1] == "" {
		return "", "", pkgerrors.Errorf("certificate id %q does not match https://{vault}/certificates/{name}", id)
	}
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), parts[1], nil
}
