package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	codepipelinetypes "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	pkgerrors "github.com/pkg/errors"
)

// codePipelineVerifier verifies an AwsCodePipeline via GetPipeline.
//
// A missing pipeline surfaces as the typed PipelineNotFoundException.
// Deletion is synchronous (DeletePipeline is a single control-plane call),
// so no transitional deleting state needs handling.
type codePipelineVerifier struct{}

func (*codePipelineVerifier) IDOutputKey() string { return "pipeline_name" }

func (*codePipelineVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := codePipelineExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscodepipeline verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscodepipeline %q not found after deploy", id)
	}
	return nil
}

func (*codePipelineVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := codePipelineExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscodepipeline verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscodepipeline %q still exists after destroy", id)
	}
	return nil
}

func codePipelineExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := codepipeline.NewFromConfig(cfg, func(o *codepipeline.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetPipeline(ctx, &codepipeline.GetPipelineInput{Name: &name})
	if err != nil {
		var notFound *codepipelinetypes.PipelineNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
