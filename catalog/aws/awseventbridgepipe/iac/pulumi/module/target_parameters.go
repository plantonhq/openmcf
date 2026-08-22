package module

import (
	awseventbridgepipev1alpha1 "github.com/plantonhq/planton/catalog/aws/awseventbridgepipe/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/pipes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildTargetParameters renders the target-family invocation shaping
// (the spec's CEL guarantees at most one family block) plus the input
// template.
func buildTargetParameters(target *awseventbridgepipev1alpha1.AwsEventBridgePipeTargetParameters) *pipes.PipeTargetParametersArgs {
	args := &pipes.PipeTargetParametersArgs{}

	if target.InputTemplate != "" {
		args.InputTemplate = pulumi.String(target.InputTemplate)
	}

	if target.Sqs != nil {
		sqs := &pipes.PipeTargetParametersSqsQueueParametersArgs{}
		if target.Sqs.MessageGroupId != "" {
			sqs.MessageGroupId = pulumi.String(target.Sqs.MessageGroupId)
		}
		if target.Sqs.MessageDeduplicationId != "" {
			sqs.MessageDeduplicationId = pulumi.String(target.Sqs.MessageDeduplicationId)
		}
		args.SqsQueueParameters = sqs
	}

	if target.Kinesis != nil {
		args.KinesisStreamParameters = &pipes.PipeTargetParametersKinesisStreamParametersArgs{
			PartitionKey: pulumi.String(target.Kinesis.PartitionKey),
		}
	}

	if target.Lambda != nil {
		args.LambdaFunctionParameters = &pipes.PipeTargetParametersLambdaFunctionParametersArgs{
			InvocationType: pulumi.String(target.Lambda.InvocationType),
		}
	}

	if target.StepFunction != nil {
		args.StepFunctionStateMachineParameters = &pipes.PipeTargetParametersStepFunctionStateMachineParametersArgs{
			InvocationType: pulumi.String(target.StepFunction.InvocationType),
		}
	}

	if target.EcsTask != nil {
		args.EcsTaskParameters = buildEcsTaskParameters(target.EcsTask)
	}

	if target.BatchJob != nil {
		args.BatchJobParameters = buildBatchJobParameters(target.BatchJob)
	}

	if target.RedshiftData != nil {
		redshift := &pipes.PipeTargetParametersRedshiftDataParametersArgs{
			Database: pulumi.String(target.RedshiftData.Database),
			Sqls:     pulumi.ToStringArray(target.RedshiftData.Sqls),
		}
		if target.RedshiftData.DbUser != "" {
			redshift.DbUser = pulumi.String(target.RedshiftData.DbUser)
		}
		if target.RedshiftData.SecretManagerArn != "" {
			redshift.SecretManagerArn = pulumi.String(target.RedshiftData.SecretManagerArn)
		}
		if target.RedshiftData.StatementName != "" {
			redshift.StatementName = pulumi.String(target.RedshiftData.StatementName)
		}
		if target.RedshiftData.WithEvent {
			redshift.WithEvent = pulumi.Bool(true)
		}
		args.RedshiftDataParameters = redshift
	}

	if target.SagemakerPipeline != nil {
		pipelineParameters := pipes.PipeTargetParametersSagemakerPipelineParametersPipelineParameterArray{}
		for _, parameter := range target.SagemakerPipeline.PipelineParameters {
			pipelineParameters = append(pipelineParameters, &pipes.PipeTargetParametersSagemakerPipelineParametersPipelineParameterArgs{
				Name:  pulumi.String(parameter.Name),
				Value: pulumi.String(parameter.Value),
			})
		}
		args.SagemakerPipelineParameters = &pipes.PipeTargetParametersSagemakerPipelineParametersArgs{
			PipelineParameters: pipelineParameters,
		}
	}

	if target.EventbridgeEventBus != nil {
		bus := &pipes.PipeTargetParametersEventbridgeEventBusParametersArgs{}
		if target.EventbridgeEventBus.DetailType != "" {
			bus.DetailType = pulumi.String(target.EventbridgeEventBus.DetailType)
		}
		if target.EventbridgeEventBus.Source != "" {
			bus.Source = pulumi.String(target.EventbridgeEventBus.Source)
		}
		if target.EventbridgeEventBus.EndpointId != "" {
			bus.EndpointId = pulumi.String(target.EventbridgeEventBus.EndpointId)
		}
		if len(target.EventbridgeEventBus.Resources) > 0 {
			bus.Resources = pulumi.ToStringArray(target.EventbridgeEventBus.Resources)
		}
		if target.EventbridgeEventBus.Time != "" {
			bus.Time = pulumi.String(target.EventbridgeEventBus.Time)
		}
		args.EventbridgeEventBusParameters = bus
	}

	if target.CloudwatchLogs != nil {
		cloudwatchLogs := &pipes.PipeTargetParametersCloudwatchLogsParametersArgs{}
		if target.CloudwatchLogs.LogStreamName != "" {
			cloudwatchLogs.LogStreamName = pulumi.String(target.CloudwatchLogs.LogStreamName)
		}
		if target.CloudwatchLogs.Timestamp != "" {
			cloudwatchLogs.Timestamp = pulumi.String(target.CloudwatchLogs.Timestamp)
		}
		args.CloudwatchLogsParameters = cloudwatchLogs
	}

	if target.Http != nil {
		httpParameters := &pipes.PipeTargetParametersHttpParametersArgs{
			HeaderParameters:      pulumi.ToStringMap(target.Http.HeaderParameters),
			QueryStringParameters: pulumi.ToStringMap(target.Http.QueryStringParameters),
		}
		if target.Http.PathParameterValue != "" {
			httpParameters.PathParameterValues = pulumi.String(target.Http.PathParameterValue)
		}
		args.HttpParameters = httpParameters
	}

	return args
}

func buildEcsTaskParameters(ecs *awseventbridgepipev1alpha1.AwsEventBridgePipeEcsTaskTargetParameters) *pipes.PipeTargetParametersEcsTaskParametersArgs {
	args := &pipes.PipeTargetParametersEcsTaskParametersArgs{
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
		strategies := pipes.PipeTargetParametersEcsTaskParametersCapacityProviderStrategyArray{}
		for _, strategy := range ecs.CapacityProviderStrategy {
			strategies = append(strategies, &pipes.PipeTargetParametersEcsTaskParametersCapacityProviderStrategyArgs{
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
		vpcConfiguration := &pipes.PipeTargetParametersEcsTaskParametersNetworkConfigurationAwsVpcConfigurationArgs{
			Subnets: subnets,
		}
		if len(ecs.NetworkConfiguration.SecurityGroups) > 0 {
			securityGroups := pulumi.StringArray{}
			for _, securityGroup := range ecs.NetworkConfiguration.SecurityGroups {
				securityGroups = append(securityGroups, pulumi.String(securityGroup.GetValue()))
			}
			vpcConfiguration.SecurityGroups = securityGroups
		}
		// The provider models assign_public_ip as a string enum.
		if ecs.NetworkConfiguration.AssignPublicIp {
			vpcConfiguration.AssignPublicIp = pulumi.String("ENABLED")
		} else {
			vpcConfiguration.AssignPublicIp = pulumi.String("DISABLED")
		}
		args.NetworkConfiguration = &pipes.PipeTargetParametersEcsTaskParametersNetworkConfigurationArgs{
			AwsVpcConfiguration: vpcConfiguration,
		}
	}

	if ecs.Overrides != nil {
		args.Overrides = buildEcsTaskOverrides(ecs.Overrides)
	}

	if len(ecs.PlacementConstraints) > 0 {
		constraints := pipes.PipeTargetParametersEcsTaskParametersPlacementConstraintArray{}
		for _, constraint := range ecs.PlacementConstraints {
			constraintArgs := &pipes.PipeTargetParametersEcsTaskParametersPlacementConstraintArgs{
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
		strategies := pipes.PipeTargetParametersEcsTaskParametersPlacementStrategyArray{}
		for _, strategy := range ecs.PlacementStrategy {
			strategyArgs := &pipes.PipeTargetParametersEcsTaskParametersPlacementStrategyArgs{
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

func buildEcsTaskOverrides(overrides *awseventbridgepipev1alpha1.AwsEventBridgePipeEcsTaskOverrides) *pipes.PipeTargetParametersEcsTaskParametersOverridesArgs {
	args := &pipes.PipeTargetParametersEcsTaskParametersOverridesArgs{}

	if overrides.Cpu != "" {
		args.Cpu = pulumi.String(overrides.Cpu)
	}
	if overrides.Memory != "" {
		args.Memory = pulumi.String(overrides.Memory)
	}
	if overrides.EphemeralStorageSizeInGib != nil {
		args.EphemeralStorage = &pipes.PipeTargetParametersEcsTaskParametersOverridesEphemeralStorageArgs{
			SizeInGib: pulumi.Int(int(*overrides.EphemeralStorageSizeInGib)),
		}
	}
	if overrides.ExecutionRoleArn.GetValue() != "" {
		args.ExecutionRoleArn = pulumi.String(overrides.ExecutionRoleArn.GetValue())
	}
	if overrides.TaskRoleArn.GetValue() != "" {
		args.TaskRoleArn = pulumi.String(overrides.TaskRoleArn.GetValue())
	}

	if len(overrides.ContainerOverrides) > 0 {
		containerOverrides := pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideArray{}
		for _, containerOverride := range overrides.ContainerOverrides {
			overrideArgs := &pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideArgs{}
			if containerOverride.Name != "" {
				overrideArgs.Name = pulumi.String(containerOverride.Name)
			}
			if len(containerOverride.Command) > 0 {
				overrideArgs.Commands = pulumi.ToStringArray(containerOverride.Command)
			}
			if containerOverride.Cpu != nil {
				overrideArgs.Cpu = pulumi.Int(int(*containerOverride.Cpu))
			}
			if containerOverride.Memory != nil {
				overrideArgs.Memory = pulumi.Int(int(*containerOverride.Memory))
			}
			if containerOverride.MemoryReservation != nil {
				overrideArgs.MemoryReservation = pulumi.Int(int(*containerOverride.MemoryReservation))
			}
			if len(containerOverride.Environment) > 0 {
				environments := pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideEnvironmentArray{}
				for _, environment := range containerOverride.Environment {
					environments = append(environments, &pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideEnvironmentArgs{
						Name:  pulumi.String(environment.Name),
						Value: pulumi.String(environment.Value),
					})
				}
				overrideArgs.Environments = environments
			}
			if len(containerOverride.EnvironmentFiles) > 0 {
				environmentFiles := pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideEnvironmentFileArray{}
				for _, environmentFile := range containerOverride.EnvironmentFiles {
					environmentFiles = append(environmentFiles, &pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideEnvironmentFileArgs{
						Type:  pulumi.String(environmentFile.Type),
						Value: pulumi.String(environmentFile.Value),
					})
				}
				overrideArgs.EnvironmentFiles = environmentFiles
			}
			if len(containerOverride.ResourceRequirements) > 0 {
				requirements := pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideResourceRequirementArray{}
				for _, requirement := range containerOverride.ResourceRequirements {
					requirements = append(requirements, &pipes.PipeTargetParametersEcsTaskParametersOverridesContainerOverrideResourceRequirementArgs{
						Type:  pulumi.String(requirement.Type),
						Value: pulumi.String(requirement.Value),
					})
				}
				overrideArgs.ResourceRequirements = requirements
			}
			containerOverrides = append(containerOverrides, overrideArgs)
		}
		args.ContainerOverrides = containerOverrides
	}

	if len(overrides.InferenceAcceleratorOverrides) > 0 {
		accelerators := pipes.PipeTargetParametersEcsTaskParametersOverridesInferenceAcceleratorOverrideArray{}
		for _, accelerator := range overrides.InferenceAcceleratorOverrides {
			acceleratorArgs := &pipes.PipeTargetParametersEcsTaskParametersOverridesInferenceAcceleratorOverrideArgs{}
			if accelerator.DeviceName != "" {
				acceleratorArgs.DeviceName = pulumi.String(accelerator.DeviceName)
			}
			if accelerator.DeviceType != "" {
				acceleratorArgs.DeviceType = pulumi.String(accelerator.DeviceType)
			}
			accelerators = append(accelerators, acceleratorArgs)
		}
		args.InferenceAcceleratorOverrides = accelerators
	}

	return args
}

func buildBatchJobParameters(batch *awseventbridgepipev1alpha1.AwsEventBridgePipeBatchJobTargetParameters) *pipes.PipeTargetParametersBatchJobParametersArgs {
	args := &pipes.PipeTargetParametersBatchJobParametersArgs{
		JobDefinition: pulumi.String(batch.JobDefinition),
		JobName:       pulumi.String(batch.JobName),
	}

	if batch.ArraySize != nil {
		args.ArrayProperties = &pipes.PipeTargetParametersBatchJobParametersArrayPropertiesArgs{
			Size: pulumi.Int(int(*batch.ArraySize)),
		}
	}
	if batch.RetryAttempts != nil {
		args.RetryStrategy = &pipes.PipeTargetParametersBatchJobParametersRetryStrategyArgs{
			Attempts: pulumi.Int(int(*batch.RetryAttempts)),
		}
	}
	if len(batch.Parameters) > 0 {
		args.Parameters = pulumi.ToStringMap(batch.Parameters)
	}

	if len(batch.DependsOn) > 0 {
		dependencies := pipes.PipeTargetParametersBatchJobParametersDependsOnArray{}
		for _, dependency := range batch.DependsOn {
			dependencyArgs := &pipes.PipeTargetParametersBatchJobParametersDependsOnArgs{}
			if dependency.JobId != "" {
				dependencyArgs.JobId = pulumi.String(dependency.JobId)
			}
			if dependency.Type != "" {
				dependencyArgs.Type = pulumi.String(dependency.Type)
			}
			dependencies = append(dependencies, dependencyArgs)
		}
		args.DependsOns = dependencies
	}

	if batch.ContainerOverrides != nil {
		containerOverrides := &pipes.PipeTargetParametersBatchJobParametersContainerOverridesArgs{}
		if len(batch.ContainerOverrides.Command) > 0 {
			containerOverrides.Commands = pulumi.ToStringArray(batch.ContainerOverrides.Command)
		}
		if batch.ContainerOverrides.InstanceType != "" {
			containerOverrides.InstanceType = pulumi.String(batch.ContainerOverrides.InstanceType)
		}
		if len(batch.ContainerOverrides.Environment) > 0 {
			environments := pipes.PipeTargetParametersBatchJobParametersContainerOverridesEnvironmentArray{}
			for _, environment := range batch.ContainerOverrides.Environment {
				environments = append(environments, &pipes.PipeTargetParametersBatchJobParametersContainerOverridesEnvironmentArgs{
					Name:  pulumi.String(environment.Name),
					Value: pulumi.String(environment.Value),
				})
			}
			containerOverrides.Environments = environments
		}
		if len(batch.ContainerOverrides.ResourceRequirements) > 0 {
			requirements := pipes.PipeTargetParametersBatchJobParametersContainerOverridesResourceRequirementArray{}
			for _, requirement := range batch.ContainerOverrides.ResourceRequirements {
				requirements = append(requirements, &pipes.PipeTargetParametersBatchJobParametersContainerOverridesResourceRequirementArgs{
					Type:  pulumi.String(requirement.Type),
					Value: pulumi.String(requirement.Value),
				})
			}
			containerOverrides.ResourceRequirements = requirements
		}
		args.ContainerOverrides = containerOverrides
	}

	return args
}
