package module

import (
	"fmt"

	"github.com/pkg/errors"
	awseventbridgerulev1alpha1 "github.com/plantonhq/planton/catalog/aws/awseventbridgerule/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// targets iterates over the spec's targets and creates an EventTarget for each.
func targets(ctx *pulumi.Context, locals *Locals, createdRule *cloudwatch.EventRule, provider *aws.Provider) error {
	if len(locals.Spec.Targets) == 0 {
		return nil
	}

	for i, target := range locals.Spec.Targets {
		resourceName := fmt.Sprintf("%s-%s", locals.Target.Metadata.Name, target.Name)

		args := &cloudwatch.EventTargetArgs{
			Rule:     createdRule.Name,
			Arn:      pulumi.String(target.Arn.GetValue()),
			TargetId: pulumi.StringPtr(target.Name),
		}

		// Event bus — must match the rule's bus.
		if locals.Spec.EventBusName.GetValue() != "" {
			args.EventBusName = pulumi.StringPtr(locals.Spec.EventBusName.GetValue())
		}

		// One teardown-forcing decision governs the rule and its
		// targets: the same spec flag forces RemoveTargets on destroy.
		if locals.Spec.ForceDestroy {
			args.ForceDestroy = pulumi.BoolPtr(true)
		}

		// IAM role for target invocation
		if target.RoleArn.GetValue() != "" {
			args.RoleArn = pulumi.StringPtr(target.RoleArn.GetValue())
		}

		// -----------------------------------------------------------
		// Input transformation (mutually exclusive)
		// -----------------------------------------------------------

		if target.Input != "" {
			args.Input = pulumi.StringPtr(target.Input)
		}

		if target.InputPath != "" {
			args.InputPath = pulumi.StringPtr(target.InputPath)
		}

		if target.InputTransformer != nil {
			transformerArgs := &cloudwatch.EventTargetInputTransformerArgs{
				InputTemplate: pulumi.String(target.InputTransformer.InputTemplate),
			}
			if len(target.InputTransformer.InputPaths) > 0 {
				transformerArgs.InputPaths = pulumi.ToStringMap(target.InputTransformer.InputPaths)
			}
			args.InputTransformer = transformerArgs
		}

		// -----------------------------------------------------------
		// Dead letter config
		// -----------------------------------------------------------

		if target.DeadLetterConfig != nil && target.DeadLetterConfig.Arn.GetValue() != "" {
			args.DeadLetterConfig = &cloudwatch.EventTargetDeadLetterConfigArgs{
				Arn: pulumi.StringPtr(target.DeadLetterConfig.Arn.GetValue()),
			}
		}

		// -----------------------------------------------------------
		// Retry policy — both fields are presence-aware (absent = AWS
		// default) because zero retry attempts is a meaningful setting
		// (fail straight to the DLQ).
		// -----------------------------------------------------------

		if target.RetryPolicy != nil {
			retryArgs := &cloudwatch.EventTargetRetryPolicyArgs{}
			if target.RetryPolicy.MaximumEventAgeInSeconds != nil {
				retryArgs.MaximumEventAgeInSeconds = pulumi.IntPtr(int(target.RetryPolicy.GetMaximumEventAgeInSeconds()))
			}
			if target.RetryPolicy.MaximumRetryAttempts != nil {
				retryArgs.MaximumRetryAttempts = pulumi.IntPtr(int(target.RetryPolicy.GetMaximumRetryAttempts()))
			}
			args.RetryPolicy = retryArgs
		}

		// -----------------------------------------------------------
		// Service-typed parameter blocks (at most one — CEL-enforced)
		// -----------------------------------------------------------

		// SQS: message group for FIFO queues.
		if target.SqsTarget != nil && target.SqsTarget.MessageGroupId != "" {
			args.SqsTarget = &cloudwatch.EventTargetSqsTargetArgs{
				MessageGroupId: pulumi.StringPtr(target.SqsTarget.MessageGroupId),
			}
		}

		// Kinesis: shard routing via a partition-key JSONPath.
		if target.KinesisTarget != nil && target.KinesisTarget.PartitionKeyPath != "" {
			args.KinesisTarget = &cloudwatch.EventTargetKinesisTargetArgs{
				PartitionKeyPath: pulumi.StringPtr(target.KinesisTarget.PartitionKeyPath),
			}
		}

		// API destination: path/query/header parameters for the HTTP
		// invocation.
		if target.HttpTarget != nil {
			httpArgs := &cloudwatch.EventTargetHttpTargetArgs{}
			if len(target.HttpTarget.PathParameterValues) > 0 {
				httpArgs.PathParameterValues = pulumi.ToStringArray(target.HttpTarget.PathParameterValues)
			}
			if len(target.HttpTarget.QueryStringParameters) > 0 {
				httpArgs.QueryStringParameters = pulumi.ToStringMap(target.HttpTarget.QueryStringParameters)
			}
			if len(target.HttpTarget.HeaderParameters) > 0 {
				httpArgs.HeaderParameters = pulumi.ToStringMap(target.HttpTarget.HeaderParameters)
			}
			args.HttpTarget = httpArgs
		}

		// AWS Batch: job submission parameters.
		if target.BatchTarget != nil {
			batchArgs := &cloudwatch.EventTargetBatchTargetArgs{
				// The job definition is a reference (an AwsBatchJobDefinition's
				// revision-carrying ARN output) or a literal name/name:revision.
				JobDefinition: pulumi.String(target.BatchTarget.JobDefinition.GetValue()),
				JobName:       pulumi.String(target.BatchTarget.JobName),
			}
			if target.BatchTarget.ArraySize != 0 {
				batchArgs.ArraySize = pulumi.IntPtr(int(target.BatchTarget.ArraySize))
			}
			if target.BatchTarget.JobAttempts != 0 {
				batchArgs.JobAttempts = pulumi.IntPtr(int(target.BatchTarget.JobAttempts))
			}
			args.BatchTarget = batchArgs
		}

		// ECS RunTask: the target arn is the CLUSTER; the task definition,
		// sizing, networking, and placement live in this block.
		if target.EcsTarget != nil {
			args.EcsTarget = ecsTargetArgs(target.EcsTarget)
		}

		// Redshift Data API: the target arn is the CLUSTER; each event runs
		// the SQL statement here, authenticated by db_user (temporary
		// credentials) or a Secrets Manager secret.
		if target.RedshiftTarget != nil {
			redshiftArgs := &cloudwatch.EventTargetRedshiftTargetArgs{
				Database: pulumi.String(target.RedshiftTarget.Database),
			}
			if target.RedshiftTarget.DbUser != "" {
				redshiftArgs.DbUser = pulumi.StringPtr(target.RedshiftTarget.DbUser)
			}
			if target.RedshiftTarget.SecretsManagerArn != "" {
				redshiftArgs.SecretsManagerArn = pulumi.StringPtr(target.RedshiftTarget.SecretsManagerArn)
			}
			if target.RedshiftTarget.Sql != "" {
				redshiftArgs.Sql = pulumi.StringPtr(target.RedshiftTarget.Sql)
			}
			if target.RedshiftTarget.StatementName != "" {
				redshiftArgs.StatementName = pulumi.StringPtr(target.RedshiftTarget.StatementName)
			}
			if target.RedshiftTarget.WithEvent {
				redshiftArgs.WithEvent = pulumi.BoolPtr(true)
			}
			args.RedshiftTarget = redshiftArgs
		}

		// SSM Run Command: the target arn is the DOCUMENT; each selector
		// narrows the instances the command dispatches to (selectors AND
		// together).
		if len(target.RunCommandTargets) > 0 {
			selectors := make(cloudwatch.EventTargetRunCommandTargetArray, 0, len(target.RunCommandTargets))
			for _, s := range target.RunCommandTargets {
				selectors = append(selectors, cloudwatch.EventTargetRunCommandTargetArgs{
					Key:    pulumi.String(s.Key),
					Values: pulumi.ToStringArray(s.Values),
				})
			}
			args.RunCommandTargets = selectors
		}

		// SageMaker Pipelines: the target arn is the PIPELINE; each event
		// starts an execution with these parameters.
		if target.SagemakerPipelineTarget != nil {
			sagemakerArgs := &cloudwatch.EventTargetSagemakerPipelineTargetArgs{}
			if len(target.SagemakerPipelineTarget.PipelineParameterList) > 0 {
				params := make(cloudwatch.EventTargetSagemakerPipelineTargetPipelineParameterListArray, 0, len(target.SagemakerPipelineTarget.PipelineParameterList))
				for _, p := range target.SagemakerPipelineTarget.PipelineParameterList {
					params = append(params, cloudwatch.EventTargetSagemakerPipelineTargetPipelineParameterListArgs{
						Name:  pulumi.String(p.Name),
						Value: pulumi.String(p.Value),
					})
				}
				sagemakerArgs.PipelineParameterLists = params
			}
			args.SagemakerPipelineTarget = sagemakerArgs
		}

		// AppSync: the target arn is the API's endpoint ARN; each event
		// invokes the GraphQL operation.
		if target.AppsyncTarget != nil {
			args.AppsyncTarget = &cloudwatch.EventTargetAppsyncTargetArgs{
				GraphqlOperation: pulumi.StringPtr(target.AppsyncTarget.GraphqlOperation),
			}
		}

		// -----------------------------------------------------------
		// Create target
		// -----------------------------------------------------------

		_, err := cloudwatch.NewEventTarget(ctx, resourceName, args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "failed to create EventBridge target %q (index %d)", target.Name, i)
		}
	}

	return nil
}

// ecsTargetArgs maps the spec's ECS RunTask block onto the provider's shape.
func ecsTargetArgs(ecs *awseventbridgerulev1alpha1.AwsEventBridgeTargetEcsConfig) *cloudwatch.EventTargetEcsTargetArgs {
	args := &cloudwatch.EventTargetEcsTargetArgs{
		TaskDefinitionArn: pulumi.String(ecs.TaskDefinitionArn.GetValue()),
	}

	if ecs.TaskCount != 0 {
		args.TaskCount = pulumi.IntPtr(int(ecs.TaskCount))
	}
	if ecs.LaunchType != "" {
		args.LaunchType = pulumi.StringPtr(ecs.LaunchType)
	}
	if ecs.PlatformVersion != "" {
		args.PlatformVersion = pulumi.StringPtr(ecs.PlatformVersion)
	}
	if ecs.Group != "" {
		args.Group = pulumi.StringPtr(ecs.Group)
	}
	if ecs.PropagateTags != "" {
		args.PropagateTags = pulumi.StringPtr(ecs.PropagateTags)
	}
	if ecs.EnableEcsManagedTags {
		args.EnableEcsManagedTags = pulumi.BoolPtr(true)
	}
	if ecs.EnableExecuteCommand {
		args.EnableExecuteCommand = pulumi.BoolPtr(true)
	}

	if len(ecs.CapacityProviderStrategy) > 0 {
		strategies := make(cloudwatch.EventTargetEcsTargetCapacityProviderStrategyArray, 0, len(ecs.CapacityProviderStrategy))
		for _, s := range ecs.CapacityProviderStrategy {
			strategies = append(strategies, cloudwatch.EventTargetEcsTargetCapacityProviderStrategyArgs{
				CapacityProvider: pulumi.String(s.CapacityProvider),
				Base:             pulumi.IntPtr(int(s.Base)),
				Weight:           pulumi.IntPtr(int(s.Weight)),
			})
		}
		args.CapacityProviderStrategies = strategies
	}

	if nc := ecs.NetworkConfiguration; nc != nil {
		subnets := make([]string, 0, len(nc.Subnets))
		for _, ref := range nc.Subnets {
			subnets = append(subnets, ref.GetValue())
		}
		securityGroups := make([]string, 0, len(nc.SecurityGroups))
		for _, ref := range nc.SecurityGroups {
			securityGroups = append(securityGroups, ref.GetValue())
		}
		ncArgs := &cloudwatch.EventTargetEcsTargetNetworkConfigurationArgs{
			Subnets: pulumi.ToStringArray(subnets),
		}
		if len(securityGroups) > 0 {
			ncArgs.SecurityGroups = pulumi.ToStringArray(securityGroups)
		}
		if nc.AssignPublicIp {
			ncArgs.AssignPublicIp = pulumi.BoolPtr(true)
		}
		args.NetworkConfiguration = ncArgs
	}

	if len(ecs.OrderedPlacementStrategy) > 0 {
		strategies := make(cloudwatch.EventTargetEcsTargetOrderedPlacementStrategyArray, 0, len(ecs.OrderedPlacementStrategy))
		for _, s := range ecs.OrderedPlacementStrategy {
			entry := cloudwatch.EventTargetEcsTargetOrderedPlacementStrategyArgs{
				Type: pulumi.String(s.Type),
			}
			if s.Field != "" {
				entry.Field = pulumi.StringPtr(s.Field)
			}
			strategies = append(strategies, entry)
		}
		args.OrderedPlacementStrategies = strategies
	}

	if len(ecs.PlacementConstraints) > 0 {
		constraints := make(cloudwatch.EventTargetEcsTargetPlacementConstraintArray, 0, len(ecs.PlacementConstraints))
		for _, c := range ecs.PlacementConstraints {
			entry := cloudwatch.EventTargetEcsTargetPlacementConstraintArgs{
				Type: pulumi.String(c.Type),
			}
			if c.Expression != "" {
				entry.Expression = pulumi.StringPtr(c.Expression)
			}
			constraints = append(constraints, entry)
		}
		args.PlacementConstraints = constraints
	}

	// Tags applied to the ECS TASKS each event launches -- distinct from
	// the rule's own resource tags.
	if len(ecs.Tags) > 0 {
		args.Tags = pulumi.ToStringMap(ecs.Tags)
	}

	return args
}
