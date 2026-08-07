package module

import (
	"strconv"

	"github.com/pkg/errors"
	awsautoscalinggroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsautoscalinggroup/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/autoscaling"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// group provisions the auto-scaling group itself. The group is a pure
// orchestrator: WHAT to launch lives in the referenced launch template,
// WHERE traffic comes from lives in the referenced target groups; this
// resource owns how many, where, and when. Only the group name is
// create-only in AWS -- everything else updates in place.
func group(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*autoscaling.Group, error) {
	spec := locals.AwsAutoScalingGroup.Spec

	groupName := locals.AwsAutoScalingGroup.Metadata.Name

	subnetIds := make(pulumi.StringArray, 0, len(spec.Subnets))
	for _, subnet := range spec.Subnets {
		subnetIds = append(subnetIds, pulumi.String(subnet.GetValue()))
	}

	args := &autoscaling.GroupArgs{
		Name:               pulumi.StringPtr(groupName),
		VpcZoneIdentifiers: subnetIds,
		MinSize:            pulumi.Int(int(spec.MinSize)),
		MaxSize:            pulumi.Int(int(spec.MaxSize)),
		// ASG tags are the native key/value/propagate-at-launch triple, not
		// a plain map: propagate_at_launch=true copies the identity tags
		// onto every launched instance, so fleet members never escape
		// cost-allocation and orphan-cleanup queries.
		Tags: groupTags(locals, groupName),
	}

	// Leaving desired_capacity unset lets scaling policies own the number:
	// a literal count here would fight the autoscaler on every apply.
	if spec.DesiredCapacity > 0 {
		args.DesiredCapacity = pulumi.IntPtr(int(spec.DesiredCapacity))
	}
	if spec.DesiredCapacityType != "" {
		args.DesiredCapacityType = pulumi.StringPtr(spec.DesiredCapacityType)
	}

	// Exactly one of launch_template / mixed_instances_policy is set (spec
	// CEL enforces it) -- mirroring AWS's ExactlyOneOf on the same fields.
	if spec.LaunchTemplate != nil {
		args.LaunchTemplate = launchTemplateArgs(spec.LaunchTemplate)
	}
	if spec.MixedInstancesPolicy != nil {
		args.MixedInstancesPolicy = mixedInstancesPolicyArgs(spec.MixedInstancesPolicy)
	}

	if spec.CapacityRebalance {
		args.CapacityRebalance = pulumi.BoolPtr(true)
	}
	if spec.DefaultCooldownSeconds > 0 {
		args.DefaultCooldown = pulumi.IntPtr(int(spec.DefaultCooldownSeconds))
	}
	// default_instance_warmup: 0 is meaningful to AWS ("metrics count
	// immediately"), but indistinguishable from unset in the proto; the
	// modules send it only when positive, matching the TF module.
	if spec.DefaultInstanceWarmupSeconds > 0 {
		args.DefaultInstanceWarmup = pulumi.IntPtr(int(spec.DefaultInstanceWarmupSeconds))
	}
	if spec.HealthCheckType != "" {
		args.HealthCheckType = pulumi.StringPtr(spec.HealthCheckType)
	}
	if spec.HealthCheckGracePeriodSeconds > 0 {
		args.HealthCheckGracePeriod = pulumi.IntPtr(int(spec.HealthCheckGracePeriodSeconds))
	}

	if len(spec.TargetGroups) > 0 {
		targetGroupArns := make(pulumi.StringArray, 0, len(spec.TargetGroups))
		for _, targetGroup := range spec.TargetGroups {
			targetGroupArns = append(targetGroupArns, pulumi.String(targetGroup.GetValue()))
		}
		args.TargetGroupArns = targetGroupArns
	}

	if len(spec.TerminationPolicies) > 0 {
		args.TerminationPolicies = pulumi.ToStringArray(spec.TerminationPolicies)
	}
	if spec.MaxInstanceLifetimeSeconds > 0 {
		args.MaxInstanceLifetime = pulumi.IntPtr(int(spec.MaxInstanceLifetimeSeconds))
	}
	if spec.ProtectFromScaleIn {
		args.ProtectFromScaleIn = pulumi.BoolPtr(true)
	}
	if spec.PlacementGroup != "" {
		args.PlacementGroup = pulumi.String(spec.PlacementGroup)
	}
	if spec.ServiceLinkedRoleArn != "" {
		args.ServiceLinkedRoleArn = pulumi.StringPtr(spec.ServiceLinkedRoleArn)
	}
	if len(spec.EnabledMetrics) > 0 {
		metrics := make(autoscaling.MetricArray, 0, len(spec.EnabledMetrics))
		for _, metric := range spec.EnabledMetrics {
			metrics = append(metrics, autoscaling.Metric(metric))
		}
		args.EnabledMetrics = metrics
	}
	if len(spec.SuspendedProcesses) > 0 {
		args.SuspendedProcesses = pulumi.ToStringArray(spec.SuspendedProcesses)
	}
	if spec.InstanceRefresh != nil {
		args.InstanceRefresh = instanceRefreshArgs(spec.InstanceRefresh)
	}
	if spec.WarmPool != nil {
		args.WarmPool = warmPoolArgs(spec.WarmPool)
	}
	if spec.InstanceMaintenancePolicy != nil {
		args.InstanceMaintenancePolicy = &autoscaling.GroupInstanceMaintenancePolicyArgs{
			MinHealthyPercentage: pulumi.Int(int(spec.InstanceMaintenancePolicy.MinHealthyPercentage)),
			MaxHealthyPercentage: pulumi.Int(int(spec.InstanceMaintenancePolicy.MaxHealthyPercentage)),
		}
	}
	if spec.CapacityDistributionStrategy != "" {
		args.AvailabilityZoneDistribution = &autoscaling.GroupAvailabilityZoneDistributionArgs{
			CapacityDistributionStrategy: pulumi.StringPtr(spec.CapacityDistributionStrategy),
		}
	}
	if spec.ForceDelete {
		args.ForceDelete = pulumi.BoolPtr(true)
	}
	if spec.WaitForCapacityTimeout != "" {
		args.WaitForCapacityTimeout = pulumi.StringPtr(spec.WaitForCapacityTimeout)
	}

	createdGroup, err := autoscaling.NewGroup(ctx, groupName, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auto-scaling group")
	}

	ctx.Export(OpAutoscalingGroupName, createdGroup.Name)
	ctx.Export(OpAutoscalingGroupArn, createdGroup.Arn)

	return createdGroup, nil
}

// groupTags converts the identity-tag map to the ASG-native tag triple with
// propagate_at_launch enabled on every entry.
func groupTags(locals *Locals, groupName string) autoscaling.GroupTagArray {
	tags := autoscaling.GroupTagArray{
		&autoscaling.GroupTagArgs{
			Key:               pulumi.String("Name"),
			Value:             pulumi.String(groupName),
			PropagateAtLaunch: pulumi.Bool(true),
		},
	}
	for key, value := range locals.AwsTags {
		tags = append(tags, &autoscaling.GroupTagArgs{
			Key:               pulumi.String(key),
			Value:             pulumi.String(value),
			PropagateAtLaunch: pulumi.Bool(true),
		})
	}
	return tags
}

// launchTemplateArgs maps the single-template reference. An empty version
// keeps AWS's "$Default" behavior -- the setup that lets a template update
// roll the fleet.
func launchTemplateArgs(ref *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupLaunchTemplateRef) *autoscaling.GroupLaunchTemplateArgs {
	args := &autoscaling.GroupLaunchTemplateArgs{
		Id: pulumi.StringPtr(ref.LaunchTemplateId.GetValue()),
	}
	if ref.Version != "" {
		args.Version = pulumi.StringPtr(ref.Version)
	}
	return args
}

// launchTemplateSpecificationArgs maps a template reference inside the
// mixed-instances policy (base or per-override).
func launchTemplateSpecificationArgs(ref *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupLaunchTemplateRef) *autoscaling.GroupMixedInstancesPolicyLaunchTemplateLaunchTemplateSpecificationArgs {
	args := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateLaunchTemplateSpecificationArgs{
		LaunchTemplateId: pulumi.StringPtr(ref.LaunchTemplateId.GetValue()),
	}
	if ref.Version != "" {
		args.Version = pulumi.StringPtr(ref.Version)
	}
	return args
}

// mixedInstancesPolicyArgs maps the On-Demand/Spot blend: a base template,
// per-type (or attribute-based) overrides, and the purchase-option split.
func mixedInstancesPolicyArgs(policy *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupMixedInstancesPolicy) *autoscaling.GroupMixedInstancesPolicyArgs {
	launchTemplate := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateArgs{
		LaunchTemplateSpecification: launchTemplateSpecificationArgs(policy.LaunchTemplate),
	}

	if len(policy.Overrides) > 0 {
		overrides := make(autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideArray, 0, len(policy.Overrides))
		for _, override := range policy.Overrides {
			overrideArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideArgs{}
			if override.InstanceType != "" {
				overrideArgs.InstanceType = pulumi.StringPtr(override.InstanceType)
			}
			// weighted_capacity is a string at AWS despite being numeric
			// (1-999); the proto keeps the honest int and converts here.
			if override.WeightedCapacity > 0 {
				overrideArgs.WeightedCapacity = pulumi.StringPtr(strconv.Itoa(int(override.WeightedCapacity)))
			}
			if override.LaunchTemplate != nil {
				overrideArgs.LaunchTemplateSpecification = &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideLaunchTemplateSpecificationArgs{
					LaunchTemplateId: pulumi.StringPtr(override.LaunchTemplate.LaunchTemplateId.GetValue()),
					Version:          overrideVersion(override.LaunchTemplate),
				}
			}
			if override.InstanceRequirements != nil {
				overrideArgs.InstanceRequirements = overrideInstanceRequirementsArgs(override.InstanceRequirements)
			}
			overrides = append(overrides, overrideArgs)
		}
		launchTemplate.Overrides = overrides
	}

	args := &autoscaling.GroupMixedInstancesPolicyArgs{
		LaunchTemplate: launchTemplate,
	}

	if policy.InstancesDistribution != nil {
		distribution := &autoscaling.GroupMixedInstancesPolicyInstancesDistributionArgs{}
		if policy.InstancesDistribution.OnDemandAllocationStrategy != "" {
			distribution.OnDemandAllocationStrategy = pulumi.StringPtr(policy.InstancesDistribution.OnDemandAllocationStrategy)
		}
		if policy.InstancesDistribution.OnDemandBaseCapacity > 0 {
			distribution.OnDemandBaseCapacity = pulumi.IntPtr(int(policy.InstancesDistribution.OnDemandBaseCapacity))
		}
		// Explicit 0 means all-Spot above the base -- the aggressive cost
		// posture -- so presence (not zero-ness) decides whether it is sent.
		if policy.InstancesDistribution.OnDemandPercentageAboveBaseCapacity != nil {
			distribution.OnDemandPercentageAboveBaseCapacity = pulumi.IntPtr(int(policy.InstancesDistribution.GetOnDemandPercentageAboveBaseCapacity()))
		}
		if policy.InstancesDistribution.SpotAllocationStrategy != "" {
			distribution.SpotAllocationStrategy = pulumi.StringPtr(policy.InstancesDistribution.SpotAllocationStrategy)
		}
		if policy.InstancesDistribution.SpotInstancePools > 0 {
			distribution.SpotInstancePools = pulumi.IntPtr(int(policy.InstancesDistribution.SpotInstancePools))
		}
		if policy.InstancesDistribution.SpotMaxPrice != "" {
			distribution.SpotMaxPrice = pulumi.StringPtr(policy.InstancesDistribution.SpotMaxPrice)
		}
		args.InstancesDistribution = distribution
	}

	return args
}

func overrideVersion(ref *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupLaunchTemplateRef) pulumi.StringPtrInput {
	if ref.Version != "" {
		return pulumi.StringPtr(ref.Version)
	}
	return nil
}

// overrideInstanceRequirementsArgs maps attribute-based selection for one
// mixed-instances override. Only set fields are sent so AWS's own defaults
// (e.g. bare metal excluded) keep applying.
func overrideInstanceRequirementsArgs(requirements *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupInstanceRequirements) *autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsArgs {
	args := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsArgs{}

	// memory_mib and vcpu_count are the two AWS-required dimensions; the
	// spec enforces their presence (and that min is set).
	memoryArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsMemoryMibArgs{
		Min: pulumi.IntPtr(int(requirements.MemoryMib.Min)),
	}
	if requirements.MemoryMib.Max > 0 {
		memoryArgs.Max = pulumi.IntPtr(int(requirements.MemoryMib.Max))
	}
	args.MemoryMib = memoryArgs

	vcpuArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsVcpuCountArgs{
		Min: pulumi.IntPtr(int(requirements.VcpuCount.Min)),
	}
	if requirements.VcpuCount.Max > 0 {
		vcpuArgs.Max = pulumi.IntPtr(int(requirements.VcpuCount.Max))
	}
	args.VcpuCount = vcpuArgs

	if len(requirements.AllowedInstanceTypes) > 0 {
		args.AllowedInstanceTypes = pulumi.ToStringArray(requirements.AllowedInstanceTypes)
	}
	if len(requirements.ExcludedInstanceTypes) > 0 {
		args.ExcludedInstanceTypes = pulumi.ToStringArray(requirements.ExcludedInstanceTypes)
	}
	if len(requirements.InstanceGenerations) > 0 {
		args.InstanceGenerations = pulumi.ToStringArray(requirements.InstanceGenerations)
	}
	if len(requirements.CpuManufacturers) > 0 {
		args.CpuManufacturers = pulumi.ToStringArray(requirements.CpuManufacturers)
	}
	if requirements.BareMetal != "" {
		args.BareMetal = pulumi.StringPtr(requirements.BareMetal)
	}
	if requirements.BurstablePerformance != "" {
		args.BurstablePerformance = pulumi.StringPtr(requirements.BurstablePerformance)
	}
	if requirements.RequireHibernateSupport {
		args.RequireHibernateSupport = pulumi.BoolPtr(true)
	}
	if requirements.SpotMaxPricePercentageOverLowestPrice > 0 {
		args.SpotMaxPricePercentageOverLowestPrice = pulumi.IntPtr(int(requirements.SpotMaxPricePercentageOverLowestPrice))
	}
	if requirements.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice > 0 {
		args.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice = pulumi.IntPtr(int(requirements.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice))
	}
	if requirements.OnDemandMaxPricePercentageOverLowestPrice > 0 {
		args.OnDemandMaxPricePercentageOverLowestPrice = pulumi.IntPtr(int(requirements.OnDemandMaxPricePercentageOverLowestPrice))
	}
	if requirements.LocalStorage != "" {
		args.LocalStorage = pulumi.StringPtr(requirements.LocalStorage)
	}
	if len(requirements.LocalStorageTypes) > 0 {
		args.LocalStorageTypes = pulumi.ToStringArray(requirements.LocalStorageTypes)
	}
	if requirements.TotalLocalStorageGb != nil {
		storageArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsTotalLocalStorageGbArgs{}
		if requirements.TotalLocalStorageGb.Min > 0 {
			storageArgs.Min = pulumi.Float64Ptr(requirements.TotalLocalStorageGb.Min)
		}
		if requirements.TotalLocalStorageGb.Max > 0 {
			storageArgs.Max = pulumi.Float64Ptr(requirements.TotalLocalStorageGb.Max)
		}
		args.TotalLocalStorageGb = storageArgs
	}
	if requirements.MemoryGibPerVcpu != nil {
		ratioArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsMemoryGibPerVcpuArgs{}
		if requirements.MemoryGibPerVcpu.Min > 0 {
			ratioArgs.Min = pulumi.Float64Ptr(requirements.MemoryGibPerVcpu.Min)
		}
		if requirements.MemoryGibPerVcpu.Max > 0 {
			ratioArgs.Max = pulumi.Float64Ptr(requirements.MemoryGibPerVcpu.Max)
		}
		args.MemoryGibPerVcpu = ratioArgs
	}
	if requirements.NetworkInterfaceCount != nil {
		countArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsNetworkInterfaceCountArgs{}
		if requirements.NetworkInterfaceCount.Min > 0 {
			countArgs.Min = pulumi.IntPtr(int(requirements.NetworkInterfaceCount.Min))
		}
		if requirements.NetworkInterfaceCount.Max > 0 {
			countArgs.Max = pulumi.IntPtr(int(requirements.NetworkInterfaceCount.Max))
		}
		args.NetworkInterfaceCount = countArgs
	}
	if requirements.NetworkBandwidthGbps != nil {
		bandwidthArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsNetworkBandwidthGbpsArgs{}
		if requirements.NetworkBandwidthGbps.Min > 0 {
			bandwidthArgs.Min = pulumi.Float64Ptr(requirements.NetworkBandwidthGbps.Min)
		}
		if requirements.NetworkBandwidthGbps.Max > 0 {
			bandwidthArgs.Max = pulumi.Float64Ptr(requirements.NetworkBandwidthGbps.Max)
		}
		args.NetworkBandwidthGbps = bandwidthArgs
	}
	if requirements.BaselineEbsBandwidthMbps != nil {
		ebsBandwidthArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsBaselineEbsBandwidthMbpsArgs{}
		if requirements.BaselineEbsBandwidthMbps.Min > 0 {
			ebsBandwidthArgs.Min = pulumi.IntPtr(int(requirements.BaselineEbsBandwidthMbps.Min))
		}
		if requirements.BaselineEbsBandwidthMbps.Max > 0 {
			ebsBandwidthArgs.Max = pulumi.IntPtr(int(requirements.BaselineEbsBandwidthMbps.Max))
		}
		args.BaselineEbsBandwidthMbps = ebsBandwidthArgs
	}
	if requirements.AcceleratorCount != nil {
		acceleratorArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsAcceleratorCountArgs{}
		if requirements.AcceleratorCount.Min > 0 {
			acceleratorArgs.Min = pulumi.IntPtr(int(requirements.AcceleratorCount.Min))
		}
		if requirements.AcceleratorCount.Max > 0 {
			acceleratorArgs.Max = pulumi.IntPtr(int(requirements.AcceleratorCount.Max))
		}
		args.AcceleratorCount = acceleratorArgs
	}
	if len(requirements.AcceleratorManufacturers) > 0 {
		args.AcceleratorManufacturers = pulumi.ToStringArray(requirements.AcceleratorManufacturers)
	}
	if len(requirements.AcceleratorNames) > 0 {
		args.AcceleratorNames = pulumi.ToStringArray(requirements.AcceleratorNames)
	}
	if len(requirements.AcceleratorTypes) > 0 {
		args.AcceleratorTypes = pulumi.ToStringArray(requirements.AcceleratorTypes)
	}
	if requirements.AcceleratorTotalMemoryMib != nil {
		acceleratorMemoryArgs := &autoscaling.GroupMixedInstancesPolicyLaunchTemplateOverrideInstanceRequirementsAcceleratorTotalMemoryMibArgs{}
		if requirements.AcceleratorTotalMemoryMib.Min > 0 {
			acceleratorMemoryArgs.Min = pulumi.IntPtr(int(requirements.AcceleratorTotalMemoryMib.Min))
		}
		if requirements.AcceleratorTotalMemoryMib.Max > 0 {
			acceleratorMemoryArgs.Max = pulumi.IntPtr(int(requirements.AcceleratorTotalMemoryMib.Max))
		}
		args.AcceleratorTotalMemoryMib = acceleratorMemoryArgs
	}

	return args
}

// instanceRefreshArgs maps the rolling-replacement behavior that turns a
// launch-template change into a zero-downtime fleet rollout.
func instanceRefreshArgs(refresh *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupInstanceRefresh) *autoscaling.GroupInstanceRefreshArgs {
	args := &autoscaling.GroupInstanceRefreshArgs{
		Strategy: pulumi.String(refresh.Strategy),
	}
	if len(refresh.Triggers) > 0 {
		args.Triggers = pulumi.ToStringArray(refresh.Triggers)
	}
	if refresh.Preferences != nil {
		preferences := &autoscaling.GroupInstanceRefreshPreferencesArgs{}
		// Explicit 0 ("replace the whole fleet at once") is meaningful and
		// distinct from unset (AWS default 90), so presence decides.
		if refresh.Preferences.MinHealthyPercentage != nil {
			preferences.MinHealthyPercentage = pulumi.IntPtr(int(refresh.Preferences.GetMinHealthyPercentage()))
		}
		if refresh.Preferences.MaxHealthyPercentage > 0 {
			preferences.MaxHealthyPercentage = pulumi.IntPtr(int(refresh.Preferences.MaxHealthyPercentage))
		}
		// The provider models warmup and checkpoint delay as strings
		// (nullable ints at AWS); the proto keeps honest ints and converts.
		if refresh.Preferences.InstanceWarmupSeconds > 0 {
			preferences.InstanceWarmup = pulumi.StringPtr(strconv.Itoa(int(refresh.Preferences.InstanceWarmupSeconds)))
		}
		if len(refresh.Preferences.CheckpointPercentages) > 0 {
			checkpoints := make(pulumi.IntArray, 0, len(refresh.Preferences.CheckpointPercentages))
			for _, percentage := range refresh.Preferences.CheckpointPercentages {
				checkpoints = append(checkpoints, pulumi.Int(int(percentage)))
			}
			preferences.CheckpointPercentages = checkpoints
		}
		if refresh.Preferences.CheckpointDelaySeconds > 0 {
			preferences.CheckpointDelay = pulumi.StringPtr(strconv.Itoa(int(refresh.Preferences.CheckpointDelaySeconds)))
		}
		if refresh.Preferences.AutoRollback {
			preferences.AutoRollback = pulumi.BoolPtr(true)
		}
		if len(refresh.Preferences.Alarms) > 0 {
			alarms := make(pulumi.StringArray, 0, len(refresh.Preferences.Alarms))
			for _, alarm := range refresh.Preferences.Alarms {
				alarms = append(alarms, pulumi.String(alarm.GetValue()))
			}
			preferences.AlarmSpecification = &autoscaling.GroupInstanceRefreshPreferencesAlarmSpecificationArgs{
				Alarms: alarms,
			}
		}
		if refresh.Preferences.ScaleInProtectedInstances != "" {
			preferences.ScaleInProtectedInstances = pulumi.StringPtr(refresh.Preferences.ScaleInProtectedInstances)
		}
		if refresh.Preferences.StandbyInstances != "" {
			preferences.StandbyInstances = pulumi.StringPtr(refresh.Preferences.StandbyInstances)
		}
		if refresh.Preferences.SkipMatching {
			preferences.SkipMatching = pulumi.BoolPtr(true)
		}
		args.Preferences = preferences
	}
	return args
}

// warmPoolArgs maps the pre-initialized capacity pool.
func warmPoolArgs(pool *awsautoscalinggroupv1alpha1.AwsAutoScalingGroupWarmPool) *autoscaling.GroupWarmPoolArgs {
	args := &autoscaling.GroupWarmPoolArgs{}
	if pool.PoolState != "" {
		args.PoolState = pulumi.StringPtr(pool.PoolState)
	}
	if pool.MinSize > 0 {
		args.MinSize = pulumi.IntPtr(int(pool.MinSize))
	}
	// Explicit 0 is meaningful (no prepared capacity beyond min_size), so
	// presence (not zero-ness) decides whether it is sent.
	if pool.MaxGroupPreparedCapacity != nil {
		args.MaxGroupPreparedCapacity = pulumi.IntPtr(int(pool.GetMaxGroupPreparedCapacity()))
	}
	if pool.ReuseOnScaleIn {
		args.InstanceReusePolicy = &autoscaling.GroupWarmPoolInstanceReusePolicyArgs{
			ReuseOnScaleIn: pulumi.BoolPtr(true),
		}
	}
	return args
}
