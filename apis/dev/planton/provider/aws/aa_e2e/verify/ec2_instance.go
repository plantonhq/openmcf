package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// ec2InstanceVerifier verifies an AwsEc2Instance via DescribeInstances.
// EC2 keeps terminated instances describable for a while after
// termination, so absence is state-aware: an instance in
// "shutting-down" or "terminated" counts as absent (the destroy
// succeeded; EC2 is just retaining the record), the same lifecycle
// class as NAT gateways and RDS clusters. A fully expired record
// returns the typed InvalidInstanceID.NotFound error.
type ec2InstanceVerifier struct{}

func (*ec2InstanceVerifier) IDOutputKey() string { return "instance_id" }

func (*ec2InstanceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ec2InstanceExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsec2instance verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsec2instance %q not found after deploy", id)
	}
	return nil
}

func (*ec2InstanceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ec2InstanceExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsec2instance verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsec2instance %q still exists after destroy", id)
	}
	return nil
}

func ec2InstanceExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{InstanceIds: []string{id}})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidInstanceID.NotFound" {
			return false, nil
		}
		return false, err
	}
	for _, reservation := range out.Reservations {
		for _, instance := range reservation.Instances {
			if instance.State == nil {
				continue
			}
			switch instance.State.Name {
			case ec2types.InstanceStateNameShuttingDown, ec2types.InstanceStateNameTerminated:
				// Terminated records linger; the instance is gone.
				continue
			default:
				return true, nil
			}
		}
	}
	return false, nil
}
