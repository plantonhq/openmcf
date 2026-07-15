package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	pkgerrors "github.com/pkg/errors"
)

// cloudFrontDistributionVerifier verifies an AwsCloudFront via
// GetDistribution, keyed on the distribution_id output. CloudFront is a
// global service (region is irrelevant to the API call). A distribution
// that is still propagating reports status "InProgress" -- it exists and is
// converging, so existence is simply "gettable"; the modules' own
// wait_for_deployment handles propagation. Deletion is terminal: once the
// disable-then-delete dance completes, GetDistribution returns NoSuchDistribution.
type cloudFrontDistributionVerifier struct{}

func (*cloudFrontDistributionVerifier) IDOutputKey() string { return "distribution_id" }

func (*cloudFrontDistributionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudFrontDistributionExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudfront verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudfront %q not found after deploy", id)
	}
	return nil
}

func (*cloudFrontDistributionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudFrontDistributionExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudfront verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudfront %q still exists after destroy", id)
	}
	return nil
}

func cloudFrontDistributionExists(ctx context.Context, cfg aws.Config, id string) (bool, error) {
	client := cloudfront.NewFromConfig(cfg)
	out, err := client.GetDistribution(ctx, &cloudfront.GetDistributionInput{Id: &id})
	if err != nil {
		var notFound *types.NoSuchDistribution
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return out.Distribution != nil, nil
}
