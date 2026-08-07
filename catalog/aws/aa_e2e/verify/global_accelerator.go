package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/globalaccelerator"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// globalAcceleratorVerifier verifies an AwsGlobalAccelerator via
// DescribeAccelerator keyed on the accelerator ARN.
//
// Global Accelerator is a global service whose control-plane API is homed in
// us-west-2 only, so the client is pinned there regardless of the scenario's
// region (the same pin the provider applies during deploys). An accelerator is
// "present" in any describable state — DEPLOYED and IN_PROGRESS both count,
// since IN_PROGRESS is the normal propagation phase after every change. A
// deleted accelerator returns the typed AcceleratorNotFoundException, which is
// the "absent" signal; any other error is a genuine failure and must surface.
type globalAcceleratorVerifier struct{}

func (*globalAcceleratorVerifier) IDOutputKey() string { return "accelerator_arn" }

func (*globalAcceleratorVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := globalAcceleratorExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsglobalaccelerator verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsglobalaccelerator %q not found after deploy", id)
	}
	return nil
}

func (*globalAcceleratorVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := globalAcceleratorExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsglobalaccelerator verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsglobalaccelerator %q still exists after destroy", id)
	}
	return nil
}

func globalAcceleratorExists(ctx context.Context, cfg aws.Config, arn string) (bool, error) {
	client := globalaccelerator.NewFromConfig(cfg, func(o *globalaccelerator.Options) {
		// The Global Accelerator API lives in us-west-2 only.
		o.Region = "us-west-2"
	})
	_, err := client.DescribeAccelerator(ctx, &globalaccelerator.DescribeAcceleratorInput{
		AcceleratorArn: &arn,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AcceleratorNotFoundException" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
