package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// managedSslCertificateVerifier probes a global managed SSL certificate by
// name. PROVISIONING status is expected for E2E domains without real DNS.
type managedSslCertificateVerifier struct{}

func (v *managedSslCertificateVerifier) IDOutputKey() string { return "self_link" }

func (v *managedSslCertificateVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["certificate_name"]
	cert, err := svc.Compute.SslCertificates.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "managed ssl certificate %s not found after deploy", name)
	}
	if cert.Type != "MANAGED" {
		return errors.Errorf("ssl certificate %s has type %q, expected MANAGED", name, cert.Type)
	}
	return nil
}

func (v *managedSslCertificateVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["certificate_name"]
	_, err := svc.Compute.SslCertificates.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing ssl certificate %s after destroy", name)
	}
	return errors.Errorf("managed ssl certificate %s still exists after destroy", name)
}
