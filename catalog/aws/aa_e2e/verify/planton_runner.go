package verify

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	pkgerrors "github.com/pkg/errors"
)

// plantonRunnerVerifier verifies an AwsPlantonRunner via the ECS service
// that keeps the runner running, keyed on the service ARN (which encodes
// both the dedicated cluster and the service name). Verification is
// deliberately service-level, independent of task health: ECS reports the
// service ACTIVE regardless of whether the container's control-plane
// handshake succeeds, so the lane proves the appliance's full
// provisioning surface (secret, roles, log group, security group,
// cluster, task definition, service) with synthetic credentials and
// without a live control plane.
type plantonRunnerVerifier struct{}

func (*plantonRunnerVerifier) IDOutputKey() string { return "service_arn" }

func (*plantonRunnerVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsServiceActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsplantonrunner verify-exists failed for %q", id)
	}
	if !active {
		return pkgerrors.Errorf("awsplantonrunner %q not ACTIVE after deploy", id)
	}
	return nil
}

func (*plantonRunnerVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsServiceActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsplantonrunner verify-absent failed for %q", id)
	}
	if active {
		return pkgerrors.Errorf("awsplantonrunner %q still ACTIVE after destroy", id)
	}
	return nil
}

// VerifyRuntimeFailureCause pins the fake-token task's designed failure to the
// refused control-plane join and nothing else, with two independent AWS
// evidence sources: the service's STOPPED task history must show a task whose
// container RAN and exited non-zero (a pull failure reads `CannotPull...` in
// the stopped reason -- the exact distinguisher, failed immediately and
// loudly), and the task definition's CloudWatch log group must carry the
// runner's join-step error line ("joining as runner ...") -- emitted only
// after the process read the token from Secrets Manager and REACHED the
// control plane. Polling is bounded: an ECS task needs a pull-run-exit cycle
// (a couple of minutes) before its history and logs attest the cause.
func (*plantonRunnerVerifier) VerifyRuntimeFailureCause(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region, cause string) error {
	if cause != "refused-join" {
		return pkgerrors.Errorf("unsupported runtime failure cause %q for the runner (supported: refused-join)", cause)
	}

	serviceArn, _ := outputs["service_arn"].(string)
	logGroup, _ := outputs["log_group_name"].(string)
	if serviceArn == "" || logGroup == "" {
		return pkgerrors.New("service_arn or log_group_name missing from outputs -- cannot pin the runtime failure cause")
	}
	clusterName, serviceName, err := parseEcsServiceArn(serviceArn)
	if err != nil {
		return err
	}

	ecsClient := ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		if region != "" {
			o.Region = region
		}
	})
	logsClient := cloudwatchlogs.NewFromConfig(cfg, func(o *cloudwatchlogs.Options) {
		if region != "" {
			o.Region = region
		}
	})

	deadline := time.Now().Add(5 * time.Minute)
	ranAndExited := false
	lastState := "no stopped task yet"
	for {
		if !ranAndExited {
			stopped, listErr := ecsClient.ListTasks(ctx, &ecs.ListTasksInput{
				Cluster:       &clusterName,
				ServiceName:   &serviceName,
				DesiredStatus: ecstypes.DesiredStatusStopped,
			})
			if listErr == nil && len(stopped.TaskArns) > 0 {
				described, descErr := ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
					Cluster: &clusterName,
					Tasks:   stopped.TaskArns,
				})
				if descErr == nil {
					var states []string
					for _, task := range described.Tasks {
						reason := aws.ToString(task.StoppedReason)
						if strings.Contains(reason, "CannotPull") {
							return pkgerrors.Errorf("the runner task cannot PULL the image -- a registry problem, not the designed refused join: %s", reason)
						}
						for _, c := range task.Containers {
							if strings.Contains(aws.ToString(c.Reason), "CannotPull") {
								return pkgerrors.Errorf("the runner container cannot PULL the image -- a registry problem, not the designed refused join: %s", aws.ToString(c.Reason))
							}
							if c.ExitCode != nil && *c.ExitCode != 0 {
								ranAndExited = true
							}
							exit := "<nil>"
							if c.ExitCode != nil {
								exit = fmt.Sprintf("%d", *c.ExitCode)
							}
							states = append(states, fmt.Sprintf("task stopped-reason=%q container reason=%q exit=%s", reason, aws.ToString(c.Reason), exit))
						}
					}
					if !ranAndExited {
						lastState = fmt.Sprintf("%d stopped task(s), none with a non-zero container exit yet: %s",
							len(described.Tasks), strings.Join(states, " | "))
					}
				}
			}
		}

		if ranAndExited {
			events, logErr := logsClient.FilterLogEvents(ctx, &cloudwatchlogs.FilterLogEventsInput{
				LogGroupName:  &logGroup,
				FilterPattern: aws.String(`"joining as runner"`),
			})
			if logErr == nil && len(events.Events) > 0 {
				fmt.Printf("  [verify] CAUSE: a task ran and exited non-zero (no pull failures) and %q carries the join-step error -- the ONLY failure is the refused join\n", logGroup)
				return nil
			}
			lastState = "a task ran and exited non-zero; waiting for the join-step log line"
		}

		if time.Now().After(deadline) {
			return pkgerrors.Errorf("the runner task never attested the refused join within the window; last state: %s", lastState)
		}
		time.Sleep(15 * time.Second)
	}
}
