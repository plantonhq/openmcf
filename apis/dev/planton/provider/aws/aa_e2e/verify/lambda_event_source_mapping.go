package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	pkgerrors "github.com/pkg/errors"
)

// lambdaEventSourceMappingVerifier verifies an AwsLambdaEventSourceMapping via
// GetEventSourceMapping, keyed on the uuid output.
type lambdaEventSourceMappingVerifier struct{}

func (*lambdaEventSourceMappingVerifier) IDOutputKey() string { return "uuid" }

func (*lambdaEventSourceMappingVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := lambdaEventSourceMappingExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambdaeventsourcemapping verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awslambdaeventsourcemapping %q not found after deploy", id)
	}
	return nil
}

func (*lambdaEventSourceMappingVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := lambdaEventSourceMappingExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambdaeventsourcemapping verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awslambdaeventsourcemapping %q still exists after destroy", id)
	}
	return nil
}

func lambdaEventSourceMappingExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := lambda.NewFromConfig(cfg, func(o *lambda.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.GetEventSourceMapping(ctx, &lambda.GetEventSourceMappingInput{UUID: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return out.UUID != nil, nil
}
