package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	pkgerrors "github.com/pkg/errors"
)

// secretsManagerSecretVerifier verifies an AwsSecretsManagerSecret via
// DescribeSecret, keyed on secret_arn. A secret with DeletedDate set is
// soft-deleted (scheduled for deletion within its recovery window) and
// counts as ABSENT -- destroy with a non-zero recovery window leaves the
// record describable until the window elapses.
type secretsManagerSecretVerifier struct{}

func (*secretsManagerSecretVerifier) IDOutputKey() string { return "secret_arn" }

func (*secretsManagerSecretVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := secretsManagerSecretExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssecretsmanagersecret verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssecretsmanagersecret %q not found after deploy", id)
	}
	return nil
}

func (*secretsManagerSecretVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := secretsManagerSecretExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssecretsmanagersecret verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssecretsmanagersecret %q still exists after destroy", id)
	}
	return nil
}

func secretsManagerSecretExists(ctx context.Context, cfg aws.Config, secretARN, region string) (bool, error) {
	client := secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretARN),
	})
	if err == nil {
		// Soft-deleted secrets stay describable through the recovery
		// window; DeletedDate set means the destroy already happened.
		return out.DeletedDate == nil, nil
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}
