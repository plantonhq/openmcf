package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// sslCertificateVerifier probes a self-managed Compute Engine SSL
// certificate by name and confirms GCP parsed the uploaded chain (a
// non-empty expiry proves real certificate material landed, not just an
// object shell). One kind maps to two API collections (global and regional
// SSL certificates), so the verifier routes on the region output.
type sslCertificateVerifier struct{}

func (v *sslCertificateVerifier) IDOutputKey() string { return "self_link" }

func (v *sslCertificateVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["certificate_name"]
	region := outputs["region"]

	if region != "" {
		cert, err := svc.Compute.RegionSslCertificates.Get(svc.Project, region, name).Context(ctx).Do()
		if err != nil {
			return errors.Wrapf(err, "regional ssl certificate %s/%s not found after deploy", region, name)
		}
		if cert.ExpireTime == "" {
			return errors.Errorf("regional ssl certificate %s/%s has no expire time — chain not parsed", region, name)
		}
		return nil
	}
	cert, err := svc.Compute.SslCertificates.Get(svc.Project, name).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "global ssl certificate %s not found after deploy", name)
	}
	if cert.ExpireTime == "" {
		return errors.Errorf("global ssl certificate %s has no expire time — chain not parsed", name)
	}
	return nil
}

func (v *sslCertificateVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	name := outputs["certificate_name"]
	region := outputs["region"]

	var err error
	if region != "" {
		_, err = svc.Compute.RegionSslCertificates.Get(svc.Project, region, name).Context(ctx).Do()
	} else {
		_, err = svc.Compute.SslCertificates.Get(svc.Project, name).Context(ctx).Do()
	}
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil
		}
		return errors.Wrapf(err, "unexpected error probing ssl certificate %s after destroy", name)
	}
	return errors.Errorf("ssl certificate %s still exists after destroy", name)
}
