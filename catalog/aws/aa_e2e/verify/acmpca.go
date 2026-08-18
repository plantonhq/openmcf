package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acmpca"
	acmpcatypes "github.com/aws/aws-sdk-go-v2/service/acmpca/types"
	pkgerrors "github.com/pkg/errors"
)

// privateCaVerifier verifies an AwsPrivateCa via
// DescribeCertificateAuthority, keyed on the CA's ARN (the provider's
// import ID). Exists accepts ACTIVE or PENDING_CERTIFICATE (an
// unactivated subordinate is legitimately created-but-pending); from
// outputs, a CA whose composed activation ran MUST be ACTIVE, and
// each issued certificate must be retrievable. Absence tolerates the
// DELETED parking state - a destroyed CA stays restorable (and
// visible) for its permanent_deletion_time_in_days window.
type privateCaVerifier struct{}

func (*privateCaVerifier) IDOutputKey() string { return "certificate_authority_arn" }

func (*privateCaVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := acmpca.NewFromConfig(cfg).DescribeCertificateAuthority(ctx, &acmpca.DescribeCertificateAuthorityInput{
		CertificateAuthorityArn: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsprivateca verify-exists failed for %q", id)
	}
	status := out.CertificateAuthority.Status
	if status != acmpcatypes.CertificateAuthorityStatusActive &&
		status != acmpcatypes.CertificateAuthorityStatusPendingCertificate {
		return pkgerrors.Errorf("awsprivateca %q status %q, want ACTIVE or PENDING_CERTIFICATE", id, status)
	}
	return nil
}

func (*privateCaVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := acmpca.NewFromConfig(cfg).DescribeCertificateAuthority(ctx, &acmpca.DescribeCertificateAuthorityInput{
		CertificateAuthorityArn: aws.String(id),
	})
	if err != nil {
		var notFound *acmpcatypes.ResourceNotFoundException
		if pkgerrors.As(err, &notFound) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awsprivateca verify-absent failed for %q", id)
	}
	// A deleted CA parks in DELETED for its restore window - that IS
	// the destroyed state (billing stopped at delete).
	if out.CertificateAuthority.Status == acmpcatypes.CertificateAuthorityStatusDeleted {
		return nil
	}
	return pkgerrors.Errorf("awsprivateca %q still exists after destroy (status %q)", id, out.CertificateAuthority.Status)
}

// VerifyExistsFromOutputs raises the bar for activated CAs: the
// composed activation implies ACTIVE, and every issued certificate
// must be retrievable from the CA.
func (v *privateCaVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	caArn, _ := outputs["certificate_authority_arn"].(string)
	if caArn == "" {
		return pkgerrors.New("awsprivateca outputs carry no certificate_authority_arn")
	}
	client := acmpca.NewFromConfig(cfg)

	described, err := client.DescribeCertificateAuthority(ctx, &acmpca.DescribeCertificateAuthorityInput{
		CertificateAuthorityArn: aws.String(caArn),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsprivateca verify-exists failed for %q", caArn)
	}
	status := described.CertificateAuthority.Status

	if activationArn, _ := outputs["activation_certificate_arn"].(string); activationArn != "" {
		if status != acmpcatypes.CertificateAuthorityStatusActive {
			return pkgerrors.Errorf("awsprivateca %q activated but status %q, want ACTIVE", caArn, status)
		}
	} else if status != acmpcatypes.CertificateAuthorityStatusActive &&
		status != acmpcatypes.CertificateAuthorityStatusPendingCertificate {
		return pkgerrors.Errorf("awsprivateca %q status %q, want ACTIVE or PENDING_CERTIFICATE", caArn, status)
	}

	if issuedArns, _ := outputs["issued_certificate_arns"].(map[string]interface{}); len(issuedArns) > 0 {
		for name, rawArn := range issuedArns {
			certificateArn, _ := rawArn.(string)
			if certificateArn == "" {
				return pkgerrors.Errorf("awsprivateca issued certificate %q carries an empty ARN", name)
			}
			if _, err := client.GetCertificate(ctx, &acmpca.GetCertificateInput{
				CertificateAuthorityArn: aws.String(caArn),
				CertificateArn:          aws.String(certificateArn),
			}); err != nil {
				return pkgerrors.Wrapf(err, "awsprivateca issued certificate %q not retrievable", name)
			}
		}
	}
	return nil
}
