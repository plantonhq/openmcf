package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	pkgerrors "github.com/pkg/errors"
)

// autoScalingGroupVerifier verifies an AwsAutoScalingGroup via
// DescribeAutoScalingGroups, keyed on the group name. Unlike most EC2 APIs,
// describing a nonexistent group is not an error -- the call succeeds with
// an empty result -- so absence is an empty group list, not a typed error
// code.
type autoScalingGroupVerifier struct{}

func (*autoScalingGroupVerifier) IDOutputKey() string { return "autoscaling_group_name" }

func (*autoScalingGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := autoScalingGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsautoscalinggroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsautoscalinggroup %q not found after deploy", id)
	}
	return nil
}

func (*autoScalingGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := autoScalingGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsautoscalinggroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsautoscalinggroup %q still exists after destroy", id)
	}
	return nil
}

func autoScalingGroupExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := autoscaling.NewFromConfig(cfg, func(o *autoscaling.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeAutoScalingGroups(ctx, &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{name},
	})
	if err != nil {
		return false, err
	}
	return len(out.AutoScalingGroups) > 0, nil
}
