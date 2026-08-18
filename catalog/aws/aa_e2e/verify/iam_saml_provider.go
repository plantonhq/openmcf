package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	pkgerrors "github.com/pkg/errors"
)

// iamSamlProviderVerifier verifies an AwsIamSamlProvider via
// GetSAMLProvider, keyed on the provider ARN (the resource's identity
// at AWS -- names never appear in the API). IAM is a global service,
// so the region parameter is ignored. A deleted provider returns the
// typed NoSuchEntity error, which is the "absent" signal.
type iamSamlProviderVerifier struct{}

func (*iamSamlProviderVerifier) IDOutputKey() string { return "provider_arn" }

func (*iamSamlProviderVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := iamSamlProviderExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsiamsamlprovider verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsiamsamlprovider %q not found after deploy", id)
	}
	return nil
}

func (*iamSamlProviderVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	exists, err := iamSamlProviderExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsiamsamlprovider verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsiamsamlprovider %q still exists after destroy", id)
	}
	return nil
}

func iamSamlProviderExists(ctx context.Context, cfg aws.Config, providerArn string) (bool, error) {
	client := iam.NewFromConfig(cfg)
	_, err := client.GetSAMLProvider(ctx, &iam.GetSAMLProviderInput{SAMLProviderArn: aws.String(providerArn)})
	if err != nil {
		if isIamNoSuchEntity(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
