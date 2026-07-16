package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	pkgerrors "github.com/pkg/errors"
)

// cloudwatchCompositeAlarmVerifier verifies an AwsCloudwatchCompositeAlarm
// via DescribeAlarms, keyed on the alarm_name output. Composite alarms share
// the metric-alarm namespace but only surface when AlarmTypes includes
// CompositeAlarm — the shared lookup handles both types.
type cloudwatchCompositeAlarmVerifier struct{}

func (*cloudwatchCompositeAlarmVerifier) IDOutputKey() string { return "alarm_name" }

func (*cloudwatchCompositeAlarmVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudwatchAlarmExists(ctx, cfg, id, region, types.AlarmTypeCompositeAlarm)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchcompositealarm verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchcompositealarm %q not found after deploy", id)
	}
	return nil
}

func (*cloudwatchCompositeAlarmVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudwatchAlarmExists(ctx, cfg, id, region, types.AlarmTypeCompositeAlarm)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchcompositealarm verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchcompositealarm %q still exists after destroy", id)
	}
	return nil
}
