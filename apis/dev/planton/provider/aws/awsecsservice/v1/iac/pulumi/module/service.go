package module

import (
	"fmt"

	"github.com/pkg/errors"
	awsecsservicev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecsservice/v1"
	"github.com/plantonhq/planton/internal/valuefrom"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// defaultDesiredCount seeds the running-task count when the spec leaves
// desired_count unset.
const defaultDesiredCount = 1

// ecsService creates the aws_ecs_service resource. The service only
// SCHEDULES: the task definition, cluster, target groups, subnets, and
// security groups are all referenced resources resolved before this module
// runs -- the module never creates or mutates a resource it merely
// references.
func ecsService(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*ecs.Service, error) {
	spec := locals.AwsEcsService.Spec
	serviceName := locals.AwsEcsService.Metadata.Name

	args := &ecs.ServiceArgs{
		Name:           pulumi.String(serviceName),
		Cluster:        pulumi.String(spec.ClusterArn.GetValue()),
		TaskDefinition: pulumi.String(spec.TaskDefinition.GetValue()),
		Tags:           pulumi.ToStringMap(locals.AwsTags),
	}

	// desired_count is a tri-state: unset seeds 1, explicit 0 deploys the
	// wiring with nothing running. When autoscaling owns the count the
	// scaler adjusts it after this initial seed (see the ignore note below).
	desiredCount := int32(defaultDesiredCount)
	if spec.DesiredCount != nil {
		desiredCount = *spec.DesiredCount
	}
	args.DesiredCount = pulumi.Int(int(desiredCount))

	// launch_type XOR capacity_provider_strategy is CEL-enforced in the
	// spec; with neither set, the module defaults to FARGATE explicitly
	// rather than inheriting the cluster's default -- the deployed result
	// should never depend on cluster-side mutable state.
	if len(spec.CapacityProviderStrategy) > 0 {
		strategies := make(ecs.ServiceCapacityProviderStrategyArray, 0, len(spec.CapacityProviderStrategy))
		for _, strategy := range spec.CapacityProviderStrategy {
			strategies = append(strategies, &ecs.ServiceCapacityProviderStrategyArgs{
				CapacityProvider: pulumi.String(strategy.CapacityProvider),
				Base:             pulumi.IntPtr(int(strategy.Base)),
				Weight:           pulumi.IntPtr(int(strategy.Weight)),
			})
		}
		args.CapacityProviderStrategies = strategies
	} else if spec.LaunchType != "" {
		args.LaunchType = pulumi.StringPtr(spec.LaunchType)
	} else {
		args.LaunchType = pulumi.StringPtr("FARGATE")
	}

	if spec.PlatformVersion != "" {
		args.PlatformVersion = pulumi.StringPtr(spec.PlatformVersion)
	}

	if spec.SchedulingStrategy != "" {
		args.SchedulingStrategy = pulumi.StringPtr(spec.SchedulingStrategy)
	}

	if spec.Network != nil {
		args.NetworkConfiguration = &ecs.ServiceNetworkConfigurationArgs{
			Subnets:        pulumi.ToStringArray(valuefrom.ToStringArray(spec.Network.Subnets)),
			SecurityGroups: pulumi.ToStringArray(valuefrom.ToStringArray(spec.Network.SecurityGroups)),
			AssignPublicIp: pulumi.BoolPtr(spec.Network.AssignPublicIp),
		}
	}

	// Load-balancer wiring registers task IPs into referenced target
	// groups; the listener/rule that routes INTO those groups is its own
	// first-class resource. AWS requires each group to already be
	// associated with a listener when the service is created.
	if len(spec.LoadBalancers) > 0 {
		loadBalancers := make(ecs.ServiceLoadBalancerArray, 0, len(spec.LoadBalancers))
		for _, loadBalancer := range spec.LoadBalancers {
			entry := &ecs.ServiceLoadBalancerArgs{
				TargetGroupArn: pulumi.StringPtr(loadBalancer.TargetGroupArn.GetValue()),
				ContainerName:  pulumi.String(loadBalancer.ContainerName),
				ContainerPort:  pulumi.Int(int(loadBalancer.ContainerPort)),
			}
			// The blue/green pair: ECS swaps the production listener rule
			// between the two target groups as deployments bake.
			if loadBalancer.AdvancedConfiguration != nil {
				advanced := &ecs.ServiceLoadBalancerAdvancedConfigurationArgs{
					AlternateTargetGroupArn: pulumi.String(loadBalancer.AdvancedConfiguration.AlternateTargetGroupArn.GetValue()),
					ProductionListenerRule:  pulumi.String(loadBalancer.AdvancedConfiguration.ProductionListenerRule.GetValue()),
					RoleArn:                 pulumi.String(loadBalancer.AdvancedConfiguration.RoleArn.GetValue()),
				}
				if loadBalancer.AdvancedConfiguration.TestListenerRule.GetValue() != "" {
					advanced.TestListenerRule = pulumi.StringPtr(loadBalancer.AdvancedConfiguration.TestListenerRule.GetValue())
				}
				entry.AdvancedConfiguration = advanced
			}
			loadBalancers = append(loadBalancers, entry)
		}
		args.LoadBalancers = loadBalancers

		if spec.HealthCheckGracePeriodSeconds != nil {
			args.HealthCheckGracePeriodSeconds = pulumi.IntPtr(int(*spec.HealthCheckGracePeriodSeconds))
		}
	}

	// Deployment bounds are proto-optional so an explicit value is
	// distinguishable from "let AWS default" (200/100).
	if spec.DeploymentMaximumPercent != nil {
		args.DeploymentMaximumPercent = pulumi.IntPtr(int(*spec.DeploymentMaximumPercent))
	}
	if spec.DeploymentMinimumHealthyPercent != nil {
		args.DeploymentMinimumHealthyPercent = pulumi.IntPtr(int(*spec.DeploymentMinimumHealthyPercent))
	}

	if spec.DeploymentCircuitBreaker != nil {
		args.DeploymentCircuitBreaker = &ecs.ServiceDeploymentCircuitBreakerArgs{
			Enable:   pulumi.Bool(spec.DeploymentCircuitBreaker.Enable),
			Rollback: pulumi.Bool(spec.DeploymentCircuitBreaker.Rollback),
		}
	}

	// Alarm gating watches CloudWatch alarms BY NAME during deployments --
	// the referenced AwsCloudwatchAlarm nodes publish their names as
	// outputs precisely for consumers like this.
	if spec.Alarms != nil {
		args.Alarms = &ecs.ServiceAlarmsArgs{
			AlarmNames: pulumi.ToStringArray(valuefrom.ToStringArray(spec.Alarms.AlarmNames)),
			Enable:     pulumi.Bool(spec.Alarms.Enable),
			Rollback:   pulumi.Bool(spec.Alarms.Rollback),
		}
	}

	if spec.DeploymentConfiguration != nil {
		deployment := &ecs.ServiceDeploymentConfigurationArgs{}
		if spec.DeploymentConfiguration.Strategy != "" {
			deployment.Strategy = pulumi.StringPtr(spec.DeploymentConfiguration.Strategy)
		}
		// AWS models the bake times as nullable integers; the provider
		// serializes them as strings, so the int fields format here.
		if spec.DeploymentConfiguration.BakeTimeInMinutes != nil {
			deployment.BakeTimeInMinutes = pulumi.StringPtr(fmt.Sprintf("%d", *spec.DeploymentConfiguration.BakeTimeInMinutes))
		}
		if spec.DeploymentConfiguration.CanaryConfiguration != nil {
			canary := &ecs.ServiceDeploymentConfigurationCanaryConfigurationArgs{
				CanaryPercent: pulumi.Float64Ptr(spec.DeploymentConfiguration.CanaryConfiguration.CanaryPercent),
			}
			if spec.DeploymentConfiguration.CanaryConfiguration.CanaryBakeTimeInMinutes != nil {
				canary.CanaryBakeTimeInMinutes = pulumi.StringPtr(fmt.Sprintf("%d", *spec.DeploymentConfiguration.CanaryConfiguration.CanaryBakeTimeInMinutes))
			}
			deployment.CanaryConfiguration = canary
		}
		if spec.DeploymentConfiguration.LinearConfiguration != nil {
			linear := &ecs.ServiceDeploymentConfigurationLinearConfigurationArgs{
				StepPercent: pulumi.Float64Ptr(spec.DeploymentConfiguration.LinearConfiguration.StepPercent),
			}
			if spec.DeploymentConfiguration.LinearConfiguration.StepBakeTimeInMinutes != nil {
				linear.StepBakeTimeInMinutes = pulumi.StringPtr(fmt.Sprintf("%d", *spec.DeploymentConfiguration.LinearConfiguration.StepBakeTimeInMinutes))
			}
			deployment.LinearConfiguration = linear
		}
		if len(spec.DeploymentConfiguration.LifecycleHooks) > 0 {
			hooks := make(ecs.ServiceDeploymentConfigurationLifecycleHookArray, 0, len(spec.DeploymentConfiguration.LifecycleHooks))
			for _, hook := range spec.DeploymentConfiguration.LifecycleHooks {
				hookArgs := &ecs.ServiceDeploymentConfigurationLifecycleHookArgs{
					HookTargetArn:   pulumi.String(hook.HookTargetArn),
					RoleArn:         pulumi.String(hook.RoleArn.GetValue()),
					LifecycleStages: pulumi.ToStringArray(hook.LifecycleStages),
				}
				if hook.HookDetails != "" {
					hookArgs.HookDetails = pulumi.StringPtr(hook.HookDetails)
				}
				hooks = append(hooks, hookArgs)
			}
			deployment.LifecycleHooks = hooks
		}
		args.DeploymentConfiguration = deployment
	}

	if spec.DeploymentController != "" {
		args.DeploymentController = &ecs.ServiceDeploymentControllerArgs{
			Type: pulumi.StringPtr(spec.DeploymentController),
		}
	}

	if spec.ServiceConnect != nil {
		serviceConnect := &ecs.ServiceServiceConnectConfigurationArgs{
			Enabled: pulumi.Bool(spec.ServiceConnect.Enabled),
		}
		if spec.ServiceConnect.Namespace != "" {
			serviceConnect.Namespace = pulumi.StringPtr(spec.ServiceConnect.Namespace)
		}
		if spec.ServiceConnect.LogConfiguration != nil {
			serviceConnect.LogConfiguration = serviceConnectLogConfiguration(spec.ServiceConnect.LogConfiguration)
		}
		if len(spec.ServiceConnect.Services) > 0 {
			services := make(ecs.ServiceServiceConnectConfigurationServiceArray, 0, len(spec.ServiceConnect.Services))
			for _, connectService := range spec.ServiceConnect.Services {
				serviceArgs := &ecs.ServiceServiceConnectConfigurationServiceArgs{
					PortName: pulumi.String(connectService.PortName),
				}
				if connectService.DiscoveryName != "" {
					serviceArgs.DiscoveryName = pulumi.StringPtr(connectService.DiscoveryName)
				}
				if connectService.IngressPortOverride != nil {
					serviceArgs.IngressPortOverride = pulumi.IntPtr(int(*connectService.IngressPortOverride))
				}
				if connectService.ClientAlias != nil {
					alias := &ecs.ServiceServiceConnectConfigurationServiceClientAliasArgs{
						Port: pulumi.Int(int(connectService.ClientAlias.Port)),
					}
					if connectService.ClientAlias.DnsName != "" {
						alias.DnsName = pulumi.StringPtr(connectService.ClientAlias.DnsName)
					}
					serviceArgs.ClientAlias = ecs.ServiceServiceConnectConfigurationServiceClientAliasArray{alias}
				}
				if connectService.Timeout != nil {
					timeout := &ecs.ServiceServiceConnectConfigurationServiceTimeoutArgs{}
					if connectService.Timeout.IdleTimeoutSeconds > 0 {
						timeout.IdleTimeoutSeconds = pulumi.IntPtr(int(connectService.Timeout.IdleTimeoutSeconds))
					}
					if connectService.Timeout.PerRequestTimeoutSeconds > 0 {
						timeout.PerRequestTimeoutSeconds = pulumi.IntPtr(int(connectService.Timeout.PerRequestTimeoutSeconds))
					}
					serviceArgs.Timeout = timeout
				}
				if connectService.Tls != nil {
					tls := &ecs.ServiceServiceConnectConfigurationServiceTlsArgs{
						IssuerCertAuthority: &ecs.ServiceServiceConnectConfigurationServiceTlsIssuerCertAuthorityArgs{
							AwsPcaAuthorityArn: pulumi.String(connectService.Tls.AwsPcaAuthorityArn),
						},
					}
					if connectService.Tls.KmsKey.GetValue() != "" {
						tls.KmsKey = pulumi.StringPtr(connectService.Tls.KmsKey.GetValue())
					}
					if connectService.Tls.RoleArn.GetValue() != "" {
						tls.RoleArn = pulumi.StringPtr(connectService.Tls.RoleArn.GetValue())
					}
					serviceArgs.Tls = tls
				}
				services = append(services, serviceArgs)
			}
			serviceConnect.Services = services
		}
		args.ServiceConnectConfiguration = serviceConnect
	}

	if spec.ServiceRegistries != nil {
		registries := &ecs.ServiceServiceRegistriesArgs{
			RegistryArn: pulumi.String(spec.ServiceRegistries.RegistryArn),
		}
		if spec.ServiceRegistries.ContainerName != "" {
			registries.ContainerName = pulumi.StringPtr(spec.ServiceRegistries.ContainerName)
		}
		if spec.ServiceRegistries.ContainerPort != nil {
			registries.ContainerPort = pulumi.IntPtr(int(*spec.ServiceRegistries.ContainerPort))
		}
		if spec.ServiceRegistries.Port != nil {
			registries.Port = pulumi.IntPtr(int(*spec.ServiceRegistries.Port))
		}
		args.ServiceRegistries = registries
	}

	if spec.VolumeConfiguration != nil {
		volume := spec.VolumeConfiguration.ManagedEbsVolume
		managedEbs := &ecs.ServiceVolumeConfigurationManagedEbsVolumeArgs{
			RoleArn: pulumi.String(volume.RoleArn.GetValue()),
		}
		if volume.SizeInGb > 0 {
			managedEbs.SizeInGb = pulumi.IntPtr(int(volume.SizeInGb))
		}
		if volume.VolumeType != "" {
			managedEbs.VolumeType = pulumi.StringPtr(volume.VolumeType)
		}
		if volume.Iops > 0 {
			managedEbs.Iops = pulumi.IntPtr(int(volume.Iops))
		}
		if volume.Throughput > 0 {
			managedEbs.Throughput = pulumi.IntPtr(int(volume.Throughput))
		}
		if volume.Encrypted != nil {
			managedEbs.Encrypted = pulumi.BoolPtr(*volume.Encrypted)
		}
		if volume.KmsKeyId.GetValue() != "" {
			managedEbs.KmsKeyId = pulumi.StringPtr(volume.KmsKeyId.GetValue())
		}
		if volume.SnapshotId != "" {
			managedEbs.SnapshotId = pulumi.StringPtr(volume.SnapshotId)
		}
		if volume.FileSystemType != "" {
			managedEbs.FileSystemType = pulumi.StringPtr(volume.FileSystemType)
		}
		args.VolumeConfiguration = &ecs.ServiceVolumeConfigurationArgs{
			Name:             pulumi.String(spec.VolumeConfiguration.Name),
			ManagedEbsVolume: managedEbs,
		}
	}

	if len(spec.OrderedPlacementStrategy) > 0 {
		strategies := make(ecs.ServiceOrderedPlacementStrategyArray, 0, len(spec.OrderedPlacementStrategy))
		for _, strategy := range spec.OrderedPlacementStrategy {
			entry := &ecs.ServiceOrderedPlacementStrategyArgs{
				Type: pulumi.String(strategy.Type),
			}
			if strategy.Field != "" {
				entry.Field = pulumi.StringPtr(strategy.Field)
			}
			strategies = append(strategies, entry)
		}
		args.OrderedPlacementStrategies = strategies
	}

	if len(spec.PlacementConstraints) > 0 {
		constraints := make(ecs.ServicePlacementConstraintArray, 0, len(spec.PlacementConstraints))
		for _, constraint := range spec.PlacementConstraints {
			entry := &ecs.ServicePlacementConstraintArgs{
				Type: pulumi.String(constraint.Type),
			}
			if constraint.Expression != "" {
				entry.Expression = pulumi.StringPtr(constraint.Expression)
			}
			constraints = append(constraints, entry)
		}
		args.PlacementConstraints = constraints
	}

	// Unset lets AWS decide (new services default to ENABLED where
	// supported) -- the spec deliberately has no default here because the
	// provider dropped its own.
	if spec.AvailabilityZoneRebalancing != "" {
		args.AvailabilityZoneRebalancing = pulumi.StringPtr(spec.AvailabilityZoneRebalancing)
	}

	if spec.PropagateTags != "" {
		args.PropagateTags = pulumi.StringPtr(spec.PropagateTags)
	}
	if spec.EnableEcsManagedTags {
		args.EnableEcsManagedTags = pulumi.BoolPtr(true)
	}
	if spec.EnableExecuteCommand {
		args.EnableExecuteCommand = pulumi.BoolPtr(true)
	}
	if spec.ForceDelete {
		args.ForceDelete = pulumi.BoolPtr(true)
	}

	// desired_count is runtime state once the service is live: the
	// autoscaler owns it when configured, and operators may scale out of
	// band. The module seeds the initial count and then leaves it alone --
	// matching the house convention for autoscaling-managed counts (the
	// Terraform module ignores desired_count changes the same way).
	createdService, err := ecs.NewService(ctx,
		"service",
		args,
		pulumi.Provider(provider),
		pulumi.IgnoreChanges([]string{"desiredCount"}))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ECS service")
	}

	ctx.Export(OpServiceArn, createdService.Arn)
	ctx.Export(OpServiceName, createdService.Name)
	ctx.Export(OpClusterArn, pulumi.String(spec.ClusterArn.GetValue()))
	ctx.Export(OpTaskDefinitionArn, pulumi.String(spec.TaskDefinition.GetValue()))

	return createdService, nil
}

// serviceConnectLogConfiguration converts the spec's log configuration for
// the Service Connect proxy. Secret options are name -> ARN pairs the ECS
// agent resolves at task start.
func serviceConnectLogConfiguration(logConfiguration *awsecsservicev1.AwsEcsServiceLogConfiguration) *ecs.ServiceServiceConnectConfigurationLogConfigurationArgs {
	args := &ecs.ServiceServiceConnectConfigurationLogConfigurationArgs{
		LogDriver: pulumi.String(logConfiguration.LogDriver),
	}
	if len(logConfiguration.Options) > 0 {
		args.Options = pulumi.ToStringMap(logConfiguration.Options)
	}
	if len(logConfiguration.SecretOptions) > 0 {
		secretOptions := make(ecs.ServiceServiceConnectConfigurationLogConfigurationSecretOptionArray, 0, len(logConfiguration.SecretOptions))
		for name, valueFrom := range logConfiguration.SecretOptions {
			secretOptions = append(secretOptions, &ecs.ServiceServiceConnectConfigurationLogConfigurationSecretOptionArgs{
				Name:      pulumi.String(name),
				ValueFrom: pulumi.String(valueFrom),
			})
		}
		args.SecretOptions = secretOptions
	}
	return args
}
