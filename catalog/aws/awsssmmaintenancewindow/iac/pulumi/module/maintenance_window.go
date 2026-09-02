package module

import (
	"github.com/pkg/errors"
	awsssmmaintenancewindowv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsssmmaintenancewindow/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// maintenanceWindow creates the window with its folded targets and
// tasks, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - Enabled is a provider-default-TRUE toggle - rendered only on an
//     explicit choice so the module never fights the default (unset =
//     enabled; the provider needs a second API call to create a window
//     paused);
//   - targets and tasks are true window satellites: window_id (and the
//     target's name/description/resource_type, and the task's
//     task_type) force replacement of the registration;
//   - rate controls (max_concurrency/max_errors) are only legal on a
//     task WITH targets - AWS rejects them on untargeted tasks, so the
//     module renders them only when targets exist;
//   - RUN_COMMAND tasks REQUIRE targets server-side (live-caught 400:
//     "you must specify at least one resource as the target"; only
//     Automation/Lambda/Step Functions tasks may run untargeted), and a
//     WindowTargetIds selector needs the CLOUD-GENERATED registration
//     ID - so the module resolves WindowTargetIds values that name an
//     in-spec target entry to the created registration's ID (the
//     catalog's name-based join convention); values naming no in-spec
//     target pass through unchanged for externally registered IDs;
//   - the invocation union renders exactly the one arm the spec set
//     (the spec's CELs guarantee one arm, matching task_type).
func maintenanceWindow(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &ssm.MaintenanceWindowArgs{
		// metadata.name is the window name on both engines.
		Name:     pulumi.String(locals.Target.Metadata.Name),
		Schedule: pulumi.String(spec.Schedule),
		Duration: pulumi.Int(int(spec.Duration)),
		Cutoff:   pulumi.Int(int(spec.Cutoff)),
		Tags:     pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.Enabled != nil {
		args.Enabled = pulumi.Bool(*spec.Enabled)
	}
	if spec.AllowUnassociatedTargets {
		args.AllowUnassociatedTargets = pulumi.Bool(true)
	}
	if spec.ScheduleTimezone != "" {
		args.ScheduleTimezone = pulumi.String(spec.ScheduleTimezone)
	}
	if spec.ScheduleOffset != 0 {
		args.ScheduleOffset = pulumi.Int(int(spec.ScheduleOffset))
	}
	if spec.StartDate != "" {
		args.StartDate = pulumi.String(spec.StartDate)
	}
	if spec.EndDate != "" {
		args.EndDate = pulumi.String(spec.EndDate)
	}

	createdWindow, err := ssm.NewMaintenanceWindow(ctx, "window", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create window")
	}

	// Folded target registrations, keyed by name (name, description,
	// and resource_type force replacement of the registration at the
	// provider).
	targetIds := pulumi.StringMap{}
	createdTargets := map[string]*ssm.MaintenanceWindowTarget{}
	for _, t := range spec.Targets {
		createdTarget, err := windowTarget(ctx, createdWindow, t, provider)
		if err != nil {
			return errors.Wrapf(err, "target %s", t.Name)
		}
		targetIds[t.Name] = createdTarget.ID().ToStringOutput()
		createdTargets[t.Name] = createdTarget
	}

	// Folded tasks, keyed by name.
	taskIds := pulumi.StringMap{}
	for _, t := range spec.Tasks {
		createdTask, err := windowTask(ctx, createdWindow, t, createdTargets, provider)
		if err != nil {
			return errors.Wrapf(err, "task %s", t.Name)
		}
		taskIds[t.Name] = createdTask.WindowTaskId
	}

	ctx.Export(OpWindowId, createdWindow.ID())
	ctx.Export(OpTargetIds, targetIds)
	ctx.Export(OpTaskIds, taskIds)
	return nil
}

// windowTarget registers one folded target set with the window.
func windowTarget(ctx *pulumi.Context, createdWindow *ssm.MaintenanceWindow,
	t *awsssmmaintenancewindowv1alpha1.AwsSsmMaintenanceWindowTargetEntry, provider *aws.Provider) (*ssm.MaintenanceWindowTarget, error) {

	var selectors ssm.MaintenanceWindowTargetTargetArray
	for _, s := range t.Targets {
		selectors = append(selectors, &ssm.MaintenanceWindowTargetTargetArgs{
			Key:    pulumi.String(s.Key),
			Values: pulumi.ToStringArray(s.Values),
		})
	}

	args := &ssm.MaintenanceWindowTargetArgs{
		WindowId:     createdWindow.ID(),
		Name:         pulumi.String(t.Name),
		ResourceType: pulumi.String(t.ResourceType),
		Targets:      selectors,
	}
	if t.Description != "" {
		args.Description = pulumi.String(t.Description)
	}
	if t.OwnerInformation != "" {
		args.OwnerInformation = pulumi.String(t.OwnerInformation)
	}

	return ssm.NewMaintenanceWindowTarget(ctx, "target-"+t.Name, args,
		pulumi.Provider(provider), pulumi.Parent(createdWindow))
}

// windowTask registers one folded task with the window. createdTargets
// carries the in-spec target registrations by name so WindowTargetIds
// selector values can join to their cloud-generated IDs.
func windowTask(ctx *pulumi.Context, createdWindow *ssm.MaintenanceWindow,
	t *awsssmmaintenancewindowv1alpha1.AwsSsmMaintenanceWindowTaskEntry,
	createdTargets map[string]*ssm.MaintenanceWindowTarget, provider *aws.Provider) (*ssm.MaintenanceWindowTask, error) {

	args := &ssm.MaintenanceWindowTaskArgs{
		WindowId: createdWindow.ID(),
		Name:     pulumi.String(t.Name),
		TaskType: pulumi.String(t.TaskType),
		TaskArn:  pulumi.String(t.TaskArn.GetValue()),
		Priority: pulumi.Int(int(t.Priority)),
	}
	if t.Description != "" {
		args.Description = pulumi.String(t.Description)
	}
	if t.ServiceRoleArn.GetValue() != "" {
		args.ServiceRoleArn = pulumi.String(t.ServiceRoleArn.GetValue())
	}
	if t.CutoffBehavior != "" {
		args.CutoffBehavior = pulumi.String(t.CutoffBehavior)
	}

	var selectors ssm.MaintenanceWindowTaskTargetArray
	for _, s := range t.Targets {
		values := pulumi.StringArray{}
		for _, v := range s.Values {
			// WindowTargetIds values naming an in-spec target resolve
			// to the created registration's cloud-generated ID; unknown
			// values pass through (externally registered target IDs).
			if created, ok := createdTargets[v]; ok && s.Key == "WindowTargetIds" {
				values = append(values, created.ID().ToStringOutput())
			} else {
				values = append(values, pulumi.String(v))
			}
		}
		selectors = append(selectors, &ssm.MaintenanceWindowTaskTargetArgs{
			Key:    pulumi.String(s.Key),
			Values: values,
		})
	}
	if len(selectors) > 0 {
		args.Targets = selectors

		// Rate controls are only legal on a task WITH targets.
		if t.MaxConcurrency != "" {
			args.MaxConcurrency = pulumi.String(t.MaxConcurrency)
		}
		if t.MaxErrors != "" {
			args.MaxErrors = pulumi.String(t.MaxErrors)
		}
	}

	if t.Invocation != nil {
		args.TaskInvocationParameters = taskInvocation(t.Invocation)
	}

	return ssm.NewMaintenanceWindowTask(ctx, "task-"+t.Name, args,
		pulumi.Provider(provider), pulumi.Parent(createdWindow))
}

// taskInvocation renders the one invocation arm the spec set.
func taskInvocation(inv *awsssmmaintenancewindowv1alpha1.AwsSsmMaintenanceWindowTaskInvocation) *ssm.MaintenanceWindowTaskTaskInvocationParametersArgs {
	invocationArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersArgs{}

	if rc := inv.RunCommand; rc != nil {
		runCommandArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersRunCommandParametersArgs{}
		if rc.Comment != "" {
			runCommandArgs.Comment = pulumi.String(rc.Comment)
		}
		if rc.DocumentHash != "" {
			runCommandArgs.DocumentHash = pulumi.String(rc.DocumentHash)
		}
		if rc.DocumentHashType != "" {
			runCommandArgs.DocumentHashType = pulumi.String(rc.DocumentHashType)
		}
		if rc.DocumentVersion != "" {
			runCommandArgs.DocumentVersion = pulumi.String(rc.DocumentVersion)
		}
		if rc.OutputS3Bucket.GetValue() != "" {
			runCommandArgs.OutputS3Bucket = pulumi.String(rc.OutputS3Bucket.GetValue())
		}
		if rc.OutputS3KeyPrefix != "" {
			runCommandArgs.OutputS3KeyPrefix = pulumi.String(rc.OutputS3KeyPrefix)
		}
		if rc.ServiceRoleArn.GetValue() != "" {
			runCommandArgs.ServiceRoleArn = pulumi.String(rc.ServiceRoleArn.GetValue())
		}
		if rc.TimeoutSeconds != 0 {
			runCommandArgs.TimeoutSeconds = pulumi.Int(int(rc.TimeoutSeconds))
		}
		var parameters ssm.MaintenanceWindowTaskTaskInvocationParametersRunCommandParametersParameterArray
		for _, p := range rc.Parameters {
			parameters = append(parameters, &ssm.MaintenanceWindowTaskTaskInvocationParametersRunCommandParametersParameterArgs{
				Name:   pulumi.String(p.Name),
				Values: pulumi.ToStringArray(p.Values),
			})
		}
		if len(parameters) > 0 {
			runCommandArgs.Parameters = parameters
		}
		if cw := rc.CloudwatchConfig; cw != nil {
			cloudwatchArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersRunCommandParametersCloudwatchConfigArgs{
				CloudwatchOutputEnabled: pulumi.Bool(cw.CloudwatchOutputEnabled),
			}
			if cw.CloudwatchLogGroupName != "" {
				cloudwatchArgs.CloudwatchLogGroupName = pulumi.String(cw.CloudwatchLogGroupName)
			}
			runCommandArgs.CloudwatchConfig = cloudwatchArgs
		}
		if nc := rc.NotificationConfig; nc != nil {
			notificationArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersRunCommandParametersNotificationConfigArgs{}
			if nc.NotificationArn.GetValue() != "" {
				notificationArgs.NotificationArn = pulumi.String(nc.NotificationArn.GetValue())
			}
			if len(nc.NotificationEvents) > 0 {
				notificationArgs.NotificationEvents = pulumi.ToStringArray(nc.NotificationEvents)
			}
			if nc.NotificationType != "" {
				notificationArgs.NotificationType = pulumi.String(nc.NotificationType)
			}
			runCommandArgs.NotificationConfig = notificationArgs
		}
		invocationArgs.RunCommandParameters = runCommandArgs
	}

	if a := inv.Automation; a != nil {
		automationArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersAutomationParametersArgs{}
		if a.DocumentVersion != "" {
			automationArgs.DocumentVersion = pulumi.String(a.DocumentVersion)
		}
		var parameters ssm.MaintenanceWindowTaskTaskInvocationParametersAutomationParametersParameterArray
		for _, p := range a.Parameters {
			parameters = append(parameters, &ssm.MaintenanceWindowTaskTaskInvocationParametersAutomationParametersParameterArgs{
				Name:   pulumi.String(p.Name),
				Values: pulumi.ToStringArray(p.Values),
			})
		}
		if len(parameters) > 0 {
			automationArgs.Parameters = parameters
		}
		invocationArgs.AutomationParameters = automationArgs
	}

	if l := inv.Lambda; l != nil {
		lambdaArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersLambdaParametersArgs{}
		if l.ClientContext != "" {
			lambdaArgs.ClientContext = pulumi.String(l.ClientContext)
		}
		if l.Payload != "" {
			lambdaArgs.Payload = pulumi.String(l.Payload)
		}
		if l.Qualifier != "" {
			lambdaArgs.Qualifier = pulumi.String(l.Qualifier)
		}
		invocationArgs.LambdaParameters = lambdaArgs
	}

	if sf := inv.StepFunctions; sf != nil {
		stepFunctionsArgs := &ssm.MaintenanceWindowTaskTaskInvocationParametersStepFunctionsParametersArgs{}
		if sf.Input != "" {
			stepFunctionsArgs.Input = pulumi.String(sf.Input)
		}
		if sf.Name != "" {
			stepFunctionsArgs.Name = pulumi.String(sf.Name)
		}
		invocationArgs.StepFunctionsParameters = stepFunctionsArgs
	}

	return invocationArgs
}
