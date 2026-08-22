package module

import (
	"github.com/pkg/errors"
	awseventbridgeschedulerv1alpha1 "github.com/plantonhq/planton/catalog/aws/awseventbridgescheduler/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/scheduler"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// schedule creates the schedule (and optionally its owned group) and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the schedule's name (metadata.name) and group are fixed for
//     life (replace-on-change);
//   - the group is a name-and-tags container (its own update path is
//     tags-only); the owned group carries the identity tags - the
//     schedule itself is UNTAGGABLE at AWS;
//   - a first deploy with a freshly created execution role is
//     eventually consistent (the provider retries the
//     role-not-assumable error for up to two minutes);
//   - with action_after_completion = DELETE, AWS deletes a completed
//     one-time schedule out from under IaC state (fire-and-forget
//     only).
func schedule(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var createdGroup *scheduler.ScheduleGroup
	if spec.Group != nil {
		var err error
		createdGroup, err = scheduler.NewScheduleGroup(ctx, "schedule_group", &scheduler.ScheduleGroupArgs{
			Name: pulumi.String(spec.Group.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create schedule group")
		}
		ctx.Export(OpGroupArn, createdGroup.Arn)
	} else {
		ctx.Export(OpGroupArn, pulumi.String(""))
	}

	args := &scheduler.ScheduleArgs{
		Name:               pulumi.String(locals.Target.Metadata.Name),
		ScheduleExpression: pulumi.String(spec.ScheduleExpression),
		FlexibleTimeWindow: buildFlexibleTimeWindow(spec.FlexibleTimeWindow),
		Target:             buildTarget(spec.Target),
	}

	// Owned group -> its name; else the joined group_name; else AWS's
	// "default" group.
	if createdGroup != nil {
		args.GroupName = createdGroup.Name
	} else if spec.GroupName != "" {
		args.GroupName = pulumi.String(spec.GroupName)
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.ScheduleExpressionTimezone != "" {
		args.ScheduleExpressionTimezone = pulumi.String(spec.ScheduleExpressionTimezone)
	}
	if spec.StartDate != "" {
		args.StartDate = pulumi.String(spec.StartDate)
	}
	if spec.EndDate != "" {
		args.EndDate = pulumi.String(spec.EndDate)
	}
	if spec.State != "" {
		args.State = pulumi.String(spec.State)
	}
	if spec.ActionAfterCompletion != "" {
		args.ActionAfterCompletion = pulumi.String(spec.ActionAfterCompletion)
	}
	if spec.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}

	createdSchedule, err := scheduler.NewSchedule(ctx, "schedule", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create schedule")
	}

	ctx.Export(OpScheduleArn, createdSchedule.Arn)
	ctx.Export(OpGroupName, createdSchedule.GroupName)
	return nil
}

func buildFlexibleTimeWindow(window *awseventbridgeschedulerv1alpha1.AwsEventBridgeScheduleFlexibleTimeWindow) *scheduler.ScheduleFlexibleTimeWindowArgs {
	args := &scheduler.ScheduleFlexibleTimeWindowArgs{
		Mode: pulumi.String(window.Mode),
	}
	if window.MaximumWindowInMinutes != nil {
		args.MaximumWindowInMinutes = pulumi.Int(int(*window.MaximumWindowInMinutes))
	}
	return args
}

func buildTarget(target *awseventbridgeschedulerv1alpha1.AwsEventBridgeScheduleTarget) *scheduler.ScheduleTargetArgs {
	args := &scheduler.ScheduleTargetArgs{
		Arn:     pulumi.String(target.Arn.GetValue()),
		RoleArn: pulumi.String(target.RoleArn.GetValue()),
	}

	if target.Input != "" {
		args.Input = pulumi.String(target.Input)
	}
	if target.DeadLetterQueueArn.GetValue() != "" {
		args.DeadLetterConfig = &scheduler.ScheduleTargetDeadLetterConfigArgs{
			Arn: pulumi.String(target.DeadLetterQueueArn.GetValue()),
		}
	}
	if target.RetryPolicy != nil {
		retryPolicy := &scheduler.ScheduleTargetRetryPolicyArgs{}
		if target.RetryPolicy.MaximumEventAgeInSeconds != nil {
			retryPolicy.MaximumEventAgeInSeconds = pulumi.Int(int(*target.RetryPolicy.MaximumEventAgeInSeconds))
		}
		if target.RetryPolicy.MaximumRetryAttempts != nil {
			retryPolicy.MaximumRetryAttempts = pulumi.Int(int(*target.RetryPolicy.MaximumRetryAttempts))
		}
		args.RetryPolicy = retryPolicy
	}
	if target.EcsParameters != nil {
		args.EcsParameters = buildEcsParameters(target.EcsParameters)
	}
	if target.EventbridgeParameters != nil {
		args.EventbridgeParameters = &scheduler.ScheduleTargetEventbridgeParametersArgs{
			DetailType: pulumi.String(target.EventbridgeParameters.DetailType),
			Source:     pulumi.String(target.EventbridgeParameters.Source),
		}
	}
	if target.KinesisParameters != nil {
		args.KinesisParameters = &scheduler.ScheduleTargetKinesisParametersArgs{
			PartitionKey: pulumi.String(target.KinesisParameters.PartitionKey),
		}
	}
	if target.SagemakerPipelineParameters != nil {
		pipelineParameters := scheduler.ScheduleTargetSagemakerPipelineParametersPipelineParameterArray{}
		for _, parameter := range target.SagemakerPipelineParameters.PipelineParameters {
			pipelineParameters = append(pipelineParameters, &scheduler.ScheduleTargetSagemakerPipelineParametersPipelineParameterArgs{
				Name:  pulumi.String(parameter.Name),
				Value: pulumi.String(parameter.Value),
			})
		}
		args.SagemakerPipelineParameters = &scheduler.ScheduleTargetSagemakerPipelineParametersArgs{
			PipelineParameters: pipelineParameters,
		}
	}
	if target.SqsParameters != nil {
		sqsParameters := &scheduler.ScheduleTargetSqsParametersArgs{}
		if target.SqsParameters.MessageGroupId != "" {
			sqsParameters.MessageGroupId = pulumi.String(target.SqsParameters.MessageGroupId)
		}
		args.SqsParameters = sqsParameters
	}

	return args
}

func buildEcsParameters(ecs *awseventbridgeschedulerv1alpha1.AwsEventBridgeScheduleEcsParameters) *scheduler.ScheduleTargetEcsParametersArgs {
	args := &scheduler.ScheduleTargetEcsParametersArgs{
		TaskDefinitionArn: pulumi.String(ecs.TaskDefinitionArn.GetValue()),
	}

	if ecs.TaskCount != nil {
		args.TaskCount = pulumi.Int(int(*ecs.TaskCount))
	}
	if ecs.LaunchType != "" {
		args.LaunchType = pulumi.String(ecs.LaunchType)
	}
	if ecs.Group != "" {
		args.Group = pulumi.String(ecs.Group)
	}
	if ecs.PlatformVersion != "" {
		args.PlatformVersion = pulumi.String(ecs.PlatformVersion)
	}
	if ecs.PropagateTags != "" {
		args.PropagateTags = pulumi.String(ecs.PropagateTags)
	}
	if ecs.ReferenceId != "" {
		args.ReferenceId = pulumi.String(ecs.ReferenceId)
	}
	if ecs.EnableEcsManagedTags {
		args.EnableEcsManagedTags = pulumi.Bool(true)
	}
	if ecs.EnableExecuteCommand {
		args.EnableExecuteCommand = pulumi.Bool(true)
	}
	if len(ecs.Tags) > 0 {
		args.Tags = pulumi.ToStringMap(ecs.Tags)
	}

	if len(ecs.CapacityProviderStrategy) > 0 {
		strategies := scheduler.ScheduleTargetEcsParametersCapacityProviderStrategyArray{}
		for _, strategy := range ecs.CapacityProviderStrategy {
			strategies = append(strategies, &scheduler.ScheduleTargetEcsParametersCapacityProviderStrategyArgs{
				CapacityProvider: pulumi.String(strategy.CapacityProvider),
				Base:             pulumi.Int(int(strategy.Base)),
				Weight:           pulumi.Int(int(strategy.Weight)),
			})
		}
		args.CapacityProviderStrategies = strategies
	}

	if ecs.NetworkConfiguration != nil {
		subnets := pulumi.StringArray{}
		for _, subnet := range ecs.NetworkConfiguration.Subnets {
			subnets = append(subnets, pulumi.String(subnet.GetValue()))
		}
		networkConfiguration := &scheduler.ScheduleTargetEcsParametersNetworkConfigurationArgs{
			Subnets:        subnets,
			AssignPublicIp: pulumi.Bool(ecs.NetworkConfiguration.AssignPublicIp),
		}
		if len(ecs.NetworkConfiguration.SecurityGroups) > 0 {
			securityGroups := pulumi.StringArray{}
			for _, securityGroup := range ecs.NetworkConfiguration.SecurityGroups {
				securityGroups = append(securityGroups, pulumi.String(securityGroup.GetValue()))
			}
			networkConfiguration.SecurityGroups = securityGroups
		}
		args.NetworkConfiguration = networkConfiguration
	}

	if len(ecs.PlacementConstraints) > 0 {
		constraints := scheduler.ScheduleTargetEcsParametersPlacementConstraintArray{}
		for _, constraint := range ecs.PlacementConstraints {
			constraintArgs := &scheduler.ScheduleTargetEcsParametersPlacementConstraintArgs{
				Type: pulumi.String(constraint.Type),
			}
			if constraint.Expression != "" {
				constraintArgs.Expression = pulumi.String(constraint.Expression)
			}
			constraints = append(constraints, constraintArgs)
		}
		args.PlacementConstraints = constraints
	}

	if len(ecs.PlacementStrategy) > 0 {
		strategies := scheduler.ScheduleTargetEcsParametersPlacementStrategyArray{}
		for _, strategy := range ecs.PlacementStrategy {
			strategyArgs := &scheduler.ScheduleTargetEcsParametersPlacementStrategyArgs{
				Type: pulumi.String(strategy.Type),
			}
			if strategy.Field != "" {
				strategyArgs.Field = pulumi.String(strategy.Field)
			}
			strategies = append(strategies, strategyArgs)
		}
		args.PlacementStrategies = strategies
	}

	return args
}
