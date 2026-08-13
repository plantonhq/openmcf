package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	pkgerrors "github.com/pkg/errors"
)

// acmCertificateVerifier verifies an AwsCertManagerCert via
// DescribeCertificate, keyed on the cert_arn output. The live scenario
// imports a self-signed certificate (ISSUED immediately); a requested
// certificate resting in PENDING_VALIDATION is equally healthy right
// after apply (the provider does not wait for validation). Existence
// therefore requires a describable certificate whose status is not one
// of the dead states (FAILED, VALIDATION_TIMED_OUT, REVOKED, EXPIRED) --
// a certificate that exists in a dead state was never a proven deploy.
// Deletion is immediate (no recovery window, unlike KMS), so absence is
// the ResourceNotFoundException.
type acmCertificateVerifier struct{}

func (*acmCertificateVerifier) IDOutputKey() string { return "cert_arn" }

func (*acmCertificateVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	cert, err := acmCertificateGet(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscertmanagercert verify-exists failed for %q", id)
	}
	if cert == nil {
		return pkgerrors.Errorf("awscertmanagercert %q not found after deploy", id)
	}
	switch cert.Status {
	case types.CertificateStatusFailed, types.CertificateStatusValidationTimedOut,
		types.CertificateStatusRevoked, types.CertificateStatusExpired:
		return pkgerrors.Errorf("awscertmanagercert %q exists but is in dead state %q", id, cert.Status)
	}
	return nil
}

func (*acmCertificateVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	cert, err := acmCertificateGet(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscertmanagercert verify-absent failed for %q", id)
	}
	if cert != nil {
		return pkgerrors.Errorf("awscertmanagercert %q still exists after destroy", id)
	}
	return nil
}

// acmCertificateGet returns the certificate detail, or nil when AWS reports
// ResourceNotFoundException.
func acmCertificateGet(ctx context.Context, cfg aws.Config, id, region string) (*types.CertificateDetail, error) {
	client := acm.NewFromConfig(cfg, func(o *acm.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return nil, nil
		}
		return nil, err
	}
	return out.Certificate, nil
}
