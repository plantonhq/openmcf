package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	pkgerrors "github.com/pkg/errors"
)

// lambdaFunctionVerifier verifies an AwsLambda function via GetFunction,
// keyed on the function_name output.
type lambdaFunctionVerifier struct{}

func (*lambdaFunctionVerifier) IDOutputKey() string { return "function_name" }

func (*lambdaFunctionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := lambdaFunctionExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambda verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awslambda function %q not found after deploy", id)
	}
	return nil
}

func (*lambdaFunctionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := lambdaFunctionExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awslambda verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awslambda function %q still exists after destroy", id)
	}
	return nil
}

func lambdaFunctionExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := lambda.NewFromConfig(cfg, func(o *lambda.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.GetFunction(ctx, &lambda.GetFunctionInput{FunctionName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.Configuration == nil {
		return false, nil
	}
	return true, nil
}
