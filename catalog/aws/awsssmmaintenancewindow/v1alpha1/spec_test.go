package awsssmmaintenancewindowv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSsmMaintenanceWindowSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSsmMaintenanceWindowSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func boolPtr(b bool) *bool { return &b }

// minimalWindow is the smallest valid instance: a weekly window with
// no targets or tasks.
func minimalWindow() *AwsSsmMaintenanceWindowSpec {
	return &AwsSsmMaintenanceWindowSpec{
		Region:   "us-west-2",
		Schedule: "cron(0 2 ? * SUN *)",
		Duration: 3,
		Cutoff:   1,
	}
}

func runCommandTask() *AwsSsmMaintenanceWindowTaskEntry {
	return &AwsSsmMaintenanceWindowTaskEntry{
		Name:     "patch-scan",
		TaskType: "RUN_COMMAND",
		TaskArn:  svr("AWS-RunPatchBaseline"),
		Targets: []*AwsSsmMaintenanceWindowTargetSelector{{
			Key:    "WindowTargetIds",
			Values: []string{"target-id-placeholder"},
		}},
		MaxConcurrency: "10%",
		MaxErrors:      "1",
		Invocation: &AwsSsmMaintenanceWindowTaskInvocation{
			RunCommand: &AwsSsmMaintenanceWindowRunCommandInvocation{
				Comment:        "weekly patch scan",
				TimeoutSeconds: 600,
				Parameters: []*AwsSsmMaintenanceWindowTaskParameter{{
					Name:   "Operation",
					Values: []string{"Scan"},
				}},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsSsmMaintenanceWindowSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal window", func() {
			gomega.Expect(protovalidate.Validate(minimalWindow())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicitly disabled window with schedule refinements", func() {
			spec := minimalWindow()
			spec.Enabled = boolPtr(false)
			spec.ScheduleTimezone = "America/Los_Angeles"
			spec.ScheduleOffset = 2
			spec.StartDate = "2026-09-01T00:00:00-07:00"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a registered target and a matching run-command task", func() {
			spec := minimalWindow()
			spec.Targets = []*AwsSsmMaintenanceWindowTargetEntry{{
				Name:         "prod-instances",
				ResourceType: "INSTANCE",
				Targets: []*AwsSsmMaintenanceWindowTargetSelector{{
					Key:    "tag:env",
					Values: []string{"prod"},
				}},
			}}
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{runCommandTask()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts automation, lambda, and step-functions tasks with matching arms", func() {
			spec := minimalWindow()
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{
				{
					Name:     "runbook",
					TaskType: "AUTOMATION",
					TaskArn:  svr("AWS-RestartEC2Instance"),
					Invocation: &AwsSsmMaintenanceWindowTaskInvocation{
						Automation: &AwsSsmMaintenanceWindowAutomationInvocation{DocumentVersion: "$LATEST"},
					},
				},
				{
					Name:     "notify-fn",
					TaskType: "LAMBDA",
					TaskArn:  svr("arn:aws:lambda:us-west-2:123456789012:function:notify"),
					Invocation: &AwsSsmMaintenanceWindowTaskInvocation{
						Lambda: &AwsSsmMaintenanceWindowLambdaInvocation{Payload: `{"source":"mw"}`},
					},
				},
				{
					Name:     "orchestrate",
					TaskType: "STEP_FUNCTIONS",
					TaskArn:  svr("arn:aws:states:us-west-2:123456789012:stateMachine:deploy"),
					Invocation: &AwsSsmMaintenanceWindowTaskInvocation{
						StepFunctions: &AwsSsmMaintenanceWindowStepFunctionsInvocation{Name: "mw-run"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an untargeted task without rate controls", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.Targets = nil
			task.MaxConcurrency = ""
			task.MaxErrors = ""
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a cutoff equal to the duration", func() {
			spec := minimalWindow()
			spec.Cutoff = 3
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a duration above 24 hours", func() {
			spec := minimalWindow()
			spec.Duration = 25
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a schedule offset above 6", func() {
			spec := minimalWindow()
			spec.ScheduleOffset = 7
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate target names", func() {
			spec := minimalWindow()
			target := &AwsSsmMaintenanceWindowTargetEntry{
				Name:         "prod-instances",
				ResourceType: "INSTANCE",
				Targets: []*AwsSsmMaintenanceWindowTargetSelector{{
					Key:    "tag:env",
					Values: []string{"prod"},
				}},
			}
			spec.Targets = []*AwsSsmMaintenanceWindowTargetEntry{target, target}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate task names", func() {
			spec := minimalWindow()
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{runCommandTask(), runCommandTask()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target entry without selectors", func() {
			spec := minimalWindow()
			spec.Targets = []*AwsSsmMaintenanceWindowTargetEntry{{
				Name:         "empty",
				ResourceType: "INSTANCE",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a task without a task_arn", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.TaskArn = nil
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invocation with two arms", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.Invocation.Automation = &AwsSsmMaintenanceWindowAutomationInvocation{}
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invocation arm that contradicts the task type", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.TaskType = "AUTOMATION"
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a run-command timeout below 30 seconds", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.Invocation.RunCommand.TimeoutSeconds = 10
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown notification event", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.Invocation.RunCommand.NotificationConfig = &AwsSsmMaintenanceWindowNotificationConfig{
				NotificationArn:    svr("arn:aws:sns:us-west-2:123456789012:mw-events"),
				NotificationEvents: []string{"Started"},
			}
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed task max_concurrency", func() {
			spec := minimalWindow()
			task := runCommandTask()
			task.MaxConcurrency = "0%"
			spec.Tasks = []*AwsSsmMaintenanceWindowTaskEntry{task}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
