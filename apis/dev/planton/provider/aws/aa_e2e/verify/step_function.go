package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/pkg/errors"
)

// stepFunctionVerifier verifies an AwsStepFunction via DescribeStateMachine,
// keyed on the state_machine_arn output. Deletion is asynchronous: a state
// machine in DELETING status is treated as absent (the NAT-gateway lifecycle
// class) because AWS keeps it describable until the deletion completes.
type stepFunctionVerifier struct{}

func (*stepFunctionVerifier) IDOutputKey() string { return "state_machine_arn" }

func (v *stepFunctionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := stepFunctionDescribe(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsstepfunction verify-exists failed for %q", id)
	}
	if out == nil || out.Status == types.StateMachineStatusDeleting {
		return errors.Errorf("awsstepfunction %q not found (or deleting) after deploy", id)
	}
	return nil
}

func (v *stepFunctionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := stepFunctionDescribe(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsstepfunction verify-absent failed for %q", id)
	}
	if out != nil && out.Status != types.StateMachineStatusDeleting {
		return errors.Errorf("awsstepfunction %q still exists after destroy", id)
	}
	return nil
}

// stepFunctionDescribe returns the describe output, or nil when the state
// machine does not exist. The SDK's typed not-found error is the absent
// signal; every other error is a real failure.
func stepFunctionDescribe(ctx context.Context, cfg aws.Config, arn, region string) (*sfn.DescribeStateMachineOutput, error) {
	client := sfn.NewFromConfig(cfg, func(o *sfn.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeStateMachine(ctx, &sfn.DescribeStateMachineInput{
		StateMachineArn: aws.String(arn),
	})
	if err != nil {
		var notFound *types.StateMachineDoesNotExist
		if errors.As(err, &notFound) {
			return nil, nil
		}
		return nil, err
	}
	return out, nil
}
