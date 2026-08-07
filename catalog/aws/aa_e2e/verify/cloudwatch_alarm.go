package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	pkgerrors "github.com/pkg/errors"
)

// cloudwatchAlarmVerifier verifies an AwsCloudwatchAlarm via DescribeAlarms,
// keyed on the alarm_name output. Alarm deletion is synchronous, so existence
// is a non-empty exact-name result. Only metric alarms are requested — the
// composite alarm has its own verifier over the CompositeAlarm alarm type.
type cloudwatchAlarmVerifier struct{}

func (*cloudwatchAlarmVerifier) IDOutputKey() string { return "alarm_name" }

func (*cloudwatchAlarmVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudwatchAlarmExists(ctx, cfg, id, region, types.AlarmTypeMetricAlarm)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchalarm verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awscloudwatchalarm %q not found after deploy", id)
	}
	return nil
}

func (*cloudwatchAlarmVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := cloudwatchAlarmExists(ctx, cfg, id, region, types.AlarmTypeMetricAlarm)
	if err != nil {
		return pkgerrors.Wrapf(err, "awscloudwatchalarm verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awscloudwatchalarm %q still exists after destroy", id)
	}
	return nil
}

// cloudwatchAlarmExists reports whether an alarm of the given type exists
// under the exact name. DescribeAlarms with AlarmNames is an exact-match
// lookup; an unknown name simply returns an empty result (no typed
// not-found error).
func cloudwatchAlarmExists(ctx context.Context, cfg aws.Config, id, region string, alarmType types.AlarmType) (bool, error) {
	client := cloudwatch.NewFromConfig(cfg, func(o *cloudwatch.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
		AlarmNames: []string{id},
		AlarmTypes: []types.AlarmType{alarmType},
	})
	if err != nil {
		return false, err
	}
	if alarmType == types.AlarmTypeCompositeAlarm {
		return len(out.CompositeAlarms) > 0, nil
	}
	return len(out.MetricAlarms) > 0, nil
}
