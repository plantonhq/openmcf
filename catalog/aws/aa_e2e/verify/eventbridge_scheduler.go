package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulertypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	pkgerrors "github.com/pkg/errors"
)

// eventBridgeSchedulerVerifier verifies an AwsEventBridgeScheduler
// arm-for-arm from its outputs: the schedule via GetSchedule (group and
// name parsed from the schedule ARN,
// "arn:...:schedule/{group}/{name}"), and the owned group - when the
// instance has one - via GetScheduleGroup. Absence asserts both are
// gone (a deleting group reports DELETING and counts as absent).
type eventBridgeSchedulerVerifier struct{}

func (*eventBridgeSchedulerVerifier) IDOutputKey() string { return "schedule_arn" }

func (v *eventBridgeSchedulerVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := scheduleExistsByArn(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgescheduler verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awseventbridgescheduler schedule %q not found after deploy", id)
	}
	return nil
}

func (v *eventBridgeSchedulerVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := scheduleExistsByArn(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awseventbridgescheduler verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awseventbridgescheduler schedule %q still exists after destroy", id)
	}
	return nil
}

func (v *eventBridgeSchedulerVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	scheduleArn := stringOutput(outputs, "schedule_arn")
	if scheduleArn == "" {
		return pkgerrors.New("awseventbridgescheduler verify-exists: no schedule_arn in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, scheduleArn, region); err != nil {
		return err
	}
	if groupArn := stringOutput(outputs, "group_arn"); groupArn != "" {
		exists, err := scheduleGroupExistsByArn(ctx, cfg, groupArn, region)
		if err != nil {
			return pkgerrors.Wrapf(err, "awseventbridgescheduler verify-exists failed for owned group %q", groupArn)
		}
		if !exists {
			return pkgerrors.Errorf("awseventbridgescheduler owned group %q not found after deploy", groupArn)
		}
	}
	return nil
}

func (v *eventBridgeSchedulerVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	scheduleArn := stringOutput(outputs, "schedule_arn")
	if scheduleArn == "" {
		return pkgerrors.New("awseventbridgescheduler verify-absent: no schedule_arn in outputs")
	}
	if err := v.VerifyAbsent(ctx, cfg, scheduleArn, region); err != nil {
		return err
	}
	if groupArn := stringOutput(outputs, "group_arn"); groupArn != "" {
		exists, err := scheduleGroupExistsByArn(ctx, cfg, groupArn, region)
		if err != nil {
			return pkgerrors.Wrapf(err, "awseventbridgescheduler verify-absent failed for owned group %q", groupArn)
		}
		if exists {
			return pkgerrors.Errorf("awseventbridgescheduler owned group %q still exists after destroy", groupArn)
		}
	}
	return nil
}

func schedulerClient(cfg aws.Config, region string) *scheduler.Client {
	return scheduler.NewFromConfig(cfg, func(o *scheduler.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// scheduleExistsByArn parses "arn:...:schedule/{group}/{name}" and calls
// GetSchedule with the two-part key.
func scheduleExistsByArn(ctx context.Context, cfg aws.Config, scheduleArn, region string) (bool, error) {
	groupName, scheduleName, err := parseSchedulerArn(scheduleArn, "schedule/")
	if err != nil {
		return false, err
	}
	_, err = schedulerClient(cfg, region).GetSchedule(ctx, &scheduler.GetScheduleInput{
		GroupName: aws.String(groupName),
		Name:      aws.String(scheduleName),
	})
	if err == nil {
		return true, nil
	}
	var notFound *schedulertypes.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}

// scheduleGroupExistsByArn parses "arn:...:schedule-group/{name}" and
// calls GetScheduleGroup. A DELETING group counts as absent.
func scheduleGroupExistsByArn(ctx context.Context, cfg aws.Config, groupArn, region string) (bool, error) {
	resource := arnResourcePart(groupArn)
	groupName := strings.TrimPrefix(resource, "schedule-group/")
	out, err := schedulerClient(cfg, region).GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
		Name: aws.String(groupName),
	})
	if err != nil {
		var notFound *schedulertypes.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.State == schedulertypes.ScheduleGroupStateDeleting {
		return false, nil
	}
	return true, nil
}

// parseSchedulerArn splits the resource part after the given prefix
// into its "{group}/{name}" halves.
func parseSchedulerArn(arn, prefix string) (string, string, error) {
	resource := strings.TrimPrefix(arnResourcePart(arn), prefix)
	parts := strings.SplitN(resource, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", pkgerrors.Errorf("unexpected scheduler ARN shape %q", arn)
	}
	return parts[0], parts[1], nil
}

// arnResourcePart returns everything after the fifth colon (the ARN's
// resource part).
func arnResourcePart(arn string) string {
	parts := strings.SplitN(arn, ":", 6)
	if len(parts) < 6 {
		return arn
	}
	return parts[5]
}
