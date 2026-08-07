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
// certificate resting in PENDING_VALIDATION is equally "describable".
// Existence is simply whether DescribeCertificate succeeds. Deletion is
// immediate (no recovery window, unlike KMS), so absence is the
// ResourceNotFoundException.
type acmCertificateVerifier struct{}

func (*acmCertificateVerifier) IDOutputKey() string { return "cert_arn" }

func (*acmCertificateVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := acmCertificateExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscertmanagercert verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscertmanagercert %q not found after deploy", id)
	}
	return nil
}

func (*acmCertificateVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := acmCertificateExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscertmanagercert verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscertmanagercert %q still exists after destroy", id)
	}
	return nil
}

func acmCertificateExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := acm.NewFromConfig(cfg, func(o *acm.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeCertificate(ctx, &acm.DescribeCertificateInput{CertificateArn: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return out.Certificate != nil, nil
}
