package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	pkgerrors "github.com/pkg/errors"
)

// cloudwatchLogGroupVerifier verifies an AwsCloudwatchLogGroup via
// DescribeLogGroups, keyed on the log_group_name output. Log group deletion
// is synchronous (no DELETING lifecycle state), so existence is simply an
// exact name match in the described page — DescribeLogGroups matches by
// prefix, which is why the result is re-checked for equality instead of
// trusting a non-empty page.
type cloudwatchLogGroupVerifier struct{}

func (*cloudwatchLogGroupVerifier) IDOutputKey() string { return "log_group_name" }

func (*cloudwatchLogGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudwatchLogGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchloggroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchloggroup %q not found after deploy", id)
	}
	return nil
}

func (*cloudwatchLogGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudwatchLogGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchloggroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchloggroup %q still exists after destroy", id)
	}
	return nil
}

func cloudwatchLogGroupExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := cloudwatchlogs.NewFromConfig(cfg, func(o *cloudwatchlogs.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &id,
	})
	if err != nil {
		return false, err
	}
	for _, group := range out.LogGroups {
		if group.LogGroupName != nil && *group.LogGroupName == id {
			return true, nil
		}
	}
	return false, nil
}
